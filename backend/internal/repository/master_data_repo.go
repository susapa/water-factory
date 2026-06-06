package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/water-factory/api/internal/domain"
	"github.com/water-factory/api/internal/dto"
)

type MasterDataRepo struct {
	pool *pgxpool.Pool
}

func NewMasterDataRepo(pool *pgxpool.Pool) *MasterDataRepo {
	return &MasterDataRepo{pool: pool}
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

// ── UOM ──────────────────────────────────────────────────────────────────────

func (r *MasterDataRepo) ListUOM(ctx context.Context) ([]*domain.UnitOfMeasure, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, code, name FROM units_of_measure ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*domain.UnitOfMeasure, 0)
	for rows.Next() {
		u := &domain.UnitOfMeasure{}
		if err := rows.Scan(&u.ID, &u.Code, &u.Name); err != nil {
			return nil, err
		}
		items = append(items, u)
	}
	return items, nil
}

// ── Raw Material Categories ───────────────────────────────────────────────────

func (r *MasterDataRepo) ListCategories(ctx context.Context) ([]*domain.RawMaterialCategory, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, name FROM raw_material_categories ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*domain.RawMaterialCategory, 0)
	for rows.Next() {
		c := &domain.RawMaterialCategory{}
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		items = append(items, c)
	}
	return items, nil
}

// ── Raw Materials ─────────────────────────────────────────────────────────────

const rmSelectCols = `
	rm.id, rm.code, rm.name,
	rm.category_id, COALESCE(c.name, ''),
	rm.uom_id, u.code,
	rm.min_stock_qty::float8, rm.reorder_qty::float8,
	COALESCE(rm.description, ''),
	rm.is_active, rm.created_at, rm.updated_at`

const rmFrom = `
	FROM raw_materials rm
	LEFT JOIN raw_material_categories c ON c.id = rm.category_id
	JOIN units_of_measure u ON u.id = rm.uom_id`

func scanRM(row interface{ Scan(...any) error }) (*domain.RawMaterial, error) {
	rm := &domain.RawMaterial{}
	err := row.Scan(
		&rm.ID, &rm.Code, &rm.Name,
		&rm.CategoryID, &rm.CategoryName,
		&rm.UOMID, &rm.UOMCode,
		&rm.MinStockQty, &rm.ReorderQty,
		&rm.Description,
		&rm.IsActive, &rm.CreatedAt, &rm.UpdatedAt,
	)
	return rm, err
}

