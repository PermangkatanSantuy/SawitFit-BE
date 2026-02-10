package user

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
	// Profil
	router.Get("/me", h.handleGetUser)
	router.Post("/me", h.handleUpsertUser)
}

func (h *Handler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user from context
    authUser, err := auth.UserFromContext(r.Context())
    if err != nil {
        utils.WriteError(w, http.StatusUnauthorized, err)
        return
    }

    utils.WriteJSON(w, http.StatusOK, authUser)
}

func (h *Handler) handleUpsertUser(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user from context
    authUser, err := auth.UserFromContext(r.Context())
    if err != nil {
        utils.WriteError(w, http.StatusUnauthorized, err)
        return
    }

    utils.WriteJSON(w, http.StatusOK, authUser)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("ID user harus berupa angka"))
		return
	}

	// Parse Body JSON
	var payload types.UpdateUserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Hitung TDEE Sederhana
	tdee := 2000 // Default
	if payload.ActivityLevel == "Active" {
		tdee = 2500
	} else if payload.ActivityLevel == "Very Active" {
		tdee = 3000
	}

	// Sesuaikan dengan Goal
	if payload.FitnessGoal == "Cut" {
		tdee = tdee - 500
	} else if payload.FitnessGoal == "Bulk" {
		tdee = tdee + 300
	}

	// Siapkan Objek User
	user := types.User{
		ID:            authUser.Sub,
		Email:         authUser.Email, // TODO: Ambil dari Token JWT nanti
		FullName:      payload.FullName,
		FitnessGoal:   payload.FitnessGoal,
		ActivityLevel: payload.ActivityLevel,
		TDEE:          tdee,
	}

	if err := h.store.UpsertUserProfile(user); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Profil berhasil disimpan"})
}