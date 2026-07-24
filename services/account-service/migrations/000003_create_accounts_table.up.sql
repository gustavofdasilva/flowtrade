CREATE TABLE IF NOT EXISTS accounts (
    "id"            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    "user_id"    UUID           NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    "balance"    NUMERIC(20, 8) NOT NULL DEFAULT 0 CHECK (balance >= 0),
    "currency"   CHAR(3)        NOT NULL DEFAULT 'BRL',
    "created_at" TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);