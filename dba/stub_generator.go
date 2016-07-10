package dba

import (
	"net/http"
)

// StubGenerator is an uninteresting generator for test purposes.
// TODO(rjk): Perhaps this should be "mock generator"?
type StubGenerator struct {
}

func MakekStubGenerator() *StubGenerator {
	return &StubGenerator{}
}

type stubGeneratorModel struct {
	Message string
}

func (sg *StubGenerator) ForRequest(req *http.Request) interface{} {
	return &stubGeneratorModel{
		Message: "hello from inside of the program",
	}
}
