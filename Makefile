.PHONY: start build build-web

NOW = $(shell date -u '+%Y%m%d%I%M%S')

SERVER_BIN = "./cmd/server/server"
RELEASE_ROOT = "release"
RELEASE_GINADMIN = "release/ginadmin"

all: start

start: 
	rm $(SERVER_BIN)
	@go build -o $(SERVER_BIN) ./cmd/server
	$(SERVER_BIN) -c ./config/config.toml -m ./config/model.conf

web:
	cd web && yarn && yarn start

test:
	go test -cover -race ./...

build:
	rm $(SERVER_BIN)
	@go build -ldflags "-w -s" -o $(SERVER_BIN) ./cmd/server

build-web:
	rm -rf ./web/dist
	cd web && yarn && yarn run build

pack: build build-web
	rm -rf $(RELEASE_ROOT)
	mkdir -p $(RELEASE_GINADMIN)
	cp -r $(SERVER_BIN) config script web/dist $(RELEASE_GINADMIN)
	cd $(RELEASE_ROOT) && zip -r ginadmin.$(NOW).zip "ginadmin"
