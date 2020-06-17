include Makefile.vars
include Makefile.common

DATABASE_HOST=localhost:5432
DATABASE_NAME=atlas_cli_test
DATABASE_USER=postgres
DATABASE_PASSWORD=postgres

migrate:
	@migrate -database 'postgres://$(DATABASE_USER):$(DATABASE_PASSWORD)@$(DATABASE_HOST)/$(DATABASE_NAME)?sslmode=disable' -path ./db/migration up

local-build:
	go build -o test cmd/server/*.go
