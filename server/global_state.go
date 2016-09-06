package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/blevesearch/bleve"
	"github.com/gorilla/securecookie"
	"github.com/sfbrigade/sfsbook/dba"
	"github.com/sfbrigade/sfsbook/dba/fieldmap"
	"github.com/sfbrigade/sfsbook/setup"
)

// GlobalState is state shared across all server requests.
type GlobalState struct {
	EmbeddableResources

	// Cache when I need one.

	// Databases
	ResourceGuide bleve.Index
	PasswordFile  bleve.Index

	// Cookie management.
	securecookie.SecureCookie

	// Flags
	Immutable bool
}

// makeCookie builds and saves a cookie.
// TODO(rjk): Add automatic cookie rotation with aging and batches.
func makeCookie(statepath, cookiename string) ([]byte, error) {
	path := filepath.Join(statepath, cookiename)
	key, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println("making new cookie", cookiename)
		key := securecookie.GenerateRandomKey(32)
		if key == nil {
			return nil, fmt.Errorf("No cookie for %s and can't make one", cookiename)
		}

		// TODO(rjk): Make sure that the umask is set appropriately.
		cookiefile, err := os.Create(path)
		if err != nil {
			return nil, fmt.Errorf("Can't create a %s to hold new cookie: %v",
				path, err)
		}

		if n, err := cookiefile.Write(key); err != nil || n != len(key) {
			return nil, fmt.Errorf("Can't write new cookie %s.  len is %d instead of %d or error: %v",
				path, n, len(key), err)
		}
	}
	return key, nil
}

// MakeGlobalRequestState builds all the global state shared between all
// requests including the contents of the persistentroot`/state/'
// directory, the database connections, the global cache and cookie
// authentication keys.
func MakeGlobalState(persistentroot string) (*GlobalState, error) {
	statepath := filepath.Join(persistentroot, "state")
	log.Println("hello from setup, creating state in", statepath)

	if err := os.MkdirAll(statepath, 0777); err != nil {
		return nil, fmt.Errorf("Couldn't make directory", statepath, "because", err)
	}

	// It is unnecessary to create the site directory because if it doesn't exist,
	sitedir := filepath.Join(persistentroot, "site")
	if _, err := os.Stat(sitedir); err != nil {
		log.Println("There is no site directory so all resources must be embedded.")
		sitedir = ""
	}

	// make keys
	if err := setup.MakeKeys(statepath); err != nil {
		return nil, fmt.Errorf("Don't have and can't make keys.", err)
	}

	resourceguide, err := dba.OpenBleve(persistentroot, fieldmap.RefGuide)
	if err != nil {
		return nil, fmt.Errorf("Can't open/create the resource guide database: %v", err)
	}

	immutable := false
	passwordfile, err := dba.OpenBleve(persistentroot, fieldmap.PasswordFile)
	if err != nil {
		log.Println("Operating in read-only mode because can't open/create user database:", err)
		immutable = true
	}

	// Make cookie keys.
	hashKey, err := makeCookie(statepath, "hashkey.dat")
	if err != nil {
		return nil, err
	}
	blockKey, err := makeCookie(statepath, "blockkey.dat")
	if err != nil {
		return nil, err
	}

	return &GlobalState{
		EmbeddableResources: *MakeEmbeddableResource(sitedir),
		SecureCookie:        *securecookie.New(hashKey, blockKey),
		ResourceGuide:       resourceguide,
		PasswordFile:        passwordfile,
		Immutable:           immutable,
	}, nil
}
