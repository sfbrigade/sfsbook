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

// GlobalState is state shared across all server requests.
type GlobalState struct {
	EmbeddableResources

	// Why do I need? Only for setup?
	Persistentroot string

	// Cache.

	// Databases
	ResourceGuide bleve.Index
	PasswordFile  bleve.Index

	// Cookie keys

	// Flags
	Immutable bool
}

// MakeGlobalRequestState builds all the global state shared between all
// requests including the contents of the persistentroot`/state/'
// directory, the database connections, the global cache and cookie
// authentication keys.
func MakeGlobalState(persistentroot string) (*GlobalState, error) {
	pth := filepath.Join(persistentroot, "state")
	log.Println("hello from setup, creating state in", pth)

	if err := os.MkdirAll(pth, 0777); err != nil {
		return nil, fmt.Errorf("Couldn't make directory", pth, "because", err)
	}

	// It is unnecessary to create the site directory because if it doesn't exist,
	sitedir := filepath.Join(persistentroot, "site")
	if _, err := os.Stat(sitedir); err != nil {
		log.Println("There is no site directory so all resources must be embedded.")
		sitedir = ""
	}

	// make keys
	if err := setup.MakeKeys(pth); err != nil {
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

	// TODO(rjk): Setup cookies. Setup the global cache.

	return &GlobalState{
		EmbeddableResources: *MakeEmbeddableResource(sitedir),
		Persistentroot: persistentroot,
		ResourceGuide:  resourceguide,
		PasswordFile:   passwordfile,
		Immutable:      immutable,
	}, nil
}
