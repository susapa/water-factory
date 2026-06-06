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

type RMInventoryRepo struct {
	pool *pgxpool.Pool
}

func NewRMInventoryRepo(pool *pgxpool.Pool) *RMInventoryRepo {
	return &RMInventoryRepo{pool: pool}
}

// ── Warehouse Locations ───────────────────────────────────────────────────────

func (r *RMInventoryRepo) ListWarehouseLocations(ctx context.Context) ([]*domain.WarehouseLocation, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, zone, COALESCE(row_no,''), COALESCE(bay_no,''), COALESCE(description,'')
		 FROM warehouse_locations ORDER BY zone, row_no, bay_no`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*domain.WarehouseLocation, 0)
	for rows.Next() {
		loc := &domain.WarehouseLocation{}
		if err := rows.Scan(&loc.ID, &loc.Zone, &loc.RowNo, &loc.BayNo, &loc.Description); err != nil {
			return nil, err
		}
		items = append(items, loc)
	}
	return items, nil
}

func locationLabel(zone, row, bay string) string {
	label := zone
	if row != "" {
		label += "-" + row
	}
	if bay != "" {
		label += "-" + bay
	}
	return label
}

// ── GRN ──────────────────────────────────────────────────────────────────────

const grnSelectCols = `
	g.id, g.grn_number,
	g.supplier_id, COALESCE(s.name, ''),
	g.received_by, COALESCE(u.full_name, ''),
	g.received_date::text,
	COALESCE(g.po_reference, ''),
	g.status,
	COALESCE(g.notes, ''),
	g.created_at, g.updated_at`

const grnFrom = `
	FROM rm_goods_receipts g
	LEFT JOIN suppliers s ON s.id = g.supplier_id
	LEFT JOIN users u ON u.id = g.received_by`

func scanGRN(row interface{ Scan(...any) error }) (*domain.GRN, error) {
	g := &domain.GRN{Lines: make([]domain.GRNLine, 0)}
	err := row.Scan(
		&g.ID, &g.GRNNumber,
		&g.SupplierID, &g.SupplierName,
		&g.ReceivedBy, &g.ReceivedByName,
		&g.ReceivedDate,
		&g.POReference,
		&g.Status,
		&g.Notes,
		&g.CreatedAt, &g.UpdatedAt,
	)
	return g, err
}

func (r *RMInventoryRepo) ListGRNs(ctx context.Context) ([]*domain.GRN, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT`+grnSelectCols+grnFrom+` ORDER BY g.created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*domain.GRN, 0)
	for rows.Next() {
		g, err := scanGRN(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, g)
	}
	return items, nil
}

func (r *RMInventoryRepo) GetGRN(ctx context.Context, id string) (*domain.GRN, error) {
	grn, err := scanGRN(r.pool.QueryRow(ctx,
		`SELECT`+grnSelectCols+grnFrom+` WHERE g.id = $1`, id))
	if err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx, `
		SELECT
			l.id, l.grn_id,
			l.raw_material_id, rm.code, rm.name,
			COALESCE(l.lot_number, ''),
			l.received_qty::float8,
			l.unit_cost,
			l.expiry_date::text,
			l.location_id,
			COALESCE(wl.zone, ''), COALESCE(wl.row_no, ''), COALESCE(wl.bay_no, ''),
			COALESCE(l.notes, '')
		FROM rm_goods_receipt_lines l
		JOIN raw_materials rm ON rm.id = l.raw_material_id
		LEFT JOIN warehouse_locations wl ON wl.id = l.location_id
		WHERE l.grn_id = $1
		ORDER BY l.id`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var line domain.GRNLine
		var zone, rowNo, bayNo string
		if err := rows.Scan(
			&line.ID, &line.GRNID,
			&line.RawMaterialID, &line.RMCode, &line.RMName,
			&line.LotNumber,
			&line.ReceivedQty,
			&line.UnitCost,
			&line.ExpiryDate,
			&line.LocationID,
			&zone, &rowNo, &bayNo,
			&line.Notes,
		); err != nil {
			return nil, err
		}
		line.LocationLabel = locationLabel(zone, rowNo, bayNo)
		grn.Lines = append(grn.Lines, line)
	}
	return grn, nil
}

func (r *RMInventoryRepo) CreateGRN(ctx context.Context, userID string, req *dto.CreateGRNRequest) (*domain.GRN, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	grnNumber, err := number.Next(ctx, r.pool, "GRN")
	if err != nil {
		return nil, err
	}

	var grnID string
	err = tx.QueryRow(ctx, `
		INSERT INTO rm_goods_receipts
			(grn_number, supplier_id, received_by, received_date, po_reference, status, notes)
		VALUES ($1, $2, $3, $4, $5, 'draft', $6)
		RETURNING id`,
		grnNumber, req.SupplierID, userID, req.ReceivedDate, req.POReference, req.Notes,
	).Scan(&grnID)
	if err != nil {
		return nil, err
	}

	for _, line := range req.Lines {
		_, err = tx.Exec(ctx, `
			INSERT INTO rm_goods_receipt_lines
				(grn_id, raw_material_id, lot_number, received_qty, unit_cost, expiry_date, location_id, notes)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			grnID, line.RawMaterialID, line.LotNumber, line.ReceivedQty,
			line.UnitCost, line.ExpiryDate, line.LocationID, line.Notes,
		)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return r.GetGRN(ctx, grnID)
}

