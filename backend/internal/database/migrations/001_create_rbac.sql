-- RBAC tables

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS roles (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    created_at  TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email         VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name     VARCHAR(100) NOT NULL,
    role_id       INTEGER NOT NULL REFERENCES roles(id),
    is_active     BOOLEAN DEFAULT TRUE,
    created_at    TIMESTAMPTZ DEFAULT NOW(),
    updated_at    TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_role_id ON users(role_id);

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash  VARCHAR(255) NOT NULL,
    expires_at  TIMESTAMPTZ NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT NOW()
);

-- Sequence table for document number generation
CREATE TABLE IF NOT EXISTS sequences (
    prefix      VARCHAR(20) PRIMARY KEY,
    last_value  BIGINT NOT NULL DEFAULT 0
);

-- Seed roles
INSERT INTO roles (name, description) VALUES
    ('admin',               'Full system access + user management'),
    ('warehouse_manager',   'Raw material and finished goods warehouse'),
    ('production_manager',  'Production orders and yield recording'),
    ('sales_admin',         'Sales orders, invoicing, and dispatch'),
    ('delivery',            'Delivery dispatch confirmation only')
ON CONFLICT (name) DO NOTHING;

-- Seed admin user: password = Admin@1234
INSERT INTO users (email, password_hash, full_name, role_id)
SELECT 'admin@water.local',
       '$2a$10$jdjXqS2z7LicSpSnwyyH3.KPpT33F0fb4f4csaQSqxXFy0Nh4x9U2',
       'System Admin',
       id
FROM roles WHERE name = 'admin'
ON CONFLICT (email) DO NOTHING;
