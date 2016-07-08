package setup

import (
	"log"
	"os"
	"path/filepath"
)

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

// ConstructNecessaryStartingState builds all the necessary state to get started,
// placing it in the persistentpath/"state"
func ConstructNecessaryStartingState(persistentroot string) {
	pth := filepath.Join(persistentroot, "state")
	log.Println("hello from setup, creating state in", pth)

	if err := os.MkdirAll(pth, 0777); err != nil {
		log.Fatalln("Couldn't make directory", pth, "because", err)
	}

	// make key
	if err := MakeKeys(pth); err != nil {
		log.Fatalln("Don't have and can't make keys.", err)
	}

	// TODO(rjk): Re-write or setup a database.

}
