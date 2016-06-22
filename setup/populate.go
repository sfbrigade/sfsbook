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
	log.Println("hello from setup, dumping stuff to", persistentroot)
	pth := filepath.Join(persistentroot, "state")

	// make key
	if err := MakeKeys(pth); err != nil {
		log.Fatalln("Don't have and can't make keys.", err)
	}

	// TODO(rjk): Re-write or setup a database.

}
