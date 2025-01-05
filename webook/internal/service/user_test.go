package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"my-go-basic-study/webook/internal/domain"
	"my-go-basic-study/webook/internal/repository"
	repomock "my-go-basic-study/webook/internal/repository/mock"
	"testing"
)

func Test_userService_Login(t *testing.T) {

	testCase := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) repository.UserRepository
		ctx      context.Context
		user     domain.User
		wantErr  error
		wantUser domain.User
	}{
		{
			name: "登陆成功",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				mockUserRepository := repomock.NewMockUserRepository(ctrl)
				//	FindByEmail(ctx context.Context, email string) (domain.User, error)
				mockUserRepository.EXPECT().
					FindByEmail(gomock.Any(), "1123@qq.com").
					Return(domain.User{
						Email:    "1123@qq.com",
						Password: "$2a$10$MNWC6FuLNk7s66vVkh3TFuYfsz.KJH7GRartbvsFmx.17.bteUdfi",
					}, nil)
				return mockUserRepository
			},
			ctx: context.Background(),
			user: domain.User{
				Email:    "1123@qq.com",
				Password: "123456#qqq",
			},
			wantErr: nil,
			wantUser: domain.User{
				Email:    "1123@qq.com",
				Password: "$2a$10$MNWC6FuLNk7s66vVkh3TFuYfsz.KJH7GRartbvsFmx.17.bteUdfi",
			},
		},
		{
			name: "用户不存在",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				mockUserRepository := repomock.NewMockUserRepository(ctrl)
				//	FindByEmail(ctx context.Context, email string) (domain.User, error)
				mockUserRepository.EXPECT().
					FindByEmail(gomock.Any(), "1123@qq.com").
					Return(domain.User{}, repository.ErrUserNotFound)
				return mockUserRepository
			},
			ctx: context.Background(),
			user: domain.User{
				Email:    "1123@qq.com",
				Password: "123456#qqq",
			},
			wantErr:  ErrInvalidUserOrPassword,
			wantUser: domain.User{},
		},
		{
			name: "DB 错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				mockUserRepository := repomock.NewMockUserRepository(ctrl)
				//	FindByEmail(ctx context.Context, email string) (domain.User, error)
				mockUserRepository.EXPECT().
					FindByEmail(gomock.Any(), "1123@qq.com").
					Return(domain.User{}, errors.New("随便一个错误"))
				return mockUserRepository
			},
			ctx: context.Background(),
			user: domain.User{
				Email:    "1123@qq.com",
				Password: "123456#qqq",
			},
			wantErr:  errors.New("随便一个错误"),
			wantUser: domain.User{},
		},
		{
			name: "密码错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				mockUserRepository := repomock.NewMockUserRepository(ctrl)
				//	FindByEmail(ctx context.Context, email string) (domain.User, error)
				mockUserRepository.EXPECT().
					FindByEmail(gomock.Any(), "1123@qq.com").
					Return(domain.User{
						Email:    "1123@qq.com",
						Password: "$2a$10$MNWC6FuLNk7s66vVkh3TFuYfsz.KJH7GRartbvsFmx.17.bteUdfi",
					}, nil)
				return mockUserRepository
			},
			ctx: context.Background(),
			user: domain.User{
				Email:    "1123@qq.com",
				Password: "1234567#qqq",
			},
			wantErr:  ErrInvalidUserOrPassword,
			wantUser: domain.User{},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockUserRepo := tc.mock(ctrl)
			userSvc := NewUserService(mockUserRepo)
			user, err := userSvc.Login(tc.ctx, tc.user)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, user)
		})
	}

}

func TestGeneratePasswrod(t *testing.T) {
	cryptedPassword, err := bcrypt.GenerateFromPassword([]byte("123456#qqq"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(cryptedPassword))
}
