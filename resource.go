package main

import (
	"encoding/json"
	"errors"
	"html/template"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

var (
	ErrResourceNotFound = errors.New("Resource not found.")
)

type Resource struct {
	Name         string
	Slug         string
	Categories   string
	Description  string
	Services     string
	Email        string
	Address      string
	Languages    string
	PopsServed   string
	Website      string
	BusinessLine string
	CrisisLine   string
	Fax          string
}

type escapedResources struct {
	Name         template.HTML
	Slug         template.HTML
	Categories   template.HTML
	Description  template.HTML
	Services     template.HTML
	Email        template.HTML
	Address      template.HTML
	Languages    template.HTML
	PopsServed   template.HTML
	Website      template.HTML
	BusinessLine template.HTML
	CrisisLine   template.HTML
	Fax          template.HTML
}

var ResourceFile = "refguides/2016-april/refguide.json"

func GetResources() (resources []*Resource, err error) {
	file, err := ioutil.ReadFile(ResourceFile)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(file, &resources); err != nil {
		return nil, err
	}

	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}

	for _, r := range resources {
		r.Slug = reg.ReplaceAllString(strings.ToLower(r.Name), "-")
		r.Description = strings.Replace(r.Description, `"`, `\"`, -1)
		r.Services = strings.Replace(r.Services, `"`, `\"`, -1)
	}

	return resources, nil
}

func GetResourceBySlug(slug string) (*Resource, error) {
	resources, err := GetResources()
	if err != nil {
		return nil, err
	}

	for _, r := range resources {
		if r.Slug == slug {
			return r, nil
		}
	}

	return nil, ErrResourceNotFound
}
