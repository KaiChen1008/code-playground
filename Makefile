IMAGES := $(shell ls cmd)

.PHONY: build run test up down delete clean
build:
	go build -o server ./cmd/server/main.go

run: build
	./server

test:
	go test ./...

up:
	docker-compose up --build -d

down:
	docker-compose down

generate:
	go generate ./...
	go mod tidy

# delete all data files
delete:
	rm -f data/*.json

clean:
	rm -f $(IMAGES)
