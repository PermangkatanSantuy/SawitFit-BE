package weight

import (
	"fmt"
	"net/http"

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
