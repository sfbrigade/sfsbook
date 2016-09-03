package server

import (
	"fmt"
	"html/template"
	"net/http"
	"sync"

	"github.com/sfbrigade/sfsbook/dba"
)

// templatedServer is a
type templatedServer struct {
	sync.Mutex
	templates map[string]*template.Template
	global        *GlobalState

	generator dba.Generator
}

func MakeTemplatedServer(global  *GlobalState, g dba.Generator) *templatedServer {
	return &templatedServer{
		templates: make(map[string]*template.Template),
		global:        global,
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

	str, err := gs.global.GetAsString(sn)
	if err != nil {
		// TODO(rjk): Rationalize error handling here. There needs to be a 404 page.
		respondWithError(w, fmt.Sprintln("Server error", err))
	}

	// TODO(rjk): I need to do something smarter about caching.
	// I removed the cache of templates pending the global cache.
	gs.ServeForStrings(str, w, req)
}

// ServeStream implementation re-parses the template each time and then
// executes it. The presumption is that in stream serving mode, a single developer
// is using the software.
func (gs *templatedServer) ServeForStrings(templatestr string, w http.ResponseWriter, req interface{}) {
	// TODO(rjk): Logs, perf measurements, etc.
	template, err := template.New("htmlbase").Parse(string(templatestr))
	if err != nil {
		respondWithError(w, fmt.Sprintln("Can't parse template", err))
		return
	}

	generatedResult := gs.generator.ForRequest(req)
	generatedResult.SetDebug(true)
	if err := template.Execute(w, generatedResult); err != nil {
		respondWithError(w, fmt.Sprintln("Can't execute template", err))
	}
}
