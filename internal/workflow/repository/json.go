package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/luis12loureiro/neurun/internal/workflow/domain"
)

const (
	FilePermUserReadWriteGroupRead = 0644
	FilePermUserReadWriteOnly      = 0600
)

type JSONRepo struct {
	path     string
	filename string
}

func NewJSONRepository(path string, filename string) domain.WorkflowRepository {
	return &JSONRepo{path: path, filename: filename}
}

func (r *JSONRepo) Create(t *domain.Worklow) error {
	filePath := filepath.Join(r.path, r.filename)
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, FilePermUserReadWriteGroupRead)
	if err != nil {
		return fmt.Errorf("failed to read workflow file: %w", err)
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(t); err != nil {
		return fmt.Errorf("failed to write workflow file: %w", err)
	}
	return nil
}

func (r *JSONRepo) Get(id string) (*domain.Worklow, error) {
	var w domain.Worklow
	filePath := filepath.Join(r.path, r.filename)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	dec := json.NewDecoder(file)
	for dec.More() {
		if err := dec.Decode(&w); err != nil {
			return nil, fmt.Errorf("failed to decode workflow: %w", err)
		}
		if w.ID == id {
			return &w, nil
		}
	}
	return nil, fmt.Errorf("workflow with id %s not found", id)
}
