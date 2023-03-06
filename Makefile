.PHONY: start build

NOW = $(shell date -u '+%Y%m%d%I%M%S')

RELEASE_VERSION = v10.0.0-beta

APP 			= ginadmin
SERVER_BIN  	= ${APP}
GIT_COUNT 		= $(shell git rev-list --all --count)
GIT_HASH        = $(shell git rev-parse --short HEAD)
RELEASE_TAG     = $(RELEASE_VERSION).$(GIT_COUNT).$(GIT_HASH)

CONFIG_DIR      = ./configs
STATIC_DIR      = ./build/dist

all: start

start:
	@go run -ldflags "-X main.VERSION=$(RELEASE_TAG)" main.go start -c $(CONFIG_DIR) -s $(STATIC_DIR)

build:
	@go build -ldflags "-w -s -X main.VERSION=$(RELEASE_TAG)" -o $(SERVER_BIN)

# go install github.com/google/wire/cmd/wire@latest
wire:
	@wire gen ./internal/library/wirex

# go install github.com/swaggo/swag/cmd/swag@latest
swagger:
	@swag init --parseDependency --generalInfo ./main.go --output ./internal/swagger

# Dependency: https://github.com/OpenAPITools/openapi-generator
# brew install openapi-generator
openapi:
	@openapi-generator generate -i ./internal/swagger/swagger.json -o ./internal/swagger/v3 -g openapi --minimal-update

clean:
	rm -rf data $(SERVER_BIN)

serve: build
	./$(SERVER_BIN) start -c $(CONFIG_DIR) -s $(STATIC_DIR)

serve-d: build
	./$(SERVER_BIN) start -c $(CONFIG_DIR) -s $(STATIC_DIR) -d

stop:
	./$(SERVER_BIN) stop