func (r *MasterDataRepo) ListRawMaterials(ctx context.Context) ([]*domain.RawMaterial, error) {
	rows, err := r.pool.Query(ctx, `SELECT`+rmSelectCols+rmFrom+` ORDER BY rm.created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*domain.RawMaterial, 0)
	for rows.Next() {
		rm, err := scanRM(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, rm)
	}
	return items, nil
}

func (r *MasterDataRepo) GetRawMaterial(ctx context.Context, id string) (*domain.RawMaterial, error) {
	row := r.pool.QueryRow(ctx, `SELECT`+rmSelectCols+rmFrom+` WHERE rm.id = $1`, id)
	rm, err := scanRM(row)
	if err != nil {
		return nil, fmt.Errorf("get raw material: %w", err)
	}
	return rm, nil
}

func (r *MasterDataRepo) CreateRawMaterial(ctx context.Context, req *dto.CreateRawMaterialRequest) (*domain.RawMaterial, error) {
	row := r.pool.QueryRow(ctx,
		`INSERT INTO raw_materials (code, name, category_id, uom_id, min_stock_qty, reorder_qty, description)
		 VALUES ($1,$2,$3,$4,$5,$6,$7)
		 RETURNING id`,
		req.Code, req.Name, req.CategoryID, req.UOMID, req.MinStockQty, req.ReorderQty, req.Description,
	)
	var id string
	if err := row.Scan(&id); err != nil {
		if isUniqueViolation(err) {
			return nil, fmt.Errorf("code '%s' already exists", req.Code)
		}
		return nil, fmt.Errorf("create raw material: %w", err)
	}
	return r.GetRawMaterial(ctx, id)
}

func (r *MasterDataRepo) UpdateRawMaterial(ctx context.Context, id string, req *dto.UpdateRawMaterialRequest) (*domain.RawMaterial, error) {
	_, err := r.pool.Exec(ctx,
		`UPDATE raw_materials
		 SET name=$1, category_id=$2, uom_id=$3, min_stock_qty=$4, reorder_qty=$5, description=$6, updated_at=NOW()
		 WHERE id=$7`,
		req.Name, req.CategoryID, req.UOMID, req.MinStockQty, req.ReorderQty, req.Description, id,
	)
	if err != nil {
		return nil, fmt.Errorf("update raw material: %w", err)
	}
	return r.GetRawMaterial(ctx, id)
}

func (r *MasterDataRepo) SetRawMaterialActive(ctx context.Context, id string, active bool) error {
	_, err := r.pool.Exec(ctx, `UPDATE raw_materials SET is_active=$1, updated_at=NOW() WHERE id=$2`, active, id)
	return err
}

// ── Finished Goods ────────────────────────────────────────────────────────────

const fgSelectCols = `
	fg.id, fg.code, fg.name,
	fg.uom_id, u.code,
	fg.shelf_life_days,
	fg.min_stock_qty::float8,
	COALESCE(fg.description, ''),
	fg.is_active, fg.created_at, fg.updated_at`

const fgFrom = `
	FROM finished_goods fg
	JOIN units_of_measure u ON u.id = fg.uom_id`

func scanFG(row interface{ Scan(...any) error }) (*domain.FinishedGood, error) {
	fg := &domain.FinishedGood{}
	err := row.Scan(
		&fg.ID, &fg.Code, &fg.Name,
		&fg.UOMID, &fg.UOMCode,
		&fg.ShelfLifeDays,
		&fg.MinStockQty,
		&fg.Description,
		&fg.IsActive, &fg.CreatedAt, &fg.UpdatedAt,
	)
	return fg, err
}

func (r *MasterDataRepo) ListFinishedGoods(ctx context.Context) ([]*domain.FinishedGood, error) {
	rows, err := r.pool.Query(ctx, `SELECT`+fgSelectCols+fgFrom+` ORDER BY fg.created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*domain.FinishedGood, 0)
	for rows.Next() {
		fg, err := scanFG(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, fg)
	}
	return items, nil
}

func (r *MasterDataRepo) GetFinishedGood(ctx context.Context, id string) (*domain.FinishedGood, error) {
	row := r.pool.QueryRow(ctx, `SELECT`+fgSelectCols+fgFrom+` WHERE fg.id = $1`, id)
	fg, err := scanFG(row)
	if err != nil {
		return nil, fmt.Errorf("get finished good: %w", err)
	}
	return fg, nil
}

func (r *MasterDataRepo) CreateFinishedGood(ctx context.Context, req *dto.CreateFinishedGoodRequest) (*domain.FinishedGood, error) {
	var id string
	err := r.pool.QueryRow(ctx,
		`INSERT INTO finished_goods (code, name, uom_id, shelf_life_days, min_stock_qty, description)
		 VALUES ($1,$2,$3,$4,$5,$6)
		 RETURNING id`,
		req.Code, req.Name, req.UOMID, req.ShelfLifeDays, req.MinStockQty, req.Description,
	).Scan(&id)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, fmt.Errorf("code '%s' already exists", req.Code)
		}
		return nil, fmt.Errorf("create finished good: %w", err)
	}
	return r.GetFinishedGood(ctx, id)
}

