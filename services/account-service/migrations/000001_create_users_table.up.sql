DO $$
BEGIN
    CREATE TYPE user_role AS ENUM (
        'ADMIN',
        'CONSUMER'
    );
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

CREATE TABLE IF NOT EXISTS public.users (
    "id"            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    "username"      TEXT        NOT NULL,
    "email"         TEXT        UNIQUE NOT NULL,
    "password_hash" TEXT        NOT NULL,
    "role"          user_role   NOT NULL DEFAULT 'CONSUMER',
    "created_at"    TIMESTAMP   DEFAULT NOW(),
    "updated_at"    TIMESTAMP   DEFAULT NOW(),
    "deleted_at"    TIMESTAMP   NULL
);

