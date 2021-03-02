BIN_DIR     ?= ./bin

COVERAGE_PROFILE ?= coverage.out

DATABASE_USER             ?= discoverrewind_admin
TEST_DATABASE_NAME        ?= discoverrewind_test
DEVELOPMENT_DATABASE_NAME ?= discoverrewind

DATABASE_DEBUG ?= true
LOCAL_GREENHOUSE ?= true

TEST_FLAGS ?=

default: start

.PHONY: build
build: install
	@echo "---> Building"
	CGO_ENABLED=0 go build -o $(BIN_DIR)/api -installsuffix cgo ./cmd/serve
	CGO_ENABLED=0 go build -o $(BIN_DIR)/migrations -installsuffix cgo ./cmd/migrations

.PHONY: clean
clean:
	@echo "---> Cleaning"
	go clean
	rm -rf $(BIN_DIR) $(COVERAGE_PROFILE) ./tmp

.PHONY: db\:migrate
db\:migrate:
	@echo "---> Migrating within Docker"
	DATABASE_DEBUG=${DATABASE_DEBUG} docker-compose exec -e DATABASE_DEBUG=false api go run cmd/migrations/*.go migrate

.PHONY: db\:migrate\:create
db\:migrate\:create:
	@echo "---> Creating new migration"
	DATABASE_DEBUG=${DATABASE_DEBUG} docker-compose exec -e DATABASE_DEBUG=false api go run cmd/migrations/*.go create $(name)

.PHONY: db\:rollback
db\:rollback:
	@echo "---> Rolling back within Docker"
	DATABASE_DEBUG=${DATABASE_DEBUG} docker-compose exec -e DATABASE_DEBUG=false api go run cmd/migrations/*.go rollback

.PHONY: install
install:
	@echo "---> Installing dependencies"
	go mod download

.PHONY: lint
lint: $(BIN_DIR)/golangci-lint
	@echo "---> Linting"
	$(BIN_DIR)/golangci-lint run

.PHONY: psql
psql:
	@echo "---> Starting psql within Docker"
	DATABASE_DEBUG=${DATABASE_DEBUG} docker-compose exec -e PAGER=less postgres psql -U $(DATABASE_USER) $(DEVELOPMENT_DATABASE_NAME)

.PHONY: start
start:
	@echo "---> Starting docker-compose stack"
	touch ~/.psqlrc ~/.inputrc
	DATABASE_DEBUG=${DATABASE_DEBUG} COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker-compose up --build --remove-orphans

$(BIN_DIR)/golangci-lint:
	@echo "--> Installing linter"
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(BIN_DIR) v1.27.0

.PHONY: test
test:
	@echo "---> Testing"
	DATABASE_DEBUG=${DATABASE_DEBUG} docker-compose exec -e PAGER=less postgres createdb -U $(DATABASE_USER) -O $(DATABASE_USER) $(TEST_DATABASE_NAME) || true
	ENVIRONMENT=test go test -race ./pkg/... -coverprofile $(COVERAGE_PROFILE) $(TEST_FLAGS)
