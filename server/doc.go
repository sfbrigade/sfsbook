package server

// The source code is arranged in the following pattern. There are
// two kinds of init: handler factory tooling and handler creation.
//  Tooling contains all error-causing init. 
// 
// A client
// of the server package creates a HandlerFactory object and can take
// action (as desired) to correct error situations.
//
// The HandlerFactory combines tooling objects. An assortment of
// methods named make*Tooling create the tooling that is kept in
// the HandlerFactory object.
//
// Once the HandlerFactory is built, then it can make handlers. 
// http.Handler creation may not return errors, only http.Handler
// implementations.
//
// Some
// handlers will mutate the request context and delegate handling to
// a different http.Handler instance.
//
// The make*Handler functions should be kept with the handler
// implementation.
