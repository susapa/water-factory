package dto

type GRNLineRequest struct {
	RawMaterialID string   `json:"raw_material_id" binding:"required"`
	LotNumber     string   `json:"lot_number"`
	ReceivedQty   float64  `json:"received_qty" binding:"required,gt=0"`
	UnitCost      *float64 `json:"unit_cost"`
	ExpiryDate    *string  `json:"expiry_date"`
	LocationID    *int     `json:"location_id"`
	Notes         string   `json:"notes"`
}

type CreateGRNRequest struct {
	SupplierID   *string          `json:"supplier_id"`
	ReceivedDate string           `json:"received_date" binding:"required"`
	POReference  string           `json:"po_reference"`
	Notes        string           `json:"notes"`
	Lines        []GRNLineRequest `json:"lines" binding:"required,min=1,dive"`
}

type UpdateGRNRequest struct {
	SupplierID   *string          `json:"supplier_id"`
	ReceivedDate string           `json:"received_date" binding:"required"`
	POReference  string           `json:"po_reference"`
	Notes        string           `json:"notes"`
	Lines        []GRNLineRequest `json:"lines" binding:"required,min=1,dive"`
}

type AdjustmentLineRequest struct {
	RMStockLotID string  `json:"rm_stock_lot_id" binding:"required"`
	CountedQty   float64 `json:"counted_qty" binding:"required,min=0"`
	Reason       string  `json:"reason"`
}

type CreateAdjustmentRequest struct {
	AdjustmentDate string                  `json:"adjustment_date" binding:"required"`
	Type           string                  `json:"type" binding:"required,oneof=cycle_count damage correction"`
	Notes          string                  `json:"notes"`
	Lines          []AdjustmentLineRequest `json:"lines" binding:"required,min=1,dive"`
}
