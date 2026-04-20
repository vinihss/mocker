-- Migrations for Mocker API
-- This file is executed automatically by storage/sqlite.go

-- migrations/001_create_tables.sql
CREATE TABLE IF NOT EXISTS mocks (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    path TEXT NOT NULL,
    method TEXT NOT NULL,
    description TEXT,
    is_active INTEGER DEFAULT 1,
    response_type TEXT NOT NULL,
    response_status_code INTEGER DEFAULT 200,
    response_headers TEXT,
    response_body TEXT NOT NULL,
    response_delay_ms INTEGER DEFAULT 0,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS instructions (
    id TEXT PRIMARY KEY,
    mock_id TEXT NOT NULL REFERENCES mocks(id) ON DELETE CASCADE,
    type TEXT NOT NULL,
    content TEXT NOT NULL,
    item_order INTEGER NOT NULL,
    created_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS mock_tests (
    id TEXT PRIMARY KEY,
    mock_id TEXT NOT NULL REFERENCES mocks(id) ON DELETE CASCADE,
    input TEXT,
    output TEXT,
    status_code INTEGER,
    took_ms INTEGER,
    created_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_mocks_path_method ON mocks(path, method);
CREATE INDEX IF NOT EXISTS idx_instructions_mock_id ON instructions(mock_id);
CREATE INDEX IF NOT EXISTS idx_mock_tests_mock_id ON mock_tests(mock_id);