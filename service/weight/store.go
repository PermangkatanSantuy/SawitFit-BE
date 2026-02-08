package weight

import (
	"context" // <-- Wajib ada buat pgx

	// Buat handle error
	"github.com/jackc/pgx/v5/pgxpool" // Library temanmu
	"github.com/types"
)

type Store struct {
	db *pgxpool.Pool // <-- Ganti jadi Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{db: db}
}

// ---------------------------------------------------------
//  Weight Entry Functions (TAMBAHAN PENTING)
// ---------------------------------------------------------

func (s *Store) CreateWeightEntry(entry types.WeightEntry) error {

	query := `
		INSERT INTO weight_entries (id_user, weight_kg, date) 
		VALUES ($1, $2, $3)
		ON CONFLICT (id_user, date) 
		DO UPDATE SET 
			weight_kg = EXCLUDED.weight_kg;
	`
	_, err := s.db.Exec(context.Background(), query, entry.UserID, entry.Weight, entry.Date)
	return err
}

func (s *Store) GetWeightHistory(userID string) ([]types.WeightEntry, error) {
	query := `
		SELECT id_weight_entry, id_user, weight_kg, date, created_at 
		FROM weight_entries 
		WHERE id_user = $1
		ORDER BY date DESC
	`

	rows, err := s.db.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []types.WeightEntry
	for rows.Next() {
		var w types.WeightEntry
		if err := rows.Scan(&w.ID, &w.UserID, &w.Weight, &w.Date, &w.CreatedAt); err != nil {
			return nil, err
		}
		history = append(history, w)
	}

	return history, nil
}

// kalo mau buat pencarian getweight berdasarkan tanggal khsusu buat function baru

func (s *Store) ResetWeightHistory(userID string) error {
	// Query untuk hapus data berat 
	query := `
		DELETE FROM weight_entries WHERE id_user = $1
	`

	// eksekusi querry nya
	_, err := s.db.Exec(context.Background(), query, userID)
	return err
}