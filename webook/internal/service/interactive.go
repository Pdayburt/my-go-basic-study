package service

import (
	"context"
	"my-go-basic-study/webook/internal/repository"
)

type InteractiveService interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
}

type interactiveService struct {
	repo repository.InteractiveRepository
}

func (i *interactiveService) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	return i.repo.IncrReadCnt(ctx, biz, bizId)
}
