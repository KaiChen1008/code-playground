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

swagger:
	-docker run --rm -v $$(pwd):/work -w /work quay.io/goswagger/swagger generate model -f api/swagger.yaml -m cmd/server/domain
	go fmt ./cmd/server/domain/...

# delete all data files
delete:
	rm -f data/*.json

clean:
	rm -f $(IMAGES)
