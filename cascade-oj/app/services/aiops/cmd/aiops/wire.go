//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"net/http"

	"cascade-oj/app/services/aiops/internal/biz"
	"cascade-oj/app/services/aiops/internal/conf"
	"cascade-oj/app/services/aiops/internal/data"
	"cascade-oj/app/services/aiops/internal/server"
	"cascade-oj/app/services/aiops/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(
	*conf.Server,
	*conf.Data,
	*conf.LLM,
	*conf.Prometheus,
	*conf.Polling,
	log.Logger,
	http.Handler,
) (*kratos.App, func(), error) {
	panic(wire.Build(
		data.ProviderSet,
		biz.ProviderSet,
		service.ProviderSet,
		server.ProviderSet,
		newApp,
	))
}
