package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync/atomic"

	"code-playground/cmd/server/domain"
	"code-playground/cmd/server/repository"
	"code-playground/pkg/config"
)

type SnippetUseCase interface {
	RunSnippet(ctx context.Context, req *domain.RunRequest) (*domain.RunResponse, error)
	GetSnippet(ctx context.Context, id string) (*domain.Snippet, error)
	DeleteSnippet(ctx context.Context, id string) error
	FormatSnippet(ctx context.Context, req *domain.FormatRequest) (*domain.FormatResponse, error)
	GetLanguages(ctx context.Context) ([]domain.LanguageInfo, error)
}

type snippetUseCase struct {
	repo                repository.SnippetRepository
	runner              CodeRunner
	maxCodeChars        int
	maxTotalSubmissions int
	languages           map[string]config.LanguageConfig
	submissionCount     int64
}

func NewSnippetUseCase(repo repository.SnippetRepository, runner CodeRunner, maxCodeChars int, maxTotalSubmissions int, languages map[string]config.LanguageConfig) SnippetUseCase {
	return &snippetUseCase{
		repo:                repo,
		runner:              runner,
		maxCodeChars:        maxCodeChars,
		maxTotalSubmissions: maxTotalSubmissions,
		languages:           languages,
	}
}

func (uc *snippetUseCase) RunSnippet(ctx context.Context, req *domain.RunRequest) (*domain.RunResponse, error) {
	if uc.maxCodeChars > 0 && req.Code != nil && len(*req.Code) > uc.maxCodeChars {
		return nil, fmt.Errorf("code too long (max %d characters)", uc.maxCodeChars)
	}

	if uc.maxTotalSubmissions > 0 {
		count := atomic.AddInt64(&uc.submissionCount, 1)
		if count > int64(uc.maxTotalSubmissions) {
			atomic.AddInt64(&uc.submissionCount, -1) // revert
			return nil, fmt.Errorf("server has reached maximum number of submissions (%d)", uc.maxTotalSubmissions)
		}
	}

	output, err := uc.runner.Run(ctx, *req.Language, *req.Code)
	if err != nil {
		output = fmt.Sprintf("Error during execution: %v", err)
	}

	id := req.ID
	if id == "" {
		// Generate a 6-character short ID (3 bytes of hex)
		b := make([]byte, 3)
		rand.Read(b)
		id = hex.EncodeToString(b)

		snippet := &domain.Snippet{
			ID:       id,
			Language: *req.Language,
			Code:     *req.Code,
			Output:   output,
		}

		if err := uc.repo.Save(snippet); err != nil {
			return nil, fmt.Errorf("failed to save snippet: %w", err)
		}
	}

	return &domain.RunResponse{
		ID:     id,
		Output: output,
	}, nil
}

func (uc *snippetUseCase) GetSnippet(ctx context.Context, id string) (*domain.Snippet, error) {
	return uc.repo.GetByID(id)
}

func (uc *snippetUseCase) DeleteSnippet(ctx context.Context, id string) error {
	return uc.repo.Delete(id)
}

func (uc *snippetUseCase) FormatSnippet(ctx context.Context, req *domain.FormatRequest) (*domain.FormatResponse, error) {
	formatted, err := uc.runner.Format(ctx, *req.Language, *req.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to format: %w", err)
	}
	return &domain.FormatResponse{Code: formatted}, nil
}

func (uc *snippetUseCase) GetLanguages(ctx context.Context) ([]domain.LanguageInfo, error) {
	var languages []domain.LanguageInfo
	for name, langCfg := range uc.languages {
		languages = append(languages, domain.LanguageInfo{
			Name:    name,
			Version: langCfg.Version,
		})
	}
	return languages, nil
}
