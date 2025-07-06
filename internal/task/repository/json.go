package repository

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	domain "github.com/luis12loureiro/neurun/internal/task/domain"
)

type JSONRepo struct {
	path     string
	filename string
}

func NewJSONRepository(path string, filename string) *JSONRepo {
	return &JSONRepo{path: path, filename: filename}
}

func (r *JSONRepo) Create(t domain.Task) error {
	file, err := os.OpenFile("big_encode.json", os.O_APPEND, os.ModePerm)
	if err != nil {
		return errors.New("failed to open task file")
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(t)
}

func (r *JSONRepo) GetTask(id string) (domain.Task, error) {
	var task domain.Task

	filePath := filepath.Join(r.path, r.filename)
	file, err := os.Open(filePath)
	if err != nil {
		return task, errors.New("task file not found")
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&task); err != nil {
		return task, errors.New("task not found")
	}

	return task, nil
}
