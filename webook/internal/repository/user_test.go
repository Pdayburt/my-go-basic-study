package repository

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"my-go-basic-study/webook/internal/domain"
	"my-go-basic-study/webook/internal/repository/cache"
	cachemock "my-go-basic-study/webook/internal/repository/cache/mock"
	"my-go-basic-study/webook/internal/repository/dao"
	daomock "my-go-basic-study/webook/internal/repository/dao/mock"
	"testing"
)

func Test_userRepository_FindByID(t *testing.T) {

	testCase := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache)
		ctx      context.Context
		id       int64
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "缓存未命中，查询成功",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {

				mockUserCache := cachemock.NewMockUserCache(ctrl)
				//Get(ctx context.Context, id int64) (domain.User, error)
				mockUserCache.EXPECT().Get(context.Background(), int64(123)).
					Return(domain.User{}, cache.ErrKeyNotExist)
				//	Set(ctx context.Context, user domain.User) error
				mockUserCache.EXPECT().Set(gomock.Any(), domain.User{
					Id:       123,
					Email:    "test@test.com",
					Password: "test",
				}).Return(nil)

				mockUserDao := daomock.NewMockUserDao(ctrl)
				//	FindById(ctx context.Context, id int64) (User, error)
				mockUserDao.EXPECT().FindById(context.Background(), int64(123)).
					Return(dao.User{
						Id:       123,
						Email:    "test@test.com",
						Password: "test",
					}, nil)
				return mockUserDao, mockUserCache
			},
			ctx: context.Background(),
			id:  123,
			wantUser: domain.User{
				Id:       123,
				Email:    "test@test.com",
				Password: "test",
			},
			wantErr: nil,
		},
		{
			name: "缓存命中",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {

				mockUserCache := cachemock.NewMockUserCache(ctrl)
				//Get(ctx context.Context, id int64) (domain.User, error)
				mockUserCache.EXPECT().Get(context.Background(), int64(123)).
					Return(domain.User{
						Id:       123,
						Email:    "test@test.com",
						Password: "test",
					}, nil)

				mockUserDao := daomock.NewMockUserDao(ctrl)
				return mockUserDao, mockUserCache
			},
			ctx: context.Background(),
			id:  123,
			wantUser: domain.User{
				Id:       123,
				Email:    "test@test.com",
				Password: "test",
			},
			wantErr: nil,
		},
		{
			name: "缓存没有，数据库查询失败",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {

				mockUserCache := cachemock.NewMockUserCache(ctrl)
				//Get(ctx context.Context, id int64) (domain.User, error)
				mockUserCache.EXPECT().Get(context.Background(), int64(123)).
					Return(domain.User{}, cache.ErrKeyNotExist)

				mockUserDao := daomock.NewMockUserDao(ctrl)
				//	FindById(ctx context.Context, id int64) (User, error)
				mockUserDao.EXPECT().FindById(context.Background(), int64(123)).
					Return(dao.User{}, errors.New("test"))
				return mockUserDao, mockUserCache
			},
			ctx:      context.Background(),
			id:       123,
			wantUser: domain.User{},
			wantErr:  errors.New("test"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			userDao, userCache := tc.mock(ctrl)
			userRepo := NewUserRepository(userDao, userCache)
			gotUser, err := userRepo.FindByID(tc.ctx, tc.id)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, gotUser)
		})
	}
}
