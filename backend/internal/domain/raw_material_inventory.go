package domain

import "time"

type WarehouseLocation struct {
	ID          int    `json:"id"`
	Zone        string `json:"zone"`
	RowNo       string `json:"row_no"`
	BayNo       string `json:"bay_no"`
	Description string `json:"description"`
}

type GRNLine struct {
	ID            string   `json:"id"`
	GRNID         string   `json:"grn_id"`
	RawMaterialID string   `json:"raw_material_id"`
	RMCode        string   `json:"rm_code"`
	RMName        string   `json:"rm_name"`
	LotNumber     string   `json:"lot_number"`
	ReceivedQty   float64  `json:"received_qty"`
	UnitCost      *float64 `json:"unit_cost"`
	ExpiryDate    *string  `json:"expiry_date"`
	LocationID    *int     `json:"location_id"`
	LocationLabel string   `json:"location_label"`
	Notes         string   `json:"notes"`
}

type GRN struct {
	ID             string    `json:"id"`
	GRNNumber      string    `json:"grn_number"`
	SupplierID     *string   `json:"supplier_id"`
	SupplierName   string    `json:"supplier_name"`
	ReceivedBy     string    `json:"received_by"`
	ReceivedByName string    `json:"received_by_name"`
	ReceivedDate   string    `json:"received_date"`
	POReference    string    `json:"po_reference"`
	Status         string    `json:"status"`
	Notes          string    `json:"notes"`
	Lines          []GRNLine `json:"lines"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type StockLot struct {
	ID            string    `json:"id"`
	RawMaterialID string    `json:"raw_material_id"`
	RMCode        string    `json:"rm_code"`
	RMName        string    `json:"rm_name"`
	UOMCode       string    `json:"uom_code"`
	LotNumber     string    `json:"lot_number"`
	GRNLineID     *string   `json:"grn_line_id"`
	ReceivedDate  string    `json:"received_date"`
	ExpiryDate    *string   `json:"expiry_date"`
	LocationID    *int      `json:"location_id"`
	LocationLabel string    `json:"location_label"`
	InitialQty    float64   `json:"initial_qty"`
	CurrentQty    float64   `json:"current_qty"`
	UnitCost      *float64  `json:"unit_cost"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type StockMovement struct {
	ID              string    `json:"id"`
	RMStockLotID    string    `json:"rm_stock_lot_id"`
	RawMaterialID   string    `json:"raw_material_id"`
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

type StockLotDetail struct {
	StockLot
	Movements []StockMovement `json:"movements"`
}

type StockSummary struct {
	RawMaterialID  string  `json:"raw_material_id"`
	RMCode         string  `json:"rm_code"`
	RMName         string  `json:"rm_name"`
	UOMCode        string  `json:"uom_code"`
	TotalQty       float64 `json:"total_qty"`
	LotCount       int     `json:"lot_count"`
	EarliestExpiry *string `json:"earliest_expiry"`
}

type AdjustmentLine struct {
	ID           string  `json:"id"`
	AdjustmentID string  `json:"adjustment_id"`
	RMStockLotID string  `json:"rm_stock_lot_id"`
	LotNumber    string  `json:"lot_number"`
	RawMaterialID string `json:"raw_material_id"`
	RMCode       string  `json:"rm_code"`
	RMName       string  `json:"rm_name"`
	SystemQty    float64 `json:"system_qty"`
	CountedQty   float64 `json:"counted_qty"`
	VarianceQty  float64 `json:"variance_qty"`
	Reason       string  `json:"reason"`
}

type Adjustment struct {
	ID              string           `json:"id"`
	AdjNumber       string           `json:"adj_number"`
	AdjustmentDate  string           `json:"adjustment_date"`
	Type            string           `json:"type"`
	PerformedBy     string           `json:"performed_by"`
	PerformedByName string           `json:"performed_by_name"`
	ApprovedBy      *string          `json:"approved_by"`
	ApprovedByName  string           `json:"approved_by_name"`
	Status          string           `json:"status"`
	Notes           string           `json:"notes"`
	Lines           []AdjustmentLine `json:"lines"`
	CreatedAt       time.Time        `json:"created_at"`
}
