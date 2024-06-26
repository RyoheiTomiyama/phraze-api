-include .env

up:
	@docker compose up -d

# GraphQL
gqlgen:
	docker compose run --rm api go run github.com/99designs/gqlgen generate

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
	atlas migrate hash --dir "file://seeds/dev"
	atlas migrate apply\
		--url "postgres://postgres:password@0.0.0.0:5432/phraze?sslmode=disable"\
		--dir "file://seeds/dev"\
		--allow-dirty
	