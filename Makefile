.PHONY: start build

NOW = $(shell date -u '+%Y%m%d%I%M%S')

RELEASE_VERSION = v9.0.0

APP 			= ginadmin
SERVER_BIN  	= ${APP}
RELEASE_ROOT 	= release
RELEASE_SERVER 	= release/${APP}
GIT_COUNT 		= $(shell git rev-list --all --count)
GIT_HASH        = $(shell git rev-parse --short HEAD)
RELEASE_TAG     = $(RELEASE_VERSION).$(GIT_COUNT).$(GIT_HASH)

all: start

start:
	@go run -ldflags "-X github.com/LyricTian/gin-admin/v9/cmd.VERSION=$(RELEASE_TAG)" main.go start --configdir ./configs

build:
	@go build -ldflags "-w -s -X github.com/LyricTian/gin-admin/v9/cmd.VERSION=$(RELEASE_TAG)" -o $(SERVER_BIN)

wire:
	@wire gen ./internal/inject

swagger:
	@swag init --parseDependency --generalInfo ./main.go --output ./internal/swagger

clean:
	rm -rf data $(SERVER_BIN) pkg/x/cachex/tmp pkg/jwtauth/tmp