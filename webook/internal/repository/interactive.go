package repository

import (
	"context"
	"my-go-basic-study/webook/internal/domain"
	"my-go-basic-study/webook/internal/repository/cache"
	"my-go-basic-study/webook/internal/repository/dao"
)

type InteractiveRepository interface {
	Get(ctx context.Context, biz string, id int64) (domain.Interactive, error)
	IncrReadCnt(ctx context.Context, biz string, id int64) error
}

type interactiveRepository struct {
	dao   dao.InteractiveDAO
	cache cache.InteractiveCache
}

func (i *interactiveRepository) Get(ctx context.Context, biz string, id int64) (domain.Interactive, error) {
	//TODO implement me
	panic("implement me")
}

func (i *interactiveRepository) IncrReadCnt(ctx context.Context, biz string, id int64) error {
	//TODO implement me
	panic("implement me")
}
