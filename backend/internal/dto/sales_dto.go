package dto

// ── Sales Order ───────────────────────────────────────────────────────────────

type CreateSOLineRequest struct {
	FinishedGoodID string  `json:"finished_good_id" binding:"required"`
	OrderedQty     float64 `json:"ordered_qty"      binding:"required,gt=0"`
	UnitPrice      float64 `json:"unit_price"       binding:"required,gte=0"`
	DiscountPct    float64 `json:"discount_pct"`
}

type CreateSalesOrderRequest struct {
	CustomerID    string              `json:"customer_id"    binding:"required"`
	OrderDate     string              `json:"order_date"     binding:"required"`
	RequestedDate string              `json:"requested_date"`
	VATRate       float64             `json:"vat_rate"`
	Notes         string              `json:"notes"`
	Lines         []CreateSOLineRequest `json:"lines" binding:"required,min=1"`
}

type UpdateSalesOrderRequest struct {
	CustomerID    string              `json:"customer_id"`
	OrderDate     string              `json:"order_date"`
	RequestedDate string              `json:"requested_date"`
	VATRate       float64             `json:"vat_rate"`
	Notes         string              `json:"notes"`
	Lines         []CreateSOLineRequest `json:"lines" binding:"min=1"`
}

// ── Delivery Order ────────────────────────────────────────────────────────────

type CreateDeliveryOrderRequest struct {
	SalesOrderID string  `json:"sales_order_id" binding:"required"`
	DeliveryDate string  `json:"delivery_date"`
	DriverUserID string  `json:"driver_user_id"`
	VehicleID    *int    `json:"vehicle_id"`
	RouteNotes   string  `json:"route_notes"`
}

// ── Invoice ───────────────────────────────────────────────────────────────────

type CreateInvoiceRequest struct {
	SalesOrderID    string `json:"sales_order_id"    binding:"required"`
	DeliveryOrderID string `json:"delivery_order_id"`
	InvoiceDate     string `json:"invoice_date"      binding:"required"`
	DueDate         string `json:"due_date"`
}

type MarkPaidRequest struct {
	PaymentDate   string `json:"payment_date"   binding:"required"`
	PaymentMethod string `json:"payment_method"`
}