func (r *MasterDataRepo) UpdateFinishedGood(ctx context.Context, id string, req *dto.UpdateFinishedGoodRequest) (*domain.FinishedGood, error) {
	_, err := r.pool.Exec(ctx,
		`UPDATE finished_goods
		 SET name=$1, uom_id=$2, shelf_life_days=$3, min_stock_qty=$4, description=$5, updated_at=NOW()
		 WHERE id=$6`,
		req.Name, req.UOMID, req.ShelfLifeDays, req.MinStockQty, req.Description, id,
	)
	if err != nil {
		return nil, fmt.Errorf("update finished good: %w", err)
	}
	return r.GetFinishedGood(ctx, id)
}

func (r *MasterDataRepo) SetFinishedGoodActive(ctx context.Context, id string, active bool) error {
	_, err := r.pool.Exec(ctx, `UPDATE finished_goods SET is_active=$1, updated_at=NOW() WHERE id=$2`, active, id)
	return err
}

// ── BOM ───────────────────────────────────────────────────────────────────────

func (r *MasterDataRepo) ListBOMs(ctx context.Context) ([]*domain.BOM, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT b.id, b.finished_good_id, fg.code, fg.name,
		       b.version, b.is_active, b.effective_date::text, COALESCE(b.notes,''), b.created_at
		FROM bill_of_materials b
		JOIN finished_goods fg ON fg.id = b.finished_good_id
		ORDER BY b.created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*domain.BOM, 0)
	for rows.Next() {
		b := &domain.BOM{Lines: []domain.BOMLine{}}
		if err := rows.Scan(&b.ID, &b.FinishedGoodID, &b.FGCode, &b.FGName,
			&b.Version, &b.IsActive, &b.EffectiveDate, &b.Notes, &b.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, b)
	}
	return items, nil
}

func (r *MasterDataRepo) GetBOM(ctx context.Context, id string) (*domain.BOM, error) {
	b := &domain.BOM{Lines: []domain.BOMLine{}}
	err := r.pool.QueryRow(ctx, `
		SELECT b.id, b.finished_good_id, fg.code, fg.name,
		       b.version, b.is_active, b.effective_date::text, COALESCE(b.notes,''), b.created_at
		FROM bill_of_materials b
		JOIN finished_goods fg ON fg.id = b.finished_good_id
		WHERE b.id = $1`, id,
	).Scan(&b.ID, &b.FinishedGoodID, &b.FGCode, &b.FGName,
		&b.Version, &b.IsActive, &b.EffectiveDate, &b.Notes, &b.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("get bom: %w", err)
	}

	lines, err := r.listBOMLines(ctx, id)
	if err != nil {
		return nil, err
	}
	b.Lines = lines
	return b, nil
}

