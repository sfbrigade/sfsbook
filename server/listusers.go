package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search/query"
	"github.com/sfbrigade/sfsbook/dba"
)

type listUsers struct {
	embr         *embeddableResources
	passwordfile dba.PasswordIndex
}

// makeListUsersHandler returns a handler for changing the user password.
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
	templates := []string{sn}
	// do the redirect?
	parseAndExecuteTemplate(gs.embr, w, req, templates, listusersresult)
}

// TODO(rjk): Note refactoring opportunity with basic search?
// Also: search could have the same "search on the results page"
func (gs *listUsers) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Println("ServeHTTP")

	listusersresult := new(listUsersResult)

	// Skip immediately to end if not authed.
	if !GetCookie(req).HasCapability(CapabilityViewUsers) {
		gs.ender(w, req, listusersresult)
		return
	}

	log.Println("ServeHTTP didn't skip to end")

	var queryop query.Query
	queryop = bleve.NewMatchAllQuery()

	// Setup a query. The query is different if we have specified it.
	if req.Method == "POST" {
		log.Println("got a post")

		if err := req.ParseForm(); err != nil {
			respondWithError(w, fmt.Sprintln("bad uploaded form data to login: ", err))
			return
		}

		// We only use the posted info.
		// This dumps passwords in the clear in the log.
		// TODO(rjk): delete before landing this code.
		log.Println("parsing form:")
		for k, v := range req.PostForm {
			log.Println("	", k, "		", v)
		}

		userquery, err := getValidatedString("userquery", req.PostForm)
		if err != nil {
			log.Println("no userquery in Post", err)
			listusersresult.Diagnosticmessage = "Ignoring an unusable user search query."
			gs.ender(w, req, listusersresult)
			return
		}
		listusersresult.Userquery = userquery

		queryop = bleve.NewWildcardQuery(userquery)

		// Is displayname indexed? I might want to do that...
		// I need to version the database...

	}

	// I need to make this search the right way. And bound the result set
	// size.
	sreq := bleve.NewSearchRequest(queryop)
	sreq.Fields = []string{"name", "role", "display_name"}

	// These two values need to come from the URL args.
	sreq.Size = 10
	sreq.From = 0

	// This is an error case (something is wrong internally)
	searchResults, err := gs.passwordfile.Search(sreq)
	if err != nil {
		respondWithError(w, fmt.Sprintln("database couldn't respond with useful results", err))
	}

	if len(searchResults.Hits) < 1 {
		// This probably means that the user has entered an invalid query.
		listusersresult.Diagnosticmessage = "Userquery matches no users."
		gs.ender(w, req, listusersresult)
		return
	}

	users := make([]map[string]interface{}, 0, len(searchResults.Hits))
	for _, sr := range searchResults.Hits {
		u := make(map[string]interface{})
		for k, v := range sr.Fields {
			// Could test and drop the unfortunate?
			u[k] = v.(string)
		}
		u["uuid"] = sr.ID
		users = append(users, u)
	}
	listusersresult.Querysuccess = true
	listusersresult.Users = users

	gs.ender(w, req, listusersresult)
}
