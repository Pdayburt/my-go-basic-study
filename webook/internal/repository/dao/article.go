package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type ArticleDao interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, article Article) error
}

type GormArticleDao struct {
	db *gorm.DB
}

func (g *GormArticleDao) UpdateById(ctx context.Context, article Article) error {
	now := time.Now().UnixMilli()
	article.Utime = now
	return g.db.WithContext(ctx).
		Model(&article).
		Where("id = ?", article.Id).Updates(map[string]interface{}{
		"title":   article.Title,
		"content": article.Content,
		"utime":   article.Utime,
	}).Error
}

func (g *GormArticleDao) Insert(ctx context.Context, art Article) (int64, error) {

	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	err := g.db.WithContext(ctx).Create(&art).Error
	return art.Id, err
}

func NewArticleDao(db *gorm.DB) ArticleDao {
	return &GormArticleDao{
		db: db,
	}
}

type Article struct {
	Id      int64  `gorm:"primaryKey;autoIncrement"`
	Title   string `gorm:"type:varchar(1024);not null"`
	Content string `gorm:"type:BLOB"`
	/*AuthorId int64  `gorm:"index:idx_authorId_ctime"`
	Ctime    int64  `gorm:"index:idx_authorId_ctime"`*/
	AuthorId int64 `gorm:"index:idx_authorId"`
	Ctime    int64
	Utime    int64
}
