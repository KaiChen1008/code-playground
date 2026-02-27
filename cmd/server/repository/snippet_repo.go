package repository

import (
	"code-playground/cmd/server/domain"
	"code-playground/pkg/errors"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type SnippetRepository interface {
	Save(snippet *domain.Snippet) error
	GetByID(id string) (*domain.Snippet, error)
	Delete(id string) error
}

type fileSnippetRepository struct {
	dataDir string
}

func NewFileSnippetRepository(dataDir string) (SnippetRepository, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, errors.New("failed to create data directory", err)
	}
	return &fileSnippetRepository{dataDir: dataDir}, nil
}

func (r *fileSnippetRepository) Save(snippet *domain.Snippet) error {
	filePath := filepath.Join(r.dataDir, fmt.Sprintf("%s.json", snippet.ID))
	data, err := json.Marshal(snippet)
	if err != nil {
		return errors.New("failed to marshal snippet", err)
	}
	return os.WriteFile(filePath, data, 0644)
}

func (r *fileSnippetRepository) GetByID(id string) (*domain.Snippet, error) {
	filePath := filepath.Join(r.dataDir, fmt.Sprintf("%s.json", id))
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("snippet not found", err)
		}
		return nil, errors.New("failed to read snippet file", err)
	}
	var snippet domain.Snippet
	if err := json.Unmarshal(data, &snippet); err != nil {
		return nil, errors.New("failed to unmarshal snippet", err)
	}
	return &snippet, nil
}

func (r *fileSnippetRepository) Delete(id string) error {
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
