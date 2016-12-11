package setup

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/crypto/acme/autocert"
)

var use_acme = flag.Bool("use_acme", false, "Obtain a real certificate from Let's Encrypt")

// CertFactory contains all state needed to construct a certificate scheme
// for the application.
type CertFactory struct {
	statepath string
	autocert  *autocert.Manager
}

// MakeCertFactory manfactures a CertFactory object that maintains the necessary
// state to manage certificates.
func MakeCertFactory(persistentroot string) (*CertFactory, error) {
	statepath := filepath.Join(persistentroot, "state")
	log.Println("hello from setup, creating state in", statepath)

	if err := os.MkdirAll(statepath, 0777); err != nil {
		return nil, fmt.Errorf("Couldn't make directory", statepath, "because", err)
	}

	var m *autocert.Manager
	if *use_acme {
		log.Println("using acme")
		m = &autocert.Manager{
			// TODO(rjk): I need to not make this hard-coded.
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist("mando.liqui.org"),
			Cache: autocert.DirCache(filepath.Join(statepath, "certcache")),
		}
	} else {
		if err := MakeKeys(statepath); err != nil {
			return nil, fmt.Errorf("Don't have and can't make keys.", err)
		}
	}

	return &CertFactory{
		statepath: statepath,
		autocert:  m,
	}, nil
}

func (cf *CertFactory) GetCertfFileName() string {
	if cf.autocert != nil {
		return ""
	}
	return filepath.Join(cf.statepath, "cert.pem")
}

func (cf *CertFactory) GetKeyFileName() string {
	if cf.autocert != nil {
		return ""
	}
	return filepath.Join(cf.statepath, "key.pem")
}

func (cf *CertFactory) GetTLSConfig() *tls.Config {
	if cf.autocert != nil {
		return &tls.Config{GetCertificate: cf.autocert.GetCertificate}
	}
	return nil
}
