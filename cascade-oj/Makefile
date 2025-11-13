GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)

.PHONY: init
# init env
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
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
	go run entgo.io/ent/cmd/ent generate ./ent/schema

# show help
help:
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@printf "\033[36m%-20s\033[0m %s\n" "init" "init env";
	@echo ''
	@echo 'Current OS:'
	@echo $(GOHOSTOS)

.DEFAULT_GOAL := help