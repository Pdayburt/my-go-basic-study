.PHONY: mock
mock:
	@mockgen -source=webook/internal/service/user.go -destination=webook/internal/service/mock/user.mock.go -package=svcmock
	@mockgen -source=webook/internal/repository/user.go -destination=webook/internal/repository/mock/user.mock.go -package=repomock
	@mockgen -source=webook/internal/repository/dao/user.go -destination=webook/internal/repository/dao/mock/user.mock.go -package=daomock
	@mockgen -source=webook/internal/repository/cache/user.go -destination=webook/internal/repository/cache/mock/user.mock.go -package=cachemock
	@mockgen -destination=webook/internal/repository/cache/redismock/cmdable.mock.go -package=redismock github.com/redis/go-redis/v9 Cmdable

	@go mod tidy