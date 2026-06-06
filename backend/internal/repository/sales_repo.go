package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/water-factory/api/internal/domain"
	"github.com/water-factory/api/internal/dto"
	"github.com/water-factory/api/pkg/number"
)

type SalesRepo struct {
	pool *pgxpool.Pool
}

func NewSalesRepo(pool *pgxpool.Pool) *SalesRepo {
	return &SalesRepo{pool: pool}
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func calcSOTotals(lines []dto.CreateSOLineRequest, vatRate float64) (subtotal, vatAmt, total float64) {
	for _, l := range lines {
		lineTotal := l.OrderedQty * l.UnitPrice * (1 - l.DiscountPct/100)
		subtotal += lineTotal
	}
	vatAmt = subtotal * (vatRate / 100)
	total = subtotal + vatAmt
	return
}

// ── Vehicles ──────────────────────────────────────────────────────────────────

func (r *SalesRepo) ListVehicles(ctx context.Context) ([]*domain.Vehicle, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, license_plate, COALESCE(type,''), COALESCE(capacity_units,0)::float8, is_active
		 FROM vehicles WHERE is_active = true ORDER BY license_plate`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*domain.Vehicle, 0)
	for rows.Next() {
		v := &domain.Vehicle{}
		if err := rows.Scan(&v.ID, &v.LicensePlate, &v.Type, &v.CapacityUnits, &v.IsActive); err != nil {
			return nil, err
		}
		items = append(items, v)
	}
	return items, nil
}

// ── Sales Orders ──────────────────────────────────────────────────────────────

const soSelectCols = `
	so.id, so.order_number,
	so.customer_id::text, COALESCE(c.name,''),
	so.order_date::text,
	CASE WHEN so.requested_date IS NULL THEN NULL ELSE so.requested_date::text END,
	so.status,
	COALESCE(so.subtotal,0)::float8, COALESCE(so.vat_amount,0)::float8, COALESCE(so.total_amount,0)::float8,
	COALESCE(so.vat_rate,7)::float8,
	COALESCE(so.notes,''),
	COALESCE(so.created_by::text,''), COALESCE(u.full_name,''),
	so.created_at, so.updated_at`

const soFrom = `
	FROM sales_orders so
	LEFT JOIN customers c ON c.id = so.customer_id
	LEFT JOIN users u ON u.id = so.created_by`

func scanSO(row interface{ Scan(...any) error }) (*domain.SalesOrder, error) {
	s := &domain.SalesOrder{Lines: make([]domain.SalesOrderLine, 0)}
	err := row.Scan(
		&s.ID, &s.OrderNumber,
		&s.CustomerID, &s.CustomerName,
		&s.OrderDate, &s.RequestedDate,
		&s.Status,
		&s.Subtotal, &s.VATAmount, &s.TotalAmount, &s.VATRate,
		&s.Notes,
		&s.CreatedBy, &s.CreatedByName,
		&s.CreatedAt, &s.UpdatedAt,
	)
	return s, err
}

func (r *SalesRepo) loadSOLines(ctx context.Context, soID string) ([]domain.SalesOrderLine, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT l.id, l.sales_order_id::text,
		       l.finished_good_id::text, fg.code, fg.name, u.code,
		       l.ordered_qty::float8, l.unit_price::float8,
		       COALESCE(l.discount_pct,0)::float8, COALESCE(l.line_total,0)::float8,
		       l.status
		FROM sales_order_lines l
		JOIN finished_goods fg ON fg.id = l.finished_good_id
		JOIN units_of_measure u ON u.id = fg.uom_id
		WHERE l.sales_order_id = $1
		ORDER BY fg.code`, soID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	lines := make([]domain.SalesOrderLine, 0)
	for rows.Next() {
		l := domain.SalesOrderLine{}
		if err := rows.Scan(&l.ID, &l.SalesOrderID,
			&l.FinishedGoodID, &l.FGCode, &l.FGName, &l.UOMCode,
			&l.OrderedQty, &l.UnitPrice, &l.DiscountPct, &l.LineTotal,
			&l.Status); err != nil {
			return nil, err
		}
		lines = append(lines, l)
	}
	return lines, nil
}

func (r *SalesRepo) ListSalesOrders(ctx context.Context) ([]*domain.SalesOrder, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT `+soSelectCols+soFrom+` ORDER BY so.created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*domain.SalesOrder, 0)
	for rows.Next() {
		s, err := scanSO(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, s)
	}
	return items, nil
}

