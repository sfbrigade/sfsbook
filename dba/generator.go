// Package dba provides the backend of the sfsbook server. There are different
// generator functions for each kind of page.
package dba

// Generator is an abstract interface implemented by each page's middle-tier.
// The server directory contains the front tier, Generators are the middle-tier, etc.
// A Generator instance knows how to convert a specific request into a structure
// that will be used to populate a template appropriate to the page.
type Generator interface {
	ForRequest(req interface{}) GeneratedResult
}

// GeneratedResult is an abstract interface implemented by the structure returned
// by ForRequest.
type GeneratedResult interface {
	// Success indicates to the consumer (e.g. page template) of the GeneratedResult
	// that the result was successfully generated.
	Success() bool

	// Debug indicates to the consumer of the GeneratedResult should it show
	// debugging output.
	Debug() bool

	// Mark this GeneratedResult to contain additional debugging contents.
	SetDebug(bool)

	// ErrorMessage returns a server error message.
	ErrorMessage() string
}

// generatedResultCore is common code for use in GeneratedResult instances.
type generatedResultCore struct {
	// Indicates if the query was successful. (i.e. that it produced data.)
	success bool

	// Show this if something went wrong.
	// TODO(rjk): Localize this string.
	failureText string

	// True if we should display additional debugging info.
	// All the resource generators need to return this flag.
	debug bool
}

func (grc *generatedResultCore) Success() bool {
	return grc.success
}

func (grc *generatedResultCore) Debug() bool {
	return grc.debug
}

func (grc *generatedResultCore) SetDebug(t bool) {
	grc.debug = t
}

func (grc *generatedResultCore) ErrorMessage() string {
	return grc.failureText
}
