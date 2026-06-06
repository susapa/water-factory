package dto

type FGAdjustmentLineRequest struct {
	FGStockLotID string  `json:"fg_stock_lot_id" binding:"required"`
	CountedQty   float64 `json:"counted_qty" binding:"required,min=0"`
	Reason       string  `json:"reason"`
}

type CreateFGAdjustmentRequest struct {
	AdjustmentDate string                    `json:"adjustment_date" binding:"required"`
	Type           string                    `json:"type" binding:"required,oneof=cycle_count write_off write_in"`
	Notes          string                    `json:"notes"`
	Lines          []FGAdjustmentLineRequest `json:"lines" binding:"required,min=1,dive"`
}
