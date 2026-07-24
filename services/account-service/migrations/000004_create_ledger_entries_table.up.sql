DO $$
BEGIN
    CREATE TYPE ledger_entry_type AS ENUM (
        'DEPOSIT',
        'WITHDRAWAL',
        'TRADE_DEBIT',
        'TRADE_CREDIT'
    );
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

CREATE TABLE IF NOT EXISTS ledger_entries (
    id          UUID              PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id  UUID              NOT NULL REFERENCES accounts(id),
    type        ledger_entry_type NOT NULL,
    amount      NUMERIC(20, 8)    NOT NULL CHECK (amount > 0),
    balance_after NUMERIC(20, 8)  NOT NULL,
    description VARCHAR(255),
    reference_id UUID,            
    created_at  TIMESTAMPTZ       NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ledger_entries_account_id ON ledger_entries(account_id);
CREATE INDEX IF NOT EXISTS idx_ledger_entries_created_at ON ledger_entries(account_id, created_at DESC);