func (r *SalesRepo) GetSalesOrder(ctx context.Context, id string) (*domain.SalesOrder, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT `+soSelectCols+soFrom+` WHERE so.id = $1`, id)
	s, err := scanSO(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("sales order not found")
		}
		return nil, err
	}
	s.Lines, err = r.loadSOLines(ctx, id)
	return s, err
}

func (r *SalesRepo) CreateSalesOrder(ctx context.Context, userID string, req *dto.CreateSalesOrderRequest) (*domain.SalesOrder, error) {
	orderNumber, err := number.Next(ctx, r.pool, "SO")
	if err != nil {
		return nil, err
	}

	if req.VATRate == 0 {
		req.VATRate = 7.0
	}
	subtotal, vatAmt, total := calcSOTotals(req.Lines, req.VATRate)

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var soID string
	var reqDate *string
	if req.RequestedDate != "" {
		reqDate = &req.RequestedDate
	}

	err = tx.QueryRow(ctx, `
		INSERT INTO sales_orders
			(order_number, customer_id, order_date, requested_date, status,
			 subtotal, vat_amount, total_amount, vat_rate, notes, created_by)
		VALUES ($1,$2,$3,$4,'draft',$5,$6,$7,$8,$9,$10)
		RETURNING id`,
		orderNumber, req.CustomerID, req.OrderDate, reqDate,
		subtotal, vatAmt, total, req.VATRate, req.Notes, userID,
	).Scan(&soID)
	if err != nil {
		return nil, fmt.Errorf("insert sales_order: %w", err)
	}

	for _, l := range req.Lines {
		lineTotal := l.OrderedQty * l.UnitPrice * (1 - l.DiscountPct/100)
		_, err = tx.Exec(ctx, `
			INSERT INTO sales_order_lines
				(sales_order_id, finished_good_id, ordered_qty, unit_price, discount_pct, line_total, status)
			VALUES ($1,$2,$3,$4,$5,$6,'pending')`,
			soID, l.FinishedGoodID, l.OrderedQty, l.UnitPrice, l.DiscountPct, lineTotal)
		if err != nil {
			return nil, fmt.Errorf("insert so_line: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return r.GetSalesOrder(ctx, soID)
}

func (r *SalesRepo) UpdateSalesOrder(ctx context.Context, id string, req *dto.UpdateSalesOrderRequest) (*domain.SalesOrder, error) {
	// check draft
	var status string
	if err := r.pool.QueryRow(ctx, `SELECT status FROM sales_orders WHERE id=$1`, id).Scan(&status); err != nil {
		return nil, fmt.Errorf("sales order not found")
	}
	if status != "draft" {
		return nil, fmt.Errorf("only draft orders can be updated")
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var vatRate float64
	if req.VATRate == 0 {
		vatRate = 7.0
	} else {
		vatRate = req.VATRate
	}

	var reqDate *string
	if req.RequestedDate != "" {
		reqDate = &req.RequestedDate
	}

	subtotal, vatAmt, total := calcSOTotals(req.Lines, vatRate)

	_, err = tx.Exec(ctx, `
		UPDATE sales_orders SET
			customer_id=$1, order_date=$2, requested_date=$3,
			subtotal=$4, vat_amount=$5, total_amount=$6, vat_rate=$7, notes=$8,
			updated_at=NOW()
		WHERE id=$9`,
		req.CustomerID, req.OrderDate, reqDate,
		subtotal, vatAmt, total, vatRate, req.Notes, id)
	if err != nil {
		return nil, err
	}

	// Replace lines
	if _, err := tx.Exec(ctx, `DELETE FROM sales_order_lines WHERE sales_order_id=$1`, id); err != nil {
		return nil, err
	}
	for _, l := range req.Lines {
		lineTotal := l.OrderedQty * l.UnitPrice * (1 - l.DiscountPct/100)
		_, err = tx.Exec(ctx, `
			INSERT INTO sales_order_lines
				(sales_order_id, finished_good_id, ordered_qty, unit_price, discount_pct, line_total, status)
			VALUES ($1,$2,$3,$4,$5,$6,'pending')`,
			id, l.FinishedGoodID, l.OrderedQty, l.UnitPrice, l.DiscountPct, lineTotal)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return r.GetSalesOrder(ctx, id)
}

func (r *SalesRepo) ConfirmSalesOrder(ctx context.Context, id string) (*domain.SalesOrder, error) {
	var status string
	if err := r.pool.QueryRow(ctx, `SELECT status FROM sales_orders WHERE id=$1`, id).Scan(&status); err != nil {
		return nil, fmt.Errorf("sales order not found")
	}
	if status != "draft" {
		return nil, fmt.Errorf("conflict: order is already %s", status)
	}

	_, err := r.pool.Exec(ctx,
		`UPDATE sales_orders SET status='confirmed', updated_at=NOW() WHERE id=$1`, id)
	if err != nil {
		return nil, err
	}
	return r.GetSalesOrder(ctx, id)
}

func (r *SalesRepo) CancelSalesOrder(ctx context.Context, id string) (*domain.SalesOrder, error) {
	var status string
	if err := r.pool.QueryRow(ctx, `SELECT status FROM sales_orders WHERE id=$1`, id).Scan(&status); err != nil {
		return nil, fmt.Errorf("sales order not found")
	}
	if status != "draft" && status != "confirmed" {
		return nil, fmt.Errorf("conflict: cannot cancel order with status %s", status)
	}

	_, err := r.pool.Exec(ctx,
		`UPDATE sales_orders SET status='cancelled', updated_at=NOW() WHERE id=$1`, id)
	if err != nil {
		return nil, err
	}
	return r.GetSalesOrder(ctx, id)
}

// ── Delivery Orders ───────────────────────────────────────────────────────────

const doSelectCols = `
	d.id, d.do_number,
	d.sales_order_id::text, COALESCE(so.order_number,''),
	COALESCE(c.name,''),
	CASE WHEN d.delivery_date IS NULL THEN NULL ELSE d.delivery_date::text END,
	CASE WHEN d.driver_user_id IS NULL THEN NULL ELSE d.driver_user_id::text END,
	COALESCE(u.full_name,''),
	d.vehicle_id,
	COALESCE(v.license_plate,''),
	d.status,
	d.dispatched_at, d.delivered_at,
	COALESCE(d.route_notes,''),
	d.created_at, d.updated_at`

const doFrom = `
	FROM delivery_orders d
	LEFT JOIN sales_orders so ON so.id = d.sales_order_id
	LEFT JOIN customers c ON c.id = so.customer_id
	LEFT JOIN users u ON u.id = d.driver_user_id
	LEFT JOIN vehicles v ON v.id = d.vehicle_id`

func scanDO(row interface{ Scan(...any) error }) (*domain.DeliveryOrder, error) {
	d := &domain.DeliveryOrder{Lines: make([]domain.DeliveryOrderLine, 0)}
	err := row.Scan(
		&d.ID, &d.DONumber,
		&d.SalesOrderID, &d.SONumber,
		&d.CustomerName,
		&d.DeliveryDate,
		&d.DriverUserID, &d.DriverName,
		&d.VehicleID, &d.LicensePlate,
		&d.Status,
		&d.DispatchedAt, &d.DeliveredAt,
		&d.RouteNotes,
		&d.CreatedAt, &d.UpdatedAt,
	)
	return d, err
}

func (r *SalesRepo) loadDOLines(ctx context.Context, doID string) ([]domain.DeliveryOrderLine, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT dl.id, dl.delivery_order_id::text,
		       COALESCE(dl.sales_order_line_id::text,''),
		       dl.fg_stock_lot_id::text,
		       COALESCE(sl.batch_number,''),
		       dl.finished_good_id::text,
		       fg.code, fg.name,
		       dl.picked_qty::float8,
		       COALESCE(dl.dispatched_qty,0)::float8
		FROM delivery_order_lines dl
		JOIN fg_stock_lots sl ON sl.id = dl.fg_stock_lot_id
		JOIN finished_goods fg ON fg.id = dl.finished_good_id
		WHERE dl.delivery_order_id = $1
		ORDER BY fg.code, sl.expiry_date`, doID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	lines := make([]domain.DeliveryOrderLine, 0)
	for rows.Next() {
		l := domain.DeliveryOrderLine{}
		if err := rows.Scan(&l.ID, &l.DeliveryOrderID,
			&l.SalesOrderLineID, &l.FGStockLotID, &l.BatchNumber,
			&l.FinishedGoodID, &l.FGCode, &l.FGName,
			&l.PickedQty, &l.DispatchedQty); err != nil {
			return nil, err
		}
		lines = append(lines, l)
	}
	return lines, nil
}

