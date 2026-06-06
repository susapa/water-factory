package domain

import "time"

type ProductionOrder struct {
	ID               string                  `json:"id"`
	OrderNumber      string                  `json:"order_number"`
	FinishedGoodID   string                  `json:"finished_good_id"`
	FGCode           string                  `json:"fg_code"`
	FGName           string                  `json:"fg_name"`
	BOMID            *string                 `json:"bom_id"`
	PlannedQty       float64                 `json:"planned_qty"`
	ActualYieldQty   float64                 `json:"actual_yield_qty"`
	DefectQty        float64                 `json:"defect_qty"`
	WasteQty         float64                 `json:"waste_qty"`
	PlannedStartDate *string                 `json:"planned_start_date"`
	PlannedEndDate   *string                 `json:"planned_end_date"`
	ActualStartDate  *time.Time              `json:"actual_start_date"`
	ActualEndDate    *time.Time              `json:"actual_end_date"`
	Status           string                  `json:"status"`
	CreatedBy        string                  `json:"created_by"`
	CreatedByName    string                  `json:"created_by_name"`
	SalesOrderID     *string                 `json:"sales_order_id"`
	Notes            string                  `json:"notes"`
	Requirements     []ProductionRequirement `json:"requirements"`
	CreatedAt        time.Time               `json:"created_at"`
	UpdatedAt        time.Time               `json:"updated_at"`
}

type ProductionRequirement struct {
	ID                string  `json:"id"`
	ProductionOrderID string  `json:"production_order_id"`
	RawMaterialID     string  `json:"raw_material_id"`
	RMCode            string  `json:"rm_code"`
	RMName            string  `json:"rm_name"`
	UOMCode           string  `json:"uom_code"`
	RequiredQty       float64 `json:"required_qty"`
	IssuedQty         float64 `json:"issued_qty"`
	Status            string  `json:"status"`
}

type ProductionRMIssue struct {
	ID                string    `json:"id"`
	ProductionOrderID string    `json:"production_order_id"`
	RMStockLotID      string    `json:"rm_stock_lot_id"`
	LotNumber         string    `json:"lot_number"`
	RawMaterialID     string    `json:"raw_material_id"`
	RMCode            string    `json:"rm_code"`
	RMName            string    `json:"rm_name"`
	IssuedQty         float64   `json:"issued_qty"`
	IssuedBy          string    `json:"issued_by"`
	IssuedByName      string    `json:"issued_by_name"`
	IssuedAt          time.Time `json:"issued_at"`
}

type ProductionYield struct {
	ID                string        `json:"id"`
	ProductionOrderID string        `json:"production_order_id"`
	OrderNumber       string        `json:"order_number"`
	FGName            string        `json:"fg_name"`
	YieldQty          float64       `json:"yield_qty"`
	DefectQty         float64       `json:"defect_qty"`
	WasteQty          float64       `json:"waste_qty"`
	BatchNumber       string        `json:"batch_number"`
	ProductionDate    string        `json:"production_date"`
	ExpiryDate        string        `json:"expiry_date"`
	RecordedBy        string        `json:"recorded_by"`
	RecordedByName    string        `json:"recorded_by_name"`
	RecordedAt        time.Time     `json:"recorded_at"`
	Notes             string        `json:"notes"`
	DefectDetails     []DefectDetail `json:"defect_details"`
}

type DefectDetail struct {
	ID          string  `json:"id"`
	YieldID     string  `json:"yield_id"`
	DefectType  string  `json:"defect_type"`
	Qty         float64 `json:"qty"`
	Description string  `json:"description"`
}
