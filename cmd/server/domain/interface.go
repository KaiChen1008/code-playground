package domain

import (
	"context"

	"code-playground/cmd/server/domain/models"
)

type Usecase interface {
	RunSnippet(ctx context.Context, req *models.RunRequest) (*models.RunResponse, error)
	GetSnippet(ctx context.Context, id, password string) (*models.Snippet, error)
	DeleteSnippet(ctx context.Context, id string) error
	FormatSnippet(ctx context.Context, req *models.FormatRequest) (*models.FormatResponse, error)
	GetLanguages(ctx context.Context) ([]models.LanguageInfo, error)
}

type CodeRunner interface {
	Run(ctx context.Context, language, code string) (string, error)
	Format(ctx context.Context, language, code string) (string, error)
}

type Repository interface {
	Save(snippet *models.Snippet) error
	GetByID(id string) (*models.Snippet, error)
	Delete(id string) error
}