func (r *SalesRepo) ListDeliveryOrders(ctx context.Context) ([]*domain.DeliveryOrder, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT `+doSelectCols+doFrom+` ORDER BY d.created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*domain.DeliveryOrder, 0)
	for rows.Next() {
		d, err := scanDO(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, d)
	}
	return items, nil
}

func (r *SalesRepo) GetDeliveryOrder(ctx context.Context, id string) (*domain.DeliveryOrder, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT `+doSelectCols+doFrom+` WHERE d.id = $1`, id)
	d, err := scanDO(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("delivery order not found")
		}
		return nil, err
	}
	d.Lines, err = r.loadDOLines(ctx, id)
	return d, err
}

// CreateDeliveryOrder auto-picks FG lots using FEFO for each SO line.
func (r *SalesRepo) CreateDeliveryOrder(ctx context.Context, req *dto.CreateDeliveryOrderRequest) (*domain.DeliveryOrder, error) {
	// Validate SO status
	var soStatus string
	if err := r.pool.QueryRow(ctx, `SELECT status FROM sales_orders WHERE id=$1`, req.SalesOrderID).Scan(&soStatus); err != nil {
		return nil, fmt.Errorf("sales order not found")
	}
	if soStatus != "confirmed" {
		return nil, fmt.Errorf("conflict: sales order must be confirmed before creating delivery (current: %s)", soStatus)
	}

	// Load SO lines
	soLines, err := r.loadSOLines(ctx, req.SalesOrderID)
	if err != nil {
		return nil, err
	}

	// For each SO line, FEFO-pick lots
	type doLine struct {
		soLineID string
		lotID    string
		fgID     string
		qty      float64
	}
	var allocations []doLine

	for _, sl := range soLines {
		if sl.Status == "dispatched" {
			continue // already fulfilled
		}
		remaining := sl.OrderedQty

		// FEFO lots
		rows, err := r.pool.Query(ctx, `
			SELECT id::text, current_qty::float8
			FROM fg_stock_lots
			WHERE finished_good_id = $1 AND status = 'available' AND current_qty > 0
			ORDER BY expiry_date ASC NULLS LAST, created_at ASC`,
			sl.FinishedGoodID)
		if err != nil {
			return nil, err
		}

		type lot struct {
			id  string
			qty float64
		}
		var lots []lot
		for rows.Next() {
			var l lot
			if err := rows.Scan(&l.id, &l.qty); err != nil {
				rows.Close()
				return nil, err
			}
			lots = append(lots, l)
		}
		rows.Close()

		for _, l := range lots {
			if remaining <= 0 {
				break
			}
			pick := l.qty
			if pick > remaining {
				pick = remaining
			}
			allocations = append(allocations, doLine{soLineID: sl.ID, lotID: l.id, fgID: sl.FinishedGoodID, qty: pick})
			remaining -= pick
		}

		if remaining > 0.0001 {
			return nil, fmt.Errorf("insufficient stock for %s (short by %.3f)", sl.FGName, remaining)
		}
	}

	doNumber, err := number.Next(ctx, r.pool, "DO")
	if err != nil {
		return nil, err
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var driverID *string
	if req.DriverUserID != "" {
		driverID = &req.DriverUserID
	}
	var delDate *string
	if req.DeliveryDate != "" {
		delDate = &req.DeliveryDate
	}

	var doID string
	err = tx.QueryRow(ctx, `
		INSERT INTO delivery_orders
			(do_number, sales_order_id, delivery_date, driver_user_id, vehicle_id, status, route_notes)
		VALUES ($1,$2,$3,$4,$5,'pending',$6)
		RETURNING id`,
		doNumber, req.SalesOrderID, delDate, driverID, req.VehicleID, req.RouteNotes,
	).Scan(&doID)
	if err != nil {
		return nil, fmt.Errorf("insert delivery_order: %w", err)
	}

	for _, a := range allocations {
		var soLineID *string
		if a.soLineID != "" {
			soLineID = &a.soLineID
		}
		_, err = tx.Exec(ctx, `
			INSERT INTO delivery_order_lines
				(delivery_order_id, sales_order_line_id, fg_stock_lot_id, finished_good_id, picked_qty, dispatched_qty)
			VALUES ($1,$2,$3,$4,$5,0)`,
			doID, soLineID, a.lotID, a.fgID, a.qty)
		if err != nil {
			return nil, fmt.Errorf("insert do_line: %w", err)
		}
	}

	// Update SO status to picking
	_, err = tx.Exec(ctx, `UPDATE sales_orders SET status='picking', updated_at=NOW() WHERE id=$1`, req.SalesOrderID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return r.GetDeliveryOrder(ctx, doID)
}

// DispatchDeliveryOrder deducts FG stock and inserts SALES_DISPATCH movements.
func (r *SalesRepo) DispatchDeliveryOrder(ctx context.Context, id string, userID string) (*domain.DeliveryOrder, error) {
	var doStatus, soID string
	err := r.pool.QueryRow(ctx,
		`SELECT d.status, d.sales_order_id::text FROM delivery_orders d WHERE d.id=$1`, id,
	).Scan(&doStatus, &soID)
	if err != nil {
		return nil, fmt.Errorf("delivery order not found")
	}
	if doStatus != "pending" && doStatus != "loading" {
		return nil, fmt.Errorf("conflict: delivery order is already %s", doStatus)
	}

	lines, err := r.loadDOLines(ctx, id)
	if err != nil {
		return nil, err
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	for _, l := range lines {
		var qtyBefore float64
		if err := tx.QueryRow(ctx,
			`SELECT current_qty FROM fg_stock_lots WHERE id=$1`, l.FGStockLotID,
		).Scan(&qtyBefore); err != nil {
			return nil, fmt.Errorf("lot %s not found", l.FGStockLotID)
		}

		qtyAfter := qtyBefore - l.PickedQty
		if qtyAfter < -0.0001 {
			return nil, fmt.Errorf("insufficient stock in lot %s", l.BatchNumber)
		}
		if qtyAfter < 0 {
			qtyAfter = 0
		}

		_, err = tx.Exec(ctx,
			`UPDATE fg_stock_lots SET current_qty=$1, updated_at=NOW() WHERE id=$2`,
			qtyAfter, l.FGStockLotID)
		if err != nil {
			return nil, err
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO fg_stock_movements
				(fg_stock_lot_id, finished_good_id, movement_type, reference_type, reference_id,
				 qty, qty_before, qty_after, performed_by, notes)
			VALUES ($1,$2,'SALES_DISPATCH','delivery_order',$3,$4,$5,$6,$7,'')`,
			l.FGStockLotID, l.FinishedGoodID, id, l.PickedQty, qtyBefore, qtyAfter, userID)
		if err != nil {
			return nil, fmt.Errorf("insert movement: %w", err)
		}

		// Update dispatched_qty on line
		_, err = tx.Exec(ctx,
			`UPDATE delivery_order_lines SET dispatched_qty=$1 WHERE id=$2`,
			l.PickedQty, l.ID)
		if err != nil {
			return nil, err
		}
	}

	now := time.Now()
	_, err = tx.Exec(ctx,
		`UPDATE delivery_orders SET status='dispatched', dispatched_at=$1, updated_at=NOW() WHERE id=$2`,
		now, id)
	if err != nil {
		return nil, err
	}

	// Update SO lines status
	_, err = tx.Exec(ctx, `
		UPDATE sales_order_lines SET status='dispatched'
		WHERE sales_order_id=$1 AND id IN (
			SELECT DISTINCT sales_order_line_id FROM delivery_order_lines WHERE delivery_order_id=$2
		)`, soID, id)
	if err != nil {
		return nil, err
	}

	// Update SO status
	_, err = tx.Exec(ctx,
		`UPDATE sales_orders SET status='dispatched', updated_at=NOW() WHERE id=$1`, soID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return r.GetDeliveryOrder(ctx, id)
}

