package listing

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Better: pass *pgxpool.Pool and open per-call connections
type PgxPool interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

type Store struct{ P PgxPool }

func (s *Store) Create(ctx context.Context, userID int64, p CreateParams) (Listing, error) {
	const q = `
    INSERT INTO listings(title, description, price, category, user_id)
    VALUES ($1,$2,$3,$4,$5)
    RETURNING id, title, description, price, category, user_id, status, created_at`
	var l Listing
	err := s.P.QueryRow(ctx, q, p.Title, p.Description, p.Price, p.Category, userID).
		Scan(&l.ID, &l.Title, &l.Description, &l.Price, &l.Category, &l.UserID, &l.Status, &l.CreatedAt)
	return l, err
}

func (s *Store) Get(ctx context.Context, id int64) (Listing, error) {
	const q = `SELECT id,title,description,price,category,user_id,status,created_at FROM listings WHERE id=$1`
	var l Listing
	err := s.P.QueryRow(ctx, q, id).
		Scan(&l.ID, &l.Title, &l.Description, &l.Price, &l.Category, &l.UserID, &l.Status, &l.CreatedAt)
	return l, err
}

type ListFilters struct {
	Q        *string
	Category *Category
	Status   *Status
	MinPrice *int64
	MaxPrice *int64
	Limit    int
	Offset   int
	Sort     string // "created_at_desc", "price_asc", etc.
}

func (s *Store) List(ctx context.Context, f ListFilters) ([]Listing, error) {
	sb := strings.Builder{}
	sb.WriteString(`SELECT id,title,description,price,category,user_id,status,created_at FROM listings`)
	var where []string
	var args []any
	i := 1

	if f.Q != nil && *f.Q != "" {
		where = append(where, fmt.Sprintf("(title ILIKE $%d OR description ILIKE $%d)", i, i+1))
		args = append(args, "%"+*f.Q+"%", "%"+*f.Q+"%")
		i += 2
	}
	if f.Category != nil {
		where = append(where, fmt.Sprintf("category = $%d", i))
		args = append(args, *f.Category)
		i++
	}
	if f.Status != nil {
		where = append(where, fmt.Sprintf("status = $%d", i))
		args = append(args, *f.Status)
		i++
	}
	if f.MinPrice != nil {
		where = append(where, fmt.Sprintf("price >= $%d", i))
		args = append(args, *f.MinPrice)
		i++
	}
	if f.MaxPrice != nil {
		where = append(where, fmt.Sprintf("price <= $%d", i))
		args = append(args, *f.MaxPrice)
		i++
	}

	if len(where) > 0 {
		sb.WriteString(" WHERE " + strings.Join(where, " AND "))
	}

	switch f.Sort {
	case "price_asc":
		sb.WriteString(" ORDER BY price ASC")
	case "price_desc":
		sb.WriteString(" ORDER BY price DESC")
	default:
		sb.WriteString(" ORDER BY created_at DESC")
	}

	if f.Limit <= 0 || f.Limit > 100 {
		f.Limit = 20
	}
	sb.WriteString(fmt.Sprintf(" LIMIT %d OFFSET %d", f.Limit, f.Offset))

	rows, err := s.P.Query(ctx, sb.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Listing
	for rows.Next() {
		var l Listing
		if err := rows.Scan(&l.ID, &l.Title, &l.Description, &l.Price, &l.Category, &l.UserID, &l.Status, &l.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, rows.Err()
}

func (s *Store) Update(ctx context.Context, id int64, userid int64, p UpdateParams) (Listing, error) {
	// build dynamic SET
	var sets []string
	var args []any
	i := 1

	if p.Title != nil {
		sets = append(sets, fmt.Sprintf("title=$%d", i))
		args = append(args, *p.Title)
		i++
	}
	if p.Description != nil {
		sets = append(sets, fmt.Sprintf("description=$%d", i))
		args = append(args, *p.Description)
		i++
	}
	if p.Price != nil {
		sets = append(sets, fmt.Sprintf("price=$%d", i))
		args = append(args, *p.Price)
		i++
	}
	if p.Category != nil {
		sets = append(sets, fmt.Sprintf("category=$%d", i))
		args = append(args, *p.Category)
		i++
	}
	if p.Status != nil {
		sets = append(sets, fmt.Sprintf("status=$%d", i))
		args = append(args, *p.Status)
		i++
	}

	if len(sets) == 0 {
		return s.Get(ctx, id)
	}

	q := fmt.Sprintf(
		`UPDATE listings SET %s WHERE id=$%d AND user_id=%d RETURNING id,title,description,price,category,user_id,status,created_at`,
		strings.Join(sets, ","),
		i,
		userid,
	)
	args = append(args, id)

	var l Listing
	err := s.P.QueryRow(ctx, q, args...).
		Scan(&l.ID, &l.Title, &l.Description, &l.Price, &l.Category, &l.UserID, &l.Status, &l.CreatedAt)
	return l, err
}

func (s *Store) Archive(ctx context.Context, id int64) error {
	_, err := s.P.Exec(ctx, `UPDATE listings SET status='ARCHIVED' WHERE id=$1`, id)
	return err
}

func (s *Store) Delete(ctx context.Context, id int64) error {
	_, err := s.P.Exec(ctx, `DELETE FROM listings WHERE id=$1`, id)
	return err
}
