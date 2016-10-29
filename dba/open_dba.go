package dba

import (
	"os"
	"path/filepath"

	"github.com/blevesearch/bleve"
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

		bi, err = bleve.New(dbpath, dxf.Mapping())
		if err != nil {
			goto cleanup
		}

		if err = dxf.LoadStartData(bi, persistentroot); err != nil {
			goto cleanup
		}
		if err = bi.Close(); err != nil {
			goto cleanup
		}

		// Re-open to work-around KV not done writing.
		if bi, err = bleve.Open(dbpath); err != nil {
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
