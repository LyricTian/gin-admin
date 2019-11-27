.PHONY: start build

NOW = $(shell date -u '+%Y%m%d%I%M%S')

SERVER_BIN = "./cmd/server/server"
RELEASE_ROOT = "release"
RELEASE_SERVER = "release/server"

all: start

build:
	@go build -ldflags "-w -s" -o $(SERVER_BIN) ./cmd/server

start: 
	go run cmd/server/main.go -c ./configs/config.toml -m ./configs/model.conf -swagger ./docs/swagger -menu ./configs/menu.json

swagger:
	swag init -g ./internal/app/routers/swagger.go -o ./docs/swagger

test:
	@go test -cover -race ./...

clean:
	rm -rf data release $(SERVER_BIN) ./internal/app/test/data ./cmd/server/data

pack: build
	rm -rf $(RELEASE_ROOT)
	mkdir -p $(RELEASE_SERVER)
	cp -r $(SERVER_BIN) configs docs $(RELEASE_SERVER)
	cd $(RELEASE_ROOT) && zip -r server.$(NOW).zip "server"
