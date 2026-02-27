package main

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"code-playground/cmd/server/delivery"
	"code-playground/cmd/server/repository"
	"code-playground/cmd/server/usecase"
	"code-playground/pkg/config"
	"code-playground/pkg/runner"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		logrus.Fatalf("failed to load config: %v", err)
	}
	dataDir := cfg.Server.DataDir
	languages := cfg.Languages
	maxCodeChars := cfg.Server.MaxCodeChars
	maxTotalSubmissions := cfg.Server.MaxTotalSubmissions
	maxConcurrentRunners := cfg.Server.MaxConcurrentRunners
	rateLimit := cfg.Server.RateLimit
	port := cfg.Server.Port

	repo, err := repository.NewFileRepo(dataDir)
	if err != nil {
		logrus.Fatalf("failed to initialize repository: %v", err)
	}

	codeRunner := runner.NewTestcontainersRunner(languages)
	uc := usecase.New(repo, codeRunner, maxCodeChars, maxTotalSubmissions, languages, maxConcurrentRunners)
	handler := delivery.NewSnippetHandler(uc)
	r := delivery.NewRouter(rateLimit, handler)

	addr := fmt.Sprintf(":%d", port)

	logrus.Infof("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		logrus.Fatalf("failed to start server: %v", err)
	}
}
