package dba

import (
	"github.com/blevesearch/bleve"
)

// StubGenerator is an uninteresting generator for test purposes.
// TODO(rjk): Perhaps this should be "mock generator"?
// TODO(rjk): Divide this into mock and full.
type StubGenerator struct {
	index bleve.Index
}

func MakeStubGenerator(bi bleve.Index) *StubGenerator {
	return &StubGenerator{}
}

// stubGeneratorModel is a GeneratedResult implementation that does nothing.
type stubGeneratorModel struct {
	generatedResultCore
	Message string
}

func (sg *StubGenerator) ForRequest(req interface{}) GeneratedResult {
	// TODO(rjk): do a query here against the index based on req.
	// TODO(rjk): write some kind of parse thing that valiadates input
	// TODO(rjk): manage cookies etc.
	// TODO(rjk): web doesn't belong in DBA. And dba doesn't belong in server so fix up your layering issues.
	return &stubGeneratorModel{
		Message: "hello from inside of the program",
	}
}
