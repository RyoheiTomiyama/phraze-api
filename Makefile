-include .env

up:
	@docker compose up -d

# GraphQL
gqlgen:
	docker compose run --rm api go run github.com/99designs/gqlgen generate

# Vulncheck
vulncheck:
	docker compose run --rm api sh -c \
	"go install golang.org/x/vuln/cmd/govulncheck@latest ; govulncheck ./..."

# Lint
lint:
	docker run --rm -v $(shell pwd)/app:/app -w /app golangci/golangci-lint:v1.59.1-alpine golangci-lint run

# DATABASE 
migrate:
	atlas schema apply\
		--url "postgres://postgres:password@0.0.0.0:5432/phraze?sslmode=disable"\
		--to "file://atlas/schema.sql"\
		--dev-url "docker://postgres"

migrate-diff:
	atlas schema diff\
		--from "postgres://postgres:password@0.0.0.0:5432/phraze?sslmode=disable"\
		--to "file://atlas/schema.sql"\
		--dev-url "docker://postgres"

seed:
	atlas migrate hash --dir "file://atlas/seeds/dev"
	atlas migrate apply\
		--url "postgres://postgres:password@0.0.0.0:5432/phraze?sslmode=disable"\
		--dir "file://atlas/seeds/dev"\
		--allow-dirty
	