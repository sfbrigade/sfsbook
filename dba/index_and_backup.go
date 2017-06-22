package dba

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/document"
)

const (
	// BackupFolderLocation is the folder where backup json dumps go to.
	BackupFolderLocation = "/tmp"

	// UUIDsIndexName is the name of index that is used to store UUIDs of the
	// resources.
	UUIDsIndexName = "UUIDsIndex"
)

// indexAndBackup wraps bleve.Index Document and Index methods. In addition it
// triggers backing up resources data to disk on Index operation.
type indexAndBackup struct {
	index bleve.Index
}

func (i *indexAndBackup) Index(id string, data interface{}) error {
	if err := i.index.Index(id, data); err != nil {
		return err
	}

	go i.backup()

	return nil
}

func (i *indexAndBackup) Document(id string) (*document.Document, error) {
	return i.index.Document(id)
}

func (i *indexAndBackup) backup() {
	uuids, err := i.Document(UUIDsIndexName)
	if err != nil {
		log.Println("indexAndBackup#backup: ERR getting index: ", UUIDsIndexName, err)
		return
	}

	log.Println("indexAndBackup#backup: About to backup: ", len(uuids.Fields), "from index: ", UUIDsIndexName)

	listOfResources := []map[string]interface{}{}

	for _, uuid := range uuids.Fields {
		resource, err := i.Document(string(uuid.Value()))
		if err != nil {
			log.Println("indexAndBackup#backup: ERR finding resource: ", string(uuid.Value()), err)
			return
		}

		dd, err := MakeMapFromDocument(resource)
		if err != nil {
			log.Println("indexAndBackup#backup: MakeMapFromDocument: ", string(uuid.Value()), err)
			return
		}

		jsonable := map[string]interface{}{}

		for k, v := range dd {
			if _, ok := immutableFields[k]; ok {
				continue
			}
			if _, ok := mustInitializeFields[k]; ok {
				continue
			}
			jsonable[k] = v
		}

		listOfResources = append(listOfResources, jsonable)
	}

	jsonString, err := json.MarshalIndent(listOfResources, "", "  ")
	if err != nil {
		log.Println("indexAndBackup#backup: ERR marshalling data: ", err)
		return
	}

	backupFileName := filepath.Join(
		BackupFolderLocation, fmt.Sprintf("%v.json", time.Now().Unix()),
	)
	if err := ioutil.WriteFile(backupFileName, jsonString, 0644); err != nil {
		log.Println("indexAndBackup#backup: ERR opening file to backup data: ", err)
		return
	}
}
