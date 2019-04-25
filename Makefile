.PHONY: start build

NOW = $(shell date -u '+%Y%m%d%I%M%S')

SERVER_BIN = "./cmd/ginadmin/ginadmin"
RELEASE_ROOT = "release"
RELEASE_SERVER = "release/ginadmin"

all: start

build:
	@go build -ldflags "-w -s" -o $(SERVER_BIN) ./cmd/ginadmin

start: 
	@go build -o $(SERVER_BIN) ./cmd/ginadmin
	$(SERVER_BIN) -c ./configs/ginadmin/config.toml -m ./configs/ginadmin/model.conf -swagger ./internal/app/ginadmin/swagger

swagger:
	swaggo -s ./internal/app/ginadmin/swagger.go -p . -o ./internal/app/ginadmin/swagger

test:
	@go test -cover -race ./...

clean:
	rm -rf data release $(SERVER_BIN) ./internal/app/ginadmin/test/data

pack: build
	rm -rf $(RELEASE_ROOT)
	mkdir -p $(RELEASE_SERVER)
	cp -r $(SERVER_BIN) configs $(RELEASE_SERVER)
	cd $(RELEASE_ROOT) && zip -r ginadmin.$(NOW).zip "ginadmin"
