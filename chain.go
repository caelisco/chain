package chain

import (
	"net/http"
)

// Mux is a http.Handler which dispatches requests to different handlers with middleware chaining.
type Mux struct {
	router      *http.ServeMux
	middlewares []func(http.Handler) http.Handler
}

// New returns a new initialized Mux instance.
func New() *Mux {
	return &Mux{
		router: http.NewServeMux(),
	}
}

// Use registers middleware with the Mux instance. Middleware must have the signature `func(http.Handler) http.Handler`.
func (m *Mux) Use(mw ...func(http.Handler) http.Handler) {
	m.middlewares = append(m.middlewares, mw...)
}

// Group is used to create 'groups' of routes in a Mux. Middleware registered inside the group will only be used on the routes in that group.
func (m *Mux) Group(fn func(*Mux)) {
	groupMux := &Mux{
		router:      m.router,
		middlewares: append([]func(http.Handler) http.Handler{}, m.middlewares...),
	}
	fn(groupMux)
}

// Handle registers a new handler for the given request pattern and HTTP method.
func (m *Mux) Handle(pattern string, handler http.Handler) {
	m.router.Handle(pattern, m.wrap(handler))
}

// HandleFunc is an adapter which allows using a http.HandlerFunc as a handler.
func (m *Mux) HandleFunc(pattern string, handlerFunc http.HandlerFunc) {
	m.router.Handle(pattern, m.wrap(handlerFunc))
}

// ServeHTTP makes the router implement the http.Handler interface.
func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.router.ServeHTTP(w, r)
}

func (m *Mux) wrap(handler http.Handler) http.Handler {
	for i := len(m.middlewares) - 1; i >= 0; i-- {
		handler = m.middlewares[i](handler)
	}
	return handler
}
