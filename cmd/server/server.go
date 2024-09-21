package server

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/nikhil-thorat/goauth/internals/user"
)

type Server struct {
	addr string
	db   *sql.DB
}

func NewServer(addr string, db *sql.DB) *Server {
	return &Server{addr: addr, db: db}
}

func (s *Server) Run() error {
	log.Printf("SERVER STARTED ON : http://localhost:%v\n", s.addr)

	router := http.NewServeMux()

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(router)

	return http.ListenAndServe(":"+s.addr, router)
}
