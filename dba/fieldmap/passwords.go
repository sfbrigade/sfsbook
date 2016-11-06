package fieldmap

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/mapping"
	"github.com/blevesearch/bleve/analysis/analyzer/keyword"
	"github.com/pborman/uuid"
	"golang.org/x/crypto/bcrypt"
)

// buildPasswordDocumentMapping makes a mapping for the password file.
func buildPasswordEntryMapping() *mapping.DocumentMapping {
	passwordFieldMapping := bleve.NewTextFieldMapping()
	passwordFieldMapping.Index = false
	passwordFieldMapping.Store = true
	passwordFieldMapping.IncludeTermVectors = false
	passwordFieldMapping.IncludeInAll = false
	passwordFieldMapping.Analyzer = keyword.Name

	// This mapping presumes that costs are less than 255. This might
	// not always be true. I had used a number field (originally) but
	// it made the test flaky.
	costFieldMapping := bleve.NewTextFieldMapping()
	costFieldMapping.Index = false
	costFieldMapping.IncludeTermVectors = false
	costFieldMapping.IncludeInAll = false
	costFieldMapping.Store = true
	costFieldMapping.Analyzer = keyword.Name

	// passwordEntryMapping is a document for each user.
	passwordEntryMapping := bleve.NewDocumentMapping()

	passwordEntryMapping.DefaultAnalyzer = keyword.Name
	passwordEntryMapping.AddFieldMappingsAt("username", KeywordFieldMapping)
	passwordEntryMapping.AddFieldMappingsAt("passwordhash", passwordFieldMapping)
	passwordEntryMapping.AddFieldMappingsAt("cost", costFieldMapping)
	passwordEntryMapping.AddFieldMappingsAt("account_created", DateTimeMapping)
	passwordEntryMapping.AddFieldMappingsAt("last_login", DateTimeMapping)
	passwordEntryMapping.AddFieldMappingsAt("display_name", KeywordFieldMapping)

	return passwordEntryMapping
}

// PasswordFileType implements IndexFactory for password data.
type PasswordFileType string

func (g PasswordFileType) Name() string {
	return string(g)
}

func (_ PasswordFileType) Mapping() *mapping.IndexMappingImpl {
	return AllDocumentMapping(IndexDocumentMap{
		"password": buildPasswordEntryMapping(),
	})
}

var PasswordFile = PasswordFileType("password.bleve")

var init_passwords = flag.Bool("init_passwords", false, "Create a set of insecure test passwords if none exist")
var init_admin_password = flag.Bool("init_admin_password", false, "Create a secure admin password.")

func (_ PasswordFileType) LoadStartData(i bleve.Index, pathroot string) error {
	var s []map[string]interface{}

	switch {
	case *init_passwords && *init_admin_password:
		return fmt.Errorf("Trying to make secure admin password and insecure test passwords is incompatible. Pick one.")
	case *init_passwords:
		s = []map[string]interface{}{
			map[string]interface{}{
				"name":         "volunteer",
				"cost":         string(bcrypt.DefaultCost),
				"passwordhash": "open",
				"role":         "volunteer",
				"display_name": "Pikachu Helper",
			},
			map[string]interface{}{
				"name":         "admin",
				"cost":         string(bcrypt.DefaultCost),
				"passwordhash": "sesame",
				"role":         "admin",
				"display_name": "Pokemon Guardian",
			},
		}
		log.Println("Setting up default password file. Warning! This is not secure. Seriously!")
	case *init_admin_password:
		// Need to generate a unique random string.
		defaultpassword := RandStringBytesRmndr(8)
		log.Printf("Admin password is %s. Please change immediately", defaultpassword)
		s = []map[string]interface{}{
			map[string]interface{}{
				"name":         "admin",
				"cost":         string(bcrypt.DefaultCost),
				"passwordhash": defaultpassword,
				"role":         "admin",
				"display_name": "Default Administrator",
			},
		}
	default:
		return fmt.Errorf("There is no password file. Cowardly refusing to create a really insecure one without a command line flag.")
	}

	// TODO(rjk): easy refactoring... can save typing and improve correctness.
	// want Create, Delete, Update right? All can batch.
	// can use a struct? Fixed structs work?
	

	batch := i.NewBatch()
	for _, r := range s {
		// Using a uuid means that identities can be assigned in parallel
		// across multiple machines without requiring consensus.
		// We need to arrange for the same uuid to be used for both this
		// record and the additional user information.
		rid := uuid.NewRandom()
		r["_type"] = "password"

		pw := []byte(r["passwordhash"].(string))
		hash, err := bcrypt.GenerateFromPassword(pw, bcrypt.DefaultCost)
		if err != nil {
			log.Println("bcrypt didn't do its thing:", err)
			continue
		}
		r["passwordhash"] = string(hash)

		r["account_created"] = time.Now()

		// NB: in Bleve, the argument map must have supported types or else
		// Bleve gives up on the field. This makes a certain amount of sense but
		// the non-support for byte[] vs string and different kinds of numbers
		// perplexed me briefly.
		batch.Index(string(rid), r)
	}

	if err := i.Batch(batch); err != nil {
		return err
	}
	return nil
}

// This code is copied from:
// http://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
// There is no need to over-engineer a solution for something that shouldn't
// be invoked that often.

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// RandStringBytesRmndr returns a random character string of length n.
func RandStringBytesRmndr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Int63()%int64(len(letters))]
	}
	return string(b)
}
