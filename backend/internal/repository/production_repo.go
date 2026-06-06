package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/water-factory/api/internal/domain"
	"github.com/water-factory/api/internal/dto"
	"github.com/water-factory/api/pkg/number"
)

type ProductionRepo struct {
	pool *pgxpool.Pool
}

func NewProductionRepo(pool *pgxpool.Pool) *ProductionRepo {
	return &ProductionRepo{pool: pool}
}

// ── helpers ───────────────────────────────────────────────────────────────────

const orderSelectCols = `
	po.id, po.order_number,
	po.finished_good_id, fg.code, fg.name,
	po.bom_id,
	po.planned_qty::float8,
	po.actual_yield_qty::float8,
	po.defect_qty::float8,
	po.waste_qty::float8,
	po.planned_start_date::text,
	po.planned_end_date::text,
	po.actual_start_date,
	po.actual_end_date,
	po.status,
	po.created_by, COALESCE(u.full_name,''),
	po.sales_order_id,
	COALESCE(po.notes,''),
	po.created_at, po.updated_at`

const orderFrom = `
	FROM production_orders po
	JOIN finished_goods fg ON fg.id = po.finished_good_id
	LEFT JOIN users u ON u.id = po.created_by`

func scanOrder(row interface{ Scan(...any) error }) (*domain.ProductionOrder, error) {
	o := &domain.ProductionOrder{Requirements: make([]domain.ProductionRequirement, 0)}
	err := row.Scan(
		&o.ID, &o.OrderNumber,
		&o.FinishedGoodID, &o.FGCode, &o.FGName,
		&o.BOMID,
		&o.PlannedQty,
		&o.ActualYieldQty,
		&o.DefectQty,
		&o.WasteQty,
		&o.PlannedStartDate,
		&o.PlannedEndDate,
		&o.ActualStartDate,
		&o.ActualEndDate,
		&o.Status,
		&o.CreatedBy, &o.CreatedByName,
		&o.SalesOrderID,
		&o.Notes,
		&o.CreatedAt, &o.UpdatedAt,
	)
	return o, err
}

