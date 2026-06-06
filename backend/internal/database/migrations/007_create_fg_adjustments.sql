-- Finished Goods Adjustment tables

CREATE TABLE IF NOT EXISTS fg_adjustments (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    adj_number      VARCHAR(50) UNIQUE NOT NULL,
    adjustment_date DATE NOT NULL,
    type            VARCHAR(30) NOT NULL
                    CHECK (type IN ('cycle_count','write_off','write_in')),
    performed_by    UUID REFERENCES users(id),
    approved_by     UUID REFERENCES users(id),
    status          VARCHAR(20) DEFAULT 'pending'
                    CHECK (status IN ('pending','approved','cancelled')),
    notes           TEXT,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS fg_adjustment_lines (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    adjustment_id    UUID NOT NULL REFERENCES fg_adjustments(id) ON DELETE CASCADE,
    fg_stock_lot_id  UUID NOT NULL REFERENCES fg_stock_lots(id),
    finished_good_id UUID NOT NULL REFERENCES finished_goods(id),
    system_qty       NUMERIC(12,3) NOT NULL,
    counted_qty      NUMERIC(12,3) NOT NULL,
    variance_qty     NUMERIC(12,3) GENERATED ALWAYS AS (counted_qty - system_qty) STORED,
    reason           TEXT
);
