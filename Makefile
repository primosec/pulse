.PHONY: run build test migrate sqlc swagger docker-up docker-down lint

run:
	air

build:
	go build -o bin/pulse ./cmd/pulse

test:
	go test -race ./...

migrate:
	goose -dir db/migrations postgres "$(DATABASE_URL)" up

sqlc:
	sqlc generate

swagger:
	swag init -g cmd/pulse/main.go -o docs

docker-up:
	docker compose up -d

docker-down:
	docker compose down

lint:
	golangci-lint run ./...
