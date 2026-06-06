package dto

// Raw Material

type CreateRawMaterialRequest struct {
	Code        string  `json:"code" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	CategoryID  *int    `json:"category_id"`
	UOMID       int     `json:"uom_id" binding:"required"`
	MinStockQty float64 `json:"min_stock_qty"`
	ReorderQty  float64 `json:"reorder_qty"`
	Description string  `json:"description"`
}

type UpdateRawMaterialRequest struct {
	Name        string  `json:"name" binding:"required"`
	CategoryID  *int    `json:"category_id"`
	UOMID       int     `json:"uom_id" binding:"required"`
	MinStockQty float64 `json:"min_stock_qty"`
	ReorderQty  float64 `json:"reorder_qty"`
	Description string  `json:"description"`
}

// Finished Good

type CreateFinishedGoodRequest struct {
	Code          string  `json:"code" binding:"required"`
	Name          string  `json:"name" binding:"required"`
	UOMID         int     `json:"uom_id" binding:"required"`
	ShelfLifeDays int     `json:"shelf_life_days" binding:"required,min=1"`
	MinStockQty   float64 `json:"min_stock_qty"`
	Description   string  `json:"description"`
}

type UpdateFinishedGoodRequest struct {
	Name          string  `json:"name" binding:"required"`
	UOMID         int     `json:"uom_id" binding:"required"`
	ShelfLifeDays int     `json:"shelf_life_days" binding:"required,min=1"`
	MinStockQty   float64 `json:"min_stock_qty"`
	Description   string  `json:"description"`
}

// BOM

type BOMLineRequest struct {
	RawMaterialID string  `json:"raw_material_id" binding:"required"`
	QtyPerUnit    float64 `json:"qty_per_unit" binding:"required,gt=0"`
	WasteFactor   float64 `json:"waste_factor"`
}

type CreateBOMRequest struct {
	FinishedGoodID string         `json:"finished_good_id" binding:"required"`
	Version        int            `json:"version" binding:"required,min=1"`
	IsActive       bool           `json:"is_active"`
	EffectiveDate  string         `json:"effective_date" binding:"required"`
	Notes          string         `json:"notes"`
	Lines          []BOMLineRequest `json:"lines" binding:"required,min=1,dive"`
}

type UpdateBOMRequest struct {
	Version       int              `json:"version" binding:"required,min=1"`
	IsActive      bool             `json:"is_active"`
	EffectiveDate string           `json:"effective_date" binding:"required"`
	Notes         string           `json:"notes"`
	Lines         []BOMLineRequest `json:"lines" binding:"required,min=1,dive"`
}

// Customer

type CreateCustomerRequest struct {
	Code        string  `json:"code" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	TaxID       string  `json:"tax_id"`
	Address     string  `json:"address"`
	Phone       string  `json:"phone"`
	Email       string  `json:"email"`
	CreditLimit float64 `json:"credit_limit"`
	CreditDays  int     `json:"credit_days"`
}

type UpdateCustomerRequest struct {
	Name        string  `json:"name" binding:"required"`
	TaxID       string  `json:"tax_id"`
	Address     string  `json:"address"`
	Phone       string  `json:"phone"`
	Email       string  `json:"email"`
	CreditLimit float64 `json:"credit_limit"`
	CreditDays  int     `json:"credit_days"`
}

// Supplier

type CreateSupplierRequest struct {
	Code    string `json:"code" binding:"required"`
	Name    string `json:"name" binding:"required"`
	TaxID   string `json:"tax_id"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
}

type UpdateSupplierRequest struct {
	Name    string `json:"name" binding:"required"`
	TaxID   string `json:"tax_id"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
}
