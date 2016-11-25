package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"strconv"

	"github.com/pborman/uuid"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search/query"
	"github.com/sfbrigade/sfsbook/dba"
)

const (
	_ = iota
	_SEARCHACTION
	_RESETPASSWORD
	_DELETEUSERS
	_ROLECHANGE_TO_ADMIN
	_ROLECHANGE_TO_VOLUNTEER
	_ROLECHANGE_TO_NOROLE
	_BADACTORREQUEST
)



type listUsers struct {
	embr         *embeddableResources
	passwordfile dba.PasswordIndex
}

// makeListUsersHandler returns a handler for editing user data.
func (hf *HandlerFactory) makeListUsersHandler() http.Handler {
	return &listUsers{
		embr:         makeEmbeddableResource(hf.sitedir),
		passwordfile: hf.passwordfile,
	}
}

type listUsersResult struct {
	Userquery string
	// TODO(rjk): Consider making this typed in some fashion.
	Users             []map[string]interface{}
	Querysuccess      bool
	Diagnosticmessage string
}

func (gs *listUsers) ender(w http.ResponseWriter, req *http.Request, listusersresult interface{}) {
	sn := req.URL.Path
	str, err := gs.embr.GetAsString(sn)
	if err != nil {
		// TODO(rjk): Rationalize error handling here. There needs to be a 404 page.
		respondWithError(w, fmt.Sprintln("Server error", err))
		return
	}

	parseAndExecuteTemplate(w, req, str, listusersresult)
}

// TODO(rjk): Note refactoring opportunity with resource search.
func (gs *listUsers) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Println("ServeHTTP")

	listusersresult := new(listUsersResult)

	// Skip immediately to end if not authed.
	if !GetCookie(req).HasCapability(CapabilityViewUsers) {
		gs.ender(w, req, listusersresult)
		return
	}

	var queryop query.Query
	queryop = bleve.NewMatchAllQuery()

	// Setup a query. The query is different if we have specified it.
	log.Println("not a post but proceeding anyway")
	log.Println("req", req)

	if err := req.ParseForm(); err != nil {
		respondWithError(w, fmt.Sprintln("bad uploaded form data to login: ", err))
		return
	}

	// Remember that nothing provided as part of the form can be
	// trusted to be valid. All must be validated. I need a general
	// mechanism to address clients that are broken.
	selecteduuids := make([]uuid.UUID, 0, 10)
	action := _SEARCHACTION
	for k, v := range req.Form {
		log.Println("listusers ServeHTTP processing form item:", k, v)
		switch {
		case strings.HasPrefix(k, "selected-"):
			uuidstring, err := getValidatedString(k, req.Form)
			if err != nil {
				log.Println("posted form contained invalid k,v pair:", k, v, err)
				action = _BADACTORREQUEST
				break
			}
			uuid := uuid.Parse(uuidstring)
			if uuid == nil {
				log.Println("posted form contained invalid uuid", v)
				action = _BADACTORREQUEST
				break
			}
			selecteduuids = append(selecteduuids, uuid)
		case k == "rolechange":
			// parse the desired rolechange
			rc := v[0]
			switch rc {
			case "nochange":
				// Don't have to do anything
			case "admin":
				action = _ROLECHANGE_TO_ADMIN
			case "norole":
				action = _ROLECHANGE_TO_ADMIN
			case "volunteer":
				action = _ROLECHANGE_TO_ADMIN
			default:
				log.Println("posted form contained invalid rolechange", k, v)
				action = _BADACTORREQUEST
				break
			}
		case k == "deleteuser":
			action = _DELETEUSERS
		case k == "userquery":
			userquery, err := getValidatedString("userquery", req.Form)
			if err != nil {
				log.Println("no userquery", err)
				listusersresult.Diagnosticmessage = "Showing all..."
			} else {
				// We had an argument to search with.
				listusersresult.Userquery = userquery
				queryop = bleve.NewWildcardQuery(userquery)
				// TODO(rjk): Improve database indexing.
			}
		case k =="resetpassword":
			action = _RESETPASSWORD
			// I don't currently have enough information in the password
			// database to do this. And adding the additional data may
			// require interaction with how we implement oauth integration.
			
		default:
			// TODO(rjk): This is a convenience to simplify implementation.
			log.Println("ignoring extra fields in form", k, v)
		}
	}

	// Do the updates
	switch action {
	case _SEARCHACTION:
		// Or do nothing if only searchin.
	case _RESETPASSWORD:
		log.Println("notimplemented: resetpassword applied to", selecteduuids)
	case _DELETEUSERS:
		log.Println("notimplemented deleteusers applied to", selecteduuids)
	case _ROLECHANGE_TO_ADMIN:
		log.Println("notimplemented rolechange to admin applied to", selecteduuids)
	case _ROLECHANGE_TO_VOLUNTEER:
		log.Println("notimplemented rolechange to volunteer applied to", selecteduuids)
	case _ROLECHANGE_TO_NOROLE:
		log.Println("notimplemented rolechange to norole applied to", selecteduuids)
	default: // includes _BADACTORREQUEST
		respondWithError(w, "client is attempting something wrong")
		return
	}

	// And now do the search.

	// I need to make this search the right way. And bound the result set
	// size.
	sreq := bleve.NewSearchRequest(queryop)
	sreq.Fields = []string{"name", "role", "display_name"}

	// These two values need to come from the URL args so that I can 
	// page through many users.
	sreq.Size = 10
	sreq.From = 0

	// This is an error case (something is wrong internally)
	searchResults, err := gs.passwordfile.Search(sreq)
	if err != nil {
		respondWithError(w, fmt.Sprintln("database couldn't respond with useful results", err))
		return
	}

	if len(searchResults.Hits) < 1 {
		// This probably means that the user has entered an invalid query.
		listusersresult.Diagnosticmessage = "Userquery matches no users."
		gs.ender(w, req, listusersresult)
		return
	}

	users := make([]map[string]interface{}, 0, len(searchResults.Hits))
	for i, sr := range searchResults.Hits {
		u := make(map[string]interface{})
		for k, v := range sr.Fields {
			// Could test and drop the unfortunate?
			u[k] = v.(string)
		}

		uuidcasted := uuid.UUID(sr.ID)
		// I thought about encrypting the UUIDs. But to get this content, one
		// must already have the admin role and that is enforced server side
		// via a strongly encrypted cookie. And they are cryptographically 
		// difficult to guess already.
		u["uuid"] = uuidcasted.String()
		u["index"] = strconv.FormatInt(int64(i), 10)
		users = append(users, u)
	}
	listusersresult.Querysuccess = true
	listusersresult.Users = users

	gs.ender(w, req, listusersresult)
}