func (r *MasterDataRepo) listBOMLines(ctx context.Context, bomID string) ([]domain.BOMLine, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT bl.id, bl.bom_id, bl.raw_material_id, rm.code, rm.name,
		       bl.qty_per_unit::float8, bl.waste_factor::float8
		FROM bom_lines bl
		JOIN raw_materials rm ON rm.id = bl.raw_material_id
		WHERE bl.bom_id = $1`, bomID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lines []domain.BOMLine
	for rows.Next() {
		l := domain.BOMLine{}
		if err := rows.Scan(&l.ID, &l.BOMID, &l.RawMaterialID, &l.RMCode, &l.RMName,
			&l.QtyPerUnit, &l.WasteFactor); err != nil {
			return nil, err
		}
		lines = append(lines, l)
	}
	if lines == nil {
		lines = []domain.BOMLine{}
	}
	return lines, nil
}

func (r *MasterDataRepo) CreateBOM(ctx context.Context, req *dto.CreateBOMRequest) (*domain.BOM, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var bomID string
	err = tx.QueryRow(ctx,
		`INSERT INTO bill_of_materials (finished_good_id, version, is_active, effective_date, notes)
		 VALUES ($1,$2,$3,$4,$5)
		 RETURNING id`,
		req.FinishedGoodID, req.Version, req.IsActive, req.EffectiveDate, req.Notes,
	).Scan(&bomID)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, fmt.Errorf("BOM version %d already exists for this finished good", req.Version)
		}
		return nil, fmt.Errorf("create bom: %w", err)
	}

	for _, l := range req.Lines {
		_, err = tx.Exec(ctx,
			`INSERT INTO bom_lines (bom_id, raw_material_id, qty_per_unit, waste_factor)
			 VALUES ($1,$2,$3,$4)`,
			bomID, l.RawMaterialID, l.QtyPerUnit, l.WasteFactor,
		)
		if err != nil {
			return nil, fmt.Errorf("insert bom line: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}
	return r.GetBOM(ctx, bomID)
}

func (r *MasterDataRepo) UpdateBOM(ctx context.Context, id string, req *dto.UpdateBOMRequest) (*domain.BOM, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`UPDATE bill_of_materials
		 SET version=$1, is_active=$2, effective_date=$3, notes=$4
		 WHERE id=$5`,
		req.Version, req.IsActive, req.EffectiveDate, req.Notes, id,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, fmt.Errorf("BOM version %d already exists for this finished good", req.Version)
		}
		return nil, fmt.Errorf("update bom: %w", err)
	}

	if _, err = tx.Exec(ctx, `DELETE FROM bom_lines WHERE bom_id=$1`, id); err != nil {
		return nil, fmt.Errorf("delete bom lines: %w", err)
	}

	for _, l := range req.Lines {
		_, err = tx.Exec(ctx,
			`INSERT INTO bom_lines (bom_id, raw_material_id, qty_per_unit, waste_factor)
			 VALUES ($1,$2,$3,$4)`,
			id, l.RawMaterialID, l.QtyPerUnit, l.WasteFactor,
		)
		if err != nil {
			return nil, fmt.Errorf("insert bom line: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}
	return r.GetBOM(ctx, id)
}

// ── Customers ─────────────────────────────────────────────────────────────────

const custSelectCols = `
	id, code, name,
	COALESCE(tax_id,''), COALESCE(address,''), COALESCE(phone,''), COALESCE(email,''),
	credit_limit::float8, credit_days,
	is_active, created_at, updated_at`

func scanCustomer(row interface{ Scan(...any) error }) (*domain.Customer, error) {
	c := &domain.Customer{}
	err := row.Scan(
		&c.ID, &c.Code, &c.Name,
		&c.TaxID, &c.Address, &c.Phone, &c.Email,
		&c.CreditLimit, &c.CreditDays,
		&c.IsActive, &c.CreatedAt, &c.UpdatedAt,
	)
	return c, err
}

func (r *MasterDataRepo) ListCustomers(ctx context.Context) ([]*domain.Customer, error) {
	rows, err := r.pool.Query(ctx, `SELECT `+custSelectCols+` FROM customers ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*domain.Customer, 0)
	for rows.Next() {
		c, err := scanCustomer(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, c)
	}
	return items, nil
}

func (r *MasterDataRepo) GetCustomer(ctx context.Context, id string) (*domain.Customer, error) {
	row := r.pool.QueryRow(ctx, `SELECT `+custSelectCols+` FROM customers WHERE id=$1`, id)
	c, err := scanCustomer(row)
	if err != nil {
		return nil, fmt.Errorf("get customer: %w", err)
	}
	return c, nil
}

func (r *MasterDataRepo) CreateCustomer(ctx context.Context, req *dto.CreateCustomerRequest) (*domain.Customer, error) {
	var id string
	err := r.pool.QueryRow(ctx,
		`INSERT INTO customers (code, name, tax_id, address, phone, email, credit_limit, credit_days)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		 RETURNING id`,
		req.Code, req.Name, req.TaxID, req.Address, req.Phone, req.Email, req.CreditLimit, req.CreditDays,
	).Scan(&id)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, fmt.Errorf("code '%s' already exists", req.Code)
		}
		return nil, fmt.Errorf("create customer: %w", err)
	}
	return r.GetCustomer(ctx, id)
}

