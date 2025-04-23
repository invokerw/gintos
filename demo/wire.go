//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github/invokerw/gintos/demo/internal/biz"
	"github/invokerw/gintos/demo/internal/conf"
	"github/invokerw/gintos/demo/internal/data"
	"github/invokerw/gintos/demo/internal/initialize"
	"github/invokerw/gintos/demo/internal/router"
	"github/invokerw/gintos/demo/internal/service"
	"github/invokerw/gintos/log"

	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, log.Logger) (*App, func(), error) {
	panic(wire.Build(
		data.ProviderSet,
		biz.ProviderSet,
		initialize.ProviderSet,
		service.ProviderSet,
		router.ProviderSet,
		newApp))
}
