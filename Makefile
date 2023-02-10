.PHONY: start build

NOW = $(shell date -u '+%Y%m%d%I%M%S')

RELEASE_VERSION = v8.1.0

APP 		= gin-admin
SERVER_BIN  	= ./cmd/${APP}/${APP}
RELEASE_ROOT 	= release
RELEASE_SERVER 	= release/${APP}
GIT_COUNT 	= $(shell git rev-list --all --count)
GIT_HASH        = $(shell git rev-parse --short HEAD)
RELEASE_TAG     = $(RELEASE_VERSION).$(GIT_COUNT).$(GIT_HASH)

all: start

build:
	@go build -ldflags "-w -s -X main.VERSION=$(RELEASE_TAG)" -o $(SERVER_BIN) ./cmd/${APP}

start:
	@go run -ldflags "-X main.VERSION=$(RELEASE_TAG)" ./cmd/${APP}/main.go web -c ./configs/config.toml -m ./configs/model.conf --menu ./configs/menu.yaml

swagger:
	@swag init -d ./cmd/${APP}/,internal/ --output ./internal/app/swagger

wire:
	@wire gen ./internal/app

test:
	cd ./internal/app/test && go test -v

clean:
	rm -rf data release $(SERVER_BIN) internal/app/test/data cmd/${APP}/data

pack: build
	rm -rf $(RELEASE_ROOT) && mkdir -p $(RELEASE_SERVER)
	cp -r $(SERVER_BIN) configs $(RELEASE_SERVER)
	cd $(RELEASE_ROOT) && tar -cvf $(APP).tar ${APP} && rm -rf ${APP}
