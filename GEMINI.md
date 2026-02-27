# GEMINI.md - code-playground

This project is a full-stack project. it's like a ideone.com, but a self-hosted version.

I only want to focus on three langauges, golang, python, javascript.

Front-end:

- put in `ui/`
- ig-style theme
- users can choose a langugae(python, golang, javascript)
- users can write/paste code to a block.
  - this block should support highlights for code (based on the language).
  - this block should also support formatting.
- press `run` to submit a request to backend to run the code and return the result and unique id for this code.
- generate a unique url(id) for this code.
  - for other people, they can use like localhost:8091/{id} to see the code
- need a buttom to delete a code

Back-end:

- [ ] use golang to implement this. put it in cmd.
- [ ] follow clean-arch.
- [ ] store user-submitted code in `./data` . store codes directly, do not use a database.
- [ ] use `testcontainer` to run the submitted codes.
- [ ] use docker & docker-compose to run the server

## Tech Stack

- [ ] **API**: spec first, generate open api spec first, and use `go-swagger` to generate code.
- [ ] **Go**: Primary language for the distributed system.
- [ ] **Python**: Used for auxiliary scripts.
- [ ] **Web**: Gin-based API and a simple UI.
- [ ] **Package Management**:
  - [ ] Go: Standard Go modules.
  - [ ] Python: `uv`.
- [ ] **Infrastructure**: Docker and Docker Compose.

## Project Structure

- `cmd/`: Application entry points.
- `pkg/`: Core libraries and shared logic.
- `scripts/`: Python-based utility scripts.
- `ui/`: Web interface assets.
- `data/`: Default directory for temporary download artifacts.

## Development Workflow

### Go (Primary)

- **Run Tests**: `go test ./...`
- **Build Docker Images**
- **Style**: Follow standard `go fmt` and `go vet`. Use `logrus` for logging as established in the codebase.
- **API**: use go-swagger to generate code.

### Python (Scripts)

- **Tooling**: Always use `uv` for managing Python dependencies and running scripts.
- **Location**: All Python logic resides in `scripts/`.

### Infrastructure

- **Full Stack**: `make up` to start the entire distributed system.

## Workspace Conventions

- **Configuration**: Managed via `config.yaml` and handled by `pkg/config`.
- **Error Handling**: Use the custom error wrappers in `pkg/errors` for consistency.
- **Testing**: Add unit tests in `_test.go` files. Use `testify` for assertions as seen in `pkg/job/queue_test.go`.


## Before Stop

- `run go mod tidy` before you stop generating code.
- make sure there's no error in each files.
