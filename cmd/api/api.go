package api

import (
	"log"
	"net/http"

	// "github.com/gorilla/mux"
	"github.com/go-chi/chi"            // untuk handler link yang dituju
	"github.com/go-chi/chi/middleware" // Middleware bawaan untuk log
	"github.com/internal/auth"
	"github.com/jackc/pgx/v5/pgxpool" // <-- Pakai library temanmu
	"github.com/service/user"
	"github.com/service/weight"
)

type APIServer struct {
	addr string
	db   *pgxpool.Pool // <-- Ubah dari *sql.DB jadi *pgxpool.Pool
	verifier *auth.Verifier
}

func NewAPIServer(addr string, db *pgxpool.Pool, verifier *auth.Verifier) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
		verifier: verifier,
	}
}

func (s *APIServer) Run() error {
	// 1. Inisialisasi Router Chi
	router := chi.NewRouter()

	// 2. Pasang Middleware (Opsional tapi sangat disarankan)
	// Logger: Agar kamu bisa lihat di terminal siapa yang akses API
	// Recoverer: Agar server tidak mati total (crash) kalau ada error fatal
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(auth.Middleware(s.verifier))

	// 3. Setup Service User
	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	
	// 3. Setup Service User
	weightStore := weight.NewStore(s.db)
	weightHandler := weight.NewHandler(weightStore)

	// Fungsi: Membuat "folder" khusus agar semua link diawali dengan "/api/v1".
	router.Route("/api/v1", func(r chi.Router) {
		
		// Kita mengirim 'r' (router khusus grup ini) ke handler.
		// HASIL: Link "/me" di dalam userHandler otomatis berubah
		// dari "localhost:8080/me" MENJADI "localhost:8080/api/v1/me"
		userHandler.RegisterRoutes(r)
		weightHandler.RegisterRoutes(r)
	})

	log.Println("Server berjalan di", s.addr)
	return http.ListenAndServe(s.addr, router)
}