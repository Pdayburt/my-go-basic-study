//go:build wireinject

package wire

import (
	"github.com/google/wire"
	"my-go-basic-study/wire/repository"
	"my-go-basic-study/wire/repository/dao"
)

func initRepository() *repository.UserRepository {
	wire.Build(dao.NewUser, repository.NewUserRepository, InitDB)
	return &repository.UserRepository{}
}
