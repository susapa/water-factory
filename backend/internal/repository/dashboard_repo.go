package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/water-factory/api/internal/domain"
)

type DashboardRepo struct {
	pool *pgxpool.Pool
}

func NewDashboardRepo(pool *pgxpool.Pool) *DashboardRepo {
	return &DashboardRepo{pool: pool}
}

func (r *DashboardRepo) GetSummary(ctx context.Context) (*domain.DashboardSummary, error) {
	s := &domain.DashboardSummary{
		SalesLast7Days: make([]domain.DailySales, 0),
		POStatusCounts: make(map[string]int),
		FGExpiringSoon: make([]domain.FGExpiringSoonItem, 0),
	}

	// Q1: KPI counts
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM sales_orders WHERE status IN ('confirmed','picking')`).Scan(&s.PendingSOCount)
	if err != nil {
		return nil, err
	}
	err = r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM production_orders WHERE status = 'in_progress'`).Scan(&s.InProgressPOCount)
	if err != nil {
		return nil, err
	}
	err = r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM delivery_orders WHERE status IN ('pending','loading')`).Scan(&s.PendingDOCount)
	if err != nil {
		return nil, err
	}
	err = r.pool.QueryRow(ctx,
		`SELECT COUNT(*), COALESCE(SUM(total_amount),0)::float8 FROM invoices WHERE status = 'issued'`,
	).Scan(&s.IssuedInvoiceCount, &s.IssuedInvoiceTotal)
	if err != nil {
		return nil, err
	}

	// Q2: daily sales last 7 days
	rows, err := r.pool.Query(ctx, `
		SELECT d.day::date::text AS date,
		       COALESCE(SUM(so.total_amount),0)::float8 AS total
		FROM generate_series(CURRENT_DATE - INTERVAL '6 days', CURRENT_DATE, '1 day') AS d(day)
		LEFT JOIN sales_orders so ON so.order_date = d.day::date AND so.status != 'cancelled'
		GROUP BY d.day ORDER BY d.day
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var ds domain.DailySales
		if err := rows.Scan(&ds.Date, &ds.Total); err != nil {
			return nil, err
		}
		s.SalesLast7Days = append(s.SalesLast7Days, ds)
	}

	// Q3: PO status counts
	rows2, err := r.pool.Query(ctx, `
		SELECT status, COUNT(*)::int FROM production_orders
		WHERE status != 'cancelled' GROUP BY status
	`)
	if err != nil {
		return nil, err
	}
	defer rows2.Close()
	for rows2.Next() {
		var status string
		var count int
		if err := rows2.Scan(&status, &count); err != nil {
			return nil, err
		}
		s.POStatusCounts[status] = count
	}

	// Q4: FG lots expiring within 7 days
	rows3, err := r.pool.Query(ctx, `
		SELECT sl.batch_number, fg.name, sl.current_qty::float8, sl.expiry_date::text
		FROM fg_stock_lots sl
		JOIN finished_goods fg ON fg.id = sl.finished_good_id
		WHERE sl.status = 'available'
		  AND sl.current_qty > 0
		  AND sl.expiry_date::date BETWEEN CURRENT_DATE AND CURRENT_DATE + INTERVAL '7 days'
		ORDER BY sl.expiry_date, fg.name
	`)
	if err != nil {
		return nil, err
	}
	defer rows3.Close()
	for rows3.Next() {
		var item domain.FGExpiringSoonItem
		if err := rows3.Scan(&item.LotNumber, &item.FGName, &item.Qty, &item.ExpiryDate); err != nil {
			return nil, err
		}
		s.FGExpiringSoon = append(s.FGExpiringSoon, item)
	}

	return s, nil
}
