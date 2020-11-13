package service

import (
	"database/sql"

	"github.com/gorilla/mux"
)

type Server struct {
	DB     *sql.DB
	Router *mux.Router
}

func New(db *sql.DB, r *mux.Router) *Server {
	s := Server{
		DB:     db,
		Router: r,
	}
	return &s
}

func (s *Server) InitializeRoutes() {
	s.Router.HandleFunc("/players", s.GetPlayers).Methods("GET")
	s.Router.HandleFunc("/player", s.CreatePlayer).Methods("POST")
	s.Router.HandleFunc("/player/{id:[0-9]+}", s.GetPlayer).Methods("GET")
	s.Router.HandleFunc("/player/{id:[0-9]+}", s.UpdatePlayer).Methods("PUT")
	s.Router.HandleFunc("/player/{id:[0-9]+}", s.DeletePlayer).Methods("DELETE")
}
