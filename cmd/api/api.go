package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool" // <-- Pakai library temanmu
	"github.com/service/user"
)

type APIServer struct {
	addr string
	db   *pgxpool.Pool // <-- Ubah dari *sql.DB jadi *pgxpool.Pool
}

func NewAPIServer(addr string, db *pgxpool.Pool) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	log.Println("Server berjalan di", s.addr)
	return http.ListenAndServe(s.addr, router)
}