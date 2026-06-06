package domain

import "time"

type UnitOfMeasure struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type RawMaterialCategory struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type RawMaterial struct {
	ID           string    `json:"id"`
	Code         string    `json:"code"`
	Name         string    `json:"name"`
	CategoryID   *int      `json:"category_id"`
	CategoryName string    `json:"category_name"`
	UOMID        int       `json:"uom_id"`
	UOMCode      string    `json:"uom_code"`
	MinStockQty  float64   `json:"min_stock_qty"`
	ReorderQty   float64   `json:"reorder_qty"`
	Description  string    `json:"description"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type FinishedGood struct {
	ID            string    `json:"id"`
	Code          string    `json:"code"`
	Name          string    `json:"name"`
	UOMID         int       `json:"uom_id"`
	UOMCode       string    `json:"uom_code"`
	ShelfLifeDays int       `json:"shelf_life_days"`
	MinStockQty   float64   `json:"min_stock_qty"`
	Description   string    `json:"description"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type BOM struct {
	ID             string    `json:"id"`
	FinishedGoodID string    `json:"finished_good_id"`
	FGCode         string    `json:"fg_code"`
	FGName         string    `json:"fg_name"`
	Version        int       `json:"version"`
	IsActive       bool      `json:"is_active"`
	EffectiveDate  string    `json:"effective_date"`
	Notes          string    `json:"notes"`
	Lines          []BOMLine `json:"lines"`
	CreatedAt      time.Time `json:"created_at"`
}

type BOMLine struct {
	ID            string  `json:"id"`
	BOMID         string  `json:"bom_id"`
	RawMaterialID string  `json:"raw_material_id"`
	RMCode        string  `json:"rm_code"`
	RMName        string  `json:"rm_name"`
	QtyPerUnit    float64 `json:"qty_per_unit"`
	WasteFactor   float64 `json:"waste_factor"`
}

type Customer struct {
	ID          string    `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	TaxID       string    `json:"tax_id"`
	Address     string    `json:"address"`
	Phone       string    `json:"phone"`
	Email       string    `json:"email"`
	CreditLimit float64   `json:"credit_limit"`
	CreditDays  int       `json:"credit_days"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Supplier struct {
	ID        string    `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	TaxID     string    `json:"tax_id"`
	Address   string    `json:"address"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
