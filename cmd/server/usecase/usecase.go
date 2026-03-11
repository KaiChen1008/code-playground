package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync/atomic"

	"golang.org/x/crypto/bcrypt"

	"code-playground/cmd/server/domain"
	"code-playground/cmd/server/domain/models"
	"code-playground/pkg/config"
)

type usecase struct {
	repo                domain.Repository
	runner              domain.CodeRunner
	maxCodeChars        int
	maxTotalSubmissions int
	languages           map[string]config.LanguageConfig
	submissionCount     int64
	semaphore           chan struct{}
}

func New(repo domain.Repository, runner domain.CodeRunner, maxCodeChars int, maxTotalSubmissions int, languages map[string]config.LanguageConfig, maxConcurrentRunners int) *usecase {
	var sem chan struct{}
	if maxConcurrentRunners > 0 {
		sem = make(chan struct{}, maxConcurrentRunners)
	}
	return &usecase{
		repo:                repo,
		runner:              runner,
		maxCodeChars:        maxCodeChars,
		maxTotalSubmissions: maxTotalSubmissions,
		languages:           languages,
		semaphore:           sem,
	}
}

func (uc *usecase) RunSnippet(ctx context.Context, req *models.RunRequest) (*models.RunResponse, error) {
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

	if uc.semaphore != nil {
		select {
		case uc.semaphore <- struct{}{}:
			defer func() {
				<-uc.semaphore
			}()
		case <-ctx.Done():
			return nil, ctx.Err()
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

		snippet := &models.Snippet{
			ID:       id,
			Language: *req.Language,
			Code:     *req.Code,
			Output:   output,
		}

		if req.Password != "" {
			hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
			if err != nil {
				return nil, fmt.Errorf("failed to hash password: %w", err)
			}
			snippet.HasPassword = true
			snippet.PasswordHash = string(hash)
		}

		if err := uc.repo.Save(snippet); err != nil {
			return nil, fmt.Errorf("failed to save snippet: %w", err)
		}
	}

	return &models.RunResponse{
		ID:     id,
		Output: output,
	}, nil
}

func (uc *usecase) GetSnippet(ctx context.Context, id, password string) (*models.Snippet, error) {
	snippet, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if snippet.HasPassword {
		if password == "" {
			return nil, fmt.Errorf("password required")
		}
		if err := bcrypt.CompareHashAndPassword([]byte(snippet.PasswordHash), []byte(password)); err != nil {
			return nil, fmt.Errorf("invalid password")
		}
	}

	// Never return the password hash to the client
	snippet.PasswordHash = ""
	return snippet, nil
}

func (uc *usecase) DeleteSnippet(ctx context.Context, id string) error {
	return uc.repo.Delete(id)
}

func (uc *usecase) FormatSnippet(ctx context.Context, req *models.FormatRequest) (*models.FormatResponse, error) {
	if uc.semaphore != nil {
		select {
		case uc.semaphore <- struct{}{}:
			defer func() {
				<-uc.semaphore
			}()
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	formatted, err := uc.runner.Format(ctx, *req.Language, *req.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to format: %w", err)
	}
	return &models.FormatResponse{Code: formatted}, nil
}

func (uc *usecase) GetLanguages(ctx context.Context) ([]models.LanguageInfo, error) {
	var languages []models.LanguageInfo
	for name, langCfg := range uc.languages {
		languages = append(languages, models.LanguageInfo{
			Name:    name,
			Version: langCfg.Version,
		})
	}
	return languages, nil
}
