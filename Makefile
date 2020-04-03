.PHONY: start build

NOW = $(shell date -u '+%Y%m%d%I%M%S')

SERVER_BIN = "./cmd/gin-admin/gin-admin"
RELEASE_ROOT = "release"
RELEASE_SERVER = "release/gin-admin"

all: start

build:
	@go build -ldflags "-w -s" -o $(SERVER_BIN) ./cmd/gin-admin

start: 
	go run cmd/gin-admin/main.go web -c ./configs/config.toml -m ./configs/model.conf --menu ./configs/menu.yaml

swagger:
	swag init --generalInfo ./internal/app/swagger/swagger.go --output ./internal/app/swagger

wire:
	wire gen ./internal/app/initialize

test:
	@go test -cover -race ./...

clean:
	rm -rf data release $(SERVER_BIN) ./internal/app/test/data ./cmd/gin-admin/data

pack: build
	rm -rf $(RELEASE_ROOT)
	mkdir -p $(RELEASE_SERVER)
	cp -r $(SERVER_BIN) configs docs $(RELEASE_SERVER)
	cd $(RELEASE_ROOT) && zip -r server.$(NOW).zip "server"
