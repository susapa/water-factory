package domain

import "time"

type FGStockLot struct {
	ID                string    `json:"id"`
	FinishedGoodID    string    `json:"finished_good_id"`
	FGCode            string    `json:"fg_code"`
	FGName            string    `json:"fg_name"`
	UOMCode           string    `json:"uom_code"`
	BatchNumber       string    `json:"batch_number"`
	ProductionOrderID *string   `json:"production_order_id"`
	ProductionDate    string    `json:"production_date"`
	ExpiryDate        string    `json:"expiry_date"`
	LocationID        *int      `json:"location_id"`
	LocationLabel     string    `json:"location_label"`
	InitialQty        float64   `json:"initial_qty"`
	CurrentQty        float64   `json:"current_qty"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type FGStockMovement struct {
	ID              string    `json:"id"`
	FGStockLotID    string    `json:"fg_stock_lot_id"`
	FinishedGoodID  string    `json:"finished_good_id"`
	MovementType    string    `json:"movement_type"`
	ReferenceType   string    `json:"reference_type"`
	ReferenceID     *string   `json:"reference_id"`
	Qty             float64   `json:"qty"`
	QtyBefore       float64   `json:"qty_before"`
	QtyAfter        float64   `json:"qty_after"`
	PerformedBy     string    `json:"performed_by"`
	PerformedByName string    `json:"performed_by_name"`
	Notes           string    `json:"notes"`
	CreatedAt       time.Time `json:"created_at"`
}

type FGStockLotDetail struct {
	FGStockLot
	Movements []FGStockMovement `json:"movements"`
}

type FGStockSummary struct {
	FinishedGoodID string  `json:"finished_good_id"`
	FGCode         string  `json:"fg_code"`
	FGName         string  `json:"fg_name"`
	UOMCode        string  `json:"uom_code"`
	TotalQty       float64 `json:"total_qty"`
	LotCount       int     `json:"lot_count"`
	NearestExpiry  *string `json:"nearest_expiry"`
}

type FGAdjustmentLine struct {
	ID             string  `json:"id"`
	AdjustmentID   string  `json:"adjustment_id"`
	FGStockLotID   string  `json:"fg_stock_lot_id"`
	BatchNumber    string  `json:"batch_number"`
	FinishedGoodID string  `json:"finished_good_id"`
	FGCode         string  `json:"fg_code"`
	FGName         string  `json:"fg_name"`
	SystemQty      float64 `json:"system_qty"`
	CountedQty     float64 `json:"counted_qty"`
	VarianceQty    float64 `json:"variance_qty"`
	Reason         string  `json:"reason"`
}

type FGAdjustment struct {
	ID              string             `json:"id"`
	AdjNumber       string             `json:"adj_number"`
	AdjustmentDate  string             `json:"adjustment_date"`
	Type            string             `json:"type"`
	PerformedBy     string             `json:"performed_by"`
	PerformedByName string             `json:"performed_by_name"`
	ApprovedBy      *string            `json:"approved_by"`
	ApprovedByName  string             `json:"approved_by_name"`
	Status          string             `json:"status"`
	Notes           string             `json:"notes"`
	Lines           []FGAdjustmentLine `json:"lines"`
	CreatedAt       time.Time          `json:"created_at"`
}