func (r *ProductionRepo) listRequirements(ctx context.Context, orderID string) ([]domain.ProductionRequirement, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT req.id, req.production_order_id,
		       req.raw_material_id, rm.code, rm.name, COALESCE(u.code,''),
		       req.required_qty::float8, req.issued_qty::float8, req.status
		FROM production_order_rm_requirements req
		JOIN raw_materials rm ON rm.id = req.raw_material_id
		LEFT JOIN units_of_measure u ON u.id = rm.uom_id
		WHERE req.production_order_id = $1`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reqs := make([]domain.ProductionRequirement, 0)
	for rows.Next() {
		r2 := domain.ProductionRequirement{}
		if err := rows.Scan(
			&r2.ID, &r2.ProductionOrderID,
			&r2.RawMaterialID, &r2.RMCode, &r2.RMName, &r2.UOMCode,
			&r2.RequiredQty, &r2.IssuedQty, &r2.Status,
		); err != nil {
			return nil, err
		}
		reqs = append(reqs, r2)
	}
	return reqs, nil
}

// ── Orders ────────────────────────────────────────────────────────────────────

func (r *ProductionRepo) ListOrders(ctx context.Context) ([]*domain.ProductionOrder, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT`+orderSelectCols+orderFrom+` ORDER BY po.created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*domain.ProductionOrder, 0)
	for rows.Next() {
		o, err := scanOrder(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, o)
	}

	// Attach requirements for each order
	for _, o := range items {
		reqs, err := r.listRequirements(ctx, o.ID)
		if err != nil {
			return nil, err
		}
		o.Requirements = reqs
	}
	return items, nil
}

func (r *ProductionRepo) GetOrder(ctx context.Context, id string) (*domain.ProductionOrder, error) {
	o, err := scanOrder(r.pool.QueryRow(ctx,
		`SELECT`+orderSelectCols+orderFrom+` WHERE po.id = $1`, id))
	if err != nil {
		return nil, fmt.Errorf("production order not found")
	}
	reqs, err := r.listRequirements(ctx, id)
	if err != nil {
		return nil, err
	}
	o.Requirements = reqs
	return o, nil
}

func (r *ProductionRepo) CreateOrder(ctx context.Context, userID string, req *dto.CreateProductionOrderRequest) (*domain.ProductionOrder, error) {
	num, err := number.Next(ctx, r.pool, "PO")
	if err != nil {
		return nil, err
	}

	var id string
	err = r.pool.QueryRow(ctx, `
		INSERT INTO production_orders
		  (order_number, finished_good_id, bom_id, planned_qty,
		   planned_start_date, planned_end_date, created_by, notes)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		RETURNING id`,
		num, req.FinishedGoodID, req.BOMID, req.PlannedQty,
		req.PlannedStartDate, req.PlannedEndDate, userID, req.Notes,
	).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("create production order: %w", err)
	}
	return r.GetOrder(ctx, id)
}

func (r *ProductionRepo) UpdateOrder(ctx context.Context, id string, req *dto.UpdateProductionOrderRequest) (*domain.ProductionOrder, error) {
	tag, err := r.pool.Exec(ctx, `
		UPDATE production_orders
		SET planned_qty=$1, planned_start_date=$2, planned_end_date=$3,
		    notes=$4, updated_at=NOW()
		WHERE id=$5 AND status='draft'`,
		req.PlannedQty, req.PlannedStartDate, req.PlannedEndDate, req.Notes, id,
	)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, fmt.Errorf("order is not in draft status")
	}
	return r.GetOrder(ctx, id)
}

func (r *ProductionRepo) ConfirmOrder(ctx context.Context, id string) (*domain.ProductionOrder, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// Lock and validate
	var status string
	var plannedQty float64
	var bomID *string
	err = tx.QueryRow(ctx,
		`SELECT status, planned_qty::float8, bom_id FROM production_orders WHERE id=$1 FOR UPDATE`, id,
	).Scan(&status, &plannedQty, &bomID)
	if err != nil {
		return nil, fmt.Errorf("production order not found")
	}
	if status != "draft" {
		return nil, fmt.Errorf("order is not in draft status")
	}
	if bomID == nil {
		return nil, fmt.Errorf("order has no BOM assigned")
	}

	// Explode BOM into requirements snapshot
	rows, err := tx.Query(ctx, `
		SELECT raw_material_id, qty_per_unit::float8, waste_factor::float8
		FROM bom_lines WHERE bom_id = $1`, *bomID)
	if err != nil {
		return nil, err
	}

	type bomLine struct {
		rmID        string
		qtyPerUnit  float64
		wasteFactor float64
	}
	var lines []bomLine
	for rows.Next() {
		var l bomLine
		if err := rows.Scan(&l.rmID, &l.qtyPerUnit, &l.wasteFactor); err != nil {
			rows.Close()
			return nil, err
		}
		lines = append(lines, l)
	}
	rows.Close()

	if len(lines) == 0 {
		return nil, fmt.Errorf("BOM has no lines")
	}

	for _, l := range lines {
		requiredQty := l.qtyPerUnit * (1 + l.wasteFactor) * plannedQty
		_, err = tx.Exec(ctx, `
			INSERT INTO production_order_rm_requirements
			  (production_order_id, raw_material_id, required_qty)
			VALUES ($1,$2,$3)`,
			id, l.rmID, requiredQty,
		)
		if err != nil {
			return nil, fmt.Errorf("insert requirement: %w", err)
		}
	}

	_, err = tx.Exec(ctx,
		`UPDATE production_orders SET status='confirmed', updated_at=NOW() WHERE id=$1`, id)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}
	return r.GetOrder(ctx, id)
}

func (r *ProductionRepo) IssueRM(ctx context.Context, id, userID string, req *dto.IssueRMRequest) (*domain.ProductionOrder, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// Validate order status
	var status string
	err = tx.QueryRow(ctx,
		`SELECT status FROM production_orders WHERE id=$1 FOR UPDATE`, id,
	).Scan(&status)
	if err != nil {
		return nil, fmt.Errorf("production order not found")
	}
	if status != "confirmed" && status != "rm_issued" {
		return nil, fmt.Errorf("order is not in confirmed status")
	}

	for _, line := range req.Lines {
		// Get requirement for this RM
		var reqID string
		var requiredQty, issuedQty float64
		err = tx.QueryRow(ctx, `
			SELECT id, required_qty::float8, issued_qty::float8
			FROM production_order_rm_requirements
			WHERE production_order_id=$1 AND raw_material_id=$2
			FOR UPDATE`, id, line.RawMaterialID,
		).Scan(&reqID, &requiredQty, &issuedQty)
		if err != nil {
			return nil, fmt.Errorf("requirement not found for raw material %s", line.RawMaterialID)
		}

		remaining := line.Qty
		// FIFO: pick available lots ordered by expiry_date ASC, received_date ASC
		lotRows, err := tx.Query(ctx, `
			SELECT id, current_qty::float8
			FROM rm_stock_lots
			WHERE raw_material_id=$1 AND status='available' AND current_qty > 0
			ORDER BY expiry_date ASC NULLS LAST, received_date ASC
			FOR UPDATE`, line.RawMaterialID)
		if err != nil {
			return nil, err
		}

		type lot struct {
			id         string
			currentQty float64
		}
		var lots []lot
		for lotRows.Next() {
			var l lot
			if err := lotRows.Scan(&l.id, &l.currentQty); err != nil {
				lotRows.Close()
				return nil, err
			}
			lots = append(lots, l)
		}
		lotRows.Close()

		totalAvailable := 0.0
		for _, l := range lots {
			totalAvailable += l.currentQty
		}
		if totalAvailable < remaining {
			return nil, fmt.Errorf("insufficient stock for raw material %s: need %.3f, available %.3f",
				line.RawMaterialID, remaining, totalAvailable)
		}

		for _, lot := range lots {
			if remaining <= 0 {
				break
			}
			take := remaining
			if take > lot.currentQty {
				take = lot.currentQty
			}

			qtyBefore := lot.currentQty
			qtyAfter := lot.currentQty - take

			newStatus := "available"
			if qtyAfter == 0 {
				newStatus = "depleted"
			}

			_, err = tx.Exec(ctx, `
				UPDATE rm_stock_lots
				SET current_qty = current_qty - $1,
				    status = $2,
				    updated_at = NOW()
				WHERE id = $3`, take, newStatus, lot.id)
			if err != nil {
				return nil, err
			}

			_, err = tx.Exec(ctx, `
				INSERT INTO production_rm_issues
				  (production_order_id, rm_stock_lot_id, raw_material_id, issued_qty, issued_by)
				VALUES ($1,$2,$3,$4,$5)`,
				id, lot.id, line.RawMaterialID, take, userID)
			if err != nil {
				return nil, err
			}

			_, err = tx.Exec(ctx, `
				INSERT INTO rm_stock_movements
				  (rm_stock_lot_id, raw_material_id, movement_type,
				   reference_type, reference_id, qty, qty_before, qty_after, performed_by)
				VALUES ($1,$2,'ISSUE_TO_PROD','production_order',$3,$4,$5,$6,$7)`,
				lot.id, line.RawMaterialID, id, take, qtyBefore, qtyAfter, userID)
			if err != nil {
				return nil, err
			}

			remaining -= take
		}

		newIssuedQty := issuedQty + line.Qty
		reqStatus := "partial"
		if newIssuedQty >= requiredQty {
			reqStatus = "fully_issued"
		}
		_, err = tx.Exec(ctx, `
			UPDATE production_order_rm_requirements
			SET issued_qty=$1, status=$2
			WHERE id=$3`, newIssuedQty, reqStatus, reqID)
		if err != nil {
			return nil, err
		}
	}

	// Check if all requirements are fully issued → update order status
	var pendingCount int
	err = tx.QueryRow(ctx, `
		SELECT COUNT(*) FROM production_order_rm_requirements
		WHERE production_order_id=$1 AND status != 'fully_issued'`, id,
	).Scan(&pendingCount)
	if err != nil {
		return nil, err
	}

	newOrderStatus := "confirmed"
	if pendingCount == 0 {
		newOrderStatus = "rm_issued"
	}
	_, err = tx.Exec(ctx,
		`UPDATE production_orders SET status=$1, updated_at=NOW() WHERE id=$2`, newOrderStatus, id)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}
	return r.GetOrder(ctx, id)
}

func (r *ProductionRepo) StartProduction(ctx context.Context, id string) (*domain.ProductionOrder, error) {
	tag, err := r.pool.Exec(ctx, `
		UPDATE production_orders
		SET status='in_progress', actual_start_date=NOW(), updated_at=NOW()
		WHERE id=$1 AND status='rm_issued'`, id)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, fmt.Errorf("order is not in rm_issued status")
	}
	return r.GetOrder(ctx, id)
}

func (r *ProductionRepo) RecordYield(ctx context.Context, id, userID string, req *dto.RecordYieldRequest) (*domain.ProductionYield, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var status string
	err = tx.QueryRow(ctx,
		`SELECT status FROM production_orders WHERE id=$1 FOR UPDATE`, id,
	).Scan(&status)
	if err != nil {
		return nil, fmt.Errorf("production order not found")
	}
	if status != "in_progress" {
		return nil, fmt.Errorf("order is not in_progress status")
	}

	var yieldID string
	err = tx.QueryRow(ctx, `
		INSERT INTO production_yields
		  (production_order_id, yield_qty, defect_qty, waste_qty,
		   batch_number, production_date, expiry_date, recorded_by, notes)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		RETURNING id`,
		id, req.YieldQty, req.DefectQty, req.WasteQty,
		req.BatchNumber, req.ProductionDate, req.ExpiryDate,
		userID, req.Notes,
	).Scan(&yieldID)
	if err != nil {
		return nil, fmt.Errorf("insert yield: %w", err)
	}

	for _, d := range req.DefectDetails {
		_, err = tx.Exec(ctx, `
			INSERT INTO production_defect_details (yield_id, defect_type, qty, description)
			VALUES ($1,$2,$3,$4)`,
			yieldID, d.DefectType, d.Qty, d.Description)
		if err != nil {
			return nil, err
		}
	}

	_, err = tx.Exec(ctx, `
		UPDATE production_orders
		SET actual_yield_qty = actual_yield_qty + $1,
		    defect_qty       = defect_qty + $2,
		    waste_qty        = waste_qty + $3,
		    updated_at       = NOW()
		WHERE id = $4`, req.YieldQty, req.DefectQty, req.WasteQty, id)
	if err != nil {
		return nil, err
	}

	// Create FG stock lot + receipt movement from yield
	var fgID string
	if err := tx.QueryRow(ctx,
		`SELECT finished_good_id::text FROM production_orders WHERE id=$1`, id,
	).Scan(&fgID); err != nil {
		return nil, fmt.Errorf("get finished_good_id: %w", err)
	}

	var fgLotID string
	err = tx.QueryRow(ctx, `
		INSERT INTO fg_stock_lots
			(finished_good_id, batch_number, production_order_id,
			 production_date, expiry_date, initial_qty, current_qty, status)
		VALUES ($1, $2, $3, $4, $5, $6, $6, 'available')
		RETURNING id`,
		fgID, req.BatchNumber, id,
		req.ProductionDate, req.ExpiryDate, req.YieldQty,
	).Scan(&fgLotID)
	if err != nil {
		return nil, fmt.Errorf("insert fg_stock_lot: %w", err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO fg_stock_movements
			(fg_stock_lot_id, finished_good_id, movement_type,
			 reference_type, reference_id,
			 qty, qty_before, qty_after, performed_by)
		VALUES ($1, $2, 'PRODUCTION_RECEIPT', 'production_yield', $3, $4, 0, $4, $5)`,
		fgLotID, fgID, yieldID, req.YieldQty, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("insert fg_stock_movement: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}
	return r.GetYield(ctx, yieldID)
}

func (r *ProductionRepo) CompleteOrder(ctx context.Context, id string) (*domain.ProductionOrder, error) {
	tag, err := r.pool.Exec(ctx, `
		UPDATE production_orders
		SET status='completed', actual_end_date=NOW(), updated_at=NOW()
		WHERE id=$1 AND status='in_progress'`, id)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, fmt.Errorf("order is not in_progress status")
	}
	return r.GetOrder(ctx, id)
}

