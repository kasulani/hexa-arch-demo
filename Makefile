#-----------------------------------------------------------------------------------------------------------------------
# Variables
# (https://www.gnu.org/software/make/manual/html_node/Using-Variables.html#Using-Variables)
#-----------------------------------------------------------------------------------------------------------------------

# Colorful outputs
RESET=\e[0m
COLOR_DEFAULT=\e[39m
COLOR_WHITE=\e[97m
COLOR_BLUE=\e[34m
TEXT_INVERSE=\e[7m

BUILD_DIR ?= $(CURDIR)/build
BINARY_CLI=shortener
BINARY_CLI_SRC=$(CURDIR)/cmd/shortener
BDD_TEST=$(CURDIR)/internal/app

GO_LINKER_FLAGS=-ldflags="-s -w"
SRC_DIRS=internal

# Db migration tool
BINARY_MIGRATE_TOOL=/bin/migrate
BINARY_MIGRATE_TOOL_VERSION=v4.4.0
BINARY_MIGRATE_TOOL_SHA256SUM=319216ab5662704cc081fb9101c325b5f8187e3af198bce2705d26b7bd42ae92
BINARY_MIGRATE_TOOL_URL=https://github.com/golang-migrate/migrate/releases/download/${BINARY_MIGRATE_TOOL_VERSION}/migrate.linux-amd64.tar.gz

#-----------------------------------------------------------------------------------------------------------------------
# Dependencies
#-----------------------------------------------------------------------------------------------------------------------
.PHONY: deps deps-dev deps-clean deps-ci

deps:
	${call print, "Installing dependencies"}
	${call go, mod vendor -v}

deps-dev:
	${call print, "Installing CompileDaemon"}
	${call go, get -v -u github.com/githubnemo/CompileDaemon}

	${call print, "Installing Linters"}
	${call go, get -v github.com/golangci/golangci-lint/cmd/golangci-lint@v1.36.0}

	${call print, "Installing Godog"}
	${call go, get -v -u github.com/cucumber/godog/cmd/godog@v0.11.0}

	${call print, "Installing gitleaks"}
	@wget https://github.com/zricethezav/gitleaks/releases/download/v4.3.1/gitleaks-linux-amd64 -O $(GOPATH)/bin/gitleaks
	@chmod +x $(GOPATH)/bin/gitleaks

	${call print, "Installing Gherkin Formatter"}
	@wget https://github.com/antham/ghokin/releases/download/v1.6.1/ghokin_linux_amd64 -O $(GOPATH)/bin/ghokin
	@chmod +x $(GOPATH)/bin/ghokin

deps-clean:
	${call print, "Cleaning dependencies"}
	@rm -rfv vendor
#-----------------------------------------------------------------------------------------------------------------------
# Building
#-----------------------------------------------------------------------------------------------------------------------
.PHONY: build-all build-cli build-migration-tool

build-all: build-cli build-migration-tool

build-cli:
	${call print, "Building cli binary"}
	${call go, build -v -o ${BUILD_DIR}/${BINARY_CLI} ${GO_LINKER_FLAGS} ${BINARY_CLI_SRC}}

build-migration-tool:
	${call print, "Downloading migration tool binary"}
	curl -L ${BINARY_MIGRATE_TOOL_URL} | tar -xvz
	echo "${BINARY_MIGRATE_TOOL_SHA256SUM} migrate.linux-amd64" | sha256sum --check
	mv migrate.linux-amd64 ${BINARY_MIGRATE_TOOL}
	chmod +x ${BINARY_MIGRATE_TOOL}
#-----------------------------------------------------------------------------------------------------------------------
# Migrations
#-----------------------------------------------------------------------------------------------------------------------
.PHONY: migrations

migrations:
	${call print, "Migrating database"}
	@migrate -database="${DATABASE_DSN}" -path=migrations up
#-----------------------------------------------------------------------------------------------------------------------
# Installing
#-----------------------------------------------------------------------------------------------------------------------
.PHONY: install-all install-cli

install-all: install-cli

install-cli:
	${call print, "Installing cli binary"}
	${call go, install -v ${BINARY_CLI_SRC}}
#-----------------------------------------------------------------------------------------------------------------------
# Development
#-----------------------------------------------------------------------------------------------------------------------
.PHONY: dev-up dev-rebuild dev-stop dev-rm dev-ssh lint

dev-up:
	${call print, "Starting development containers"}
	@docker-compose up -d dev

dev-rebuild:
	${call print, "Rebuilding development containers"}
	@docker-compose down -v
	@docker-compose build --no-cache
	@docker rmi -f $$(docker images | grep '<none>' | awk '{print $$3}')

dev-stop:
	${call print, "Stopping development containers"}
	@docker-compose stop

dev-rm:
	${call print, "Deleting containers"}
	@docker-compose rm -f

dev-ssh:
	@docker-compose exec dev bash

lint:
	${call print, "linting"}
	@golangci-lint -v run
#-----------------------------------------------------------------------------------------------------------------------
# Testing
#-----------------------------------------------------------------------------------------------------------------------
.PHONY: test-unit test-unit-coverage test-behavior

test-all: test-unit

test-unit:
	${call print, "Running unit tests"}
	${call go, test -v -race -tags unit $$(for d in $(SRC_DIRS); do echo ./$$d/...; done)}

test-unit-coverage:
	${call print, "Running unit tests with coverage"}
	${call go, test -v -race -tags unit $$(for d in $(SRC_DIRS); do echo ./$$d/...; done) -cover -coverprofile=coverage_unit.txt -covermode=atomic}
#-----------------------------------------------------------------------------------------------------------------------
# Helpers
#-----------------------------------------------------------------------------------------------------------------------
define print
	@printf "${TEXT_INVERSE}${COLOR_WHITE}%3s${COLOR_BLUE} %-75s ${COLOR_WHITE} ${COLOR_DEFAULT}${RESET}\n" " " $(1)
endef

# Run on host
define go
	@go $(1)
endef

define godog
	@godog $(1)
endef