func (r *RMInventoryRepo) UpdateGRN(ctx context.Context, id string, req *dto.UpdateGRNRequest) (*domain.GRN, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var status string
	if err := tx.QueryRow(ctx,
		`SELECT status FROM rm_goods_receipts WHERE id = $1 FOR UPDATE`, id,
	).Scan(&status); err != nil {
		return nil, err
	}
	if status != "draft" {
		return nil, fmt.Errorf("GRN is not in draft status")
	}

	_, err = tx.Exec(ctx, `
		UPDATE rm_goods_receipts
		SET supplier_id=$1, received_date=$2, po_reference=$3, notes=$4, updated_at=NOW()
		WHERE id=$5`,
		req.SupplierID, req.ReceivedDate, req.POReference, req.Notes, id,
	)
	if err != nil {
		return nil, err
	}

	if _, err = tx.Exec(ctx, `DELETE FROM rm_goods_receipt_lines WHERE grn_id=$1`, id); err != nil {
		return nil, err
	}

	for _, line := range req.Lines {
		_, err = tx.Exec(ctx, `
			INSERT INTO rm_goods_receipt_lines
				(grn_id, raw_material_id, lot_number, received_qty, unit_cost, expiry_date, location_id, notes)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			id, line.RawMaterialID, line.LotNumber, line.ReceivedQty,
			line.UnitCost, line.ExpiryDate, line.LocationID, line.Notes,
		)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return r.GetGRN(ctx, id)
}

func (r *RMInventoryRepo) ConfirmGRN(ctx context.Context, id string, userID string) (*domain.GRN, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var status string
	if err := tx.QueryRow(ctx,
		`SELECT status FROM rm_goods_receipts WHERE id = $1 FOR UPDATE`, id,
	).Scan(&status); err != nil {
		return nil, err
	}
	if status != "draft" {
		return nil, fmt.Errorf("GRN is not in draft status")
	}

	type lineInfo struct {
		id            string
		rawMaterialID string
		lotNumber     string
		receivedQty   float64
		unitCost      *float64
		expiryDate    *string
		locationID    *int
	}

	lineRows, err := tx.Query(ctx, `
		SELECT id, raw_material_id, COALESCE(lot_number,''), received_qty::float8, unit_cost, expiry_date::text, location_id
		FROM rm_goods_receipt_lines WHERE grn_id = $1 ORDER BY id`, id)
	if err != nil {
		return nil, err
	}
	var lines []lineInfo
	for lineRows.Next() {
		var l lineInfo
		if err := lineRows.Scan(&l.id, &l.rawMaterialID, &l.lotNumber, &l.receivedQty, &l.unitCost, &l.expiryDate, &l.locationID); err != nil {
			lineRows.Close()
			return nil, err
		}
		lines = append(lines, l)
	}
	lineRows.Close()

	var receivedDate string
	if err := tx.QueryRow(ctx, `SELECT received_date::text FROM rm_goods_receipts WHERE id=$1`, id).Scan(&receivedDate); err != nil {
		return nil, err
	}

	for _, line := range lines {
		lotNum := line.lotNumber
		if lotNum == "" {
			generated, err := number.Next(ctx, r.pool, "LOT")
			if err != nil {
				return nil, err
			}
			lotNum = generated
			if _, err := tx.Exec(ctx, `UPDATE rm_goods_receipt_lines SET lot_number=$1 WHERE id=$2`, lotNum, line.id); err != nil {
				return nil, err
			}
		}

		var lotID string
		err = tx.QueryRow(ctx, `
			INSERT INTO rm_stock_lots
				(raw_material_id, lot_number, grn_line_id, received_date, expiry_date, location_id,
				 initial_qty, current_qty, unit_cost, status)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $7, $8, 'available')
			RETURNING id`,
			line.rawMaterialID, lotNum, line.id, receivedDate,
			line.expiryDate, line.locationID, line.receivedQty, line.unitCost,
		).Scan(&lotID)
		if err != nil {
			if isUniqueViolation(err) {
				return nil, fmt.Errorf("lot number '%s' already exists for this raw material", lotNum)
			}
			return nil, err
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO rm_stock_movements
				(rm_stock_lot_id, raw_material_id, movement_type, reference_type, reference_id,
				 qty, qty_before, qty_after, performed_by)
			VALUES ($1, $2, 'GRN', 'GRN', $3, $4, 0, $4, $5)`,
			lotID, line.rawMaterialID, id, line.receivedQty, userID,
		)
		if err != nil {
			return nil, err
		}
	}

	if _, err = tx.Exec(ctx,
		`UPDATE rm_goods_receipts SET status='confirmed', updated_at=NOW() WHERE id=$1`, id,
	); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return r.GetGRN(ctx, id)
}

