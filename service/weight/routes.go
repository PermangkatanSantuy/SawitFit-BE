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

	var payload types.CreateWeightEntryPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	history, err := h.store.GetWeightHistory(payload.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, history)
}