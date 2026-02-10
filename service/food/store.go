package food

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/types"
)

type Store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{db: db}
}

func (s *Store) CreateFoodEntry(entry types.FoodEntry) error {
	query := `
		INSERT INTO food_entries (id_user, food_name, date, time, calories, protein, carbs, fats, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
	`
	_, err := s.db.Exec(context.Background(), query,
		entry.IDUser, 
		entry.FoodName,
		entry.Date,
		entry.Time,
		entry.Calories,
		entry.Protein,
		entry.Carbs,
		entry.Fats,
	)
	if err != nil {
		return fmt.Errorf("gagal input makanan: %v", err)
	}
	return nil
}

func (s *Store) GetFoodEntriesByDate(userID string, date string) ([]types.FoodEntry, error) {
	query := `
		SELECT id_food_entry, id_user, food_name, date::text, time::text, calories, protein, carbs, fats, created_at
		FROM food_entries
		WHERE id_user = $1 AND date = $2
		ORDER BY time ASC
	`
	rows, err := s.db.Query(context.Background(), query, userID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entries := []types.FoodEntry{}
	for rows.Next() {
		var e types.FoodEntry
		err := rows.Scan(&e.ID, &e.IDUser, &e.FoodName, &e.Date, &e.Time, &e.Calories, &e.Protein, &e.Carbs, &e.Fats, &e.CreatedAt)
		if err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}