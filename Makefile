.PHONY: start build

NOW = $(shell date -u '+%Y%m%d%I%M%S')

RELEASE_VERSION = v9.0.1

APP 			= ginadmin
SERVER_BIN  	= ${APP}
RELEASE_ROOT 	= release
RELEASE_SERVER 	= release/${APP}
GIT_COUNT 		= $(shell git rev-list --all --count)
GIT_HASH        = $(shell git rev-parse --short HEAD)
RELEASE_TAG     = $(RELEASE_VERSION).$(GIT_COUNT).$(GIT_HASH)

all: start

start:
	@go run -ldflags "-X main.VERSION=$(RELEASE_TAG)" main.go start --configdir ./configs

build:
	@go build -ldflags "-w -s -X main.VERSION=$(RELEASE_TAG)" -o $(SERVER_BIN)

wire:
	@wire gen ./internal/inject

swagger:
	@swag init --parseDependency --generalInfo ./main.go --output ./internal/swagger

# Dependency: https://github.com/OpenAPITools/openapi-generator
# brew install openapi-generator
openapi:
	@openapi-generator generate -i ./internal/swagger/swagger.json -o ./internal/swagger/v3 -g openapi --minimal-update && cp ./internal/swagger/v3/openapi.json ./configs/openapi.json

clean:
	rm -rf data $(SERVER_BIN) pkg/x/cachex/tmp pkg/jwtauth/tmp