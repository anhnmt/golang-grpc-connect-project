// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/xdorro/golang-grpc-base-project/internal/interceptor"
	"github.com/xdorro/golang-grpc-base-project/internal/module/auth/biz"
	"github.com/xdorro/golang-grpc-base-project/internal/module/auth/service"
	"github.com/xdorro/golang-grpc-base-project/internal/module/user/biz"
	"github.com/xdorro/golang-grpc-base-project/internal/module/user/repo"
	"github.com/xdorro/golang-grpc-base-project/internal/module/user/service"
	"github.com/xdorro/golang-grpc-base-project/internal/server"
	"github.com/xdorro/golang-grpc-base-project/internal/service"
	"github.com/xdorro/golang-grpc-base-project/pkg/casbin"
	"github.com/xdorro/golang-grpc-base-project/pkg/repo"
	"net/http"
)

// Injectors from wire.go:

func initServer() server.IServer {
	serveMux := http.NewServeMux()
	iRepo := repo.NewRepo()
	option := &casbin.Option{
		Repo: iRepo,
	}
	iCasbin := casbin.NewCasbin(option)
	interceptorOption := &interceptor.Option{
		Casbin: iCasbin,
	}
	iInterceptor := interceptor.NewInterceptor(interceptorOption)
	userrepoOption := &userrepo.Option{
		Repo: iRepo,
	}
	userrepoIRepo := userrepo.NewRepo(userrepoOption)
	userbizOption := &userbiz.Option{
		UserRepo: userrepoIRepo,
	}
	iUserBiz := userbiz.NewBiz(userbizOption)
	userserviceOption := &userservice.Option{
		UserBiz: iUserBiz,
	}
	iUserService := userservice.NewService(userserviceOption)
	authbizOption := &authbiz.Option{
		UserRepo: userrepoIRepo,
	}
	iAuthBiz := authbiz.NewBiz(authbizOption)
	authserviceOption := &authservice.Option{
		AuthBiz: iAuthBiz,
	}
	iAuthService := authservice.NewService(authserviceOption)
	serviceOption := &service.Option{
		Mux:         serveMux,
		Interceptor: iInterceptor,
		Repo:        iRepo,
		UserService: iUserService,
		AuthService: iAuthService,
	}
	iService := service.NewService(serviceOption)
	serverOption := &server.Option{
		Mux:     serveMux,
		Service: iService,
	}
	iServer := server.NewServer(serverOption)
	return iServer
}
