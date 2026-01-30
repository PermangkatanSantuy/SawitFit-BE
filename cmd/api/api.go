package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	
	// Import Service User & Food
	"github.com/service/user"
	"github.com/service/food" 
)

type APIServer struct {
	addr string
	db   *pgxpool.Pool
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

	// --- 1. FITUR USER (Hasil Merge db/crud) ---
	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	// --- 2. FITUR FOOD (Baru Ditambahkan) ---
	foodStore := food.NewStore(s.db)
	foodHandler := food.NewHandler(foodStore)
	foodHandler.RegisterRoutes(subrouter)
	// ----------------------------------------

	log.Println("Server berjalan di", s.addr)
	return http.ListenAndServe(s.addr, router)
}