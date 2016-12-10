package setup

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// HandlerFactory contains all state needed to construct the various
// specialized http.Handler instances provided by the server. Once
// a HandlerFactory exists, it can vend handlers without errors.
type CertFactory struct {
	statepath string
}

func MakeCertFactory(persistentroot string) (*CertFactory, error) {
	statepath := filepath.Join(persistentroot, "state")
	log.Println("hello from setup, creating state in", statepath)

	if err := os.MkdirAll(statepath, 0777); err != nil {
		return nil, fmt.Errorf("Couldn't make directory", statepath, "because", err)
	}

	if err := MakeKeys(statepath); err != nil {
		return nil, fmt.Errorf("Don't have and can't make keys.", err)
	}

	return &CertFactory{
		statepath: statepath,
	}, nil
}

func (cf *CertFactory) GetCertfFileName() string {
	return filepath.Join(cf.statepath, "cert.pem")
}

func (cf *CertFactory) GetKeyFileName() string {
	return filepath.Join(cf.statepath, "key.pem")
}
