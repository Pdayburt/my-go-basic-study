//go:build wireinject

package startup

import (
	"github.com/google/wire"
	"my-go-basic-study/webook/interactive/repository"
	"my-go-basic-study/webook/interactive/repository/cache"
	"my-go-basic-study/webook/interactive/repository/dao"
	"my-go-basic-study/webook/interactive/service"
)

var thirdProvider = wire.NewSet(InitRedis, InitDb)

var interactiveSvcProvider = wire.NewSet(
	service.NewInteractiveService,
	repository.NewCachedInteractiveRepository,
	dao.NewGORMInteractiveDAO,
	cache.NewInteractiveRedisCache,
)

func InitInteractiveService() service.InteractiveService {
	wire.Build(thirdProvider, interactiveSvcProvider)
	return service.NewInteractiveService(nil)
}
