package number

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Next generates a document number like GRN-20260605-00001.
// It uses a sequences table to guarantee uniqueness across concurrent requests.
func Next(ctx context.Context, pool *pgxpool.Pool, prefix string) (string, error) {
	var seq int64
	err := pool.QueryRow(ctx,
		`INSERT INTO sequences (prefix, last_value)
		 VALUES ($1, 1)
		 ON CONFLICT (prefix) DO UPDATE SET last_value = sequences.last_value + 1
		 RETURNING last_value`,
		prefix,
	).Scan(&seq)
	if err != nil {
		return "", fmt.Errorf("next sequence: %w", err)
	}

	date := time.Now().Format("20060102")
	return fmt.Sprintf("%s-%s-%05d", prefix, date, seq), nil
}
