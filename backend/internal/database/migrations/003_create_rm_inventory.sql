-- Raw Material Inventory tables

CREATE TABLE IF NOT EXISTS rm_goods_receipts (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    grn_number    VARCHAR(50) UNIQUE NOT NULL,
    supplier_id   UUID REFERENCES suppliers(id),
    received_by   UUID REFERENCES users(id),
    received_date DATE NOT NULL DEFAULT CURRENT_DATE,
    po_reference  VARCHAR(100),
    status        VARCHAR(20) DEFAULT 'draft'
                  CHECK (status IN ('draft','confirmed','cancelled')),
    notes         TEXT,
    created_at    TIMESTAMPTZ DEFAULT NOW(),
    updated_at    TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS rm_goods_receipt_lines (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    grn_id          UUID NOT NULL REFERENCES rm_goods_receipts(id) ON DELETE CASCADE,
    raw_material_id UUID NOT NULL REFERENCES raw_materials(id),
    lot_number      VARCHAR(100) NOT NULL,
    received_qty    NUMERIC(12,3) NOT NULL CHECK (received_qty > 0),
    unit_cost       NUMERIC(14,4),
    expiry_date     DATE,
    location_id     INTEGER REFERENCES warehouse_locations(id),
    notes           TEXT
);

CREATE TABLE IF NOT EXISTS rm_stock_lots (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    raw_material_id UUID NOT NULL REFERENCES raw_materials(id),
    lot_number      VARCHAR(100) NOT NULL,
    grn_line_id     UUID REFERENCES rm_goods_receipt_lines(id),
    received_date   DATE NOT NULL,
    expiry_date     DATE,
    location_id     INTEGER REFERENCES warehouse_locations(id),
    initial_qty     NUMERIC(12,3) NOT NULL,
    current_qty     NUMERIC(12,3) NOT NULL,
    unit_cost       NUMERIC(14,4),
    status          VARCHAR(20) DEFAULT 'available'
                    CHECK (status IN ('available','reserved','depleted','quarantine')),
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(raw_material_id, lot_number)
);

CREATE INDEX IF NOT EXISTS idx_rm_stock_lots_rm_id ON rm_stock_lots(raw_material_id);
CREATE INDEX IF NOT EXISTS idx_rm_stock_lots_status ON rm_stock_lots(status);
CREATE INDEX IF NOT EXISTS idx_rm_stock_lots_received_date ON rm_stock_lots(received_date);

CREATE TABLE IF NOT EXISTS rm_stock_movements (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rm_stock_lot_id UUID NOT NULL REFERENCES rm_stock_lots(id),
    raw_material_id UUID NOT NULL REFERENCES raw_materials(id),
    movement_type   VARCHAR(30) NOT NULL
                    CHECK (movement_type IN ('GRN','ISSUE_TO_PROD','RETURN_FROM_PROD','ADJUSTMENT_IN','ADJUSTMENT_OUT','CYCLE_COUNT')),
    reference_type  VARCHAR(30),
    reference_id    UUID,
    qty             NUMERIC(12,3) NOT NULL,
    qty_before      NUMERIC(12,3) NOT NULL,
    qty_after       NUMERIC(12,3) NOT NULL,
    performed_by    UUID REFERENCES users(id),
    notes           TEXT,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_rm_movements_lot ON rm_stock_movements(rm_stock_lot_id);
CREATE INDEX IF NOT EXISTS idx_rm_movements_created ON rm_stock_movements(created_at);

CREATE TABLE IF NOT EXISTS rm_adjustments (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    adj_number      VARCHAR(50) UNIQUE NOT NULL,
    adjustment_date DATE NOT NULL DEFAULT CURRENT_DATE,
    type            VARCHAR(20) NOT NULL CHECK (type IN ('cycle_count','damage','correction')),
    performed_by    UUID REFERENCES users(id),
    approved_by     UUID REFERENCES users(id),
    status          VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending','approved','cancelled')),
    notes           TEXT,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS rm_adjustment_lines (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    adjustment_id   UUID NOT NULL REFERENCES rm_adjustments(id) ON DELETE CASCADE,
    rm_stock_lot_id UUID NOT NULL REFERENCES rm_stock_lots(id),
    raw_material_id UUID NOT NULL REFERENCES raw_materials(id),
    system_qty      NUMERIC(12,3) NOT NULL,
    counted_qty     NUMERIC(12,3) NOT NULL,
    variance_qty    NUMERIC(12,3) GENERATED ALWAYS AS (counted_qty - system_qty) STORED,
    reason          TEXT
);
