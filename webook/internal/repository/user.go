package repository

import (
	"context"
	"log"
	"my-go-basic-study/webook/internal/domain"
	"my-go-basic-study/webook/internal/repository/cache"
	"my-go-basic-study/webook/internal/repository/dao"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByID(ctx context.Context, id int64) (domain.User, error)
}

type userRepository struct {
	dao   dao.UserDao
	cache cache.UserCache
}

func NewUserRepository(dao dao.UserDao, cache cache.UserCache) UserRepository {
	return &userRepository{
		dao:   dao,
		cache: cache,
	}
}

func (r *userRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	user, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       user.Id,
		Email:    user.Email,
		Password: user.Password,
	}, nil

}

// FindByID

func (r *userRepository) FindByID(ctx context.Context, id int64) (domain.User, error) {
	u, err := r.cache.Get(ctx, id)
	//缓存里 有数据
	if err == nil {
		return u, nil
	}
	/*//缓存里 没数据 去数据库加载
	if errors.Is(err, cache.ErrKeyNotExist) {
		return domain.User{}, err
	}*/
	user, err := r.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	uR := domain.User{
		Id:       user.Id,
		Email:    user.Email,
		Password: user.Password,
	}
	err = r.cache.Set(ctx, uR)
	if err != nil {
		log.Println(err)
	}

	return uR, nil

}
