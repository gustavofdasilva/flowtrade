CREATE TABLE IF NOT EXISTS asset_prices (
    "id"         UUID           PRIMARY KEY DEFAULT gen_random_uuid(),
    "asset_id"   UUID           NOT NULL REFERENCES assets(id),
    "price"      NUMERIC(20, 8) NOT NULL CHECK (price > 0),
    "source"     VARCHAR(50)    NOT NULL DEFAULT 'simulated',
    "created_at" TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_asset_prices_asset_id_created_at ON asset_prices(asset_id, created_at DESC);