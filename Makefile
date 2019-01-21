.PHONY: start build-server build-web

NOW = $(shell date -u '+%Y%m%d%I%M%S')

SERVER_BIN = "./cmd/server/server"
RELEASE_ROOT = "release"
RELEASE_SERVER = "release/ginadmin"

all: start

build-server:
	@go build -ldflags "-w -s" -o $(SERVER_BIN) ./cmd/server

build-web:
	cd web && yarn && yarn run build

build: build-server build-web

start: build-server build-web
	$(SERVER_BIN) -c ./config/config.toml -m ./config/model.conf -www ./web/dist

start-dev-server: 
	@go build -o $(SERVER_BIN) ./cmd/server
	$(SERVER_BIN) -c ./config/config.toml -m ./config/model.conf -swagger ./src/web/swagger

start-dev-web:
	cd web && yarn && yarn start

swagger:
	swaggo -s ./src/web/swagger.go -p ./src -o ./src/web/swagger

test:
	@go test -cover -race ./...

pack: build-server build-web
	rm -rf $(RELEASE_ROOT)
	mkdir -p $(RELEASE_SERVER)
	cp -r $(SERVER_BIN) config web/dist $(RELEASE_SERVER)
	cd $(RELEASE_ROOT) && zip -r ginadmin.$(NOW).zip "ginadmin"
