package dba

import (
	"github.com/blevesearch/bleve"
	"github.com/sfbrigade/sfsbook/dba/fieldmap"
)

// PasswordIndex is a minimal set of entry points into bleve.Index
// that we need to operate. Brought out to simplify mocking.
// See bleve.Index for documentation and usage.
type PasswordIndex interface {
	Index(id string, data interface{}) error
	Search(req *bleve.SearchRequest) (*bleve.SearchResult, error)
	MapForDocument(id string) (map[string]interface{}, error)
}

type BlevePasswordIndex struct {
	idx bleve.Index
}

func (pdoc *BlevePasswordIndex) MapForDocument(id string) (map[string]interface{}, error) {
	idx := pdoc.idx
	doc, err := idx.Document(id)
	if err != nil {
		return nil, err
	}
	return MakeMapFromDocument(doc)
}

func (pdoc *BlevePasswordIndex) Search(req *bleve.SearchRequest) (*bleve.SearchResult, error) {
	return pdoc.idx.Search(req)
}

func (pdoc *BlevePasswordIndex) Index(id string, data interface{}) error {
	return pdoc.idx.Index(id, data)
}

func OpenPassword(persistentroot string) (PasswordIndex, error) {
	passwordfile, err := OpenBleve(persistentroot, fieldmap.PasswordFile)
	if err != nil {
		return nil, err
	}
	return &BlevePasswordIndex{passwordfile}, err
}
