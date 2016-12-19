package server

import (
	"fmt"
	"html/template"
	"log"
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

	// The req contains the cookie info. And so we can bound viewability
	// in the database.
	results := gs.generator.ForRequest(req)
	templates := []string{sn, "/head.html", "/header.html", "/footer.html"}
	// TODO(rjk): I need to do something smarter about caching.
	// I removed the cache of templates pending the global cache.
	parseAndExecuteTemplate(gs.embr, w, req, templates, results)
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
func parseAndExecuteTemplate(embr *embeddableResources, w http.ResponseWriter, req *http.Request, templateNames []string, result interface{}) {
	// TODO(rjk): Logs, perf measurements, etc.
	templateStrings, err := getTemplateStrings(embr, templateNames)
	if err != nil {
		respondWithError(w, fmt.Sprintln("getTemplateStrings failed", err))
		return
	}

	template, err := template.New("htmlbase").Parse(templateStrings[0])
	if err != nil {
		respondWithError(w, fmt.Sprintln("Can't parse template", err))
		return
	}
	for _, t := range templateStrings[1:] {
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

// getTemplateStrings returns a new slice containing the template string for
// each template name in the original slice or an error if something is wrong.
func getTemplateStrings(embr *embeddableResources, templateNames []string) ([]string, error) {
	templateStrings := make([]string, len(templateNames))
	for i, v := range templateNames {
		log.Println("this is the template", v)
		str, err := embr.GetAsString(v)
		if err != nil {
			return []string{}, err
		}
		templateStrings[i] = str
	}
	return templateStrings, nil
}
