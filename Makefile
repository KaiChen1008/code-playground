IMAGES := $(shell ls cmd)

.PHONY: build run test up down delete clean ui-install ui-dev ui-build ui-lint

ui-install:
	cd ui && yarn install

ui-dev:
	cd ui && yarn dev

ui-build:
	cd ui && yarn build

ui-lint:
	cd ui && yarn lint

build: ui-build
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
