package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/luis12loureiro/neurun/internal/workflow/domain"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteRepo struct {
	db *sql.DB
}

func NewSQLiteRepository(dbPath string, name string) (domain.WorkflowRepository, error) {
	filePath := filepath.Join(dbPath, name)
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	repo := &SQLiteRepo{db: db}
	if err := repo.createTables(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}
	return repo, nil
}

func (r *SQLiteRepo) Create(w *domain.Workflow) error {
	// start transaction
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()
	query := `
		INSERT INTO workflow (id, name, description, status)
		VALUES (?, ?, ?, ?)`
	_, err = tx.Exec(query, w.ID, w.Name, w.Description, w.Status)
	if err != nil {
		return fmt.Errorf("failed to insert workflow: %w", err)
	}
	// commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (r *SQLiteRepo) Get(id string) (*domain.Workflow, error) {
	query := `
		SELECT id, name, description, status
		FROM workflow
		WHERE id = ?`
	var w domain.Workflow
	err := r.db.QueryRow(query, id).Scan(
		&w.ID,
		&w.Name,
		&w.Description,
		&w.Status,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("workflow with id %s not found", id)
		}
		return nil, fmt.Errorf("failed to get workflow: %w", err)
	}
	return &w, nil
}

func (r *SQLiteRepo) Update(w *domain.Workflow) error {
	tasksJSON, err := json.Marshal(w.Tasks)
	if err != nil {
		return fmt.Errorf("failed to marshal tasks: %w", err)
	}

	query := `
		UPDATE workflow
		SET name = ?, description = ?, status = ?, tasks_json = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`

	result, err := r.db.Exec(query, w.Name, w.Description, w.Status, string(tasksJSON), w.ID)
	if err != nil {
		return fmt.Errorf("failed to update workflow: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("workflow with id %s not found for update", w.ID)
	}

	return nil
}

func (r *SQLiteRepo) createTables() error {
	createDbTables := `
	CREATE TABLE IF NOT EXISTS workflow (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		status TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS task (
        id TEXT PRIMARY KEY,
        name TEXT NOT NULL,
        type TEXT NOT NULL,
        status TEXT NOT NULL,
        retries INTEGER DEFAULT 0,
        retry_delay INTEGER DEFAULT 0,
        condition TEXT,
        workflow_id TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (workflow_id) REFERENCES workflow(id) ON DELETE CASCADE
    );
	CREATE TABLE IF NOT EXISTS task_next (
        task_id TEXT NOT NULL,
        next_task_id TEXT NOT NULL,
        PRIMARY KEY (task_id, next_task_id),
        FOREIGN KEY (task_id) REFERENCES task(id) ON DELETE CASCADE,
        FOREIGN KEY (next_task_id) REFERENCES task(id) ON DELETE CASCADE
    );
	CREATE TABLE IF NOT EXISTS log_payload (
        task_id TEXT PRIMARY KEY,
        message TEXT NOT NULL,
        FOREIGN KEY (task_id) REFERENCES task(id) ON DELETE CASCADE
    );
	CREATE TABLE IF NOT EXISTS http_payload (
        task_id TEXT PRIMARY KEY,
        url TEXT NOT NULL,
        method TEXT NOT NULL,
        body BLOB,
        headers TEXT, -- JSON object for headers map
        query_params TEXT, -- JSON object for query params map
        timeout INTEGER NOT NULL DEFAULT 3,
        follow_redirects BOOLEAN NOT NULL DEFAULT FALSE,
        verify_ssl BOOLEAN NOT NULL DEFAULT FALSE,
        expected_status_code INTEGER NOT NULL DEFAULT 200,
        FOREIGN KEY (task_id) REFERENCES task(id) ON DELETE CASCADE
    );
	CREATE TABLE IF NOT EXISTS http_auth (
        task_id TEXT PRIMARY KEY,
        auth_type TEXT NOT NULL, -- 'basic', 'bearer', 'apikey'
        auth_data TEXT NOT NULL, -- JSON object containing auth details
        FOREIGN KEY (task_id) REFERENCES task(id) ON DELETE CASCADE
    );
	`
	_, err := r.db.Exec(createDbTables)
	if err != nil {
		return fmt.Errorf("failed to create db tables: %w", err)
	}
	return nil
}

func (r *SQLiteRepo) Close() error {
	return r.db.Close()
}
