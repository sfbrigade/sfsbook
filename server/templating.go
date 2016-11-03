package server

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/sfbrigade/sfsbook/dba"
)

// templatedServer is a server instance that uses results from generator
// to populate a Go template.
type templatedServer struct {
	embr      *embeddableResources
	generator dba.Generator
}

func (hf *HandlerFactory) makeTemplatedHandler(g dba.Generator) *templatedServer {
	return &templatedServer{
		embr:      makeEmbeddableResource(hf.sitedir),
		generator: g,
	}
}

func (gs *templatedServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	sn := req.URL.Path
	if sn == "" {
		sn = "index.html"
	}

	if err := req.ParseForm(); err != nil {
		respondWithError(w, fmt.Sprintln("invalid form parameters", err))
	}

	str, err := gs.embr.GetAsString(sn)
	if err != nil {
		// TODO(rjk): Rationalize error handling here. There needs to be a 404 page.
		respondWithError(w, fmt.Sprintln("Server error", err))
	}
	hdr, err := gs.embr.GetAsString("/header.html")
	if err != nil {
		// TODO(rjk): Rationalize error handling here. There needs to be a 404 page.
		respondWithError(w, fmt.Sprintln("Server error", err))
		return
	}
	ftr, err := gs.embr.GetAsString("/footer.html")
	if err != nil {
		// TODO(rjk): Rationalize error handling here. There needs to be a 404 page.
		respondWithError(w, fmt.Sprintln("Server error", err))
		return
	}
	// The req contains the cookie info. And so we can bound viewability
	// in the database.
	results := gs.generator.ForRequest(req)
	templates := []string{str,hdr,ftr}
	// TODO(rjk): I need to do something smarter about caching.
	// I removed the cache of templates pending the global cache.
	parseAndExecuteTemplate(w, req, templates, results)
}

// TODO(rjk): I think that this is not quite right code structure.
// instead, there needs to be a dbareq re-writing layer.

type templateParameters struct {
	Results       interface{}
	DecodedCookie *UserCookie
}

// parseAndExecuteTemplate implementation re-parses the template from templatestr
// and executes it with a templateParameters object. This is a utility method to be
// used from anywhere that the provided template is to be used.
// TODO(rjk): add caching of results.
// TODO(rjk): permit many arguments. They need to get bundled into a kv-store
// that keeps things more flexible.
func parseAndExecuteTemplate(w http.ResponseWriter, req *http.Request, templatestrings []string, result interface{}) {
	// TODO(rjk): Logs, perf measurements, etc.
	template, err := template.New("htmlbase").Parse(templatestrings[0])
	if err != nil {
		respondWithError(w, fmt.Sprintln("Can't parse template", err))
		return
	}
	for _, t := range templatestrings[1:] {
		_, err = template.Parse(t)
		if err != nil {
			respondWithError(w, fmt.Sprintln("Can't parse template", err))
			return
		}
	}

	// Also, I need to make result optional so this is the wrong way to proceed.
	tp := &templateParameters{
		Results:       result,
		DecodedCookie: GetCookie(req),
	}

	if err := template.Execute(w, tp); err != nil {
		respondWithError(w, fmt.Sprintln("Can't execute template", err))
	}
}
