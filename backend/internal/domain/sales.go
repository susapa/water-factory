package domain

import "time"

// ── Sales Order ───────────────────────────────────────────────────────────────

type SalesOrderLine struct {
	ID             string  `json:"id"`
	SalesOrderID   string  `json:"sales_order_id"`
	FinishedGoodID string  `json:"finished_good_id"`
	FGCode         string  `json:"fg_code"`
	FGName         string  `json:"fg_name"`
	UOMCode        string  `json:"uom_code"`
	OrderedQty     float64 `json:"ordered_qty"`
	UnitPrice      float64 `json:"unit_price"`
	DiscountPct    float64 `json:"discount_pct"`
	LineTotal      float64 `json:"line_total"`
	Status         string  `json:"status"` // pending | allocated | dispatched
}

type SalesOrder struct {
	ID            string          `json:"id"`
	OrderNumber   string          `json:"order_number"`
	CustomerID    string          `json:"customer_id"`
	CustomerName  string          `json:"customer_name"`
	OrderDate     string          `json:"order_date"`
	RequestedDate *string         `json:"requested_date"`
	Status        string          `json:"status"` // draft|confirmed|picking|dispatched|invoiced|cancelled
	Subtotal      float64         `json:"subtotal"`
	VATAmount     float64         `json:"vat_amount"`
	TotalAmount   float64         `json:"total_amount"`
	VATRate       float64         `json:"vat_rate"`
	Notes         string          `json:"notes"`
	CreatedBy     string          `json:"created_by"`
	CreatedByName string          `json:"created_by_name"`
	Lines         []SalesOrderLine `json:"lines"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

// ── Vehicle ───────────────────────────────────────────────────────────────────

type Vehicle struct {
	ID            int     `json:"id"`
	LicensePlate  string  `json:"license_plate"`
	Type          string  `json:"type"`
	CapacityUnits float64 `json:"capacity_units"`
	IsActive      bool    `json:"is_active"`
}

// ── Delivery Order ────────────────────────────────────────────────────────────

type DeliveryOrderLine struct {
	ID               string  `json:"id"`
	DeliveryOrderID  string  `json:"delivery_order_id"`
	SalesOrderLineID string  `json:"sales_order_line_id"`
	FGStockLotID     string  `json:"fg_stock_lot_id"`
	BatchNumber      string  `json:"batch_number"`
	FinishedGoodID   string  `json:"finished_good_id"`
	FGCode           string  `json:"fg_code"`
	FGName           string  `json:"fg_name"`
	PickedQty        float64 `json:"picked_qty"`
	DispatchedQty    float64 `json:"dispatched_qty"`
}

type DeliveryOrder struct {
	ID           string              `json:"id"`
	DONumber     string              `json:"do_number"`
	SalesOrderID string              `json:"sales_order_id"`
	SONumber     string              `json:"so_number"`
	CustomerName string              `json:"customer_name"`
	DeliveryDate *string             `json:"delivery_date"`
	DriverUserID *string             `json:"driver_user_id"`
	DriverName   string              `json:"driver_name"`
	VehicleID    *int                `json:"vehicle_id"`
	LicensePlate string              `json:"license_plate"`
	Status       string              `json:"status"` // pending|loading|dispatched|delivered|failed
	DispatchedAt *time.Time          `json:"dispatched_at"`
	DeliveredAt  *time.Time          `json:"delivered_at"`
	RouteNotes   string              `json:"route_notes"`
	Lines        []DeliveryOrderLine `json:"lines"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}

// ── Invoice ───────────────────────────────────────────────────────────────────

type Invoice struct {
	ID              string    `json:"id"`
	InvoiceNumber   string    `json:"invoice_number"`
	SalesOrderID    string    `json:"sales_order_id"`
	SONumber        string    `json:"so_number"`
	CustomerName    string    `json:"customer_name"`
	DeliveryOrderID *string   `json:"delivery_order_id"`
	DONumber        string    `json:"do_number"`
	InvoiceDate     string    `json:"invoice_date"`
	DueDate         *string   `json:"due_date"`
	Subtotal        float64   `json:"subtotal"`
	VATAmount       float64   `json:"vat_amount"`
	TotalAmount     float64   `json:"total_amount"`
	Status          string    `json:"status"` // issued|paid|overdue|cancelled
	PaymentDate     *string   `json:"payment_date"`
	PaymentMethod   string    `json:"payment_method"`
	CreatedBy       string    `json:"created_by"`
	CreatedByName   string    `json:"created_by_name"`
	CreatedAt       time.Time `json:"created_at"`
}
