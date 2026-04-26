ifneq ($(wildcard .env),)
include .env
export
else
$(warning WARNING: .env file not found! Using .env.example)
include .env.example
export
endif

BASE_STACK = docker compose -f docker-compose.yml
INTEGRATION_TEST_STACK = $(BASE_STACK) -f docker-compose-integration-test.yml
ALL_STACK = $(INTEGRATION_TEST_STACK)

# HELP =================================================================================================================
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

db-up: ### Start database
	$(BASE_STACK) up -d db
.PHONY: db-up

db-down: ### Stop database
	$(BASE_STACK) stop db
.PHONY: db-down

compose-up: ### Run docker compose (without backend and reverse proxy)
	$(BASE_STACK) up --build -d db rabbitmq nats && docker compose logs -f
.PHONY: compose-up

compose-up-all: ### Run docker compose (with backend and reverse proxy)
	$(BASE_STACK) up --build -d
.PHONY: compose-up-all

compose-down-all: ### Stop docker compose (with backend and reverse proxy)
	$(BASE_STACK) down
.PHONY: compose-down-all

compose-up-integration-test: ### Run docker compose with integration test
	$(INTEGRATION_TEST_STACK) up --build --abort-on-container-exit --exit-code-from integration-test
.PHONY: compose-up-integration-test

compose-down: ### Down docker compose
	$(ALL_STACK) down --remove-orphans
.PHONY: compose-down

swag-v1: ### swag init
	swag init -g internal/controller/restapi/router.go
.PHONY: swag-v1

sqlc: ### generate source files from sql
	sqlc generate
.PHONY: sqlc

ent: ### generate ent source files
	ent generate ./internal/repo/persistent/ent/schema
.PHONY: ent

proto-v1: ### generate source files from proto
	protoc --go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		docs/proto/v1/*.proto
.PHONY: proto-v1

deps: ### deps tidy + verify
	go mod tidy && go mod verify
.PHONY: deps

deps-audit: ### check dependencies vulnerabilities
	govulncheck ./...
.PHONY: deps-audit

fix-diff: ### Show code changes by `go fix`
	go fix -diff ./...
.PHONY: fix-diff

format: ### Run code formatter
	go fix ./...
	gofumpt -l -w .
	gci write . --skip-generated -s standard -s default
.PHONY: format

run: deps sqlc ent swag-v1 proto-v1 ### swag run for API v1
	go mod download && \
	CGO_ENABLED=0 go run -tags migrate ./cmd/app
.PHONY: run

build: deps sqlc ent swag-v1 proto-v1 ### build the application
	go mod download && \
	CGO_ENABLED=0 go build -o ./main ./cmd/app
.PHONY: build

docker-rm-volume: ### remove docker volume
	docker volume rm go-clean-template-gin_pg-data
.PHONY: docker-rm-volume

linter-golangci: ### check by golangci linter
	golangci-lint run
.PHONY: linter-golangci

linter-hadolint: ### check by hadolint linter
	find . -name "Dockerfile*" | xargs hadolint
.PHONY: linter-hadolint

linter-dotenv: ### check by dotenv linter
	dotenv-linter
.PHONY: linter-dotenv

test: ### run test
	go test -v -race -covermode atomic -coverprofile=coverage.txt ./internal/... ./pkg/...
.PHONY: test

integration-test: ### run integration-test
	go clean -testcache && go test -v ./integration-test/...
.PHONY: integration-test

mock: ### run mockgen
	mockgen -source ./internal/repo/contracts.go -package usecase_test > ./internal/usecase/mocks_repo_test.go
	mockgen -source ./internal/usecase/contracts.go -package usecase_test > ./internal/usecase/mocks_usecase_test.go
.PHONY: mock

migrate-create:  ### create new migration with name="$(name)" from ent schema
	atlas migrate diff "$(name)" --dir "file://migrations" --to "ent://internal/repo/persistent/ent/schema" --dev-url "docker://postgres/16/dev?search_path=public"
.PHONY: migrate-create

migrate-up: ### apply all pending migrations
	atlas migrate apply --url '$(PG_URL)?sslmode=disable' --dir "file://migrations"
.PHONY: migrate-up

migrate-down: ### rollback last migration
	atlas migrate down --url '$(PG_URL)?sslmode=disable' --dir "file://migrations" --amount 1
.PHONY: migrate-down

migrate-redo: ### rollback and reapply last migration
	atlas migrate down --url '$(PG_URL)?sslmode=disable' --dir "file://migrations" --amount 1
	atlas migrate apply --url '$(PG_URL)?sslmode=disable' --dir "file://migrations"
.PHONY: migrate-redo

migrate-status: ### show migration version
	atlas migrate status --url '$(PG_URL)?sslmode=disable' --dir "file://migrations"
.PHONY: migrate-status

migrate-list: ### list migrations, order by modified date
	ls -l migrations/*.up.sql
.PHONY: migrate-list

bin-deps: ### install tools
	go install tool
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install entgo.io/ent/cmd/ent@latest
	go install github.com/air-verse/air@latest
	curl -sSf https://atlasgo.sh | sh
	pnpm install -g swagger-typescript-api
.PHONY: bin-deps

gen-ts:
	# can change location, as well as axios/ky instead of native fetch
	swagger-typescript-api generate \
		--path docs/swagger.yaml \
		--output ./ \
		--name api.ts
.PHONY: gen-ts

schema: ### Generate database schema
	docker exec db pg_dump --schema-only --no-owner --no-privileges $(PG_URL) > docs/schema.sql
.PHONY: schema

pre-commit: swag-v1 proto-v1 mock format linter-golangci test gen-ts schema ### run pre-commit
.PHONY: pre-commit
