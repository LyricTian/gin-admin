.PHONY: start build

NOW = $(shell date -u '+%Y%m%d%I%M%S')

RELEASE_VERSION = v10.0.2

APP 			= ginadmin
SERVER_BIN  	= ${APP}
GIT_COUNT 		= $(shell git rev-list --all --count)
GIT_HASH        = $(shell git rev-parse --short HEAD)
RELEASE_TAG     = $(RELEASE_VERSION).$(GIT_COUNT).$(GIT_HASH)

CONFIG_DIR       = ./configs
CONFIG_FILES     = dev
STATIC_DIR       = ./build/dist
START_ARGS       = -d $(CONFIG_DIR) -c $(CONFIG_FILES) -s $(STATIC_DIR)

all: start

start:
	@go run -ldflags "-X main.VERSION=$(RELEASE_TAG)" main.go start $(START_ARGS)

build:
	@go build -ldflags "-w -s -X main.VERSION=$(RELEASE_TAG)" -o $(SERVER_BIN)

# go install github.com/google/wire/cmd/wire@latest
wire: installCli
	@wire gen ./internal/wirex

# go install github.com/swaggo/swag/cmd/swag@latest
swagger: installCli
	@swag init --parseDependency --generalInfo ./main.go --output ./internal/swagger

# https://github.com/OpenAPITools/openapi-generator
openapi:
	docker run --rm -v ${PWD}:/local openapitools/openapi-generator-cli generate -i /local/internal/swagger/swagger.yaml -g openapi -o /local/internal/swagger/v3

clean:
	rm -rf data $(SERVER_BIN)

serve: build
	./$(SERVER_BIN) start $(START_ARGS)

serve-d: build
	./$(SERVER_BIN) start $(START_ARGS) -d

stop:
	./$(SERVER_BIN) stop

installCli:
	go install github.com/gin-admin/gin-admin-cli/v10@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/google/wire/cmd/wire@latest

addStruct: installCli
	@if [ ! -n "${m}" ]; then \
		echo "补充参数： m=[模块名]"; \
		exit 1; \
	fi;
	@if [ ! -n "${s}" ]; then \
		echo "补充参数： s=[结构名(大驼峰命名)]"; \
		exit 1; \
	fi;
	@if [ ! -n "${c}" ]; then \
		echo "补充参数： c=[结构注释]"; \
		exit 1; \
	fi;
	gin-admin-cli gen -d . -m ${m} --structs ${s} --structs-comment '${s} ${c}'
	git add .

rmStruct: installCli
	@if [ ! -n "${m}" ]; then \
		echo "补充参数： m=[模块名]"; \
		exit 1; \
	fi;
	@if [ ! -n "${s}" ]; then \
		echo "补充参数： s=[结构名]"; \
		exit 1; \
	fi;
	gin-admin-cli rm -d . -m ${m} --structs ${s}

rmModule:
	@if [ ! -n "${m}" ]; then \
		echo "补充参数： m=[模块名]"; \
		exit 1; \
	fi;
	@if [ ! -d "internal/mods/${m}/api" ]; then \
		rm -rf internal/mods/${m}; \
	else \
		read -p "模块不为空，确认删除？[y/n]" input; \
		if [ "$${input}" = "y" ]; then \
			rm -rf internal/mods/${m}; \
		fi; \
	fi;
	read -p  "手动删除 internal/mods/mods.go 中的模块，然后回车执行 make wire:"
	make wire;

gitToDev:
	@git add . && \
	git commit ; \
	git pull origin main && \
	current_branch=`git symbolic-ref --short -q HEAD` && \
	echo "当前分支：$${current_branch}" && \
	git push --set-upstream origin $${current_branch} && \
	git checkout dev && \
	git pull origin dev && \
	git merge $${current_branch} && \
	git push origin dev && \
	git checkout $${current_branch};


gitToMain:
	@git add . && \
	git commit ; \
	git pull origin main && \
	current_branch=`git symbolic-ref --short -q HEAD` && \
	echo "当前分支：$${current_branch}" && \
	git push --set-upstream origin $${current_branch} && \
	git checkout main && \
	git pull origin main && \
	git merge $${current_branch} && \
	git push origin main && \
	git checkout $${current_branch};

