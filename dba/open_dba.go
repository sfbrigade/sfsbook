package dba

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/pborman/uuid"
)

// IndexFactory specifies how to name and populate a database.
type IndexFactory interface {
	// The file name for the database.
	Name() string

	// LoadStartData populates a newly-created empty database from
	// pre-existing data.
	LoadStartData(idx bleve.Index, root string) error

	// Mapping returns IndexMapping for this database.
	Mapping() *bleve.IndexMapping
}

// OpenBleve opens the backing database or builds it if it doesn't exist.
func OpenBleve(persistentroot string, dxf IndexFactory) (bleve.Index, error) {
	dbpath := filepath.Join(persistentroot, "state", dxf.Name())
	bi, err := bleve.Open(dbpath)
	if err == bleve.ErrorIndexPathDoesNotExist {
		// TODO(rjkroege): Might consider making the path configurable. Or something.
		// At present, tries to be smart for development. I'll worry about deployment
		// later.

		log.Printf("Indexing the provided datafile...")
		bi, err = bleve.New(dbpath, dxf.Mapping())
		if err != nil {
			goto cleanup
		}

		if err = dxf.LoadStartData(bi, persistentroot); err != nil {
			goto cleanup
		}
	} else if err != nil {
		goto cleanup
	}
	return bi, nil
cleanup:
	os.RemoveAll(dbpath)
	return nil, err
}

const sourcefile = "refguide.json"

type RefGuideType string

func (g RefGuideType) Name() string {
	return string(g)
}

func (_ RefGuideType) Mapping() *bleve.IndexMapping {
	return allDocumentMapping(IndexDocumentMap{
		"resource": buildResourceDocumentMapping(),
	})
}

var RefGuide = RefGuideType("sfsbook.bleve")

func (_ RefGuideType) LoadStartData(i bleve.Index, pathroot string) error {
	log.Println("Indexing... now")

	jsonBytes, err := ioutil.ReadFile(filepath.Join(pathroot, sourcefile))
	if err != nil {
		return err
	}

	log.Println("read the database")

	// parse bytes as json
	var parsedResources []map[string]interface{}
	err = json.Unmarshal(jsonBytes, &parsedResources)
	if err != nil {
		return err
	}

	log.Println("parsed the database record")

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

	log.Println("built a batch")

	err = i.Batch(batch)
	if err != nil {
		return err
	}
	log.Println("done Indexing...")

	// TODO(rjk): add some comments and default password setup here.

	return nil
}
