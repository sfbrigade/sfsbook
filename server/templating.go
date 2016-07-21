package server

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/sfbrigade/sfsbook/dba"
)

// templatedServer is a 
type templatedServer struct {
	sync.Mutex
	templates map[string]*template.Template
	ff *fileFinder

	generator dba.Generator
}

func MakeTemplatedServer(ff *fileFinder, g dba.Generator) *templatedServer {
	return &templatedServer{ 
		templates: make(map[string]*template.Template),
		ff: ff,
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

	if err := gs.ff.StreamOrString(sn, gs, w, req); err != nil {
		respondWithError(w, fmt.Sprintln("Server error", err))
	}
}


// ServeString caches the parsed templates. 
func (gs *templatedServer) ServeForString(s string, w http.ResponseWriter, req interface{}) {
	gs.Lock()
	template, ok := gs.templates[s]
	gs.Unlock()

	if !ok {
		var err error
		template, err = template.New("htmlbase").Parse(s)
		if err != nil {
			respondWithError(w, fmt.Sprintln("Can't parse template", err))
			return
		}
		gs.Lock()
		if _, ok := gs.templates[s]; !ok {
			gs.templates[s] = template
		}
		gs.Unlock()
	}

	// TODO(rjk): plumb the state into here and wire it up in some way.
	// The basic idea: there's a different staticServer instance for each of the
	// the various files.
	if err := template.Execute(w, gs.generator.ForRequest(req)); err != nil {
		respondWithError(w, fmt.Sprintln("Can't execute template", err))
	}
}

// ServeStream implementation re-parses the template each time and then
// executes it. The presumption is that in stream serving mode, a single developer
// is using the software.
func (gs *templatedServer) ServeForStream(reader io.Reader, w http.ResponseWriter, req interface{}) {
	templatestr, err := ioutil.ReadAll(reader)
	if err != nil {
		respondWithError(w, fmt.Sprintln("Can't read source file", err))
		return
	}
	// TODO(rjk): Logs, perf measurements, etc.
	template, err := template.New("htmlbase").Parse(string(templatestr))
	if err != nil {
		respondWithError(w, fmt.Sprintln("Can't parse template", err))
		return
	}

	if err := template.Execute(w, gs.generator.ForRequest(req)); err != nil {
		respondWithError(w, fmt.Sprintln("Can't execute template", err))
	}
}
