package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

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
	for _, task := range w.Tasks {
		if err := r.createTask(tx, task, w.ID); err != nil {
			return fmt.Errorf("failed to insert task %s: %w", task.ID, err)
		}
	}
	// commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (r *SQLiteRepo) Get(id string) (*domain.Workflow, error) {
	query := `
		SELECT 
			w.id, w.name, w.description, w.status,
			t.id, t.name, t.type, t.status, t.retries, t.retry_delay, t.condition,
			lp.message,
			hp.url, hp.method, hp.body, hp.headers, hp.query_params, 
			hp.timeout, hp.follow_redirects, hp.verify_ssl, hp.expected_status_code,
			ha.auth_type, ha.auth_data
		FROM workflow w
		LEFT JOIN task t ON w.id = t.workflow_id
		LEFT JOIN log_payload lp ON t.id = lp.task_id
		LEFT JOIN http_payload hp ON t.id = hp.task_id
		LEFT JOIN http_auth ha ON t.id = ha.task_id
		WHERE w.id = ?
		ORDER BY t.id`

	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to query workflow: %w", err)
	}
	defer rows.Close()

	var workflow *domain.Workflow
	tasksMap := make(map[string]*domain.Task)

	for rows.Next() {
		var (
			// workflow fields
			wID, wName, wDescription, wStatus string
			// task fields (nullable)
			tID, tName, tType, tStatus, tCondition sql.NullString
			tRetries                               sql.NullInt32
			tRetryDelayMs                          sql.NullInt64
			// log payload (nullable)
			logMessage sql.NullString
			// HTTP payload (nullable)
			httpURL, httpMethod, httpBody, httpHeaders, httpQueryParams sql.NullString
			httpTimeoutMs                                               sql.NullInt64
			httpFollowRedirects, httpVerifySSL                          sql.NullBool
			httpExpectedStatusCode                                      sql.NullInt32
			// HTTP auth (nullable)
			authType, authDataJSON sql.NullString
		)

		err := rows.Scan(
			&wID, &wName, &wDescription, &wStatus,
			&tID, &tName, &tType, &tStatus, &tRetries, &tRetryDelayMs, &tCondition,
			&logMessage,
			&httpURL, &httpMethod, &httpBody, &httpHeaders, &httpQueryParams,
			&httpTimeoutMs, &httpFollowRedirects, &httpVerifySSL, &httpExpectedStatusCode,
			&authType, &authDataJSON,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// create workflow if not exists
		if workflow == nil {
			workflow = &domain.Workflow{
				ID:          wID,
				Name:        wName,
				Description: wDescription,
				Status:      domain.WorklowStatus(wStatus),
			}
		}

		// skip if no task (workflow with no tasks)
		if !tID.Valid {
			continue
		}

		taskID := tID.String
		_, exists := tasksMap[taskID]
		if !exists {
			// create new task
			task := &domain.Task{
				ID:         taskID,
				Name:       tName.String,
				Type:       domain.TaskType(tType.String),
				Status:     domain.TaskStatus(tStatus.String),
				Retries:    uint8(tRetries.Int32),
				RetryDelay: time.Duration(tRetryDelayMs.Int64) * time.Second,
				Condition:  tCondition.String,
			}

			// set payload based on task type
			switch task.Type {
			case domain.TaskTypeLog:
				if logMessage.Valid {
					task.Payload = &domain.LogPayload{
						Message: logMessage.String,
					}
				}
			case domain.TaskTypeHTTP:
				if httpURL.Valid {
					httpPayload := &domain.HTTPPayload{
						URL:                httpURL.String,
						Method:             httpMethod.String,
						Body:               []byte(httpBody.String),
						Timeout:            time.Duration(httpTimeoutMs.Int64) * time.Second,
						FollowRedirects:    httpFollowRedirects.Bool,
						VerifySSL:          httpVerifySSL.Bool,
						ExpectedStatusCode: httpExpectedStatusCode.Int32,
					}
					// unmarshal JSON fields
					if httpHeaders.Valid && httpHeaders.String != "" {
						if err := json.Unmarshal([]byte(httpHeaders.String), &httpPayload.Headers); err != nil {
							return nil, fmt.Errorf("failed to unmarshal headers: %w", err)
						}
					}
					if httpQueryParams.Valid && httpQueryParams.String != "" {
						if err := json.Unmarshal([]byte(httpQueryParams.String), &httpPayload.QueryParams); err != nil {
							return nil, fmt.Errorf("failed to unmarshal query params: %w", err)
						}
					}
					// set auth if present
					if authType.Valid && authDataJSON.Valid {
						auth, err := r.unmarshalAuth(authType.String, authDataJSON.String)
						if err != nil {
							return nil, fmt.Errorf("failed to unmarshal auth: %w", err)
						}
						httpPayload.Auth = auth
					}
					task.Payload = httpPayload
				}
			}
			tasksMap[taskID] = task
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}
	if workflow == nil {
		return nil, fmt.Errorf("workflow with id %s not found", id)
	}

	// load task relationships
	if len(tasksMap) > 0 {
		if err := r.loadTaskRelationships(tasksMap, id); err != nil {
			return nil, fmt.Errorf("failed to load task relationships: %w", err)
		}
	}
	// find root tasks
	workflow.Tasks = r.findRootTasks(tasksMap)
	return workflow, nil
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

func (r *SQLiteRepo) createTask(tx *sql.Tx, task *domain.Task, workflowID string) error {
	taskQuery := `
        INSERT INTO task (id, name, type, status, retries, retry_delay, condition, workflow_id)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := tx.Exec(taskQuery, task.ID, task.Name, task.Type, task.Status,
		task.Retries, task.RetryDelay.Seconds(), task.Condition, workflowID)
	if err != nil {
		return fmt.Errorf("failed to insert task: %w", err)
	}
	switch task.Type {
	case domain.TaskTypeLog:
		if err := r.insertLogPayload(tx, task); err != nil {
			return err
		}
	case domain.TaskTypeHTTP:
		if err := r.insertHTTPPayload(tx, task); err != nil {
			return err
		}
	}
	for _, nextTask := range task.Next {
		if err := r.createTask(tx, nextTask, workflowID); err != nil {
			return err
		}
		nextQuery := `
            INSERT INTO task_next (task_id, next_task_id)
            VALUES (?, ?)`
		_, err := tx.Exec(nextQuery, task.ID, nextTask.ID)
		if err != nil {
			return fmt.Errorf("failed to insert task next: %w", err)
		}
	}
	return nil
}

func (r *SQLiteRepo) insertLogPayload(tx *sql.Tx, task *domain.Task) error {
	logPayload, ok := task.Payload.(*domain.LogPayload)
	if !ok {
		return fmt.Errorf("invalid payload type for log task")
	}
	query := `
        INSERT INTO log_payload (task_id, message)
        VALUES (?, ?)`
	_, err := tx.Exec(query, task.ID, logPayload.Message)
	if err != nil {
		return fmt.Errorf("failed to insert log payload: %w", err)
	}
	return nil
}

func (r *SQLiteRepo) insertHTTPPayload(tx *sql.Tx, task *domain.Task) error {
	httpPayload, ok := task.Payload.(*domain.HTTPPayload)
	if !ok {
		return fmt.Errorf("invalid payload type for HTTP task")
	}
	// marshal headers and query params to JSON
	headersJSON, err := json.Marshal(httpPayload.Headers)
	if err != nil {
		return fmt.Errorf("failed to marshal headers: %w", err)
	}
	queryParamsJSON, err := json.Marshal(httpPayload.QueryParams)
	if err != nil {
		return fmt.Errorf("failed to marshal query params: %w", err)
	}
	query := `
        INSERT INTO http_payload (task_id, url, method, body, headers, query_params, 
            timeout, follow_redirects, verify_ssl, expected_status_code)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = tx.Exec(query, task.ID, httpPayload.URL, httpPayload.Method, httpPayload.Body,
		string(headersJSON), string(queryParamsJSON), httpPayload.Timeout.Seconds(),
		httpPayload.FollowRedirects, httpPayload.VerifySSL, httpPayload.ExpectedStatusCode)
	if err != nil {
		return fmt.Errorf("failed to insert HTTP payload: %w", err)
	}
	if httpPayload.Auth != nil {
		if err := r.insertHTTPAuth(tx, task.ID, httpPayload.Auth); err != nil {
			return err
		}
	}
	return nil
}

func (r *SQLiteRepo) insertHTTPAuth(tx *sql.Tx, taskID string, auth domain.HTTPAuthType) error {
	authType := auth.Type()
	authDataJSON, err := json.Marshal(auth)
	if err != nil {
		return fmt.Errorf("failed to marshal auth data: %w", err)
	}
	query := `
        INSERT INTO http_auth (task_id, auth_type, auth_data)
        VALUES (?, ?, ?)`
	_, err = tx.Exec(query, taskID, authType, string(authDataJSON))
	if err != nil {
		return fmt.Errorf("failed to insert HTTP auth: %w", err)
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

func (r *SQLiteRepo) loadTaskRelationships(tasksMap map[string]*domain.Task, workflowID string) error {
	// get all relationships for this workflow
	query := `
		SELECT tn.task_id, tn.next_task_id
		FROM task_next tn
		INNER JOIN task t ON tn.task_id = t.id
		WHERE t.workflow_id = ?`
	rows, err := r.db.Query(query, workflowID)
	if err != nil {
		return fmt.Errorf("failed to query task relationships: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var taskID, nextTaskID string
		if err := rows.Scan(&taskID, &nextTaskID); err != nil {
			return fmt.Errorf("failed to scan task relationship: %w", err)
		}
		task, taskExists := tasksMap[taskID]
		nextTask, nextExists := tasksMap[nextTaskID]
		if taskExists && nextExists {
			task.Next = append(task.Next, nextTask)
		}
	}
	return rows.Err()
}

func (r *SQLiteRepo) unmarshalAuth(authType, authDataJSON string) (domain.HTTPAuthType, error) {
	switch authType {
	case domain.HTTPBasicAuthType:
		var auth domain.HTTPBasicAuth
		if err := json.Unmarshal([]byte(authDataJSON), &auth); err != nil {
			return nil, fmt.Errorf("failed to unmarshal basic auth: %w", err)
		}
		return &auth, nil
	case domain.HTTPBearerAuthType:
		var auth domain.HTTPBearerAuth
		if err := json.Unmarshal([]byte(authDataJSON), &auth); err != nil {
			return nil, fmt.Errorf("failed to unmarshal bearer auth: %w", err)
		}
		return &auth, nil
	case domain.HTTPApiKeyAuthType:
		var auth domain.HTTPApiKeyAuth
		if err := json.Unmarshal([]byte(authDataJSON), &auth); err != nil {
			return nil, fmt.Errorf("failed to unmarshal API key auth: %w", err)
		}
		return &auth, nil
	default:
		return nil, fmt.Errorf("unknown auth type: %s", authType)
	}
}

func (r *SQLiteRepo) findRootTasks(tasksMap map[string]*domain.Task) []*domain.Task {
	// find tasks that are not referenced as 'next' by any other task
	referencedTasks := make(map[string]bool)
	for _, task := range tasksMap {
		for _, nextTask := range task.Next {
			referencedTasks[nextTask.ID] = true
		}
	}
	var rootTasks []*domain.Task
	for _, task := range tasksMap {
		if !referencedTasks[task.ID] {
			rootTasks = append(rootTasks, task)
		}
	}
	return rootTasks
}

func (r *SQLiteRepo) Close() error {
	return r.db.Close()
}
