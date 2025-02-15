package article

import (
	"context"
	"my-go-basic-study/webook/internal/domain"
	"my-go-basic-study/webook/internal/repository/dao"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, article domain.Article) error
	FindById(ctx context.Context, id int64) (domain.Article, error)
	Sync(ctx context.Context, art domain.Article) (int64, error)
	StatusSync(ctx context.Context, id int64, authorId int64, status domain.ArticleStatus) error
}

type CachedArticleRepository struct {
	dao dao.ArticleDao
	/*	readerDao dao.ReaderDAO
		authorDao dao.AuthorDAO*/
}

func (c *CachedArticleRepository) StatusSync(ctx context.Context, id int64, authorId int64, status domain.ArticleStatus) error {

	//return c.dao.StatusSync(ctx, id, authorId, uint8(status))pa
	panic("implement me")
}

func (c *CachedArticleRepository) Sync(ctx context.Context, art domain.Article) (int64, error) {
	//return c.dao.Sync(ctx, art)
	panic("implement me")
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
