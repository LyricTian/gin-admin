.PHONY: start build

NOW = $(shell date -u '+%Y%m%d%I%M%S')

SERVER_BIN = "./cmd/server/server"
RELEASE_ROOT = "release"
RELEASE_SERVER = "release/server"

all: start

build:
	@go build -ldflags "-w -s" -o $(SERVER_BIN) ./cmd/server

start: 
	@go build -o $(SERVER_BIN) ./cmd/server
	$(SERVER_BIN) -c ./configs/config.toml -m ./configs/model.conf -swagger ./internal/app/swagger

swagger:
	swaggo -s ./internal/app/routers/api/swagger.go -p . -o ./internal/app/swagger

test:
	@go test -cover -race ./...

clean:
	rm -rf data release $(SERVER_BIN) ./internal/app/test/data ./cmd/server/data

pack: build
	rm -rf $(RELEASE_ROOT)
	mkdir -p $(RELEASE_SERVER)
	cp -r $(SERVER_BIN) configs $(RELEASE_SERVER)
	cd $(RELEASE_ROOT) && zip -r server.$(NOW).zip "server"
