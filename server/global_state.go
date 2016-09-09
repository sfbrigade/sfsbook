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

	// Cache when I need one.

	// Databases
	ResourceGuide bleve.Index
	PasswordFile  bleve.Index

	UserState

	// Flags
	Immutable bool
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

	// This is unnecessary. The auth scheme should take care of it.
	immutable := false
	passwordfile, err := dba.OpenBleve(persistentroot, fieldmap.PasswordFile)
	if err != nil {
		log.Println("Operating in read-only mode because can't open/create user database:", err)
		immutable = true
	}

	userstate, err := MakeUserState(statepath)
	if err != nil {
		return nil, err
	}

	return &GlobalState{
		EmbeddableResources: *MakeEmbeddableResource(sitedir),
		UserState:        *userstate,
		ResourceGuide:       resourceguide,
		PasswordFile:        passwordfile,
		Immutable:           immutable,
	}, nil
}
