package repository

import (
	"context"
	"fmt"
	"math"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/water-factory/api/internal/domain"
	"github.com/water-factory/api/internal/dto"
	"github.com/water-factory/api/pkg/number"
)

type FGInventoryRepo struct {
	pool *pgxpool.Pool
}

func NewFGInventoryRepo(pool *pgxpool.Pool) *FGInventoryRepo {
	return &FGInventoryRepo{pool: pool}
}

// ── Stock Summary ─────────────────────────────────────────────────────────────

func (r *FGInventoryRepo) ListStockSummary(ctx context.Context) ([]*domain.FGStockSummary, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT
			sl.finished_good_id::text,
			fg.code,
			fg.name,
			u.code,
			SUM(sl.current_qty)::float8,
			COUNT(*)::int,
			MIN(sl.expiry_date)::text
		FROM fg_stock_lots sl
		JOIN finished_goods fg ON fg.id = sl.finished_good_id
		JOIN units_of_measure u ON u.id = fg.uom_id
		WHERE sl.status = 'available'
		GROUP BY sl.finished_good_id, fg.code, fg.name, u.code
		ORDER BY fg.code`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*domain.FGStockSummary, 0)
	for rows.Next() {
		s := &domain.FGStockSummary{}
		if err := rows.Scan(&s.FinishedGoodID, &s.FGCode, &s.FGName, &s.UOMCode,
			&s.TotalQty, &s.LotCount, &s.NearestExpiry); err != nil {
			return nil, err
		}
		items = append(items, s)
	}
	return items, nil
}

// ── Stock Lots ────────────────────────────────────────────────────────────────

const fgLotSelectCols = `
	sl.id, sl.finished_good_id::text, fg.code, fg.name, u.code,
	sl.batch_number,
	sl.production_order_id::text,
	sl.production_date::text,
	sl.expiry_date::text,
	sl.location_id,
	COALESCE(wl.zone,''), COALESCE(wl.row_no,''), COALESCE(wl.bay_no,''),
	sl.initial_qty::float8, sl.current_qty::float8,
	sl.status,
	sl.created_at, sl.updated_at`

const fgLotFrom = `
	FROM fg_stock_lots sl
	JOIN finished_goods fg ON fg.id = sl.finished_good_id
	JOIN units_of_measure u ON u.id = fg.uom_id
	LEFT JOIN warehouse_locations wl ON wl.id = sl.location_id`

func scanFGLotRow(row interface{ Scan(...any) error }) (*domain.FGStockLot, string, string, string, error) {
	sl := &domain.FGStockLot{}
	var zone, rowNo, bayNo string
	err := row.Scan(
		&sl.ID, &sl.FinishedGoodID, &sl.FGCode, &sl.FGName, &sl.UOMCode,
		&sl.BatchNumber,
		&sl.ProductionOrderID,
		&sl.ProductionDate,
		&sl.ExpiryDate,
		&sl.LocationID,
		&zone, &rowNo, &bayNo,
		&sl.InitialQty, &sl.CurrentQty,
		&sl.Status,
		&sl.CreatedAt, &sl.UpdatedAt,
	)
	return sl, zone, rowNo, bayNo, err
}

func (r *FGInventoryRepo) ListStockLots(ctx context.Context, finishedGoodID string) ([]*domain.FGStockLot, error) {
	var (
		rows interface {
			Next() bool
			Scan(...any) error
			Close()
			Err() error
		}
		err error
	)

	if finishedGoodID != "" {
		rows, err = r.pool.Query(ctx,
			`SELECT`+fgLotSelectCols+fgLotFrom+` WHERE sl.finished_good_id=$1 ORDER BY sl.expiry_date ASC`,
			finishedGoodID)
	} else {
		rows, err = r.pool.Query(ctx,
			`SELECT`+fgLotSelectCols+fgLotFrom+` ORDER BY fg.code, sl.expiry_date ASC`)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*domain.FGStockLot, 0)
	for rows.Next() {
		sl, zone, rowNo, bayNo, err := scanFGLotRow(rows)
		if err != nil {
			return nil, err
		}
		sl.LocationLabel = locationLabel(zone, rowNo, bayNo)
		items = append(items, sl)
	}
	return items, nil
}

func (r *FGInventoryRepo) GetStockLotDetail(ctx context.Context, id string) (*domain.FGStockLotDetail, error) {
	sl, zone, rowNo, bayNo, err := scanFGLotRow(r.pool.QueryRow(ctx,
		`SELECT`+fgLotSelectCols+fgLotFrom+` WHERE sl.id=$1`, id))
	if err != nil {
		return nil, err
	}
	sl.LocationLabel = locationLabel(zone, rowNo, bayNo)

	detail := &domain.FGStockLotDetail{
		FGStockLot: *sl,
		Movements:  make([]domain.FGStockMovement, 0),
	}

	mvRows, err := r.pool.Query(ctx, `
		SELECT
			m.id, m.fg_stock_lot_id::text, m.finished_good_id::text,
			m.movement_type,
			COALESCE(m.reference_type,''),
			m.reference_id::text,
			m.qty::float8, m.qty_before::float8, m.qty_after::float8,
			COALESCE(m.performed_by::text,''), COALESCE(u.full_name,''),
			COALESCE(m.notes,''),
			m.created_at
		FROM fg_stock_movements m
		LEFT JOIN users u ON u.id = m.performed_by
		WHERE m.fg_stock_lot_id = $1
		ORDER BY m.created_at`, id)
	if err != nil {
		return nil, err
	}
	defer mvRows.Close()

	for mvRows.Next() {
		var mv domain.FGStockMovement
		if err := mvRows.Scan(
			&mv.ID, &mv.FGStockLotID, &mv.FinishedGoodID,
			&mv.MovementType,
			&mv.ReferenceType,
			&mv.ReferenceID,
			&mv.Qty, &mv.QtyBefore, &mv.QtyAfter,
			&mv.PerformedBy, &mv.PerformedByName,
			&mv.Notes,
			&mv.CreatedAt,
		); err != nil {
			return nil, err
		}
		detail.Movements = append(detail.Movements, mv)
	}
	return detail, nil
}

// ── Adjustments ───────────────────────────────────────────────────────────────

const fgAdjSelectCols = `
	a.id, a.adj_number,
	a.adjustment_date::text,
	a.type,
	COALESCE(a.performed_by::text,''), COALESCE(pu.full_name,''),
	a.approved_by::text, COALESCE(au.full_name,''),
	a.status,
	COALESCE(a.notes,''),
	a.created_at`

const fgAdjFrom = `
	FROM fg_adjustments a
	LEFT JOIN users pu ON pu.id = a.performed_by
	LEFT JOIN users au ON au.id = a.approved_by`

func scanFGAdjustment(row interface{ Scan(...any) error }) (*domain.FGAdjustment, error) {
	a := &domain.FGAdjustment{Lines: make([]domain.FGAdjustmentLine, 0)}
	err := row.Scan(
		&a.ID, &a.AdjNumber,
		&a.AdjustmentDate,
		&a.Type,
		&a.PerformedBy, &a.PerformedByName,
		&a.ApprovedBy, &a.ApprovedByName,
		&a.Status,
		&a.Notes,
		&a.CreatedAt,
	)
	return a, err
}

func (r *FGInventoryRepo) ListAdjustments(ctx context.Context) ([]*domain.FGAdjustment, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT`+fgAdjSelectCols+fgAdjFrom+` ORDER BY a.created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*domain.FGAdjustment, 0)
	for rows.Next() {
		a, err := scanFGAdjustment(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, a)
	}
	return items, nil
}

func (r *FGInventoryRepo) GetAdjustment(ctx context.Context, id string) (*domain.FGAdjustment, error) {
	adj, err := scanFGAdjustment(r.pool.QueryRow(ctx,
		`SELECT`+fgAdjSelectCols+fgAdjFrom+` WHERE a.id=$1`, id))
	if err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx, `
		SELECT
			al.id, al.adjustment_id::text,
			al.fg_stock_lot_id::text, sl.batch_number,
			al.finished_good_id::text, fg.code, fg.name,
			al.system_qty::float8, al.counted_qty::float8, al.variance_qty::float8,
			COALESCE(al.reason,'')
		FROM fg_adjustment_lines al
		JOIN fg_stock_lots sl ON sl.id = al.fg_stock_lot_id
		JOIN finished_goods fg ON fg.id = al.finished_good_id
		WHERE al.adjustment_id = $1
		ORDER BY al.id`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var line domain.FGAdjustmentLine
		if err := rows.Scan(
			&line.ID, &line.AdjustmentID,
			&line.FGStockLotID, &line.BatchNumber,
			&line.FinishedGoodID, &line.FGCode, &line.FGName,
			&line.SystemQty, &line.CountedQty, &line.VarianceQty,
			&line.Reason,
		); err != nil {
			return nil, err
		}
		adj.Lines = append(adj.Lines, line)
	}
	return adj, nil
}

func (r *FGInventoryRepo) CreateAdjustment(ctx context.Context, userID string, req *dto.CreateFGAdjustmentRequest) (*domain.FGAdjustment, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	adjNumber, err := number.Next(ctx, r.pool, "FGADJ")
	if err != nil {
		return nil, err
	}

	var adjID string
	err = tx.QueryRow(ctx, `
		INSERT INTO fg_adjustments
			(adj_number, adjustment_date, type, performed_by, status, notes)
		VALUES ($1, $2, $3, $4, 'pending', $5)
		RETURNING id`,
		adjNumber, req.AdjustmentDate, req.Type, userID, req.Notes,
	).Scan(&adjID)
	if err != nil {
		return nil, err
	}

	for _, line := range req.Lines {
		var systemQty float64
		var finishedGoodID string
		if err := tx.QueryRow(ctx,
			`SELECT current_qty::float8, finished_good_id::text FROM fg_stock_lots WHERE id=$1 FOR UPDATE`,
			line.FGStockLotID,
		).Scan(&systemQty, &finishedGoodID); err != nil {
			return nil, fmt.Errorf("fg stock lot not found: %s", line.FGStockLotID)
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO fg_adjustment_lines
				(adjustment_id, fg_stock_lot_id, finished_good_id, system_qty, counted_qty, reason)
			VALUES ($1, $2, $3, $4, $5, $6)`,
			adjID, line.FGStockLotID, finishedGoodID, systemQty, line.CountedQty, line.Reason,
		)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return r.GetAdjustment(ctx, adjID)
}

