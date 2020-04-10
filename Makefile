.PHONY: start build

NOW = $(shell date -u '+%Y%m%d%I%M%S')

APP = gin-admin
SERVER_BIN = ./cmd/${APP}/${APP}
RELEASE_ROOT = release
RELEASE_SERVER = release/${APP}

all: start

build:
	@go build -ldflags "-w -s" -o $(SERVER_BIN) ./cmd/${APP}

start: 
	go run cmd/${APP}/main.go web -c ./configs/config.toml -m ./configs/model.conf --menu ./configs/menu.yaml

swagger:
	swag init --generalInfo ./internal/app/swagger.go --output ./internal/app/swagger

wire:
	wire gen ./internal/app/initialize

test:
	@go test -v ./internal/app/test

clean:
	rm -rf data release $(SERVER_BIN) ./internal/app/test/data ./cmd/${APP}/data

pack: build
	rm -rf $(RELEASE_ROOT)
	mkdir -p $(RELEASE_SERVER)
	cp -r $(SERVER_BIN) configs docs $(RELEASE_SERVER)
	cd $(RELEASE_ROOT) && zip -r ${APP}.$(NOW).zip ${APP}
	rm -rf $(RELEASE_SERVER)
