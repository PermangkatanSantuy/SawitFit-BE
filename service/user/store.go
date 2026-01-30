package user

import (
	"context" // <-- Wajib ada buat pgx
	"fmt"

	"github.com/jackc/pgx/v5" // Buat handle error
	"github.com/jackc/pgx/v5/pgxpool" // Library temanmu
	"github.com/types"
)

type Store struct {
	db *pgxpool.Pool // <-- Ganti jadi Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{db: db}
}

func (s *Store) GetUserByID(userID int64) (*types.User, error) {
	// PERBAIKAN: Gunakan COALESCE(kolom, '')
	// Artinya: Kalau datanya NULL, ganti jadi string kosong ('').
	// Ini berlaku buat email juga kalau-kalau dia NULL.

	query := `
		SELECT 
			id_user, 
			COALESCE(email, ''), 
			full_name, 
			fitness_goal, 
			activity_level, 
			tdee, 
			COALESCE(profile_picture, ''), 
			created_at 
		FROM users 
		WHERE id_user = $1`

	// Perhatikan: Scan-nya harus urut sesuai SELECT di atas
	row := s.db.QueryRow(context.Background(), query, userID)

	var u types.User
	err := row.Scan(
		&u.ID,
		&u.Email,          // Sekarang aman karena sudah di-COALESCE
		&u.FullName,
		&u.FitnessGoal,
		&u.ActivityLevel,
		&u.TDEE,
		&u.ProfilePicture, // Sekarang aman karena sudah di-COALESCE
		&u.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return &types.User{ID: userID}, nil
		}
		return nil, err
	}

	return &u, nil
}

func (s *Store) UpsertUserProfil(u types.User) error {
	// PERBAIKAN: Kita tambahkan kolom 'password' dan isi dengan string random ("-")
	// Ini supaya database tidak error "violates not-null".
	// Password asli nanti diurus fitur Register/Login, bukan di sini.

	query := `
		INSERT INTO users (id_user, email, password, full_name, fitness_goal, activity_level, tdee, created_at)
		VALUES ($1, $2, 'apa_baelah', $3, $4, $5, $6, NOW())
		ON CONFLICT (id_user) DO UPDATE SET
			full_name = EXCLUDED.full_name,
			fitness_goal = EXCLUDED.fitness_goal,
			activity_level = EXCLUDED.activity_level,
			tdee = EXCLUDED.tdee;
	`

	// Perhatikan urutan parameter ($1 s/d $6)
	_, err := s.db.Exec(context.Background(), query,
		u.ID,            // $1
		u.Email,         // $2
		u.FullName,      // $3
		u.FitnessGoal,   // $4
		u.ActivityLevel, // $5
		u.TDEE,          // $6
	)

	if err != nil {
		return fmt.Errorf("gagal upsert user: %v", err)
	}

	return nil
}