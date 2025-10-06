package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Nikhil1169/search-agent/internal/model"
	_ "github.com/lib/pq"
)

type Store struct {
	db *sql.DB
}

func NewStore(connStr string) (*Store, error) {
	// --- THIS IS THE FIXED LINE ---
	db, err := sql.Open("postgres", connStr) // Changed conn_str to connStr
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s *Store) FindListingsByParams(ctx context.Context, params *model.SearchParams) ([]model.Listing, error) {
	var args []interface{}
	query := `SELECT id, title, description, price, category, user_id, status, created_at
			  FROM listings WHERE status = 'AVAILABLE'`
	orderBy := "ORDER BY created_at DESC"

	if len(params.Keywords) > 0 {
		var keywordConditions []string
		for _, keyword := range params.Keywords {
			args = append(args, "%"+keyword+"%")
			argIndex := len(args)
			keywordConditions = append(keywordConditions, fmt.Sprintf("title ILIKE $%d", argIndex), fmt.Sprintf("description ILIKE $%d", argIndex))
		}
		query += " AND (" + strings.Join(keywordConditions, " OR ") + ")"

		if params.Category != "" {
			args = append(args, params.Category)
			orderBy = fmt.Sprintf("ORDER BY CASE WHEN category = $%d THEN 0 ELSE 1 END, created_at DESC", len(args))
		}

	} else if params.Category != "" {
		args = append(args, params.Category)
		query += fmt.Sprintf(" AND category = $%d", len(args))
	}

	if params.MaxPrice > 0 {
		args = append(args, params.MaxPrice*100)
		query += fmt.Sprintf(" AND price <= $%d", len(args))
	}
	if params.MinPrice > 0 {
		args = append(args, params.MinPrice*100)
		query += fmt.Sprintf(" AND price >= $%d", len(args))
	}

	query += " " + orderBy + " LIMIT 20"

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("database query failed: %w", err)
	}
	defer rows.Close()

	var listings []model.Listing
	for rows.Next() {
		var l model.Listing
		if err := rows.Scan(&l.ID, &l.Title, &l.Description, &l.Price, &l.Category, &l.UserID, &l.Status, &l.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		listings = append(listings, l)
	}

	return listings, nil
}
