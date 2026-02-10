package food

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi" 
	"github.com/internal/auth" 
	"github.com/types"
	"github.com/utils"
)

type Handler struct {
	store *Store
}

func NewHandler(store *Store) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router chi.Router) {
	
	router.Route("/food", func(r chi.Router) {
		r.Post("/", h.handleAddFood)
		r.Get("/{date}", h.handleGetFoodByDate)

		r.Delete("/reset", h.handleResetFoodEntries)
	})
}

func (h *Handler) handleAddFood(w http.ResponseWriter, r *http.Request) {
	// 1. AMBIL USER DARI CONTEXT 
	authUser, err := auth.UserFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("akses ditolak: %v", err))
		return
	}
	
	// ID User ada di authUser.Sub (String UUID)
	userID := authUser.Sub 

	// 2. Parse Data Makanan
	var payload types.FoodEntryPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Validasi
	if payload.FoodName == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("nama makanan wajib diisi"))
		return
	}

	// 3. Masukkan ke Struct (IDUser pake userID dari token)
	entry := types.FoodEntry{
		IDUser:   userID, 
		FoodName: payload.FoodName,
		Date:     payload.Date,
		Time:     payload.Time,
		Calories: payload.Calories,
		Protein:  payload.Protein,
		Carbs:    payload.Carbs,
		Fats:     payload.Fats,
	}

	if err := h.store.CreateFoodEntry(entry); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]string{"message": "Makanan berhasil dicatat!"})
}

func (h *Handler) handleGetFoodByDate(w http.ResponseWriter, r *http.Request) {
	// 1. AMBIL USER DARI CONTEXT
	authUser, err := auth.UserFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("akses ditolak: %v", err))
		return
	}
	userID := authUser.Sub

	// 2. Ambil Tanggal dari URL
	date := chi.URLParam(r, "date")
	if date == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("tanggal wajib diisi"))
		return
	}

	// 3. Panggil Store (Store kamu udah nerima string kan?)
	foods, err := h.store.GetFoodEntriesByDate(userID, date)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, foods)
}

func (h *Handler) handleResetFoodEntries(w http.ResponseWriter, r *http.Request) {
	authUser, err := auth.UserFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("akses ditolak: %v", err))
		return
	}
	userID := authUser.Sub

	if err := h.store.resetFoodEntries(userID); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]string{"message": "Entry makanan berhasil di reset"})
}