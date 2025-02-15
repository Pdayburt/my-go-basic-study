package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type ArticleDao interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, entity Article) error
	Sync(ctx context.Context, entity Article) (int64, error)
	SyncStatus(ctx context.Context, uid int64, id int64, status uint8) error
	GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]Article, error)
	GetById(ctx context.Context, id int64) (Article, error)
	GetPubById(ctx context.Context, id int64) (PublishedArticle, error)
}

type GormArticleDao struct {
	db *gorm.DB
}

func (g *GormArticleDao) SyncStatus(ctx context.Context, uid int64, id int64, status uint8) error {
	//TODO implement me
	panic("implement me")
}

func (g *GormArticleDao) GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]Article, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GormArticleDao) GetById(ctx context.Context, id int64) (Article, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GormArticleDao) GetPubById(ctx context.Context, id int64) (PublishedArticle, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GormArticleDao) StatusSync(ctx context.Context, id int64, authorId int64, status uint8) error {
	now := time.Now()
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&Article{}).
			Where("id = ? AND author_id = ?", id, authorId).
			Updates(map[string]interface{}{
				"status": status,
				"utime":  now,
			})

		if res.Error != nil {
			//数据库有问题呢
			return res.Error
		}
		if res.RowsAffected == 0 {
			//没有记录，即 要么ID是错的 ，要么作者不对
			return gorm.ErrRecordNotFound
		}

		return tx.Model(&Article{}).
			Where("id = ? ", id).
			Updates(map[string]interface{}{
				"status": status,
				"utime":  now,
			}).Error

	})
}

func (g *GormArticleDao) Sync(ctx context.Context, art Article) (int64, error) {

	g.db.Transaction(func(tx *gorm.DB) error {

		return ErrUserDuplicateEmail
	})
	return 0, nil
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
	Id      int64  `gorm:"primaryKey,autoIncrement" bson:"id,omitempty"`
	Title   string `gorm:"type=varchar(4096)" bson:"title,omitempty"`
	Content string `gorm:"type=BLOB" bson:"content,omitempty"`
	// 我要根据创作者ID来查询
	AuthorId int64 `gorm:"index" bson:"author_id,omitempty"`
	Status   uint8 `bson:"status,omitempty"`
	Ctime    int64 `bson:"ctime,omitempty"`
	// 更新时间
	Utime int64 `bson:"utime,omitempty"`
}

type PublishedArticle Article

type PublishedArticleV1 struct {
	Article
}
