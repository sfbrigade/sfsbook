package server

import (
	"fmt"
	"log"

	"github.com/rjkroege/mocking"
	"github.com/sfbrigade/sfsbook/dba"
	"golang.org/x/crypto/bcrypt"
	"github.com/pborman/uuid"
)

// TODO(rjk): Complete the higher level of semantic abstraction with
// some kind of appropriate wrapper for ListUsers, Delete etc.
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

// You can search for anything,
// but only (username,password := "test_user", "password") return a nonempty result
func (tape *mockPasswordIndex) Search(username string) (dba.PasswordSearchResult, error) {

	testUsername := "test_user"
	testPassHash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	match := dba.PasswordSearchResult{
		Found: true,
		DisplayName: testUsername,
		Role: "volunteer",
		PasswordHash: string(testPassHash),
		ID: uuid.NewRandom(),
	}

	miss := dba.PasswordSearchResult{Found: false}


	if username == testUsername {
		return match, nil
	}
	return miss, nil
}

func (tape *mockPasswordIndex) ListUsers(userquery string, size, from int) ([]map[string]interface{}, error) {
	return nil, fmt.Errorf("not-implemented")
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
