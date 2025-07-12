package repository

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/luis12loureiro/neurun/internal/workflow/domain"
)

type JSONRepo struct {
	path     string
	filename string
}

func NewJSONRepository(path string, filename string) domain.WorkflowRepository {
	return &JSONRepo{path: path, filename: filename}
}

func (r *JSONRepo) Create(t domain.Worklow) error {
	filePath := filepath.Join(r.path, r.filename)
	file, err := os.OpenFile(filePath, os.O_APPEND, os.ModePerm)
	if err != nil {
		return errors.New("failed to open workflow file")
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(t)
}

func (r *JSONRepo) Get(id string) (domain.Worklow, error) {
	var w domain.Worklow

	filePath := filepath.Join(r.path, r.filename)
	file, err := os.Open(filePath)
	if err != nil {
		return w, errors.New("task file not found")
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&w); err != nil {
		return w, errors.New("task not found")
	}

	return w, nil
}
