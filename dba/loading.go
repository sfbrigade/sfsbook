package dba

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/pborman/uuid"
)

// OpenBleve opens the backing database or builds it if it doesn't exist.
func OpenBleve(persistentroot string) (bleve.Index, error) {
	dbpath := filepath.Join(persistentroot, "state", "sfsbook.bleve")
	bi, err := bleve.Open(dbpath)
	if err == bleve.ErrorIndexPathDoesNotExist {
		// TODO(rjkroege): Might consider making the path configurable. Or something.
		// At present, tries to be smart for development. I'll worry about deployment
		// later.

		log.Printf("Indexing the provided datafile...")
		// create a mapping
		indexMapping, err := buildIndexMapping()
		if err != nil {
			return nil, err
		}
		bi, err = bleve.New(dbpath, indexMapping)
		if err != nil {
			return nil, err
		}

		if err = indexDatabase(bi, persistentroot); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	return bi, nil
}

const sourcefile = "refguide.json"

func indexDatabase(i bleve.Index, pathroot string) error {
	log.Printf("Indexing...")

	jsonBytes, err := ioutil.ReadFile(filepath.Join(pathroot, sourcefile))
	if err != nil {
		return err
	}

	// parse bytes as json
	var parsedResources []map[string]interface{}
	err = json.Unmarshal(jsonBytes, &parsedResources)
	if err != nil {
		return err
	}

	// So: how do I maintain flexibility in the handling of the fields?
	// Can unmarshal into a map of interface{}
	// I can set reasonable defaults.
	// Documents can have sub-documents...

	batch := i.NewBatch()
	for _, r := range parsedResources {
		rid := uuid.NewRandom().String()
		r["reviewed"] = false
		// This can be adapted to specify different types.
		r["_type"] = "resource"
		r["date_indexed"] = time.Now()
		batch.Index(rid, r)
	}

	err = i.Batch(batch)
	if err != nil {
		return err
	}
	log.Println("done Indexing...")

	// TODO(rjk): add some comments and default password setup here.

	return nil
}
