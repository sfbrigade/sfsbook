package fieldmap

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzers/keyword_analyzer"
	"github.com/pborman/uuid"
	"golang.org/x/crypto/bcrypt"
)

// buildPasswordDocumentMapping makes a mapping for the password file.
func buildPasswordEntryMapping() *bleve.DocumentMapping {
	passwordFieldMapping := bleve.NewTextFieldMapping()
	passwordFieldMapping.Index = false
	passwordFieldMapping.Store = true
	passwordFieldMapping.IncludeTermVectors = false
	passwordFieldMapping.IncludeInAll = false
	passwordFieldMapping.Analyzer = keyword_analyzer.Name

	// This mapping presumes that costs are less than 255. This might
	// not always be true. I had used a number field (originally) but
	// it made the test flaky.
	costFieldMapping := bleve.NewTextFieldMapping()
	costFieldMapping.Index = false
	costFieldMapping.IncludeTermVectors = false
	costFieldMapping.IncludeInAll = false
	costFieldMapping.Store = true
	costFieldMapping.Analyzer = keyword_analyzer.Name

	// passwordEntryMapping is a document for each user.
	passwordEntryMapping := bleve.NewDocumentMapping()

	passwordEntryMapping.DefaultAnalyzer = keyword_analyzer.Name
	passwordEntryMapping.AddFieldMappingsAt("username", keywordFieldMapping)
	passwordEntryMapping.AddFieldMappingsAt("passwordhash", passwordFieldMapping)
	passwordEntryMapping.AddFieldMappingsAt("cost", costFieldMapping)
	passwordEntryMapping.AddFieldMappingsAt("account_created", dateTimeMapping)
	passwordEntryMapping.AddFieldMappingsAt("last_login", dateTimeMapping)

	return passwordEntryMapping
}

// PasswordFileType implements IndexFactory for password data.
type PasswordFileType string

func (g PasswordFileType) Name() string {
	return string(g)
}

func (_ PasswordFileType) Mapping() *bleve.IndexMapping {
	return allDocumentMapping(IndexDocumentMap{
		"password": buildPasswordEntryMapping(),
	})
}

var PasswordFile = PasswordFileType("password.bleve")

var init_passwords = flag.Bool("init_passwords", false, "create a set of test passwords if none exist")

func (_ PasswordFileType) LoadStartData(i bleve.Index, pathroot string) error {
	if !*init_passwords {
		return fmt.Errorf("There is no password file. Cowardly refusing to create a really insecure one without a command line flag.")
	}

	log.Println("Setting up default password file. Warning! This is not secure. Seriously!")

	s := []map[string]interface{}{
		map[string]interface{}{
			"name":         "volunteer",
			"cost":         string(bcrypt.DefaultCost),
			"passwordhash": "open",
			"role":         "volunteer",
		},
		map[string]interface{}{
			"name":         "admin",
			"cost":         string(bcrypt.DefaultCost),
			"passwordhash": "sesame",
			"role":         "admin",
		},
	}

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
