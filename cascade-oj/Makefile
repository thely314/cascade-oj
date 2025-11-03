GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
GATEWAYDIR:=app/gateway
SERVICESDIR:=app/services
SUBDIRS:=public user judge admin
SUBPROJECTS:=$(addprefix $(SERVICESDIR)/,$(SUBDIRS)) $(GATEWAYDIR)

# ifeq ($(GOHOSTOS), windows)
# 	# the `find.exe` is different from `find` in bash/shell.
# 	# to see https://docs.microsoft.com/en-us/windows-server/administration/windows-commands/find.
# 	# changed to use git-bash.exe to run find cli or other cli friendly, caused of every developer has a Git.
# 	# Git_Bash= $(subst cmd\,bin\bash.exe,$(dir $(shell where git)))
# 	Git_Bash=$(subst \,/,$(subst cmd\,bin\bash.exe,$(dir $(shell where git))))
# 	INTERNAL_PROTO_FILES=$(shell $(Git_Bash) -c "find app -name *.proto -type f -print0 | xargs -0 echo")
# 	API_PROTO_FILES=$(shell $(Git_Bash) -c "find api -name *.proto -type f -print0 | xargs -0 echo")
# else
# 	# Unix shell
# 	# if you use Msys2, Cygwin or unix-like terminal in windows, you can use the below code directly.
	INTERNAL_PROTO_FILES=$(shell find app -name *.proto -type f -print0 | xargs -0 echo)
	API_PROTO_FILES=$(shell find api -name *.proto -type f -print0 | xargs -0 echo)
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
# generate internal proto at sub project internal/conf
config:
	for dir in $(SUBPROJECTS); do  \
		$(MAKE) -C $$dir config; \
	done

# show help
help:
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@printf "\033[36m%-20s\033[0m %s\n" "init" "init env";
	@printf "\033[36m%-20s\033[0m %s\n" "config" "generate internal proto";
	@echo ''
	@echo 'Sub projects:'
	@for dir in $(SUBPROJECTS); do  \
		echo $$dir; \
	done
	@echo ''
	@echo 'Current OS:'
	@echo $(GOHOSTOS)

.DEFAULT_GOAL := help