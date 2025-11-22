//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

// import (
// "cascade-oj/app/gateway/internal/biz"
// "cascade-oj/app/gateway/internal/conf"
// "cascade-oj/app/gateway/internal/data"
// "cascade-oj/app/gateway/internal/server"
// "cascade-oj/app/gateway/internal/service"

// "github.com/go-kratos/kratos/v2"
// "github.com/go-kratos/kratos/v2/log"
// "github.com/google/wire"
// )

// wireApp init kratos application.
// func wireApp(*conf.Server, *conf.Data, *conf.Jwt, log.Logger) (*kratos.App, func(), error) {
// 	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
// }
