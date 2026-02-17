package weight

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/internal/auth"
	"github.com/types" // Import Types
	"github.com/utils" // Import Utils (Asumsi folder utils ada)
)

type Handler struct {
	store *Store
}

func NewHandler(store *Store) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router chi.Router) {

	router.Get("/me/weight", h.handleGetWeightPaginated)

	// Weight (Grouping)
	router.Route("/weight", func(r chi.Router) {
		r.Post("/", h.handleCreateWeight)
		r.Get("/", h.handleGetWeightHistory)

		// Rute untuk menghapus data berat badan
		r.Delete("/reset", h.HandlerResetWeight)
	})
}

// --- Handler Weight ---

func (h *Handler) handleCreateWeight(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user from context
	authUser, err := auth.UserFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, authUser)

	var payload types.CreateWeightEntryPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if payload.Weight < 0 {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("berat badan tidak valid"))
		return
	}

	// Panggil Store langsung
	err = h.store.CreateWeightEntry(types.WeightEntry{
		UserID: authUser.Sub,
		Weight: payload.Weight,
		Date:   payload.Date,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]string{"message": "Berat berhasil dicatat"})
}

func (h *Handler) handleGetWeightHistory(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user from context
	authUser, err := auth.UserFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, authUser)

	// untuk mendapatkan user id
	userID := authUser.Sub

	history, err := h.store.GetWeightHistory(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, history)
}

func (h *Handler) HandlerResetWeight(w http.ResponseWriter, r *http.Request) {
	// 1. mendapatkan data pengguna
	authUser, err := auth.UserFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
	}

	// Infoin identitas user
	utils.WriteJSON(w, http.StatusOK, authUser)

	// 2. mengambil id dari user ke string
	userID := authUser.Sub

	// 3. menjalankan function hapus/reset berat
	if err := h.store.ResetWeightHistory(userID); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]string{"message": "Berat berhasil di reset"})
}

func (h *Handler) handleGetWeightPaginated(w http.ResponseWriter, r *http.Request) {
	// Mengambil identitas user dari token
	authUser, err := auth.UserFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("Unauthorized: %v", err))
		return
	}
	// utils.WriteJSON(w, http.StatusOK, authUser)

	// Tangkap limit dari url akan diambil sebagai teks
	limitStr := r.URL.Query().Get("limit")

	// nilai default untuk limit dari data nya, 7 hari terakhir
	limit := 7

	if limitStr != "" {

		// konversi dari teks ke integer
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil || parsedLimit <= 0 {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Invalid limit parameter: %v", err))
			return
		}
		limit = parsedLimit
	}

	if limit > 30 {
		limit = 30
	}

	// Menangkap cursor dari url (opsional)
	cursor := r.URL.Query().Get("cursor")
	if cursor != "" {

		// karena dalam go tanda + pada string itu adalah spasi jadinya akan eror pada server
		cursor = strings.Replace(cursor, " 00", "+00", 1)

		// PERBAIKAN: Kita buatkan layout khusus yang persis dengan gaya tulisan PostgreSQL
		// 2006-01-02 = Format Tahun-Bulan-Tanggal standar Go
		// 15:04:05.999999 = Jam-Menit-Detik-Milidetik
		// -07 = Format Zona Waktu (+00)
		layoutPostgres := "2006-01-02 15:04:05.999999-07"

		// Memastikan bahwa cursor merupakan timestamp
		_, err := time.Parse(layoutPostgres, cursor)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Invalid cursor parameter: %v", err))
			return
		}
	}

	entries, err := h.store.GetWeightHistoryPaginated(authUser.Sub, limit, cursor)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("Gagal mengambil data: %v", err))
		return
	}

	// Buat cursor baru untuk halaman selanjutnya
	var nextCursor string
	if len(entries) > 0 { // jika datanya lebih dari 1
		nextCursor = entries[len(entries)-1].CreatedAt // maka next nya akan menjadi data terakhir
	}

	// Membuat balasan json
	response := map[string]interface{}{
		"data": entries,
	}

	if nextCursor != "" {
		response["next cursor"] = nextCursor
	}

	utils.WriteJSON(w, http.StatusOK, response)
}
