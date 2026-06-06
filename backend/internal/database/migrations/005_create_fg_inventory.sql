-- Finished Goods Inventory tables

CREATE TABLE IF NOT EXISTS fg_stock_lots (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    finished_good_id    UUID NOT NULL REFERENCES finished_goods(id),
    batch_number        VARCHAR(100) NOT NULL,
    production_order_id UUID REFERENCES production_orders(id),
    production_date     DATE NOT NULL,
    expiry_date         DATE NOT NULL,
    location_id         INTEGER REFERENCES warehouse_locations(id),
    initial_qty         NUMERIC(12,3) NOT NULL,
    current_qty         NUMERIC(12,3) NOT NULL,
    status              VARCHAR(20) DEFAULT 'available'
                        CHECK (status IN ('available','reserved','dispatched','expired','quarantine')),
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    updated_at          TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(finished_good_id, batch_number)
);

CREATE INDEX IF NOT EXISTS idx_fg_lots_fg_id ON fg_stock_lots(finished_good_id);
CREATE INDEX IF NOT EXISTS idx_fg_lots_expiry ON fg_stock_lots(expiry_date);
CREATE INDEX IF NOT EXISTS idx_fg_lots_status ON fg_stock_lots(status);

CREATE TABLE IF NOT EXISTS fg_stock_movements (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fg_stock_lot_id  UUID NOT NULL REFERENCES fg_stock_lots(id),
    finished_good_id UUID NOT NULL REFERENCES finished_goods(id),
    movement_type    VARCHAR(30) NOT NULL
                     CHECK (movement_type IN ('PRODUCTION_RECEIPT','SALES_DISPATCH','RETURN','ADJUSTMENT_IN','ADJUSTMENT_OUT')),
    reference_type   VARCHAR(30),
    reference_id     UUID,
    qty              NUMERIC(12,3) NOT NULL,
    qty_before       NUMERIC(12,3) NOT NULL,
    qty_after        NUMERIC(12,3) NOT NULL,
    performed_by     UUID REFERENCES users(id),
    notes            TEXT,
    created_at       TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_fg_movements_lot ON fg_stock_movements(fg_stock_lot_id);
CREATE INDEX IF NOT EXISTS idx_fg_movements_created ON fg_stock_movements(created_at);

CREATE TABLE IF NOT EXISTS notifications (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type          VARCHAR(50) NOT NULL,
    reference_id  UUID,
    message       TEXT NOT NULL,
    target_roles  TEXT[] NOT NULL,
    is_read       BOOLEAN DEFAULT FALSE,
    created_at    TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_notifications_unread ON notifications(is_read, created_at);
