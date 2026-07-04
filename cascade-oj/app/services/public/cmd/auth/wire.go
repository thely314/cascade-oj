//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"net/http"

	"cascade-oj/app/services/public/internal/biz"
	"cascade-oj/app/services/public/internal/conf"
	"cascade-oj/app/services/public/internal/data"
	"cascade-oj/app/services/public/internal/server"
	"cascade-oj/app/services/public/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, *conf.Jwt, log.Logger, http.Handler) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
