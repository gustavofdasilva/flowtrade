DO $$
BEGIN
    CREATE TYPE asset_type AS ENUM (
        'STOCK',
        'ETF',
        'CRYPTO'
    );
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

CREATE TABLE IF NOT EXISTS assets (
    "id"          UUID       PRIMARY KEY DEFAULT gen_random_uuid(),
    "ticker"      VARCHAR(20)  NOT NULL UNIQUE,
    "name"        VARCHAR(100) NOT NULL,
    "type"        asset_type   NOT NULL,
    "currency"    CHAR(3)      NOT NULL DEFAULT 'BRL',
    "is_active"   BOOLEAN      NOT NULL DEFAULT TRUE,
    "created_at"  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    "updated_at"  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_assets_ticker   ON assets(ticker);
CREATE INDEX idx_assets_is_active ON assets(is_active);