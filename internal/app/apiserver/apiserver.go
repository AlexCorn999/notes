package apiserver

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/AlexCorn999/notes/internal/app/model"
	"github.com/AlexCorn999/notes/internal/app/note"
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
	s.router.Post("/user", s.handleUsersCreate())
	s.router.Post("/join", s.handleSession())
	s.router.Post("/create", s.AddNote())
	s.router.Get("/notes", s.GetList())
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

// handleUsersCreate adds the user to the database
func (s *APIServer) handleUsersCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := &model.User{
			Email:    req.Email,
			Password: req.Password,
		}

		if err := u.Validate(); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		if _, err := s.store.User().Create(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusCreated, u)
	}
}

func (s *APIServer) handleSession() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		user, err := s.store.User().FindByEmail(req.Email)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, err)
		}

		if req.Password != user.Password {
			s.error(w, r, http.StatusUnauthorized, errors.New("incorrect password"))
		}

		s.respond(w, r, http.StatusOK, nil)

	}
}

func (s *APIServer) AddNote() http.HandlerFunc {
	type request struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		n := &note.Note{
			Title: req.Title,
			Text:  req.Text,
		}

		// поменять id
		s.store.Note().CreateNote(n, 1)

		s.respond(w, r, http.StatusCreated, n)
	}
}

func (s *APIServer) GetList() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// поменять id
		result, err := s.store.Note().GetList(1)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		s.respond(w, r, http.StatusOK, result)
	}
}

func (s *APIServer) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *APIServer) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