func (r *ProductionRepo) CancelOrder(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE production_orders SET status='cancelled', updated_at=NOW()
		WHERE id=$1 AND status IN ('draft','confirmed')`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("order cannot be cancelled (not in draft or confirmed status)")
	}
	return nil
}

// ── Yields ────────────────────────────────────────────────────────────────────

func (r *ProductionRepo) GetYield(ctx context.Context, id string) (*domain.ProductionYield, error) {
	y := &domain.ProductionYield{DefectDetails: make([]domain.DefectDetail, 0)}
	err := r.pool.QueryRow(ctx, `
		SELECT y.id, y.production_order_id,
		       po.order_number, fg.name,
		       y.yield_qty::float8, y.defect_qty::float8, y.waste_qty::float8,
		       y.batch_number,
		       y.production_date::text,
		       y.expiry_date::text,
		       y.recorded_by, COALESCE(u.full_name,''),
		       y.recorded_at,
		       COALESCE(y.notes,'')
		FROM production_yields y
		JOIN production_orders po ON po.id = y.production_order_id
		JOIN finished_goods fg ON fg.id = po.finished_good_id
		LEFT JOIN users u ON u.id = y.recorded_by
		WHERE y.id = $1`, id,
	).Scan(
		&y.ID, &y.ProductionOrderID,
		&y.OrderNumber, &y.FGName,
		&y.YieldQty, &y.DefectQty, &y.WasteQty,
		&y.BatchNumber,
		&y.ProductionDate,
		&y.ExpiryDate,
		&y.RecordedBy, &y.RecordedByName,
		&y.RecordedAt,
		&y.Notes,
	)
	if err != nil {
		return nil, fmt.Errorf("yield not found")
	}

	details, err := r.listDefectDetails(ctx, id)
	if err != nil {
		return nil, err
	}
	y.DefectDetails = details
	return y, nil
}

func (r *ProductionRepo) listDefectDetails(ctx context.Context, yieldID string) ([]domain.DefectDetail, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, yield_id, defect_type, qty::float8, COALESCE(description,'')
		FROM production_defect_details
		WHERE yield_id = $1`, yieldID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	details := make([]domain.DefectDetail, 0)
	for rows.Next() {
		d := domain.DefectDetail{}
		if err := rows.Scan(&d.ID, &d.YieldID, &d.DefectType, &d.Qty, &d.Description); err != nil {
			return nil, err
		}
		details = append(details, d)
	}
	return details, nil
}

