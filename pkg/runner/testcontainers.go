package runner

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type Language struct {
	Image   string
	Version string
}

type TestcontainersRunner struct {
	languages        map[string]Language
	executionTimeout time.Duration
}

func NewTestcontainersRunner(languages map[string]Language, executionTimeout time.Duration) *TestcontainersRunner {
	return &TestcontainersRunner{
		languages:        languages,
		executionTimeout: executionTimeout,
	}
}

func (r *TestcontainersRunner) Run(ctx context.Context, language, code string) (string, error) {
	if r.executionTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, r.executionTimeout)
		defer cancel()
	}

	langCfg, ok := r.languages[language]
	if !ok {
		return "", fmt.Errorf("unsupported language: %s", language)
	}

	var entrypoint []string
	var filename string

	switch language {
	case "golang":
		filename = "main.go"
		entrypoint = []string{"go", "run", "/main.go"}
	case "python":
		filename = "main.py"
		entrypoint = []string{"python", "/main.py"}
	case "javascript":
		filename = "main.js"
		entrypoint = []string{"node", "/main.js"}
	case "rust":
		filename = "main.rs"
		// Using rustc to compile and then run for simple single-file scripts
		entrypoint = []string{"sh", "-c", "rustc /main.rs -o /main && /main"}
	case "cpp":
		filename = "main.cpp"
		// Using g++ to compile and then run
		entrypoint = []string{"sh", "-c", "g++ /main.cpp -o /main && /main"}
	default:
		return "", fmt.Errorf("unsupported language: %s", language)
	}

	// Create a temp directory to hold the code file
	tempDir, err := os.MkdirTemp("", "code-run-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	if err := os.WriteFile(filepath.Join(tempDir, filename), []byte(code), 0644); err != nil {
		return "", fmt.Errorf("failed to write code to file: %w", err)
	}

	req := testcontainers.ContainerRequest{
		Image: langCfg.Image,
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      filepath.Join(tempDir, filename),
				ContainerFilePath: "/" + filename,
				FileMode:          0644,
			},
		},
		Cmd:        entrypoint,
		WorkingDir: "/",
		// Using Stdout and Stderr to capture output.
		// Wait for the container to finish.
		WaitingFor: wait.ForExit(),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return "", fmt.Errorf("failed to start container: %w", err)
	}
	defer container.Terminate(ctx)

	// Wait for the container to finish.
	// Actually GenericContainer with Started: true and WaitingFor: wait.ForExit()
	// should wait until it exits.

	logs, err := container.Logs(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get container logs: %w", err)
	}
	defer logs.Close()

	output, err := io.ReadAll(logs)
	if err != nil {
		return "", fmt.Errorf("failed to read container logs: %w", err)
	}

	// container.Logs() often includes some prefix bytes (like docker log stream header)
	// but for simplicity we'll just trim it or assume it's clean if we use the right reader.
	// Actually, the logs from GenericContainer.Logs() are often multiplexed.
	// Let's use a cleaner way if possible, or just take the whole thing.

	res := string(output)
	// Trimming the first 8 bytes of each line if it's multiplexed (docker format)
	// But let's see. If it's a simple run, maybe it's fine.

	return cleanLogs(res), nil
}

func (r *TestcontainersRunner) Format(ctx context.Context, language, code string) (string, error) {
	if r.executionTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, r.executionTimeout)
		defer cancel()
	}

	langCfg, ok := r.languages[language]
	if !ok {
		return code, nil
	}

	if language != "golang" {
		// Fallback for others - we might just return the same code if no backend formatter exists
		// or let Monaco handle it if possible.
		return code, nil
	}

	image := langCfg.Image
	filename := "main.go"
	entrypoint := []string{"gofmt", "/main.go"}

	tempDir, err := os.MkdirTemp("", "code-format-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	if err := os.WriteFile(filepath.Join(tempDir, filename), []byte(code), 0644); err != nil {
		return "", fmt.Errorf("failed to write code to file: %w", err)
	}

	req := testcontainers.ContainerRequest{
		Image: image,
		HostConfigModifier: func(hc *container.HostConfig) {
			hc.Resources = container.Resources{
				Memory:   128 * 1024 * 1024, // 128 MB
				NanoCPUs: 1 * 1e9,           // 0.1 cpu
			}
		},
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      filepath.Join(tempDir, filename),
				ContainerFilePath: "/" + filename,
				FileMode:          0644,
			},
		},
		Cmd:        entrypoint,
		WorkingDir: "/",
		WaitingFor: wait.ForExit(),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return "", fmt.Errorf("failed to start container: %w", err)
	}
	defer container.Terminate(ctx)

	logs, err := container.Logs(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get container logs: %w", err)
	}
	defer logs.Close()

	output, err := io.ReadAll(logs)
	if err != nil {
		return "", fmt.Errorf("failed to read container logs: %w", err)
	}

	return cleanLogs(string(output)), nil
}

// cleanLogs handles docker log multiplexing headers (8 bytes) if present.
func cleanLogs(s string) string {
	if len(s) < 8 {
		return s
	}

	var sb strings.Builder
	reader := strings.NewReader(s)

	for {
		header := make([]byte, 8)
		_, err := reader.Read(header)
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}

		// First byte is stream type (1=stdout, 2=stderr)
		// Bytes 4-7 are payload size (BigEndian)
		if header[0] != 1 && header[0] != 2 {
			// Not a docker header or garbled, just return original as fallback
			// But let's check if it's the very first packet.
			if sb.Len() == 0 {
				return s
			}
			break
		}

		size := binary.BigEndian.Uint32(header[4:])
		payload := make([]byte, size)
		_, err = reader.Read(payload)
		if err != nil {
			break
		}
		sb.Write(payload)
	}

	if sb.Len() == 0 {
		return s
	}

	return strings.TrimSpace(sb.String())
}
