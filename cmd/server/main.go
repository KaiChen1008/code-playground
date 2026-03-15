package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

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
	executionTimeout := cfg.Server.ExecutionTimeout
	runnerLangs := make(map[string]runner.Language)
	for k, v := range languages {
		runnerLangs[k] = runner.Language{
			Image:   v.Image,
			Version: v.Version,
		}
	}

	codeRunner := runner.NewTestcontainersRunner(runnerLangs, executionTimeout)

	repo, err := repository.NewFileRepo(dataDir)
	if err != nil {
		logrus.Fatalf("failed to initialize repository: %v", err)
	}
	uc := usecase.New(repo, codeRunner, maxCodeChars, maxTotalSubmissions, languages, maxConcurrentRunners)
	handler := delivery.NewSnippetHandler(uc)
	r := delivery.NewRouter(rateLimit, handler, uc)

	addr := fmt.Sprintf(":%d", port)
	svr := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// ref: github.com/gin-gonic/examples/tree/master/graceful-shutdown/graceful-shutdown/notify-with-context
	// create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	svr.BaseContext = func(net.Listener) context.Context {
		return ctx // for canceling running jobs
	}

	go func() {
		logrus.Infof("Server starting on %s", addr)
		if err := svr.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("Failed to listen and serve: %v", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()
	logrus.Info("Shutting down server...")

	// this context is used to inform the server that it has 5 seconds to finish
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := svr.Shutdown(ctx); err != nil {
		logrus.Fatalf("Failed to shutdown server: %v", err)
	}
	logrus.Println("Server shutdown")
}