func (r *ProductionRepo) ListYields(ctx context.Context, orderID string) ([]*domain.ProductionYield, error) {
	query := `
		SELECT y.id, y.production_order_id,
		       po.order_number, fg.name,
		       y.yield_qty::float8, y.defect_qty::float8, y.waste_qty::float8,
		       y.batch_number,
		       y.production_date::text,
		       y.expiry_date::text,
		       y.recorded_by, COALESCE(u.full_name,''),
		       y.recorded_at,
		       COALESCE(y.notes,'')
		FROM production_yields y
		JOIN production_orders po ON po.id = y.production_order_id
		JOIN finished_goods fg ON fg.id = po.finished_good_id
		LEFT JOIN users u ON u.id = y.recorded_by`

	var rows interface {
		Next() bool
		Scan(...any) error
		Close()
	}
	var qErr error
	if orderID != "" {
		r2, err := r.pool.Query(ctx, query+` WHERE y.production_order_id=$1 ORDER BY y.recorded_at DESC`, orderID)
		rows, qErr = r2, err
	} else {
		r2, err := r.pool.Query(ctx, query+` ORDER BY y.recorded_at DESC`)
		rows, qErr = r2, err
	}
	if qErr != nil {
		return nil, qErr
	}
	defer rows.Close()

	items := make([]*domain.ProductionYield, 0)
	for rows.Next() {
		y := &domain.ProductionYield{DefectDetails: make([]domain.DefectDetail, 0)}
		if err := rows.Scan(
			&y.ID, &y.ProductionOrderID,
			&y.OrderNumber, &y.FGName,
			&y.YieldQty, &y.DefectQty, &y.WasteQty,
			&y.BatchNumber,
			&y.ProductionDate,
			&y.ExpiryDate,
			&y.RecordedBy, &y.RecordedByName,
			&y.RecordedAt,
			&y.Notes,
		); err != nil {
			return nil, err
		}
		items = append(items, y)
	}

	// Attach defect details
	for _, y := range items {
		details, err := r.listDefectDetails(ctx, y.ID)
		if err != nil {
			return nil, err
		}
		y.DefectDetails = details
	}
	return items, nil
}
