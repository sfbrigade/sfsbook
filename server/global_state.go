package server

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/blevesearch/bleve"
	"github.com/sfbrigade/sfsbook/dba"
	"github.com/sfbrigade/sfsbook/dba/fieldmap"
	"github.com/sfbrigade/sfsbook/setup"
)

// HandlerFactory contains all state needed to construct the various
// specialized http.Handler instances provided by the server. Once
// a HandlerFactory exists, it can vend handlers without errors.
type HandlerFactory struct {
	statepath string
	sitedir   string

	// Cache when I need one.

	// Databases
	resourceguide bleve.Index
	passwordfile  bleve.Index

	cookietool *cookieTooling

	// Flags
	Immutable bool
}

// MakeHandlerFactory does possibly error-generating setup for all
// asepcts of the global state including the contents of the persistentroot`/state/'
// directory, the database connections, the global cache and cookie
// authentication keys. The provided HandlerFactory can then vends various
// specialized http.Handler instances using this state without errors.
func MakeHandlerFactory(persistentroot string) (*HandlerFactory, error) {
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

	cookietool, err := makeCookieTooling(statepath)
	if err != nil {
		return nil, err
	}

	// This is unnecessary. The auth scheme should take care of it.
	immutable := false
	passwordfile, err := dba.OpenBleve(persistentroot, fieldmap.PasswordFile)
	if err != nil {
		log.Println("Operating in read-only mode because can't open/create user database:", err)
		immutable = true
	}

	return &HandlerFactory{
		statepath:     statepath,
		sitedir:       sitedir,
		resourceguide: resourceguide,
		passwordfile:  passwordfile,
		Immutable:     immutable,
		cookietool:    cookietool,
	}, nil
}
