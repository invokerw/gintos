// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github/invokerw/gintos/demo/internal/biz"
	"github/invokerw/gintos/demo/internal/conf"
	"github/invokerw/gintos/demo/internal/data"
	"github/invokerw/gintos/demo/internal/initialize"
	"github/invokerw/gintos/demo/internal/router"
	"github/invokerw/gintos/demo/internal/service"
	"github/invokerw/gintos/log"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(server *conf.Server, confData *conf.Data, file *conf.File, logger log.Logger) (*App, func(), error) {
	dataData, cleanup, err := data.NewData(confData, logger)
	if err != nil {
		return nil, nil, err
	}
	userRepo := data.NewUserRepo(dataData, logger)
	roleRepo := data.NewRoleRepo(dataData, logger)
	adapter, err := data.NewCasbinAdapter(dataData)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	enforcer, err := biz.NewCasbinEnforcer(adapter, logger)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	initRet := initialize.DoInit(userRepo, roleRepo, enforcer, logger)
	greeterRepo := data.NewGreeterRepo(dataData, logger)
	greeterUsecase := biz.NewGreeterUsecase(greeterRepo, logger)
	greeterService := service.NewGreeterService(greeterUsecase, logger)
	oss, err := biz.NewOSS(file)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	userUsecase := biz.NewUserUsecase(userRepo, oss, logger)
	authService := service.NewAuthService(userUsecase, logger)
	roleUsecase := biz.NewRoleUsecase(roleRepo, logger)
	adminService := service.NewAdminService(userUsecase, roleUsecase, enforcer, logger)
	baseService := service.NewBaseService(userUsecase, logger)
	engine := router.NewGinHttpServer(server, greeterService, authService, adminService, baseService, enforcer, file, logger)
	app := newApp(initRet, engine)
	return app, func() {
		cleanup()
	}, nil
}
