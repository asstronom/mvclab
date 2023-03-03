package webview

import (
	"net/http"

	"example.com/pkg/domain"
	"github.com/gorilla/mux"
)

type Server struct {
	service domain.UserService
	router  *mux.Router
}

func NewServer(service domain.UserService) (*Server, error) {
	server := Server{service: service}
	initTemplates()
	server.router = mux.NewRouter()
	server.initRoutes()
	return &server, nil
}

func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(":9091", s.router)
}
