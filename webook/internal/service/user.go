package service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"my-go-basic-study/webook/internal/domain"
	"my-go-basic-study/webook/internal/repository"
)

var (
	ErrUserDuplicateEmail    = repository.ErrUserDuplicateEmail
	ErrInvalidUserOrPassword = errors.New("～～邮箱或者密码错误～～")
)

type UserService interface {
	Signup(ctx context.Context, u domain.User) error
	Login(ctx context.Context, user domain.User) (domain.User, error)
	Profile(ctx context.Context, id int64) (domain.User, error)
}

var _ UserService = (*userService)(nil)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (svc *userService) Signup(ctx context.Context, u domain.User) error {
	cryptedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(cryptedPassword)
	return svc.repo.Create(ctx, u)
}

func (svc *userService) Login(ctx context.Context, user domain.User) (domain.User, error) {
	userByEmail, err := svc.repo.FindByEmail(ctx, user.Email)
	if errors.Is(err, repository.ErrUserNotFound) {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(userByEmail.Password), []byte(user.Password)); err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return userByEmail, nil

}

func (svc *userService) Profile(ctx context.Context, id int64) (domain.User, error) {
	user, err := svc.repo.FindByID(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}
