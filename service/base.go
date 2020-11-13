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
	s.Router.HandleFunc("/products", s.GetProducts).Methods("GET")
	s.Router.HandleFunc("/product", s.CreateProduct).Methods("POST")
	s.Router.HandleFunc("/product/{id:[0-9]+}", s.GetProduct).Methods("GET")
	s.Router.HandleFunc("/product/{id:[0-9]+}", s.UpdateProduct).Methods("PUT")
	s.Router.HandleFunc("/product/{id:[0-9]+}", s.DeleteProduct).Methods("DELETE")
}
