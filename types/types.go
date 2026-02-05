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