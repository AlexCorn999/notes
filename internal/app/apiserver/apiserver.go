package apiserver

import (
	"encoding/json"
	"errors"
	"fmt"

	"net/http"

	"github.com/AlexCorn999/notes/internal/app/logger"
	"github.com/AlexCorn999/notes/internal/app/model"
	"github.com/AlexCorn999/notes/internal/app/note"
	"github.com/AlexCorn999/notes/internal/app/store"
	"github.com/AlexCorn999/notes/internal/app/tokens"
	"github.com/AlexCorn999/notes/internal/app/yandex"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt"
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

// configureRouter sets endpoints
func (s *APIServer) configureRouter() {
	s.router.Use(logger.WithLogging)
	s.router.Post("/user", s.userCreate())
	s.router.HandleFunc("/join", s.validateJWT(s.login()))
	s.router.HandleFunc("/create", s.validateJWT(s.addNote()))
	s.router.HandleFunc("/notes", s.validateJWT(s.getList()))
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

// UserCreate adds the user to the database
func (s *APIServer) userCreate() http.HandlerFunc {
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

		u, err := s.store.User().Create(u)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		tok, err := tokens.GenerateToken(req.Email, req.Password)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}

		result := fmt.Sprintf("token - %s       user_id - %d", tok, u.ID)

		s.respond(w, r, http.StatusCreated, result)
	}
}

// login authorizes the user
func (s *APIServer) login() http.HandlerFunc {
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
			return
		}

		if req.Password != user.Password {
			s.error(w, r, http.StatusUnauthorized, errors.New("incorrect password"))
			return
		}

		s.respond(w, r, http.StatusOK, "Успешная авторизация :)")

	}
}

// addNote adds a note
func (s *APIServer) addNote() http.HandlerFunc {
	type request struct {
		User_id int    `json:"user_id"`
		Title   string `json:"title"`
		Text    string `json:"text"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if err := s.store.User().GetUser(req.User_id); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		n := &note.Note{
			Title: req.Title,
			Text:  req.Text,
		}

		if err := n.Validate(); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		if err := yandex.CheckAll(n); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
		}

		s.store.Note().CreateNote(n, req.User_id)

		s.respond(w, r, http.StatusCreated, n)
	}
}

// getList displays all user notes
func (s *APIServer) getList() http.HandlerFunc {
	type request struct {
		User_id int `json:"user_id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if err := s.store.User().GetUser(req.User_id); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		result, err := s.store.Note().GetList(req.User_id)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		titles := make([]string, 0)
		for _, note := range result {
			titles = append(titles, note.Title)
		}

		s.respond(w, r, http.StatusOK, titles)
	}
}

// ValidateJWT checks the validity of the token
func (s *APIServer) validateJWT(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {

			token, err := jwt.Parse(r.Header["Token"][0], func(t *jwt.Token) (interface{}, error) {
				_, ok := t.Method.(*jwt.SigningMethodHMAC)
				if !ok {
					s.error(w, r, http.StatusUnauthorized, nil)
				}
				return tokens.MySignKey, nil
			})

			if err != nil {
				s.error(w, r, http.StatusUnauthorized, err)
			}

			if token.Valid {
				next(w, r)
			}

		} else {
			s.error(w, r, http.StatusUnauthorized, errors.New("no autorized"))
		}
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
