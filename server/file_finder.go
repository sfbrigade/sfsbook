package server

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

//go:generate go run ../generator/tool/generateresources.go -output embedded_resources.go -prefix ../site/ ../site

// This file contains the code needed to either serve a resource from the
// embedded resources or find in the site directory. The above generate
// directory actually constructs the embedded resources using tooling
// from the generator package.

type embeddableResources struct {
	sitedir string
}

// makeEmbeddableResource returns a new EmbeddableResource. If
// sitedir is empty, all resources will be taken from the internal resource
// source.
func makeEmbeddableResource(sitedir string) *embeddableResources {
	return &embeddableResources{
		sitedir: sitedir,
	}
}

func (er *embeddableResources) alwaysGetEmbedded(upath string) (string, error) {
	// TODO(rjk): Resources should be compressed.
	res, ok := Resources[upath]
	if !ok {
		return "", Error(ErrorNoSuchEmbeddedResource)
	}
	return res, nil
}

// GetAsString retrieves file upath from either the embedded
// storage or from disk. It returns either a string containing the
// resource or an error if the file could not be retrieved.
func (er *embeddableResources) GetAsString(upath string) (string, error) {
	if er.sitedir == "" {
		return er.alwaysGetEmbedded(upath)
	}
	fpath := filepath.Join(er.sitedir, upath)

	if _, err := os.Stat(fpath); err != nil {
		return er.alwaysGetEmbedded(upath)
	}

	fileasbytes, err := ioutil.ReadFile(fpath)
	if err != nil {
		log.Println("EmbeddableResource.GetAsString: problem reading file", fpath, err)
		return "", Error(ErrorNoSuchFileResource)
	}
	return string(fileasbytes), nil
}
