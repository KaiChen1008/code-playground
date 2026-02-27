package main

import (
	"code-playground/cmd/server/delivery/http"
	"code-playground/cmd/server/repository"
	"code-playground/cmd/server/usecase"
	"code-playground/pkg/config"
	"code-playground/pkg/runner"
	"fmt"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	repo, err := repository.NewFileSnippetRepository(cfg.Server.DataDir)
	if err != nil {
		log.Fatalf("failed to initialize repository: %v", err)
	}

	codeRunner := runner.NewTestcontainersRunner(cfg)
	uc := usecase.NewSnippetUseCase(repo, codeRunner, cfg)
	handler := http.NewSnippetHandler(uc)
	r := http.NewRouter(handler)

	port := fmt.Sprintf(":%d", cfg.Server.Port)

	log.Infof("Server starting on %s", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
