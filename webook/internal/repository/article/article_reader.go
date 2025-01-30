package article

import (
	"context"
	"my-go-basic-study/webook/internal/domain"
)

type ArticleReaderRepository interface {
	Create(ctx context.Context, article domain.Article) (int64, error)
}
