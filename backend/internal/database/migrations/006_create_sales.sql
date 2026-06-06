-- Sales & Distribution tables

CREATE TABLE IF NOT EXISTS sales_orders (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_number   VARCHAR(50) UNIQUE NOT NULL,
    customer_id    UUID NOT NULL REFERENCES customers(id),
    order_date     DATE NOT NULL DEFAULT CURRENT_DATE,
    requested_date DATE,
    status         VARCHAR(20) DEFAULT 'draft'
                   CHECK (status IN ('draft','confirmed','picking','dispatched','invoiced','cancelled')),
    subtotal       NUMERIC(14,2) DEFAULT 0,
    vat_amount     NUMERIC(14,2) DEFAULT 0,
    total_amount   NUMERIC(14,2) DEFAULT 0,
    vat_rate       NUMERIC(5,2) DEFAULT 7.00,
    notes          TEXT,
    created_by     UUID REFERENCES users(id),
    created_at     TIMESTAMPTZ DEFAULT NOW(),
    updated_at     TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_sales_orders_status ON sales_orders(status);
CREATE INDEX IF NOT EXISTS idx_sales_orders_customer ON sales_orders(customer_id);
CREATE INDEX IF NOT EXISTS idx_sales_orders_date ON sales_orders(order_date);

CREATE TABLE IF NOT EXISTS sales_order_lines (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sales_order_id   UUID NOT NULL REFERENCES sales_orders(id) ON DELETE CASCADE,
    finished_good_id UUID NOT NULL REFERENCES finished_goods(id),
    ordered_qty      NUMERIC(12,3) NOT NULL CHECK (ordered_qty > 0),
    unit_price       NUMERIC(14,4) NOT NULL,
    discount_pct     NUMERIC(5,2) DEFAULT 0,
    line_total       NUMERIC(14,2),
    status           VARCHAR(20) DEFAULT 'pending'
                     CHECK (status IN ('pending','allocated','dispatched'))
);

CREATE TABLE IF NOT EXISTS vehicles (
    id             SERIAL PRIMARY KEY,
    license_plate  VARCHAR(20) UNIQUE NOT NULL,
    type           VARCHAR(50),
    capacity_units NUMERIC(10,2),
    is_active      BOOLEAN DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS delivery_orders (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    do_number      VARCHAR(50) UNIQUE NOT NULL,
    sales_order_id UUID NOT NULL REFERENCES sales_orders(id),
    delivery_date  DATE,
    driver_user_id UUID REFERENCES users(id),
    vehicle_id     INTEGER REFERENCES vehicles(id),
    status         VARCHAR(20) DEFAULT 'pending'
                   CHECK (status IN ('pending','loading','dispatched','delivered','failed')),
    dispatched_at  TIMESTAMPTZ,
    delivered_at   TIMESTAMPTZ,
    route_notes    TEXT,
    created_at     TIMESTAMPTZ DEFAULT NOW(),
    updated_at     TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS delivery_order_lines (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    delivery_order_id   UUID NOT NULL REFERENCES delivery_orders(id) ON DELETE CASCADE,
    sales_order_line_id UUID REFERENCES sales_order_lines(id),
    fg_stock_lot_id     UUID NOT NULL REFERENCES fg_stock_lots(id),
    finished_good_id    UUID NOT NULL REFERENCES finished_goods(id),
    picked_qty          NUMERIC(12,3) NOT NULL,
    dispatched_qty      NUMERIC(12,3) DEFAULT 0
);

CREATE TABLE IF NOT EXISTS invoices (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_number    VARCHAR(50) UNIQUE NOT NULL,
    sales_order_id    UUID NOT NULL REFERENCES sales_orders(id),
    delivery_order_id UUID REFERENCES delivery_orders(id),
    invoice_date      DATE NOT NULL DEFAULT CURRENT_DATE,
    due_date          DATE,
    subtotal          NUMERIC(14,2),
    vat_amount        NUMERIC(14,2),
    total_amount      NUMERIC(14,2),
    status            VARCHAR(20) DEFAULT 'issued'
                      CHECK (status IN ('issued','paid','overdue','cancelled')),
    payment_date      DATE,
    payment_method    VARCHAR(50),
    created_by        UUID REFERENCES users(id),
    created_at        TIMESTAMPTZ DEFAULT NOW()
);

-- FK from production_orders back to sales_orders (only if not already added)
DO $$ BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM pg_constraint WHERE conname = 'fk_prod_orders_sales_order'
  ) THEN
    ALTER TABLE production_orders
      ADD CONSTRAINT fk_prod_orders_sales_order
      FOREIGN KEY (sales_order_id) REFERENCES sales_orders(id);
  END IF;
END $$;