func (r *SalesRepo) DeliverDeliveryOrder(ctx context.Context, id string) (*domain.DeliveryOrder, error) {
	var status string
	if err := r.pool.QueryRow(ctx, `SELECT status FROM delivery_orders WHERE id=$1`, id).Scan(&status); err != nil {
		return nil, fmt.Errorf("delivery order not found")
	}
	if status != "dispatched" {
		return nil, fmt.Errorf("conflict: delivery order must be dispatched first (current: %s)", status)
	}

	now := time.Now()
	_, err := r.pool.Exec(ctx,
		`UPDATE delivery_orders SET status='delivered', delivered_at=$1, updated_at=NOW() WHERE id=$2`,
		now, id)
	if err != nil {
		return nil, err
	}
	return r.GetDeliveryOrder(ctx, id)
}

// ── Invoices ──────────────────────────────────────────────────────────────────

const invSelectCols = `
	i.id, i.invoice_number,
	i.sales_order_id::text, COALESCE(so.order_number,''),
	COALESCE(c.name,''),
	CASE WHEN i.delivery_order_id IS NULL THEN NULL ELSE i.delivery_order_id::text END,
	COALESCE(d.do_number,''),
	i.invoice_date::text,
	CASE WHEN i.due_date IS NULL THEN NULL ELSE i.due_date::text END,
	COALESCE(i.subtotal,0)::float8, COALESCE(i.vat_amount,0)::float8, COALESCE(i.total_amount,0)::float8,
	i.status,
	CASE WHEN i.payment_date IS NULL THEN NULL ELSE i.payment_date::text END,
	COALESCE(i.payment_method,''),
	COALESCE(i.created_by::text,''), COALESCE(u.full_name,''),
	i.created_at`

