GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)

# ifeq ($(GOHOSTOS), windows)
# 	# the `find.exe` is different from `find` in bash/shell.
# 	# to see https://docs.microsoft.com/en-us/windows-server/administration/windows-commands/find.
# 	# changed to use git-bash.exe to run find cli or other cli friendly, caused of every developer has a Git.
# 	# Git_Bash= $(subst cmd\,bin\bash.exe,$(dir $(shell where git)))
# 	Git_Bash=$(subst \,/,$(subst cmd\,bin\bash.exe,$(dir $(shell where git))))
# 	INTERNAL_PROTO_FILES=$(shell $(Git_Bash) -c "find internal -name *.proto -type f -print0 | xargs -0 echo")
# 	WIRE_DIR=$(shell $(Git_Bash) -c "find . -name wire.go -type f -exec dirname {} \; | xargs -0 echo")
# else
# 	# Unix shell
# 	# if you use Msys2, Cygwin or unix-like terminal in windows, you can use the below code directly.
	INTERNAL_PROTO_FILES=$(shell find internal -name *.proto -type f -print0 | xargs -0 echo)
	WIRE_DIR=$(shell find . -name wire.go -type f -exec dirname {} \; | xargs -0 echo)
# endif

.PHONY: init
# init env
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	go install github.com/google/wire/cmd/wire@latest

.PHONY: config
# generate internal proto
config:
	echo "Using $(INTERNAL_PROTO_FILES)"
	protoc --proto_path=./internal \
	       --proto_path=./third_party \
 	       --go_out=./internal \
	       $(INTERNAL_PROTO_FILES)

.PHONY: api
# generate api proto
api:
	protoc --proto_path=./api \
	       --proto_path=./third_party \
 	       --go_out=paths=source_relative:./api \
 	       --go-http_out=paths=source_relative:./api \
 	       --go-grpc_out=paths=source_relative:./api \
	       --openapi_out=fq_schema_naming=true,default_response=false:. \
	       $(API_PROTO_FILES)

.PHONY: build
# build
build:
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...

.PHONY: wire
# generate wire code
wire:
	cd $(WIRE_DIR) && go run github.com/google/wire/cmd/wire

.PHONY: generate
# generate
generate:
	go mod tidy
	go get github.com/google/wire/cmd/wire@latest
	go generate ./...

.PHONY: all
# generate all
all:
	make api;
	make config;
	make generate;

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
