package dba

import (
	"github.com/blevesearch/bleve"
)

// PasswordIndex is a minimal set of entry points into bleve.Index
// that we need to operate. Brought out to simplify mocking.
// See bleve.Index for documentation and usage.
type PasswordIndex interface {
	// Inserts the provided data bundle into the password database.
	Index(id string, data interface{}) error

	// Searches the password database for the specified search string
	// and returns a search result.
	Search(req *bleve.SearchRequest) (*bleve.SearchResult, error)

	// Returns the document map for a given document (i.e. password entry)
	// when presented with the entry's UUID
	// TODO(rjk): Make sure that the returned map can be indexed. 
	// And maybe give these better names. And use UUIDs for type clarity.
	MapForDocument(id string) (map[string]interface{}, error)

	// Deletes the specificed password entry.
	Delete (id string) error
}

type blevePasswordIndex struct {
	idx bleve.Index
}

func (pdoc *blevePasswordIndex) MapForDocument(id string) (map[string]interface{}, error) {
	idx := pdoc.idx
	doc, err := idx.Document(id)
	if err != nil {
		return nil, err
	}
	return MakeMapFromDocument(doc)
}

// TODO(rjk): Fix the layering violation of leaking bleve into the interface.
func (pdoc *blevePasswordIndex) Search(req *bleve.SearchRequest) (*bleve.SearchResult, error) {
	return pdoc.idx.Search(req)
}

func (pdoc *blevePasswordIndex) Delete(id string) error {
	return pdoc.idx.Delete(id)
}


func (pdoc *blevePasswordIndex) Index(id string, data interface{}) error {
	return pdoc.idx.Index(id, data)
}

func OpenPassword(persistentroot string) (PasswordIndex, error) {
	passwordfile, err := OpenBleve(persistentroot, PasswordFile)
	if err != nil {
		return nil, err
	}
	return &blevePasswordIndex{passwordfile}, err
}
