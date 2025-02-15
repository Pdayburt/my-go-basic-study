package dao

import (
	"context"
	"errors"
	"github.com/bwmarrin/snowflake"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"time"
)

type MongoBD struct {
	/*	client   *mongo.Client
		database *mongo.Database*/
	//制作库
	col *mongo.Collection
	//线上库
	liveCol *mongo.Collection
	node    snowflake.Node
}

func NewMongoBD(database *mongo.Database, node snowflake.Node) *MongoBD {
	return &MongoBD{
		col:     database.Collection("articles"),
		liveCol: database.Collection("publish_articles"),
		node:    node,
	}
}
func (m *MongoBD) Insert(ctx context.Context, art Article) (int64, error) {
	id := m.node.Generate().Int64()
	art.Id = id
	_, err := m.col.InsertOne(ctx, art)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (m *MongoBD) UpdateById(ctx context.Context, entity Article) error {
	updateRes, err := m.col.UpdateOne(ctx, Article{Id: entity.Id}, bson.D{
		{"$set", Article{
			Title:    entity.Title,
			Content:  entity.Content,
			AuthorId: entity.AuthorId,
			Status:   entity.Status,
			Ctime:    entity.Ctime,
			Utime:    entity.Utime,
		}},
	})
	if err != nil {
		return err
	}
	if updateRes.ModifiedCount == 0 {
		return errors.New("update article fail")
	}
	return nil
}

func (m *MongoBD) Sync(ctx context.Context, art Article) (int64, error) {
	//mongoBD中没有事务的概念 所以 分步执行
	var (
		id  = art.Id
		err error
	)
	if id > 0 {
		err = m.UpdateById(ctx, art)
	} else {
		id, err = m.Insert(ctx, art)
	}
	if err != nil {
		return 0, err
	}
	art.Id = id
	now := time.Now().UnixMilli()
	art.Utime = now
	// liveCol 是 INSERT or Update 语义
	filter := bson.D{bson.E{"id", art.Id},
		bson.E{"author_id", art.AuthorId}}
	set := bson.D{bson.E{"$set", art},
		bson.E{"$setOnInsert",
			bson.D{bson.E{"ctime", now}}}}
	_, err = m.liveCol.UpdateOne(ctx,
		filter, set)
	return id, err

}

func (m *MongoBD) SyncStatus(ctx context.Context, uid int64, id int64, status uint8) error {
	//TODO implement me
	panic("implement me")
}

func (m *MongoBD) GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]Article, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoBD) GetById(ctx context.Context, id int64) (Article, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoBD) GetPubById(ctx context.Context, id int64) (PublishedArticle, error) {
	//TODO implement me
	panic("implement me")
}
