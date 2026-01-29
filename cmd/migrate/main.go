package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load env (Sudah terbukti berhasil di root)
	if err := godotenv.Load(); err != nil {
		log.Println("Info: File .env tidak ditemukan, menggunakan environment system...")
	}

	// 2. Ambil URL (Kita sudah tahu namanya DATABASE_URL)
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("❌ Error: DATABASE_URL tidak ditemukan di .env")
	}

	// 3. Konfigurasi & Koneksi
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatal("❌ Error format URL:", err)
	}

	dbPool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatal("❌ Gagal connect ke database:", err)
	}
	defer dbPool.Close()

	// Cek koneksi sekilas
	if err := dbPool.Ping(context.Background()); err != nil {
		log.Fatal("❌ Gagal ping:", err)
	}
	fmt.Println("✅ Koneksi aman. Memulai proses pembuatan tabel...")

	// // 4. SQL Script (Sesuai tugas KAN-20)
	query := `
	-- Tabel Users
	CREATE TABLE IF NOT EXISTS users (
		id_user BIGSERIAL PRIMARY KEY,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		full_name TEXT,
		fitness_goal TEXT,
		activity_level TEXT,
		tdee INT,
		profile_picture TEXT,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);

	-- Tabel Food Entries
	CREATE TABLE IF NOT EXISTS food_entries (
		id_food_entry BIGSERIAL PRIMARY KEY,
		id_user BIGINT REFERENCES users(id_user) ON DELETE CASCADE,
		food_name TEXT NOT NULL,
		date DATE DEFAULT CURRENT_DATE,
		time TIME DEFAULT CURRENT_TIME,
		calories FLOAT DEFAULT 0,
		protein FLOAT DEFAULT 0,
		carbs FLOAT DEFAULT 0,
		fats FLOAT DEFAULT 0,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);

	-- Tabel Roles
	CREATE TABLE IF NOT EXISTS roles (
		id_role BIGSERIAL PRIMARY KEY,
		name TEXT UNIQUE NOT NULL,
		deskripsi TEXT,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);

	-- Tabel User Roles
	CREATE TABLE IF NOT EXISTS user_roles (
		id_user BIGINT REFERENCES users(id_user) ON DELETE CASCADE,
		id_role BIGINT REFERENCES roles(id_role) ON DELETE CASCADE,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		PRIMARY KEY (id_user, id_role)
	);

	-- Tabel Weight Entries
	CREATE TABLE IF NOT EXISTS weight_entries (
		id_weight_entry BIGSERIAL PRIMARY KEY,
		id_user BIGINT REFERENCES users(id_user) ON DELETE CASCADE,
		date DATE DEFAULT CURRENT_DATE,
		weight_kg FLOAT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);

	-- Tabel Workout Sessions
	CREATE TABLE IF NOT EXISTS workout_sessions (
		id_workout_session BIGSERIAL PRIMARY KEY,
		id_user BIGINT REFERENCES users(id_user) ON DELETE CASCADE,
		session_name TEXT,
		date DATE DEFAULT CURRENT_DATE,
		duration_minutes INT,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);

	-- Tabel Exercises
	CREATE TABLE IF NOT EXISTS exercises (
		id_exercise BIGSERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		muscle_group TEXT,
		equipment_type TEXT,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);

	-- Tabel Workout Sets
	CREATE TABLE IF NOT EXISTS workout_sets (
		id_workout_sets BIGSERIAL PRIMARY KEY,
		id_workout_session BIGINT REFERENCES workout_sessions(id_workout_session) ON DELETE CASCADE,
		id_exercise BIGINT REFERENCES exercises(id_exercise) ON DELETE CASCADE,
		set_number INT NOT NULL,
		weight_kg FLOAT NOT NULL,
		reps INT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);

	-- Index (Sesuai requirements)
	CREATE INDEX IF NOT EXISTS idx_food_user_date ON food_entries(id_user, created_at);
	CREATE INDEX IF NOT EXISTS idx_weight_user_date ON weight_entries(id_user, created_at);
	`

	// 5. Eksekusi
	_, err = dbPool.Exec(context.Background(), query)
	if err != nil {
		log.Fatal("❌ Gagal migrasi tabel:", err)
	}

	fmt.Println("🎉 SUKSES! Semua tabel (profiles, food, weight) dan index berhasil dibuat di Supabase.")
}
