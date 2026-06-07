package domain

type DashboardSummary struct {
	PendingSOCount     int                  `json:"pending_so_count"`
	InProgressPOCount  int                  `json:"in_progress_po_count"`
	PendingDOCount     int                  `json:"pending_do_count"`
	IssuedInvoiceCount int                  `json:"issued_invoice_count"`
	IssuedInvoiceTotal float64              `json:"issued_invoice_total"`
	SalesLast7Days     []DailySales         `json:"sales_last_7_days"`
	POStatusCounts     map[string]int       `json:"po_status_counts"`
	FGExpiringSoon     []FGExpiringSoonItem `json:"fg_expiring_soon"`
}

type DailySales struct {
	Date  string  `json:"date"`
	Total float64 `json:"total"`
}

type FGExpiringSoonItem struct {
	LotNumber  string  `json:"lot_number"`
	FGName     string  `json:"fg_name"`
	Qty        float64 `json:"qty"`
	ExpiryDate string  `json:"expiry_date"`
}
