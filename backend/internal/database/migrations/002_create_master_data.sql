-- Master Data tables

CREATE TABLE IF NOT EXISTS units_of_measure (
    id    SERIAL PRIMARY KEY,
    code  VARCHAR(20) UNIQUE NOT NULL,
    name  VARCHAR(50) NOT NULL
);

INSERT INTO units_of_measure (code, name) VALUES
    ('PCS',   'ชิ้น'),
    ('BOX',   'กล่อง'),
    ('PACK',  'แพ็ค'),
    ('LITER', 'ลิตร'),
    ('KG',    'กิโลกรัม'),
    ('ROLL',  'ม้วน'),
    ('SHEET', 'แผ่น')
ON CONFLICT (code) DO NOTHING;

CREATE TABLE IF NOT EXISTS raw_material_categories (
    id   SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

INSERT INTO raw_material_categories (name) VALUES
    ('บรรจุภัณฑ์'),
    ('สารกรอง'),
    ('สารเคมี'),
    ('อุปกรณ์')
ON CONFLICT DO NOTHING;

CREATE TABLE IF NOT EXISTS raw_materials (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code            VARCHAR(50) UNIQUE NOT NULL,
    name            VARCHAR(200) NOT NULL,
    category_id     INTEGER REFERENCES raw_material_categories(id),
    uom_id          INTEGER NOT NULL REFERENCES units_of_measure(id),
    min_stock_qty   NUMERIC(12,3) DEFAULT 0,
    reorder_qty     NUMERIC(12,3) DEFAULT 0,
    description     TEXT,
    is_active       BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS finished_goods (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code            VARCHAR(50) UNIQUE NOT NULL,
    name            VARCHAR(200) NOT NULL,
    uom_id          INTEGER NOT NULL REFERENCES units_of_measure(id),
    shelf_life_days INTEGER NOT NULL DEFAULT 365,
    min_stock_qty   NUMERIC(12,3) DEFAULT 0,
    description     TEXT,
    is_active       BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS bill_of_materials (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    finished_good_id UUID NOT NULL REFERENCES finished_goods(id) ON DELETE CASCADE,
    version          INTEGER NOT NULL DEFAULT 1,
    is_active        BOOLEAN DEFAULT TRUE,
    effective_date   DATE NOT NULL,
    notes            TEXT,
    created_at       TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(finished_good_id, version)
);

CREATE TABLE IF NOT EXISTS bom_lines (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bom_id          UUID NOT NULL REFERENCES bill_of_materials(id) ON DELETE CASCADE,
    raw_material_id UUID NOT NULL REFERENCES raw_materials(id),
    qty_per_unit    NUMERIC(12,4) NOT NULL,
    waste_factor    NUMERIC(5,4) DEFAULT 0
);

CREATE TABLE IF NOT EXISTS customers (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code         VARCHAR(50) UNIQUE NOT NULL,
    name         VARCHAR(200) NOT NULL,
    tax_id       VARCHAR(20),
    address      TEXT,
    phone        VARCHAR(30),
    email        VARCHAR(100),
    credit_limit NUMERIC(14,2) DEFAULT 0,
    credit_days  INTEGER DEFAULT 30,
    is_active    BOOLEAN DEFAULT TRUE,
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS suppliers (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code        VARCHAR(50) UNIQUE NOT NULL,
    name        VARCHAR(200) NOT NULL,
    tax_id      VARCHAR(20),
    address     TEXT,
    phone       VARCHAR(30),
    email       VARCHAR(100),
    is_active   BOOLEAN DEFAULT TRUE,
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    updated_at  TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS warehouse_locations (
    id          SERIAL PRIMARY KEY,
    zone        VARCHAR(20) NOT NULL,
    row_no      VARCHAR(10),
    bay_no      VARCHAR(10),
    description TEXT
);
