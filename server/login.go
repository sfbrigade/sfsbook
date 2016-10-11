package server


import (
	"fmt"
	"log"
	"net/http"

	"github.com/blevesearch/bleve"
)

// This is written possibly incorrectly. Refactor later.
// I need more tests and more refactoring. my code arrangement leaves
// much to be desired.
type loginServer struct {
	embr      *embeddableResources
	passwordfile  bleve.Index
	cookietool *cookieTooling
}

// makeLoginHandler returns a handler for login actions.
func (hf *HandlerFactory) makeLoginHandler() http.Handler {
	return &loginServer{
		embr:      makeEmbeddableResource(hf.sitedir),
		passwordfile: hf.passwordfile,
		cookietool: hf.cookietool,
	}
}

func (gs *loginServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	sn := req.URL.Path

	log.Println("loginServer is handling request for", sn)

	if req.Method == "POST" {
		log.Println("handling Post of resource")
		if err := req.ParseForm(); err != nil {
			respondWithError(w, fmt.Sprintln("bad uploaded form data to login: ", err))
			return
		}

		// We only use the posted info.
		log.Println("parsing form:")
		for k, v := range req.PostForm {
			log.Println("	", k, "		", v)
		}

		// dbreq.IsPost = true
		// dbreq.PostArgs = req.PostForm


		// search for name
		// get all the info
		// if hash matches
		// update the cookie
		// if hash doesn't match, indicate that it's failed.

		// debug flag should not be part of the results database.
	}

	str, err := gs.embr.GetAsString(sn)
	if err != nil {
		// TODO(rjk): Rationalize error handling here. There needs to be a 404 page.
		respondWithError(w, fmt.Sprintln("Server error", err))
	}

	// Interim. Clearly wrong.
	//results := dba.MakeStubGenerator(gs.passwordfile).ForRequest(req)
	parseAndExecuteTemplate(w, req, str, nil)
}
