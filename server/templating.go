package server

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
)

// TemplatedProcessor implements RuntimeProcessor for files
// that can be templated.
type templatedProcessor struct {
	sync.Mutex
	templates map[string]*template.Template
}


func MakeTemplatedProcessor() RuntimeProcessor {
	return &templatedProcessor{ 
		templates: make(map[string]*template.Template),
	}
}


// BasicStateVariables is a placeholder to demonstrate that the
// template code is doing the right thing.
type BasicStateVariables struct {
	Message string
	
}

func MakeBasicStateVariables() interface{} {
	return &BasicStateVariables{
		Message: "hello from inside of the program",
	}
}

// ServeString caches the parsed templates. 
func (gs *templatedProcessor) ServeString(s string, w http.ResponseWriter, req *http.Request) {
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
	if err := template.Execute(w, MakeBasicStateVariables()); err != nil {
		respondWithError(w, fmt.Sprintln("Can't execute template", err))
	}
}

// ServeStream implementation re-parses the template each time and then
// executes it. The presumption is that in stream serving mode, a single developer
// is using the software.
func (gs *templatedProcessor) ServeStream(reader io.Reader, w http.ResponseWriter, req *http.Request) {
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

	if err := template.Execute(w, MakeBasicStateVariables()); err != nil {
		respondWithError(w, fmt.Sprintln("Can't execute template", err))
	}
	
}
