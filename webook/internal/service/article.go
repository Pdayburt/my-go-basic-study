package service

import (
	"context"
	"errors"
	"my-go-basic-study/webook/internal/domain"
	"my-go-basic-study/webook/internal/repository/article"
)

type ArticleService interface {
	Save(ctx context.Context, article domain.Article) (int64, error)
	Publish(ctx context.Context, article domain.Article) (int64, error)
	PublishV1(ctx context.Context, article domain.Article) (int64, error)
}

type articleService struct {
	repo article.ArticleRepository
	//V1
	author article.ArticleAuthorRepository
	reader article.ArticleReaderRepository
}

/*func NewArticleService(repo article.ArticleRepository) ArticleService {
	return &articleService{
		repo: repo,
	}
}
*/

func NewArticleService(author article.ArticleAuthorRepository, reader article.ArticleReaderRepository) ArticleService {
	return &articleService{
		author: author,
		reader: reader,
	}
}

func (a *articleService) PublishV1(ctx context.Context, article domain.Article) (int64, error) {
	id, err := a.author.Create(ctx, article)
	if err != nil {
		return 0, err
	}
	article.Id = id
	return a.reader.Create(ctx, article)
}

func (a *articleService) Publish(ctx context.Context, article domain.Article) (int64, error) {

	return a.repo.Sync(ctx, article)
}

func (a *articleService) Save(ctx context.Context, article domain.Article) (int64, error) {
	if article.Id > 0 {
		err := a.repo.Update(ctx, article)
		return article.Id, err
	}
	return a.repo.Create(ctx, article)

}

func (a *articleService) update(ctx context.Context, art domain.Article) error {
	articleInDb, err := a.repo.FindById(ctx, art.Id)
	if err != nil {
		return err
	}
	if art.Author.Id != articleInDb.Author.Id {
		return errors.New("无权限更新别人的数据")
	}
	return a.repo.Update(ctx, art)
}
