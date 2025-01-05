//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"my-go-basic-study/webook/internal/repository"
	"my-go-basic-study/webook/internal/repository/cache"
	"my-go-basic-study/webook/internal/repository/dao"
	"my-go-basic-study/webook/internal/service"
	"my-go-basic-study/webook/internal/web"
	"my-go-basic-study/webook/ioc"
)

func InitWebService() *gin.Engine {
	wire.Build(

		ioc.InitDb, ioc.InitRedis,

		dao.NewUserDao, cache.NewUserCache,

		repository.NewUserRepository,

		service.NewUserService,

		web.NewUserHandler,
		ioc.InitGin, ioc.InitMiddleware,
	)

	return new(gin.Engine)
}