func (r *FGInventoryRepo) ApproveAdjustment(ctx context.Context, id string, approverID string) (*domain.FGAdjustment, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var status string
	if err := tx.QueryRow(ctx,
		`SELECT status FROM fg_adjustments WHERE id=$1 FOR UPDATE`, id,
	).Scan(&status); err != nil {
		return nil, err
	}
	if status != "pending" {
		return nil, fmt.Errorf("adjustment is not in pending status")
	}

	rows, err := tx.Query(ctx, `
		SELECT fg_stock_lot_id::text, system_qty::float8, counted_qty::float8, variance_qty::float8
		FROM fg_adjustment_lines WHERE adjustment_id=$1`, id)
	if err != nil {
		return nil, err
	}

	type adjLine struct {
		lotID       string
		systemQty   float64
		countedQty  float64
		varianceQty float64
	}
	var adjLines []adjLine
	for rows.Next() {
		var l adjLine
		if err := rows.Scan(&l.lotID, &l.systemQty, &l.countedQty, &l.varianceQty); err != nil {
			rows.Close()
			return nil, err
		}
		adjLines = append(adjLines, l)
	}
	rows.Close()

	for _, line := range adjLines {
		if _, err := tx.Exec(ctx,
			`UPDATE fg_stock_lots SET current_qty=$1, updated_at=NOW() WHERE id=$2`,
			line.countedQty, line.lotID,
		); err != nil {
			return nil, err
		}

		if line.countedQty == 0 {
			if _, err := tx.Exec(ctx,
				`UPDATE fg_stock_lots SET status='dispatched' WHERE id=$1`, line.lotID,
			); err != nil {
				return nil, err
			}
		}

		if line.varianceQty != 0 {
			mvType := "ADJUSTMENT_IN"
			if line.varianceQty < 0 {
				mvType = "ADJUSTMENT_OUT"
			}
			absVariance := math.Abs(line.varianceQty)

			var fgID string
			if err := tx.QueryRow(ctx,
				`SELECT finished_good_id::text FROM fg_stock_lots WHERE id=$1`, line.lotID,
			).Scan(&fgID); err != nil {
				return nil, err
			}

			if _, err := tx.Exec(ctx, `
				INSERT INTO fg_stock_movements
					(fg_stock_lot_id, finished_good_id, movement_type, reference_type, reference_id,
					 qty, qty_before, qty_after, performed_by)
				VALUES ($1, $2, $3, 'ADJUSTMENT', $4, $5, $6, $7, $8)`,
				line.lotID, fgID, mvType, id,
				absVariance, line.systemQty, line.countedQty, approverID,
			); err != nil {
				return nil, err
			}
		}
	}

	if _, err := tx.Exec(ctx,
		`UPDATE fg_adjustments SET status='approved', approved_by=$1 WHERE id=$2`,
		approverID, id,
	); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return r.GetAdjustment(ctx, id)
}

func (r *FGInventoryRepo) CancelAdjustment(ctx context.Context, id string) error {
	var status string
	if err := r.pool.QueryRow(ctx,
		`SELECT status FROM fg_adjustments WHERE id=$1`, id,
	).Scan(&status); err != nil {
		return err
	}
	if status != "pending" {
		return fmt.Errorf("only pending adjustments can be cancelled")
	}
	_, err := r.pool.Exec(ctx,
		`UPDATE fg_adjustments SET status='cancelled' WHERE id=$1`, id)
	return err
}