func (r *MasterDataRepo) UpdateCustomer(ctx context.Context, id string, req *dto.UpdateCustomerRequest) (*domain.Customer, error) {
	_, err := r.pool.Exec(ctx,
		`UPDATE customers
		 SET name=$1, tax_id=$2, address=$3, phone=$4, email=$5, credit_limit=$6, credit_days=$7, updated_at=NOW()
		 WHERE id=$8`,
		req.Name, req.TaxID, req.Address, req.Phone, req.Email, req.CreditLimit, req.CreditDays, id,
	)
	if err != nil {
		return nil, fmt.Errorf("update customer: %w", err)
	}
	return r.GetCustomer(ctx, id)
}

func (r *MasterDataRepo) SetCustomerActive(ctx context.Context, id string, active bool) error {
	_, err := r.pool.Exec(ctx, `UPDATE customers SET is_active=$1, updated_at=NOW() WHERE id=$2`, active, id)
	return err
}

// ── Suppliers ─────────────────────────────────────────────────────────────────

const suppSelectCols = `
	id, code, name,
	COALESCE(tax_id,''), COALESCE(address,''), COALESCE(phone,''), COALESCE(email,''),
	is_active, created_at, updated_at`

func scanSupplier(row interface{ Scan(...any) error }) (*domain.Supplier, error) {
	s := &domain.Supplier{}
	err := row.Scan(
		&s.ID, &s.Code, &s.Name,
		&s.TaxID, &s.Address, &s.Phone, &s.Email,
		&s.IsActive, &s.CreatedAt, &s.UpdatedAt,
	)
	return s, err
}

func (r *MasterDataRepo) ListSuppliers(ctx context.Context) ([]*domain.Supplier, error) {
	rows, err := r.pool.Query(ctx, `SELECT `+suppSelectCols+` FROM suppliers ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*domain.Supplier, 0)
	for rows.Next() {
		s, err := scanSupplier(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, s)
	}
	return items, nil
}

func (r *MasterDataRepo) GetSupplier(ctx context.Context, id string) (*domain.Supplier, error) {
	row := r.pool.QueryRow(ctx, `SELECT `+suppSelectCols+` FROM suppliers WHERE id=$1`, id)
	s, err := scanSupplier(row)
	if err != nil {
		return nil, fmt.Errorf("get supplier: %w", err)
	}
	return s, nil
}

func (r *MasterDataRepo) CreateSupplier(ctx context.Context, req *dto.CreateSupplierRequest) (*domain.Supplier, error) {
	var id string
	err := r.pool.QueryRow(ctx,
		`INSERT INTO suppliers (code, name, tax_id, address, phone, email)
		 VALUES ($1,$2,$3,$4,$5,$6)
		 RETURNING id`,
		req.Code, req.Name, req.TaxID, req.Address, req.Phone, req.Email,
	).Scan(&id)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, fmt.Errorf("code '%s' already exists", req.Code)
		}
		return nil, fmt.Errorf("create supplier: %w", err)
	}
	return r.GetSupplier(ctx, id)
}

func (r *MasterDataRepo) UpdateSupplier(ctx context.Context, id string, req *dto.UpdateSupplierRequest) (*domain.Supplier, error) {
	_, err := r.pool.Exec(ctx,
		`UPDATE suppliers
		 SET name=$1, tax_id=$2, address=$3, phone=$4, email=$5, updated_at=NOW()
		 WHERE id=$6`,
		req.Name, req.TaxID, req.Address, req.Phone, req.Email, id,
	)
	if err != nil {
		return nil, fmt.Errorf("update supplier: %w", err)
	}
	return r.GetSupplier(ctx, id)
}

func (r *MasterDataRepo) SetSupplierActive(ctx context.Context, id string, active bool) error {
	_, err := r.pool.Exec(ctx, `UPDATE suppliers SET is_active=$1, updated_at=NOW() WHERE id=$2`, active, id)
	return err
}
