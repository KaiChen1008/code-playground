package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"code-playground/cmd/server/domain/models"
	"code-playground/pkg/errors"
)

type fileSnippetRepo struct {
	dataDir string
}

func NewFileSnippetRepo(dataDir string) (*fileSnippetRepo, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, errors.New("failed to create data directory", err)
	}
	return &fileSnippetRepo{dataDir: dataDir}, nil
}

func (r *fileSnippetRepo) Save(snippet *models.Snippet) error {
	filePath := filepath.Join(r.dataDir, fmt.Sprintf("%s.json", snippet.ID))
	data, err := json.Marshal(snippet)
	if err != nil {
		return errors.New("failed to marshal snippet", err)
	}
	return os.WriteFile(filePath, data, 0644)
}

func (r *fileSnippetRepo) GetByID(id string) (*models.Snippet, error) {
	filePath := filepath.Join(r.dataDir, fmt.Sprintf("%s.json", id))
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("snippet not found", err)
		}
		return nil, errors.New("failed to read snippet file", err)
	}
	var snippet models.Snippet
	if err := json.Unmarshal(data, &snippet); err != nil {
		return nil, errors.New("failed to unmarshal snippet", err)
	}
	return &snippet, nil
}

func (r *fileSnippetRepo) Delete(id string) error {
	filePath := filepath.Join(r.dataDir, fmt.Sprintf("%s.json", id))
	err := os.Remove(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("snippet not found", err)
		}
		return errors.New("failed to delete snippet file", err)
	}
	return nil
}
