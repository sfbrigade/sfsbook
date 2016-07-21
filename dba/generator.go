// Package dba provides the backend of the sfsbook server. There are different
// generator functions for each kind of page.
package dba

// Generator is an abstract interface implemented by each page's middle-tier.
// The server directory contains the front tier, Generators are the middle-tier, etc.
// A Generator instance knows how to convert a specific request into a structure
// that will be used to populate a template appropriate to the page.
type Generator interface {
	ForRequest(req interface{}) interface{}
}
