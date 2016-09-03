package setup

// Copied from the Go 1.6.2 source tree. Builds keys for testing.

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Generate a self-signed X.509 certificate for a TLS server. Outputs to
// 'cert.pem' and 'key.pem' and will overwrite existing files.

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}

func pemBlockForKey(priv interface{}) (*pem.Block, error) {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}, nil
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			return nil, fmt.Errorf("Unable to marshal ECDSA private key: %v", err)
		}
		return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}, nil
	default:
		return nil, nil
	}
}

func makeTestKeys(certpth, keypth string) error {
	var priv interface{}
	var err error
	priv, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to generate private key: %s", err)
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return fmt.Errorf("failed to generate serial number: %s", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"sfsbook.org"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// This only generates a key for testing. A real key will need to be generated and
	// managed separately.
	hosts := strings.Split("localhost", ",")
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	template.IsCA = true
	template.KeyUsage |= x509.KeyUsageCertSign

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(priv), priv)
	if err != nil {
		return fmt.Errorf("Failed to create certificate: %s", err)
	}

	certOut, err := os.Create(certpth)
	if err != nil {
		return fmt.Errorf("failed to open %s for writing: %s", certpth, err)
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()

	//	log.Println("written", certpth)

	keyOut, err := os.OpenFile(keypth, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to open %s for writing:", keypth, err)
	}
	defer keyOut.Close()

	blk, err := pemBlockForKey(priv)
	if err != nil {
		return err
	}
	pem.Encode(keyOut, blk)
	//	log.Println("written", keypth)
	return nil
}

// Code here down is not from the Go distribution.

// MakeKeys creates the https keys in the persistent server root
// specificed by pth..
func MakeKeys(pth string) error {
	certpth := filepath.Join(pth, "cert.pem")
	keypth := filepath.Join(pth, "key.pem")

	makeNew := false
	if _, err := os.Stat(certpth); err != nil {
		makeNew = true
	}
	if _, err := os.Stat(keypth); err != nil {
		makeNew = true
	}

	if makeNew {
		return makeTestKeys(certpth, keypth)
	}
	return nil
}