func (r *RMInventoryRepo) CancelGRN(ctx context.Context, id string) error {
	var status string
	if err := r.pool.QueryRow(ctx,
		`SELECT status FROM rm_goods_receipts WHERE id=$1`, id,
	).Scan(&status); err != nil {
		return err
	}
	if status != "draft" {
		return fmt.Errorf("only draft GRNs can be cancelled")
	}
	_, err := r.pool.Exec(ctx,
		`UPDATE rm_goods_receipts SET status='cancelled', updated_at=NOW() WHERE id=$1`, id)
	return err
}

// ── Stock ─────────────────────────────────────────────────────────────────────

func (r *RMInventoryRepo) ListStockSummary(ctx context.Context) ([]*domain.StockSummary, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT
			sl.raw_material_id::text,
			rm.code,
			rm.name,
			u.code,
			SUM(sl.current_qty)::float8,
			COUNT(*)::int,
			MIN(sl.expiry_date)::text
		FROM rm_stock_lots sl
		JOIN raw_materials rm ON rm.id = sl.raw_material_id
		JOIN units_of_measure u ON u.id = rm.uom_id
		WHERE sl.status = 'available'
		GROUP BY sl.raw_material_id, rm.code, rm.name, u.code
		ORDER BY rm.code`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*domain.StockSummary, 0)
	for rows.Next() {
		s := &domain.StockSummary{}
		if err := rows.Scan(&s.RawMaterialID, &s.RMCode, &s.RMName, &s.UOMCode,
			&s.TotalQty, &s.LotCount, &s.EarliestExpiry); err != nil {
			return nil, err
		}
		items = append(items, s)
	}
	return items, nil
}

func (r *RMInventoryRepo) ListStockLots(ctx context.Context, rawMaterialID string) ([]*domain.StockLot, error) {
	query := `
		SELECT
			sl.id, sl.raw_material_id::text, rm.code, rm.name, u.code,
			sl.lot_number,
			sl.grn_line_id::text,
			sl.received_date::text,
			sl.expiry_date::text,
			sl.location_id,
			COALESCE(wl.zone,''), COALESCE(wl.row_no,''), COALESCE(wl.bay_no,''),
			sl.initial_qty::float8, sl.current_qty::float8,
			sl.unit_cost,
			sl.status,
			sl.created_at, sl.updated_at
		FROM rm_stock_lots sl
		JOIN raw_materials rm ON rm.id = sl.raw_material_id
		JOIN units_of_measure u ON u.id = rm.uom_id
		LEFT JOIN warehouse_locations wl ON wl.id = sl.location_id`

	var rows interface {
		Next() bool
		Scan(...any) error
		Close()
		Err() error
	}
	var err error

	if rawMaterialID != "" {
		rows, err = r.pool.Query(ctx, query+` WHERE sl.raw_material_id=$1 ORDER BY sl.received_date`, rawMaterialID)
	} else {
		rows, err = r.pool.Query(ctx, query+` ORDER BY rm.code, sl.received_date`)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*domain.StockLot, 0)
	for rows.Next() {
		sl, zone, rowNo, bayNo, err := scanStockLotRow(rows)
		if err != nil {
			return nil, err
		}
		sl.LocationLabel = locationLabel(zone, rowNo, bayNo)
		items = append(items, sl)
	}
	return items, nil
}

func scanStockLotRow(row interface{ Scan(...any) error }) (*domain.StockLot, string, string, string, error) {
	sl := &domain.StockLot{}
	var zone, rowNo, bayNo string
	err := row.Scan(
		&sl.ID, &sl.RawMaterialID, &sl.RMCode, &sl.RMName, &sl.UOMCode,
		&sl.LotNumber,
		&sl.GRNLineID,
		&sl.ReceivedDate,
		&sl.ExpiryDate,
		&sl.LocationID,
		&zone, &rowNo, &bayNo,
		&sl.InitialQty, &sl.CurrentQty,
		&sl.UnitCost,
		&sl.Status,
		&sl.CreatedAt, &sl.UpdatedAt,
	)
	return sl, zone, rowNo, bayNo, err
}

func (r *RMInventoryRepo) GetStockLotDetail(ctx context.Context, id string) (*domain.StockLotDetail, error) {
	sl, zone, rowNo, bayNo, err := scanStockLotRow(r.pool.QueryRow(ctx, `
		SELECT
			sl.id, sl.raw_material_id::text, rm.code, rm.name, u.code,
			sl.lot_number,
			sl.grn_line_id::text,
			sl.received_date::text,
			sl.expiry_date::text,
			sl.location_id,
			COALESCE(wl.zone,''), COALESCE(wl.row_no,''), COALESCE(wl.bay_no,''),
			sl.initial_qty::float8, sl.current_qty::float8,
			sl.unit_cost,
			sl.status,
			sl.created_at, sl.updated_at
		FROM rm_stock_lots sl
		JOIN raw_materials rm ON rm.id = sl.raw_material_id
		JOIN units_of_measure u ON u.id = rm.uom_id
		LEFT JOIN warehouse_locations wl ON wl.id = sl.location_id
		WHERE sl.id = $1`, id))
	if err != nil {
		return nil, err
	}
	sl.LocationLabel = locationLabel(zone, rowNo, bayNo)

	detail := &domain.StockLotDetail{
		StockLot:  *sl,
		Movements: make([]domain.StockMovement, 0),
	}

	mvRows, err := r.pool.Query(ctx, `
		SELECT
			m.id, m.rm_stock_lot_id::text, m.raw_material_id::text,
			m.movement_type,
			COALESCE(m.reference_type,''),
			m.reference_id::text,
			m.qty::float8, m.qty_before::float8, m.qty_after::float8,
			COALESCE(m.performed_by::text,''), COALESCE(u.full_name,''),
			COALESCE(m.notes,''),
			m.created_at
		FROM rm_stock_movements m
		LEFT JOIN users u ON u.id = m.performed_by
		WHERE m.rm_stock_lot_id = $1
		ORDER BY m.created_at`, id)
	if err != nil {
		return nil, err
	}
	defer mvRows.Close()

	for mvRows.Next() {
		var mv domain.StockMovement
		if err := mvRows.Scan(
			&mv.ID, &mv.RMStockLotID, &mv.RawMaterialID,
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

const adjSelectCols = `
	a.id, a.adj_number,
	a.adjustment_date::text,
	a.type,
	COALESCE(a.performed_by::text,''), COALESCE(pu.full_name,''),
	a.approved_by::text, COALESCE(au.full_name,''),
	a.status,
	COALESCE(a.notes,''),
	a.created_at`

const adjFrom = `
	FROM rm_adjustments a
	LEFT JOIN users pu ON pu.id = a.performed_by
	LEFT JOIN users au ON au.id = a.approved_by`

func scanAdjustment(row interface{ Scan(...any) error }) (*domain.Adjustment, error) {
	a := &domain.Adjustment{Lines: make([]domain.AdjustmentLine, 0)}
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

func (r *RMInventoryRepo) ListAdjustments(ctx context.Context) ([]*domain.Adjustment, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT`+adjSelectCols+adjFrom+` ORDER BY a.created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*domain.Adjustment, 0)
	for rows.Next() {
		a, err := scanAdjustment(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, a)
	}
	return items, nil
}

func (r *RMInventoryRepo) GetAdjustment(ctx context.Context, id string) (*domain.Adjustment, error) {
	adj, err := scanAdjustment(r.pool.QueryRow(ctx,
		`SELECT`+adjSelectCols+adjFrom+` WHERE a.id=$1`, id))
	if err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx, `
		SELECT
			al.id, al.adjustment_id::text,
			al.rm_stock_lot_id::text, sl.lot_number,
			al.raw_material_id::text, rm.code, rm.name,
			al.system_qty::float8, al.counted_qty::float8, al.variance_qty::float8,
			COALESCE(al.reason,'')
		FROM rm_adjustment_lines al
		JOIN rm_stock_lots sl ON sl.id = al.rm_stock_lot_id
		JOIN raw_materials rm ON rm.id = al.raw_material_id
		WHERE al.adjustment_id = $1
		ORDER BY al.id`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var line domain.AdjustmentLine
		if err := rows.Scan(
			&line.ID, &line.AdjustmentID,
			&line.RMStockLotID, &line.LotNumber,
			&line.RawMaterialID, &line.RMCode, &line.RMName,
			&line.SystemQty, &line.CountedQty, &line.VarianceQty,
			&line.Reason,
		); err != nil {
			return nil, err
		}
		adj.Lines = append(adj.Lines, line)
	}
	return adj, nil
}

func (r *RMInventoryRepo) CreateAdjustment(ctx context.Context, userID string, req *dto.CreateAdjustmentRequest) (*domain.Adjustment, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	adjNumber, err := number.Next(ctx, r.pool, "ADJ")
	if err != nil {
		return nil, err
	}

	var adjID string
	err = tx.QueryRow(ctx, `
		INSERT INTO rm_adjustments
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
		var rawMaterialID string
		if err := tx.QueryRow(ctx,
			`SELECT current_qty::float8, raw_material_id::text FROM rm_stock_lots WHERE id=$1 FOR UPDATE`,
			line.RMStockLotID,
		).Scan(&systemQty, &rawMaterialID); err != nil {
			return nil, fmt.Errorf("stock lot not found: %s", line.RMStockLotID)
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO rm_adjustment_lines
				(adjustment_id, rm_stock_lot_id, raw_material_id, system_qty, counted_qty, reason)
			VALUES ($1, $2, $3, $4, $5, $6)`,
			adjID, line.RMStockLotID, rawMaterialID, systemQty, line.CountedQty, line.Reason,
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

func (r *RMInventoryRepo) ApproveAdjustment(ctx context.Context, id string, approverID string) (*domain.Adjustment, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var status string
	if err := tx.QueryRow(ctx,
		`SELECT status FROM rm_adjustments WHERE id=$1 FOR UPDATE`, id,
	).Scan(&status); err != nil {
		return nil, err
	}
	if status != "pending" {
		return nil, fmt.Errorf("adjustment is not in pending status")
	}

	rows, err := tx.Query(ctx, `
		SELECT rm_stock_lot_id::text, system_qty::float8, counted_qty::float8, variance_qty::float8
		FROM rm_adjustment_lines WHERE adjustment_id=$1`, id)
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
			`UPDATE rm_stock_lots SET current_qty=$1, updated_at=NOW() WHERE id=$2`,
			line.countedQty, line.lotID,
		); err != nil {
			return nil, err
		}

		if line.countedQty == 0 {
			if _, err := tx.Exec(ctx,
				`UPDATE rm_stock_lots SET status='depleted' WHERE id=$1`, line.lotID,
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

			var rmID string
			if err := tx.QueryRow(ctx,
				`SELECT raw_material_id::text FROM rm_stock_lots WHERE id=$1`, line.lotID,
			).Scan(&rmID); err != nil {
				return nil, err
			}

			if _, err := tx.Exec(ctx, `
				INSERT INTO rm_stock_movements
					(rm_stock_lot_id, raw_material_id, movement_type, reference_type, reference_id,
					 qty, qty_before, qty_after, performed_by)
				VALUES ($1, $2, $3, 'ADJUSTMENT', $4, $5, $6, $7, $8)`,
				line.lotID, rmID, mvType, id,
				absVariance, line.systemQty, line.countedQty, approverID,
			); err != nil {
				return nil, err
			}
		}
	}

	if _, err := tx.Exec(ctx,
		`UPDATE rm_adjustments SET status='approved', approved_by=$1 WHERE id=$2`,
		approverID, id,
	); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return r.GetAdjustment(ctx, id)
}

func (r *RMInventoryRepo) CancelAdjustment(ctx context.Context, id string) error {
	var status string
	if err := r.pool.QueryRow(ctx,
		`SELECT status FROM rm_adjustments WHERE id=$1`, id,
	).Scan(&status); err != nil {
		return err
	}
	if status != "pending" {
		return fmt.Errorf("only pending adjustments can be cancelled")
	}
	_, err := r.pool.Exec(ctx,
		`UPDATE rm_adjustments SET status='cancelled' WHERE id=$1`, id)
	return err
}
