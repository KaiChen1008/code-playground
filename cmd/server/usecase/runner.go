package usecase

import (
	"context"
)

type CodeRunner interface {
	Run(ctx context.Context, language, code string) (string, error)
	Format(ctx context.Context, language, code string) (string, error)
}
