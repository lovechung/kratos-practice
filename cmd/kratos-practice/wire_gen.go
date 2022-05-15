// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"kratos-practice/internal/biz"
	"kratos-practice/internal/conf"
	"kratos-practice/internal/data"
	"kratos-practice/internal/server"
	"kratos-practice/internal/service"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confServer *conf.Server, confData *conf.Data, logger log.Logger) (*kratos.App, func(), error) {
	client := data.NewDB(confData, logger)
	redisClient := data.NewRedis(confData, logger)
	dataData, cleanup, err := data.NewData(client, redisClient, logger)
	if err != nil {
		return nil, nil, err
	}
	userRepo := data.NewUserRepo(dataData, logger)
	transaction := data.NewTransaction(dataData)
	userUseCase := biz.NewUserUseCase(userRepo, transaction, logger)
	userService := service.NewUserService(userUseCase, logger)
	carRepo := data.NewCarRepo(dataData, logger)
	carUseCase := biz.NewCarUseCase(carRepo, transaction, logger)
	carService := service.NewCarService(carUseCase, logger)
	httpServer := server.NewHTTPServer(confServer, userService, carService, logger)
	grpcServer := server.NewGRPCServer(confServer, userService, carService, logger)
	app := newApp(logger, httpServer, grpcServer)
	return app, func() {
		cleanup()
	}, nil
}