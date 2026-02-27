package repository

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"code-playground/cmd/server/domain/models"
)

func TestFileRepo(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "repo-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	repo, err := NewFileRepo(tempDir)
	assert.NoError(t, err)

	snippet := &models.Snippet{
		ID:       "test-id",
		Language: "golang",
		Code:     "fmt.Println(1)",
	}

	// Save
	err = repo.Save(snippet)
	assert.NoError(t, err)

	// Get
	got, err := repo.GetByID("test-id")
	assert.NoError(t, err)
	assert.Equal(t, snippet, got)

	// Delete
	err = repo.Delete("test-id")
	assert.NoError(t, err)

	// Get again
	_, err = repo.GetByID("test-id")
	assert.Error(t, err)
}
