package webserver

import (
	"github.com/go-chi/chi/v5/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Route struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
}

type WebServer struct {
	Router        chi.Router
	Handlers      []Route
	WebServerPort string
}

func NewWebServer(serverPort string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      make([]Route, 0),
		WebServerPort: serverPort,
	}
}

func (s *WebServer) AddHandler(method, path string, handler http.HandlerFunc) {
	s.Handlers = append(s.Handlers, Route{Method: method, Path: path, Handler: handler})
}

// loop through the handlers and add them to the router
// register middeleware logger
// start the server
func (s *WebServer) Start() {
	s.Router.Use(middleware.Logger)

	for _, route := range s.Handlers {
		switch route.Method {
		case http.MethodGet:
			s.Router.Get(route.Path, route.Handler)
		case http.MethodPost:
			s.Router.Post(route.Path, route.Handler)
		// Adicione outros métodos se precisar
		case http.MethodPut:
			s.Router.Put(route.Path, route.Handler)
		case http.MethodDelete:
			s.Router.Delete(route.Path, route.Handler)
		}
	}
	http.ListenAndServe(s.WebServerPort, s.Router)
}
