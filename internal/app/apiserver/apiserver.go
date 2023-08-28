package apiserver

import (
	"io"
	"net/http"

	"github.com/AlexCorn999/notes/internal/app/store"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type APIServer struct {
	logger *logrus.Logger
	router *chi.Mux
	store  *store.Store
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

	if err := s.configureStore(); err != nil {
		return err
	}

	s.logger.Info("starting api server")
	return http.ListenAndServe(":8080", s.router)
}

// configureLogger sets the logging level
func (s *APIServer) configureLogger() {
	s.logger.SetLevel(logrus.DebugLevel)
}

// configureRouter
func (s *APIServer) configureRouter() {
	s.router.HandleFunc("/hello", s.handleHello())
}

// configureStore opens a connection to the database and assigns a value to the server database
func (s *APIServer) configureStore() error {
	st := store.NewStore()
	if err := st.Open(); err != nil {
		return err
	}
	s.store = st
	return nil
}

func (s *APIServer) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello from handle")
	}
}