const invFrom = `
	FROM invoices i
	LEFT JOIN sales_orders so ON so.id = i.sales_order_id
	LEFT JOIN customers c ON c.id = so.customer_id
	LEFT JOIN delivery_orders d ON d.id = i.delivery_order_id
	LEFT JOIN users u ON u.id = i.created_by`

func scanInvoice(row interface{ Scan(...any) error }) (*domain.Invoice, error) {
	inv := &domain.Invoice{}
	err := row.Scan(
		&inv.ID, &inv.InvoiceNumber,
		&inv.SalesOrderID, &inv.SONumber,
		&inv.CustomerName,
		&inv.DeliveryOrderID, &inv.DONumber,
		&inv.InvoiceDate, &inv.DueDate,
		&inv.Subtotal, &inv.VATAmount, &inv.TotalAmount,
		&inv.Status,
		&inv.PaymentDate, &inv.PaymentMethod,
		&inv.CreatedBy, &inv.CreatedByName,
		&inv.CreatedAt,
	)
	return inv, err
}

func (r *SalesRepo) ListInvoices(ctx context.Context) ([]*domain.Invoice, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT `+invSelectCols+invFrom+` ORDER BY i.created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*domain.Invoice, 0)
	for rows.Next() {
		inv, err := scanInvoice(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, inv)
	}
	return items, nil
}

