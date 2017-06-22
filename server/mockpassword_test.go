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

type deleteStim struct {
	name string
	uuid string
}

type listUsersStim struct {
	name  string
	query string
	size  int
	from  int
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
func (tape *mockPasswordIndex) Search(username string) (*dba.PasswordSearchResult, error) {

	testUsername := "test_user"
	testPassHash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	match := dba.NewPasswordSearchResult(
		string(testPassHash),
		"volunteer",
		testUsername,
		uuid.NewRandom(),
	)

	miss := (*dba.PasswordSearchResult)(nil)

	if username == testUsername {
		return match, nil
	}
	return miss, nil
}

func (mpi *mockPasswordIndex) ListUsers(userquery string, size, from int) ([]map[string]interface{}, error) {
	tape := (*mocking.Tape)(mpi)
	res := tape.Record(listUsersStim{
		"PasswordIndex.ListUsers",
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

func (mpi *mockPasswordIndex) Delete(id string) error {
	tape := (*mocking.Tape)(mpi)

	res := tape.Record(deleteStim{
		"PasswordIndex.Delete",
		id,
	})

	switch v := res.(type) {
	case nil:
		return nil
	case error:
		return v
	}
	log.Fatalf("%s (mock) response %v is of bad type", "PasswordIndex.Delete", res)
	return nil
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
