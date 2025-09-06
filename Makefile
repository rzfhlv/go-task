.PHONY: build run

build:
	go build -o build/main cmd/api/main.go

run: build
	./build/main start

migrate-create: build
	./build/main migrate create --name=$(name)

migrate-up: build
	./build/main migrate up

migrate-down: build
	./build/main migrate down

deps-up:
	docker compose up -d postgres redis

deps-down:
	docker compose down

generate-mock:
	mockery --config .mockery.yml

test:
	go test ./... -cover -race -coverprofile=coverage.out -covermode=atomic
	./script/coverage.sh

test-analyze:
	go test ./... -cover -race -coverprofile=coverage.out -covermode=atomic
	go tool cover -html=coverage.out -o coverage.html
	open coverage.html
