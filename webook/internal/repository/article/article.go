package article

import (
	"context"
	"my-go-basic-study/webook/internal/domain"
	"my-go-basic-study/webook/internal/repository/dao"
	"my-go-basic-study/webook/internal/repository/dao/article"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, article domain.Article) error
	FindById(ctx context.Context, id int64) (domain.Article, error)
	Sync(ctx context.Context, art domain.Article) (int64, error)
}

type CachedArticleRepository struct {
	dao       dao.ArticleDao
	readerDao article.ReaderDAO
	authorDao article.AuthorDAO
}

func (c *CachedArticleRepository) Sync(ctx context.Context, art domain.Article) (int64, error) {
	//v1 操作两个DAO
	var (
		id  = art.Id
		err error
	)
	
}

func (c *CachedArticleRepository) FindById(ctx context.Context, id int64) (domain.Article, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CachedArticleRepository) Update(ctx context.Context, article domain.Article) error {
	return c.dao.UpdateById(ctx, dao.Article{
		Id:       article.Id,
		Title:    article.Title,
		Content:  article.Content,
		AuthorId: article.Author.Id,
	})
}

func (c *CachedArticleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {

	return c.dao.Insert(ctx, dao.Article{
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
	})
}

func NewArticleRepository(dao dao.ArticleDao) ArticleRepository {
	return &CachedArticleRepository{
		dao: dao,
	}
}
