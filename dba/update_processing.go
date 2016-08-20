package dba

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

func (rr *ResourceRequest) validateIfNecessary() error {
	if !rr.IsPost {
		// Nothing to do.
		return nil
	}

	// TODO(rjk): Must validate here.	 Many validations are possible.
	// Recall that this is user-generated content and must be treated
	// appropriately. 
	return nil
}


func (qr *ResourceResultsGenerator) mergeAndUpdateIfNecessary(rr *ResourceRequest, updatedresults map[string]interface{}) error {
	if !rr.IsPost {
		// Nothing to do.
		return nil
	}

	uuid := rr.Uuid
	postdata := rr.PostArgs
	needtoupdate := false

	// I will bundle all the errors together but it might be preferable to handle update
	// errors specially so that I can bring out to the UI in some way. i.e.: errors need to
	// be per-field. 
	// TODO(rjk): deliver errors per-field and carry through to the resource template.
	allerrors := make([]error, 0)

	for k, v:= range updatedresults {
		if _, ok := immutableFields[k]; ok {
			continue
		}

		postedarray, ok := postdata[k]
		if !ok {
			continue
		}

		if len(postedarray) != 1 {
			log.Println(k, "is weird and has a multi-valued array", postedarray)
			continue
		}
		posted := postedarray[0]

		switch v.(type) {
		case bool:
			log.Println("bool value: need to exercise")
			// The database value is bool.
			if posted == "checked" {
				// This may not be handling the boolean field properly.
				updatedresults[k] = true
				needtoupdate = true
			} else {
				updatedresults[k] = false				
			}
		case float64:
			if n, err := strconv.ParseFloat(posted, 64); err == nil {
				updatedresults[k] = n
				needtoupdate = true
			} else {
				allerrors = append(allerrors, fmt.Errorf("%s failed to convert to float: %v", err))
			}
		case time.Time:
			// TODO(rjk): Obviously, I'll want to support many date and time formats so
			// try more of them.
			if date, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", posted); err == nil {
				updatedresults[k] = date
				needtoupdate = true
			} else {
				allerrors = append(allerrors, fmt.Errorf("%s failed to convert to date: %v", err))
			}
		case string:
			updatedresults[k] = posted
			needtoupdate = true
		}			
	}

	if needtoupdate {
		updatedresults["date_last_modified"] = time.Now()
		if err := qr.index.Index(uuid, updatedresults); err != nil {
			allerrors = append(allerrors, err)
		}
	}

	// TODO(rjk): Further error rationalization possible.
	if len(allerrors) > 0 {
		return fmt.Errorf("all update errors: %v", allerrors)
	} else {
		return nil
	}
}