func (r *SalesRepo) GetInvoice(ctx context.Context, id string) (*domain.Invoice, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT `+invSelectCols+invFrom+` WHERE i.id = $1`, id)
	inv, err := scanInvoice(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("invoice not found")
		}
		return nil, err
	}
	return inv, nil
}

func (r *SalesRepo) CreateInvoice(ctx context.Context, userID string, req *dto.CreateInvoiceRequest) (*domain.Invoice, error) {
	// Validate SO exists and not already invoiced
	var soStatus, subtotal, vatAmt, total string
	var vat float64
	err := r.pool.QueryRow(ctx,
		`SELECT status, COALESCE(subtotal,0)::text, COALESCE(vat_amount,0)::text, COALESCE(total_amount,0)::text, COALESCE(vat_rate,7)
		 FROM sales_orders WHERE id=$1`, req.SalesOrderID,
	).Scan(&soStatus, &subtotal, &vatAmt, &total, &vat)
	if err != nil {
		return nil, fmt.Errorf("sales order not found")
	}
	if soStatus == "invoiced" || soStatus == "cancelled" {
		return nil, fmt.Errorf("conflict: sales order is already %s", soStatus)
	}

	// Check not already invoiced
	var existingCount int
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM invoices WHERE sales_order_id=$1 AND status != 'cancelled'`, req.SalesOrderID).Scan(&existingCount)
	if existingCount > 0 {
		return nil, fmt.Errorf("conflict: invoice already exists for this sales order")
	}

	invNumber, err := number.Next(ctx, r.pool, "INV")
	if err != nil {
		return nil, err
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var dueDate, doID *string
	if req.DueDate != "" {
		dueDate = &req.DueDate
	}
	if req.DeliveryOrderID != "" {
		doID = &req.DeliveryOrderID
	}

	var invID string
	err = tx.QueryRow(ctx, `
		INSERT INTO invoices
			(invoice_number, sales_order_id, delivery_order_id, invoice_date, due_date,
			 subtotal, vat_amount, total_amount, status, created_by)
		SELECT $1,$2,$3,$4,$5,
			   COALESCE(subtotal,0), COALESCE(vat_amount,0), COALESCE(total_amount,0),
			   'issued', $6
		FROM sales_orders WHERE id=$2
		RETURNING id`,
		invNumber, req.SalesOrderID, doID, req.InvoiceDate, dueDate, userID,
	).Scan(&invID)
	if err != nil {
		return nil, fmt.Errorf("insert invoice: %w", err)
	}

	// Update SO to invoiced
	_, err = tx.Exec(ctx,
		`UPDATE sales_orders SET status='invoiced', updated_at=NOW() WHERE id=$1`, req.SalesOrderID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return r.GetInvoice(ctx, invID)
}

func (r *SalesRepo) MarkInvoicePaid(ctx context.Context, id string, req *dto.MarkPaidRequest) (*domain.Invoice, error) {
	var status string
	if err := r.pool.QueryRow(ctx, `SELECT status FROM invoices WHERE id=$1`, id).Scan(&status); err != nil {
		return nil, fmt.Errorf("invoice not found")
	}
	if status == "paid" {
		return nil, fmt.Errorf("conflict: invoice is already paid")
	}
	if status == "cancelled" {
		return nil, fmt.Errorf("conflict: invoice is cancelled")
	}

	_, err := r.pool.Exec(ctx, `
		UPDATE invoices SET status='paid', payment_date=$1, payment_method=$2 WHERE id=$3`,
		req.PaymentDate, req.PaymentMethod, id)
	if err != nil {
		return nil, err
	}
	return r.GetInvoice(ctx, id)
}
