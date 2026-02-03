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
	ID            int64  `json:"id_user"`
}

// Types untuk Input (Request)
type CreateWeightEntryPayload struct {
	Weight float64 `json:"weight_kg"`
	Date   string  `json:"date"` // Format: "YYYY-MM-DD" (Contoh: "2026-02-02")
	ID     string   `json:"id_user"`
}

// Types untuk Output (Response dari Database)
type WeightEntry struct {
	ID        int64   `json:"id"`
	UserID    string  `json:"user_id"`
	Weight    float64 `json:"weight_kg"`
	Date      string  `json:"date"`
	CreatedAt string  `json:"created_at"`
}
