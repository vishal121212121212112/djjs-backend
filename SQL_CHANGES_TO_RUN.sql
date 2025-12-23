-- ============================================================
-- SQL Migration for Authentication System
-- Run this directly in your PostgreSQL database
-- Database: postgres (or your database name)
-- ============================================================

-- ============================================
-- STEP 1: Add columns to existing users table
-- ============================================
DO $$ 
BEGIN
    -- Add email_verified_at if it doesn't exist
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name='users' AND column_name='email_verified_at'
    ) THEN
        ALTER TABLE users ADD COLUMN email_verified_at TIMESTAMPTZ NULL;
        RAISE NOTICE 'Added email_verified_at column to users table';
    ELSE
        RAISE NOTICE 'Column email_verified_at already exists in users table';
    END IF;

    -- Add disabled_at if it doesn't exist
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name='users' AND column_name='disabled_at'
    ) THEN
        ALTER TABLE users ADD COLUMN disabled_at TIMESTAMPTZ NULL;
        RAISE NOTICE 'Added disabled_at column to users table';
    ELSE
        RAISE NOTICE 'Column disabled_at already exists in users table';
    END IF;
END $$;

-- ============================================
-- STEP 2: Create sessions table
-- ============================================
CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY, -- UUID as text (matches uuid.New().String() format)
    user_id BIGINT NOT NULL,
    refresh_token_hash BYTEA NOT NULL,
    user_agent TEXT,
    ip TEXT, -- Store as TEXT for flexibility (handles IPv4, IPv6, and proxy headers)
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_used_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at TIMESTAMPTZ NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    
    -- Foreign key - matches users.id (BIGSERIAL/BIGINT)
    CONSTRAINT fk_sessions_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create indexes for sessions
CREATE UNIQUE INDEX IF NOT EXISTS idx_sessions_refresh_token_hash ON sessions(refresh_token_hash);
CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions(expires_at);
CREATE INDEX IF NOT EXISTS idx_sessions_user_id_revoked ON sessions(user_id, revoked_at) WHERE revoked_at IS NULL;

-- ============================================
-- STEP 3: Create verification_tokens table
-- ============================================
CREATE TABLE IF NOT EXISTS verification_tokens (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT, -- UUID as text (database generates)
    user_id BIGINT NOT NULL,
    token_hash BYTEA NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT fk_verification_tokens_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_verification_tokens_token_hash ON verification_tokens(token_hash);
CREATE INDEX IF NOT EXISTS idx_verification_tokens_user_id ON verification_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_verification_tokens_expires_at ON verification_tokens(expires_at);

-- ============================================
-- STEP 4: Create password_reset_tokens table
-- ============================================
CREATE TABLE IF NOT EXISTS password_reset_tokens (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT, -- UUID as text (database generates)
    user_id BIGINT NOT NULL,
    token_hash BYTEA NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT fk_password_reset_tokens_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_password_reset_tokens_token_hash ON password_reset_tokens(token_hash);
CREATE INDEX IF NOT EXISTS idx_password_reset_tokens_user_id ON password_reset_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_password_reset_tokens_expires_at ON password_reset_tokens(expires_at);

-- ============================================
-- STEP 5: Create auth_audit_events table
-- ============================================
CREATE TABLE IF NOT EXISTS auth_audit_events (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT, -- UUID as text (database generates)
    user_id BIGINT NULL,
    type TEXT NOT NULL,
    ip TEXT,
    user_agent TEXT,
    metadata JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT fk_auth_audit_events_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_auth_audit_events_user_id ON auth_audit_events(user_id);
CREATE INDEX IF NOT EXISTS idx_auth_audit_events_created_at ON auth_audit_events(created_at);
CREATE INDEX IF NOT EXISTS idx_auth_audit_events_type ON auth_audit_events(type);
CREATE INDEX IF NOT EXISTS idx_auth_audit_events_user_id_created_at ON auth_audit_events(user_id, created_at);

-- ============================================
-- VERIFICATION (optional - run to check)
-- ============================================
-- SELECT table_name FROM information_schema.tables 
-- WHERE table_schema = 'public' 
-- AND table_name IN ('sessions', 'verification_tokens', 'password_reset_tokens', 'auth_audit_events');

-- ============================================
-- IMPORTANT NOTES:
-- ============================================
-- 1. All id columns use TEXT type (not UUID) to match Go's uuid.New().String() format
-- 2. sessions.ip uses TEXT (not INET) for flexibility with proxy headers
-- 3. All foreign keys reference users(id) which should be BIGINT/BIGSERIAL
-- 4. The migration uses IF NOT EXISTS to be idempotent (safe to run multiple times)
-- 5. If tables already exist with UUID type, you may need to drop and recreate them
