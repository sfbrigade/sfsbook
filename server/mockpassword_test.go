package server

import (
	"fmt"
	"log"

	"github.com/blevesearch/bleve"
	"github.com/rjkroege/mocking"
)

// TODO(rjk): Complete the higher level of semantic abstraction with
// some kind of appropriate wrapper for Search.
type mockPasswordIndex mocking.Tape

type indexStim struct {
	fun  string
	id   string
	data interface{}
}

type docStim struct {
	name string
	uuid string
}

type listUsersStim struct {
	query string
	size int
	from int
}

func (mpi *mockPasswordIndex) Index(id string, data interface{}) error {
	tape := (*mocking.Tape)(mpi)
	res := tape.Record(indexStim{
		"PasswordIndex.Index",
		id,
		data,
	})
	switch v := res.(type) {
	case nil:
		return nil
	case error:
		return v
	}
	log.Fatalf("%s (mock) response %v is of bad type", "PasswordIndex.Index", res)
	return nil
}

func (tape *mockPasswordIndex) Search(_ *bleve.SearchRequest) (*bleve.SearchResult, error) {
	return nil, fmt.Errorf("not-implemented")
}

func (mpi *mockPasswordIndex) ListUsers(userquery string, size, from int) ([]map[string]interface{}, error) {
	tape := (*mocking.Tape)(mpi)
	res := tape.Record(listUsersStim{
		userquery,
		size,
		from,
	})
	switch v := res.(type) {
	case []map[string]interface{}:
		return v, nil
	case error:
		return nil, v
	}
	log.Fatalf("%s (mock) response %v is of bad type", "PasswordIndex.Document", res)
	return nil, nil
}

func (tape *mockPasswordIndex) Delete(id string) error {
	return fmt.Errorf("not-implemented")
}

func (mpi *mockPasswordIndex) MapForDocument(id string) (map[string]interface{}, error) {
	tape := (*mocking.Tape)(mpi)
	res := tape.Record(docStim{
		"PasswordIndex.Document",
		id,
	})
	switch v := res.(type) {
	case map[string]interface{}:
		return v, nil
	case error:
		return nil, v
	}
	log.Fatalf("%s (mock) response %v is of bad type", "PasswordIndex.Document", res)
	return nil, nil
}
