package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"

	"github.com/vinicius/mocker/internal/config"
	"github.com/vinicius/mocker/internal/models"
)

// SQLite implements the storage interface for SQLite
type SQLite struct {
	db *sql.DB
}

// NewSQLite creates a new SQLite storage instance
func NewSQLite(cfg config.DatabaseConfig) (*SQLite, error) {
	db, err := sql.Open("sqlite", cfg.GetDatabaseDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.MaxLifetime)

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	// Run migrations
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return &SQLite{db: db}, nil
}

// Close closes the database connection
func (s *SQLite) Close() error {
	return s.db.Close()
}

// Ping checks the database connection
func (s *SQLite) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

// Mock operations

// CreateMock creates a new mock
func (s *SQLite) CreateMock(ctx context.Context, mock *models.Mock) error {
	query := `
		INSERT INTO mocks (
			id, name, path, method, description, is_active,
			response_type, response_status_code, response_headers, response_body, response_delay_ms,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.ExecContext(ctx, query,
		mock.ID,
		mock.Name,
		mock.Path,
		mock.Method,
		mock.Description,
		mock.IsActive,
		mock.ResponseType,
		mock.ResponseStatus,
		mock.ResponseHeaders,
		mock.ResponseBody,
		mock.ResponseDelayMs,
		mock.CreatedAt.Format(time.RFC3339),
		mock.UpdatedAt.Format(time.RFC3339),
	)

	return err
}

// GetMockByID retrieves a mock by ID
func (s *SQLite) GetMockByID(ctx context.Context, id string) (*models.Mock, error) {
	query := `
		SELECT id, name, path, method, description, is_active,
			response_type, response_status_code, response_headers, response_body, response_delay_ms,
			created_at, updated_at
		FROM mocks WHERE id = ?
	`

	mock := &models.Mock{}
	var createdAt, updatedAt string

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&mock.ID,
		&mock.Name,
		&mock.Path,
		&mock.Method,
		&mock.Description,
		&mock.IsActive,
		&mock.ResponseType,
		&mock.ResponseStatus,
		&mock.ResponseHeaders,
		&mock.ResponseBody,
		&mock.ResponseDelayMs,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}

	mock.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	mock.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

	return mock, nil
}

// GetMockByPathAndMethod retrieves a mock by path and method
func (s *SQLite) GetMockByPathAndMethod(ctx context.Context, path, method string) (*models.Mock, error) {
	query := `
		SELECT id, name, path, method, description, is_active,
			response_type, response_status_code, response_headers, response_body, response_delay_ms,
			created_at, updated_at
		FROM mocks WHERE path = ? AND method = ? AND is_active = 1
	`

	mock := &models.Mock{}
	var createdAt, updatedAt string

	err := s.db.QueryRowContext(ctx, query, path, method).Scan(
		&mock.ID,
		&mock.Name,
		&mock.Path,
		&mock.Method,
		&mock.Description,
		&mock.IsActive,
		&mock.ResponseType,
		&mock.ResponseStatus,
		&mock.ResponseHeaders,
		&mock.ResponseBody,
		&mock.ResponseDelayMs,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}

	mock.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	mock.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

	return mock, nil
}

// ListMocks lists all mocks with pagination
func (s *SQLite) ListMocks(ctx context.Context, page, pageSize int) ([]models.Mock, int, error) {
	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM mocks`
	if err := s.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	query := `
		SELECT id, name, path, method, description, is_active,
			response_type, response_status_code, response_headers, response_body, response_delay_ms,
			created_at, updated_at
		FROM mocks
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := s.db.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var mocks []models.Mock
	for rows.Next() {
		mock := models.Mock{}
		var createdAt, updatedAt string

		err := rows.Scan(
			&mock.ID,
			&mock.Name,
			&mock.Path,
			&mock.Method,
			&mock.Description,
			&mock.IsActive,
			&mock.ResponseType,
			&mock.ResponseStatus,
			&mock.ResponseHeaders,
			&mock.ResponseBody,
			&mock.ResponseDelayMs,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		mock.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		mock.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

		mocks = append(mocks, mock)
	}

	return mocks, total, rows.Err()
}

// UpdateMock updates an existing mock
func (s *SQLite) UpdateMock(ctx context.Context, mock *models.Mock) error {
	query := `
		UPDATE mocks SET
			name = ?, path = ?, method = ?, description = ?, is_active = ?,
			response_type = ?, response_status_code = ?, response_headers = ?, response_body = ?, response_delay_ms = ?,
			updated_at = ?
		WHERE id = ?
	`

	_, err := s.db.ExecContext(ctx, query,
		mock.Name,
		mock.Path,
		mock.Method,
		mock.Description,
		mock.IsActive,
		mock.ResponseType,
		mock.ResponseStatus,
		mock.ResponseHeaders,
		mock.ResponseBody,
		mock.ResponseDelayMs,
		mock.UpdatedAt.Format(time.RFC3339),
		mock.ID,
	)

	return err
}

// DeleteMock deletes a mock by ID
func (s *SQLite) DeleteMock(ctx context.Context, id string) error {
	query := `DELETE FROM mocks WHERE id = ?`
	_, err := s.db.ExecContext(ctx, query, id)
	return err
}

// ActivateMock activates a mock
func (s *SQLite) ActivateMock(ctx context.Context, id string) error {
	query := `UPDATE mocks SET is_active = 1, updated_at = ? WHERE id = ?`
	_, err := s.db.ExecContext(ctx, query, time.Now().Format(time.RFC3339), id)
	return err
}

// DeactivateMock deactivates a mock
func (s *SQLite) DeactivateMock(ctx context.Context, id string) error {
	query := `UPDATE mocks SET is_active = 0, updated_at = ? WHERE id = ?`
	_, err := s.db.ExecContext(ctx, query, time.Now().Format(time.RFC3339), id)
	return err
}

// MockTest operations

// CreateMockTest creates a new mock test
func (s *SQLite) CreateMockTest(ctx context.Context, test *models.MockTest) error {
	query := `
		INSERT INTO mock_tests (id, mock_id, input, output, status_code, took_ms, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.ExecContext(ctx, query,
		test.ID,
		test.MockID,
		test.Input,
		test.Output,
		test.StatusCode,
		test.TookMs,
		test.CreatedAt.Format(time.RFC3339),
	)

	return err
}

// GetMockTestsByMockID retrieves tests for a mock
func (s *SQLite) GetMockTestsByMockID(ctx context.Context, mockID string, page, pageSize int) ([]models.MockTest, int, error) {
	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM mock_tests WHERE mock_id = ?`
	if err := s.db.QueryRowContext(ctx, countQuery, mockID).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	query := `
		SELECT id, mock_id, input, output, status_code, took_ms, created_at
		FROM mock_tests
		WHERE mock_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := s.db.QueryContext(ctx, query, mockID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var tests []models.MockTest
	for rows.Next() {
		test := models.MockTest{}
		var createdAt string

		err := rows.Scan(
			&test.ID,
			&test.MockID,
			&test.Input,
			&test.Output,
			&test.StatusCode,
			&test.TookMs,
			&createdAt,
		)
		if err != nil {
			return nil, 0, err
		}

		test.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)

		tests = append(tests, test)
	}

	return tests, total, rows.Err()
}

// Instruction operations

// CreateInstruction creates a new instruction
func (s *SQLite) CreateInstruction(ctx context.Context, instr *models.Instruction) error {
	query := `
		INSERT INTO instructions (id, mock_id, type, content, item_order, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.ExecContext(ctx, query,
		instr.ID,
		instr.MockID,
		instr.Type,
		instr.Content,
		instr.Order,
		instr.CreatedAt.Format(time.RFC3339),
	)

	return err
}

// GetInstructionsByMockID retrieves instructions for a mock
func (s *SQLite) GetInstructionsByMockID(ctx context.Context, mockID string) ([]models.Instruction, error) {
	query := `
		SELECT id, mock_id, type, content, item_order, created_at
		FROM instructions
		WHERE mock_id = ?
		ORDER BY item_order ASC
	`

	rows, err := s.db.QueryContext(ctx, query, mockID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var instructions []models.Instruction
	for rows.Next() {
		instr := models.Instruction{}
		var createdAt string

		err := rows.Scan(
			&instr.ID,
			&instr.MockID,
			&instr.Type,
			&instr.Content,
			&instr.Order,
			&createdAt,
		)
		if err != nil {
			return nil, err
		}

		instr.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)

		instructions = append(instructions, instr)
	}

	return instructions, rows.Err()
}

// DeleteInstructionsByMockID deletes all instructions for a mock
func (s *SQLite) DeleteInstructionsByMockID(ctx context.Context, mockID string) error {
	query := `DELETE FROM instructions WHERE mock_id = ?`
	_, err := s.db.ExecContext(ctx, query, mockID)
	return err
}

// runMigrations runs database migrations
func runMigrations(db *sql.DB) error {
	migration := `
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
	`

	_, err := db.Exec(migration)
	return err
}
