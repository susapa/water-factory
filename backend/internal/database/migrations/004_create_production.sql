-- Production tables

CREATE TABLE IF NOT EXISTS production_orders (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_number        VARCHAR(50) UNIQUE NOT NULL,
    finished_good_id    UUID NOT NULL REFERENCES finished_goods(id),
    bom_id              UUID REFERENCES bill_of_materials(id),
    planned_qty         NUMERIC(12,3) NOT NULL CHECK (planned_qty > 0),
    actual_yield_qty    NUMERIC(12,3) DEFAULT 0,
    defect_qty          NUMERIC(12,3) DEFAULT 0,
    waste_qty           NUMERIC(12,3) DEFAULT 0,
    planned_start_date  DATE,
    planned_end_date    DATE,
    actual_start_date   TIMESTAMPTZ,
    actual_end_date     TIMESTAMPTZ,
    status              VARCHAR(20) DEFAULT 'draft'
                        CHECK (status IN ('draft','confirmed','rm_issued','in_progress','completed','cancelled')),
    created_by          UUID REFERENCES users(id),
    sales_order_id      UUID,
    notes               TEXT,
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    updated_at          TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_prod_orders_status ON production_orders(status);
CREATE INDEX IF NOT EXISTS idx_prod_orders_fg_id ON production_orders(finished_good_id);

-- BOM requirements snapshot (locked at confirm time)
CREATE TABLE IF NOT EXISTS production_order_rm_requirements (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    production_order_id UUID NOT NULL REFERENCES production_orders(id) ON DELETE CASCADE,
    raw_material_id     UUID NOT NULL REFERENCES raw_materials(id),
    required_qty        NUMERIC(12,3) NOT NULL,
    issued_qty          NUMERIC(12,3) DEFAULT 0,
    status              VARCHAR(20) DEFAULT 'pending'
                        CHECK (status IN ('pending','partial','fully_issued'))
);

-- Records which specific Lots were issued (FIFO)
CREATE TABLE IF NOT EXISTS production_rm_issues (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    production_order_id UUID NOT NULL REFERENCES production_orders(id),
    rm_stock_lot_id     UUID NOT NULL REFERENCES rm_stock_lots(id),
    raw_material_id     UUID NOT NULL REFERENCES raw_materials(id),
    issued_qty          NUMERIC(12,3) NOT NULL,
    issued_by           UUID REFERENCES users(id),
    issued_at           TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS production_yields (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    production_order_id UUID NOT NULL REFERENCES production_orders(id),
    yield_qty           NUMERIC(12,3) NOT NULL,
    defect_qty          NUMERIC(12,3) DEFAULT 0,
    waste_qty           NUMERIC(12,3) DEFAULT 0,
    batch_number        VARCHAR(100) NOT NULL,
    production_date     DATE NOT NULL DEFAULT CURRENT_DATE,
    expiry_date         DATE NOT NULL,
    recorded_by         UUID REFERENCES users(id),
    recorded_at         TIMESTAMPTZ DEFAULT NOW(),
    notes               TEXT
);

CREATE TABLE IF NOT EXISTS production_defect_details (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    yield_id    UUID NOT NULL REFERENCES production_yields(id) ON DELETE CASCADE,
    defect_type VARCHAR(100) NOT NULL,
    qty         NUMERIC(12,3) NOT NULL,
    description TEXT
);
