package dba

import (
	"fmt"

	"github.com/blevesearch/bleve/document"
)

// resultsMapFromDocument builds a results object.
func resultsMapFromDocument(doc *document.Document) (map[string]interface{}, error) {
	resultsMap := make(map[string]interface{})
	allerrors := make([]error, 0)

	for _, f := range doc.Fields {
		switch t := f.(type) {
		default:
			allerrors = append(allerrors, fmt.Errorf("%s: as unexpected type", f.Name()))
		case *document.NumericField:
			v, err := t.Number()
			if err != nil {
				allerrors = append(allerrors, fmt.Errorf("%s failed to convert to number: %v", t.Name(), err))
				continue
			}
			resultsMap[t.Name()] = v
		case *document.BooleanField:
			v, err := t.Boolean()
			if err != nil {
				allerrors = append(allerrors, fmt.Errorf("%s failed to convert to boolean: %v", t.Name(), err))
				continue
			}
			resultsMap[t.Name()] = v
		case *document.TextField:
			resultsMap[t.Name()] = string(t.Value())
		case *document.DateTimeField:
			v, err := t.DateTime()
			if err != nil {
				allerrors = append(allerrors, fmt.Errorf("%s failed to convert to time: %v", t.Name(), err))
				continue
			}
			resultsMap[t.Name()] = v
		}
	}

	var err error
	if len(allerrors) > 0 {
		err = fmt.Errorf("field errors: %v", allerrors)
	}
	return resultsMap, err
}
