package server

import (
	"context"
	"golang-jwt-example/pkg/handler"
	"golang-jwt-example/pkg/middleware"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Config struct {
	Log *zap.Logger
}

type Server struct {
	Mux        *http.ServeMux
	Handler    http.Handler
	server     *http.Server
	handler    *handler.Handler
	middleware *middleware.Middleware
	log        *zap.Logger
}

func NewServer(handler *handler.Handler, middleware *middleware.Middleware, cfg *Config) *Server {
	s := &Server{
		Mux:        http.NewServeMux(),
		handler:    handler,
		middleware: middleware,
	}

	if cfg != nil {
		if log := cfg.Log; log != nil {
			s.log = log
		}
	}

	s.registerHandler()
	return s
}

func (s *Server) Serve(listener net.Listener) error {
	server := &http.Server{
		Handler: s.Mux,
	}
	s.server = server
	if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *Server) GracefulShutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) registerHandler() {
	r := mux.NewRouter()

	st := r.PathPrefix("/student/").Subrouter()
	st.Handle("/classes", s.middleware.User.Auth(s.handler.User.GetUser())).Methods(http.MethodGet)
	//st.Handle("/login", s.handler.Student.PostLogin()).Methods(http.MethodPost)
	//st.Handle("/logout", s.middleware.Student.Auth(s.handler.Student.DeleteLogout())).Methods(http.MethodDelete)
	//st.Handle("/teachers", s.middleware.Student.Auth(s.handler.Student.GetTeachers())).Methods(http.MethodGet)
	//st.Handle("/me", s.middleware.Student.Auth(s.handler.Student.GetMe())).Methods(http.MethodGet)
	//st.Handle("/me", s.middleware.Student.Auth(s.handler.Student.PutMe())).Methods(http.MethodPut)
	//st.Handle("/me/class", s.middleware.Student.Auth(s.handler.Student.GetMeClass())).Methods(http.MethodGet)
	//st.Handle("/me/grade", s.middleware.Student.Auth(s.handler.Student.GetMeGrade())).Methods(http.MethodGet)
	//st.Handle("/me/school", s.middleware.Student.Auth(s.handler.Student.GetMeSchool())).Methods(http.MethodGet)
	//st.Handle("/me/corporation", s.middleware.Student.Auth(s.handler.Student.GetMeCorporation())).Methods(http.MethodGet)
	//st.Handle("/students", s.middleware.Student.Auth(s.handler.Student.GetStudents())).Methods(http.MethodGet)
	//st.Handle("/students/{id}", s.middleware.Student.Auth(s.handler.Student.GetStudent())).Methods(http.MethodGet)
	//st.Handle("/token", s.handler.Student.PostToken()).Methods(http.MethodPost)

	//r.Handle("/healthz", s.handler.General.HealthCheck())
	//r.Handle("/version", s.handler.General.Version())
	//r.NotFoundHandler = s.handler.General.NotFound()

	s.Mux.Handle("/", r)
}
