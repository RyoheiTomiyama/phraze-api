-include .env

up:
	@docker compose up -d

migrate:
	atlas schema apply\
		--url "postgres://postgres:password@0.0.0.0:5432/phraze?sslmode=disable"\
		--to "file://atlas/schema.hcl"\
		--dev-url "docker://postgres"

migrate-diff:
	atlas schema diff\
		--from "postgres://postgres:password@0.0.0.0:5432/phraze?sslmode=disable"\
		--to "file://atlas/schema.hcl"\
		--dev-url "docker://postgres"
