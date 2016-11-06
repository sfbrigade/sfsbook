package dba

import (
	"github.com/blevesearch/bleve"
)

// PasswordIndex is a minimal set of entry points into bleve.Index
// that we need to operate. Brought out to simplify mocking.
// See bleve.Index for documentation and usage.
type PasswordIndex interface {
	Index(id string, data interface{}) error
	Search(req *bleve.SearchRequest) (*bleve.SearchResult, error)
	MapForDocument(id string) (map[string]interface{}, error)
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

func (pdoc *blevePasswordIndex) Search(req *bleve.SearchRequest) (*bleve.SearchResult, error) {
	return pdoc.idx.Search(req)
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
