package types

import "time"

// User: Sesuai tabel 'users' di Supabase
// ID User menggunakan int64 karena di database tipenya int8 (BigInt)
type User struct {
	ID             int64     `json:"id_user"`
	Email          string    `json:"email"`
	Password       string    `json:"-"` // Tidak dikirim ke JSON
	FullName       string    `json:"full_name"`
	FitnessGoal    string    `json:"fitness_goal"`
	ActivityLevel  string    `json:"activity_level"`
	TDEE           int       `json:"tdee"`
	ProfilePicture string    `json:"profile_picture"`
	CreatedAt      time.Time `json:"created_at"`
}

// Payload: Data yang dikirim Frontend saat update profil
type UpdateUserPayload struct {
	FullName      string `json:"full_name" validate:"required"`
	FitnessGoal   string `json:"fitness_goal"`
	ActivityLevel string `json:"activity_level"`
}

type FoodEntry struct {
	ID        int64     `json:"id_food_entry"`
	IDUser    int64     `json:"id_user"`
	FoodName  string    `json:"food_name"`
	Date      string    `json:"date"`     // Format YYYY-MM-DD
	Time      string    `json:"time"`     // Format HH:MM:SS
	Calories  float64   `json:"calories"` // float8 di DB = float64 di Go
	Protein   float64   `json:"protein"`
	Carbs     float64   `json:"carbs"`
	Fats      float64   `json:"fats"`     // Pake 's' sesuai screenshot
	CreatedAt time.Time `json:"created_at"`
}

// Data yang dikirim User lewat Postman/Android
type FoodEntryPayload struct {
	FoodName string  `json:"food_name"`
	Date     string  `json:"date"`
	Time     string  `json:"time"`
	Calories float64 `json:"calories"`
	Protein  float64 `json:"protein"`
	Carbs    float64 `json:"carbs"`
	Fats     float64 `json:"fats"`
}