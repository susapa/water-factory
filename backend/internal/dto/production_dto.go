package dto

type CreateProductionOrderRequest struct {
	FinishedGoodID   string   `json:"finished_good_id" binding:"required"`
	BOMID            *string  `json:"bom_id"`
	PlannedQty       float64  `json:"planned_qty" binding:"required,gt=0"`
	PlannedStartDate *string  `json:"planned_start_date"`
	PlannedEndDate   *string  `json:"planned_end_date"`
	Notes            string   `json:"notes"`
}

type UpdateProductionOrderRequest struct {
	PlannedQty       float64  `json:"planned_qty" binding:"required,gt=0"`
	PlannedStartDate *string  `json:"planned_start_date"`
	PlannedEndDate   *string  `json:"planned_end_date"`
	Notes            string   `json:"notes"`
}

type IssueRMLine struct {
	RawMaterialID string  `json:"raw_material_id" binding:"required"`
	Qty           float64 `json:"qty" binding:"required,gt=0"`
}

type IssueRMRequest struct {
	Lines []IssueRMLine `json:"lines" binding:"required,min=1"`
}

type RecordYieldRequest struct {
	YieldQty       float64              `json:"yield_qty" binding:"required,gt=0"`
	DefectQty      float64              `json:"defect_qty"`
	WasteQty       float64              `json:"waste_qty"`
	BatchNumber    string               `json:"batch_number" binding:"required"`
	ProductionDate string               `json:"production_date" binding:"required"`
	ExpiryDate     string               `json:"expiry_date" binding:"required"`
	Notes          string               `json:"notes"`
	DefectDetails  []DefectDetailRequest `json:"defect_details"`
}

type DefectDetailRequest struct {
	DefectType  string  `json:"defect_type" binding:"required"`
	Qty         float64 `json:"qty" binding:"required,gt=0"`
	Description string  `json:"description"`
}
