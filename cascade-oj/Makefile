GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
GATEWAYDIR:=app/gateway
SERVICESDIR:=app/services
SUBDIRS:=public user judge admin
SUBPROJECTS:=$(addprefix $(SERVICESDIR)/,$(SUBDIRS)) $(GATEWAYDIR)

ifeq ($(GOHOSTOS), windows)
# 	# the `find.exe` is different from `find` in bash/shell.
# 	# to see https://docs.microsoft.com/en-us/windows-server/administration/windows-commands/find.
# 	# changed to use git-bash.exe to run find cli or other cli friendly, caused of every developer has a Git.
# 	# Git_Bash= $(subst cmd\,bin\bash.exe,$(dir $(shell where git)))
	Git_Bash=$(subst \,/,$(subst cmd\,bin\bash.exe,$(dir $(shell where git))))
	INTERNAL_PROTO_FILES=$(shell $(Git_Bash) -c "find app -name *.proto -type f -print0 | xargs -0 echo")
	API_PROTO_FILES=$(shell $(Git_Bash) -c "find api -name *.proto -type f -print0 | xargs -0 echo")
else
# 	# Unix shell
# 	# if you use Msys2, Cygwin or unix-like terminal in windows, you can use the below code directly.
	INTERNAL_PROTO_FILES=$(shell find app -name *.proto -type f -print0 | xargs -0 echo)
	API_PROTO_FILES=$(shell find api -name *.proto -type f -print0 | xargs -0 echo)
endif

.PHONY: init
# init env
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-errors@latest
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/bufbuild/buf/cmd/buf@v1.59.0

.PHONY: config
# generate internal proto at sub project internal/conf
config:
	@for dir in $(SUBPROJECTS); do  \
		$(MAKE) -C $$dir config; \
	done

.PHONY: api
# generate protobuf api go code: make api [module]
api:
	@if [ -z "$(word 2,$(MAKECMDGOALS))" ]; then \
		cd api && \
		buf generate -o cascade; \
	else \
		MODULE=$(word 2,$(MAKECMDGOALS)); \
		case "$$MODULE" in \
			auth) \
				echo "Generating API code for auth module..."; \
				cd api && buf generate -o cascade cascade/public/auth/v1/auth.proto; \
				;; \
			*) \
				echo "Error: Unknown module '$$MODULE'"; \
				echo "Available modules: auth"; \
				exit 1; \
				;; \
		esac; \
	fi

.PHONY: ent-init
# init ent
ent-init:
	go get entgo.io/ent/cmd/ent
	go mod download entgo.io/ent
	go get entgo.io/ent/dialect/sql/schema@v0.14.5
	go get entgo.io/ent/dialect@v0.14.5
	go get entgo.io/ent/dialect/sql/schema@v0.14.5
	go get entgo.io/ent/entc/load@v0.14.5
	go get entgo.io/ent/cmd/internal/base@v0.14.5
	go get entgo.io/ent/entc/gen@v0.14.5
	go get entgo.io/ent/cmd/internal/printer@v0.14.5
	go get entgo.io/ent/cmd/ent@v0.14.5
	go get github.com/go-sql-driver/mysql@v1.9.3
	go mod download github.com/go-sql-driver/mysql

.PHONY: ent
# generate ent code
ent:
	go run entgo.io/ent/cmd/ent generate ./ent/schema

.PHONY: db-init
# initialize database
# ensure you have started docker-compose service before running this command
db-init:
	go run ./table_creator.go

.PHONY: wire
# generate wire code
wire:
	@for dir in $(SUBPROJECTS); do  \
		$(MAKE) -C $$dir wire; \
	done

.PHONY: generate
# generate
generate:
	go mod tidy
	go get github.com/google/wire/cmd/wire@latest
	go generate ./...

# show help
help:
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@printf "\033[36m%-20s\033[0m %s\n" "init" "init env";
	@printf "\033[36m%-20s\033[0m %s\n" "config" "generate internal proto";
	@printf "\033[36m%-20s\033[0m %s\n" "api [module]" "generate api from proto";
	@printf "\033[36m%-20s\033[0m %s\n" "ent-init" "initialize ent environment";
	@printf "\033[36m%-20s\033[0m %s\n" "ent" "generate ent schema";
	@printf "\033[36m%-20s\033[0m %s\n" "db-init" "initialize database, docker-compose service needed";
	@printf "\033[36m%-20s\033[0m %s\n" "wire" "generate wire code";
	@printf "\033[36m%-20s\033[0m %s\n" "generate" "go mod tidy && go generate";
	@echo ''
	@echo 'Sub projects:'
	@for dir in $(SUBPROJECTS); do  \
		echo $$dir; \
	done
	@echo ''
	@echo 'Current OS:'
	@echo $(GOHOSTOS)

.DEFAULT_GOAL := help