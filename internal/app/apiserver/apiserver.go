package apiserver

import (
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type APIServer struct {
	logger *logrus.Logger
	router *chi.Mux
}

// NewAPI returns a new API server
func NewAPI() *APIServer {
	return &APIServer{
		logger: logrus.New(),
		router: chi.NewRouter(),
	}
}

// Start starts the server
func (s *APIServer) Start() error {
	s.configureLogger()
	s.configureRouter()

	s.logger.Info("starting api server")
	return http.ListenAndServe(":8080", s.router)
}

// configureLogger sets the logging level
func (s *APIServer) configureLogger() {
	s.logger.SetLevel(logrus.DebugLevel)
}

func (s *APIServer) configureRouter() {
	s.router.HandleFunc("/hello", s.handleHello())
}

func (s *APIServer) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello from handle")
	}
}
