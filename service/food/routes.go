package food

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/types"
	"github.com/utils"
)

type Handler struct {
	store *Store
}

func NewHandler(store *Store) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	// Endpoint: Tambah Makanan
	// Contoh: POST /api/v1/users/1001/food
	router.HandleFunc("/users/{userID}/food", h.handleAddFood).Methods("POST")

	// Endpoint: Lihat Makanan per Tanggal
	// Contoh: GET /api/v1/users/1001/food/2026-01-30
	router.HandleFunc("/users/{userID}/food/{date}", h.handleGetFoodByDate).Methods("GET")
}

func (h *Handler) handleAddFood(w http.ResponseWriter, r *http.Request) {
	// 1. Ambil UserID dari URL
	vars := mux.Vars(r)
	userIDStr := vars["userID"]
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("ID user tidak valid"))
		return
	}

	// 2. Baca Data JSON yang dikirim
	var payload types.FoodEntryPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Validasi sederhana
	if payload.FoodName == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("nama makanan wajib diisi"))
		return
	}

	// 3. Masukkan ke Struct Database
	entry := types.FoodEntry{
		IDUser:   userID,
		FoodName: payload.FoodName,
		Date:     payload.Date,     // Format: "2026-01-30"
		Time:     payload.Time,     // Format: "12:30:00"
		Calories: payload.Calories,
		Protein:  payload.Protein,
		Carbs:    payload.Carbs,
		Fats:     payload.Fats,
	}

	// 4. Panggil Store buat simpan ke DB
	if err := h.store.CreateFoodEntry(entry); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]string{"message": "Makanan berhasil dicatat!"})
}

func (h *Handler) handleGetFoodByDate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	
	// 1. Ambil User ID
	userIDStr := vars["userID"]
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("ID user tidak valid"))
		return
	}

	// 2. Ambil Tanggal (Format YYYY-MM-DD)
	date := vars["date"]
	if date == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("tanggal wajib diisi"))
		return
	}

	// 3. Panggil Store
	foods, err := h.store.GetFoodEntriesByDate(userID, date)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// 4. Kirim Balikan JSON
	utils.WriteJSON(w, http.StatusOK, foods)
}