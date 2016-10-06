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

	results := gs.generator.ForRequest(req)

	// TODO(rjk): I need to do something smarter about caching.
	// I removed the cache of templates pending the global cache.
	gs.serveForStrings(w, req, str, results)
}

// TODO(rjk): I think that this is not quite right code structure.
// instead, there needs to be a dbareq re-writing layer.

type templateParameters struct {
	Results       dba.GeneratedResult
	DecodedCookie *UserCookie
}

// serveForStrings implementation re-parses the template from templatestr
// and executes it with a templateParameters object.
// TODO(rjk): figure out how to cache this as necessary.
func (gs *templatedServer) serveForStrings(w http.ResponseWriter, req *http.Request, templatestr string, result dba.GeneratedResult) {
	// TODO(rjk): Logs, perf measurements, etc.
	template, err := template.New("htmlbase").Parse(string(templatestr))
	if err != nil {
		respondWithError(w, fmt.Sprintln("Can't parse template", err))
		return
	}

	// TODO(rjk): This needs to be set differently. Instead of always.
	result.SetDebug(true)
	tp := &templateParameters{
		Results:       result,
		DecodedCookie: GetCookie(req),
	}

	if err := template.Execute(w, tp); err != nil {
		respondWithError(w, fmt.Sprintln("Can't execute template", err))
	}
}
