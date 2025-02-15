package mongo

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/event"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"testing"
)

func TestMongo_insert(t *testing.T) {

	monitor := event.CommandMonitor{
		Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
			fmt.Println(startedEvent.Command)
		},
		Succeeded: func(ctx context.Context, succeededEvent *event.CommandSucceededEvent) {

		},
		Failed: func(ctx context.Context, failedEvent *event.CommandFailedEvent) {

		},
	}
	ops := options.Client().
		ApplyURI("mongodb://root:root@localhost:27017").
		SetMonitor(&monitor)
	client, err := mongo.Connect(ops)
	if err != nil {
		panic(err)
	}
	wbDb := client.Database("webook")
	articleCol := wbDb.Collection("articles")
	/*result, err := articleCol.InsertOne(context.Background(), Article{
		Id:      123,
		Title:   "我的标题",
		Content: "我的内容",
	})*/
	result, err := articleCol.InsertMany(context.Background(), []Article{
		{
			Id:      123,
			Title:   "我的标题",
			Content: "我的内容",
		},
		{
			Id:      1231,
			Title:   "我的标题1",
			Content: "我的内容1",
		},
		{
			Id:      12313,
			Title:   "我的标题13",
			Content: "我的内容13",
		},
		{
			Id:      123112,
			Title:   "我的标题112",
			Content: "我的内容112",
		},
		{
			Id:      123222,
			Title:   "我的标题222",
			Content: "我的内容222",
		},
	})

	fmt.Println(result.InsertedIDs)

}

func TestMongo_find(t *testing.T) {
	monitor := event.CommandMonitor{
		Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
			fmt.Println(startedEvent.Command)
		},
		Succeeded: func(ctx context.Context, succeededEvent *event.CommandSucceededEvent) {

		},
		Failed: func(ctx context.Context, failedEvent *event.CommandFailedEvent) {

		},
	}
	ops := options.Client().
		ApplyURI("mongodb://root:root@localhost:27017").
		SetMonitor(&monitor)
	client, err := mongo.Connect(ops)
	if err != nil {
		panic(err)
	}
	wbDb := client.Database("webook")
	articleCol := wbDb.Collection("articles")

	filter := bson.D{
		bson.E{
			Key:   "id",
			Value: 123,
		},
	}
	var art Article
	err = articleCol.FindOne(context.Background(), filter).Decode(&art)
	assert.NoError(t, err)
	fmt.Printf("filter res: %+v\n", art)
	art = Article{}
	err = articleCol.FindOne(context.Background(), Article{Id: 1232}).Decode(&art)
	if errors.Is(err, mongo.ErrNoDocuments) {
		fmt.Println("查不到数据～～")
	}
	assert.NoError(t, err)
	fmt.Printf("Article res: %+v\n", art)

}

func TestMongo_update(t *testing.T) {
	monitor := event.CommandMonitor{
		Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
			fmt.Println(startedEvent.Command)
		},
		Succeeded: func(ctx context.Context, succeededEvent *event.CommandSucceededEvent) {

		},
		Failed: func(ctx context.Context, failedEvent *event.CommandFailedEvent) {

		},
	}
	ops := options.Client().
		ApplyURI("mongodb://root:root@localhost:27017").
		SetMonitor(&monitor)
	client, err := mongo.Connect(ops)
	if err != nil {
		panic(err)
	}
	wbDb := client.Database("webook")
	articleCol := wbDb.Collection("articles")
	/*	filter := bson.D{
		bson.E{
			Key:   "id",
			Value: 123222,
		},
	}*/
	/*sets := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				bson.E{
					Key:   "title",
					Value: "我的新的标题",
				},
			},
		},
	}
	updateResult, err := articleCol.UpdateMany(context.Background(), filter, sets)
	if err != nil {
		panic(err)
	}
	fmt.Println("affected: ", updateResult.ModifiedCount)*/
	/*updateResult, err := articleCol.UpdateMany(context.Background(), filter, bson.D{
		{
			Key: "$set",
			Value: Article{
				Title: "我的新的标题222",
			},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("affected: ", updateResult.ModifiedCount)*/

	updateResult, err := articleCol.UpdateMany(context.Background(), Article{
		Id: 123,
	}, bson.D{
		bson.E{
			Key: "$set",
			Value: Article{
				Title:   "我的新的标题123",
				Content: "我的新的内容123",
			},
		},
	})

	if err != nil {
		panic(err)
	}
	fmt.Println(updateResult.MatchedCount)

}

func TestMongo_delete(t *testing.T) {
	monitor := event.CommandMonitor{
		Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
			fmt.Println(startedEvent.Command)
		},
		Succeeded: func(ctx context.Context, succeededEvent *event.CommandSucceededEvent) {

		},
		Failed: func(ctx context.Context, failedEvent *event.CommandFailedEvent) {

		},
	}
	ops := options.Client().
		ApplyURI("mongodb://root:root@localhost:27017").
		SetMonitor(&monitor)
	client, err := mongo.Connect(ops)
	if err != nil {
		panic(err)
	}
	wbDb := client.Database("webook")
	articleCol := wbDb.Collection("articles")
	//deleteResult, err := articleCol.DeleteOne(context.Background(), bson.D{
	//	bson.E{
	//		Key:   "id",
	//		Value: 123,
	//	},
	//})
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(deleteResult.DeletedCount)
	deleteResult, err := articleCol.DeleteMany(context.Background(), Article{Id: 123222})
	if err != nil {
		panic(err)
	}
	fmt.Println(deleteResult.DeletedCount)
}

func TestMongo_findWithOr(t *testing.T) {
	monitor := event.CommandMonitor{
		Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
			fmt.Println(startedEvent.Command)
		},
		Succeeded: func(ctx context.Context, succeededEvent *event.CommandSucceededEvent) {

		},
		Failed: func(ctx context.Context, failedEvent *event.CommandFailedEvent) {

		},
	}
	ops := options.Client().
		ApplyURI("mongodb://root:root@localhost:27017").
		SetMonitor(&monitor)
	client, err := mongo.Connect(ops)
	if err != nil {
		panic(err)
	}
	wbDb := client.Database("webook")
	articleCol := wbDb.Collection("articles")

	cursor, err := articleCol.Find(context.Background(), bson.D{
		bson.E{
			Key: "or$",
			Value: bson.A{
				bson.D{
					bson.E{
						Key:   "id",
						Value: 1231,
					},
				},
				bson.D{
					bson.E{
						Key:   "title",
						Value: "我的标题",
					},
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}
	var arts []Article
	err = cursor.All(context.Background(), &arts)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", arts)

}

func TestMongo_findWithAnd(t *testing.T) {
	monitor := event.CommandMonitor{
		Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
			fmt.Println(startedEvent.Command)
		},
		Succeeded: func(ctx context.Context, succeededEvent *event.CommandSucceededEvent) {

		},
		Failed: func(ctx context.Context, failedEvent *event.CommandFailedEvent) {

		},
	}
	ops := options.Client().
		ApplyURI("mongodb://root:root@localhost:27017").
		SetMonitor(&monitor)
	client, err := mongo.Connect(ops)
	if err != nil {
		panic(err)
	}
	wbDb := client.Database("webook")
	articleCol := wbDb.Collection("articles")
	find, err := articleCol.Find(context.Background(), bson.D{
		bson.E{
			Key: "$and",
			Value: bson.A{
				bson.D{
					bson.E{
						Key:   "id",
						Value: 123112,
					},
				},
				bson.D{
					bson.E{
						Key:   "title",
						Value: "我的标题112",
					},
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}
	var arts []Article
	err = find.All(context.Background(), &arts)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", arts)

}

func TestMongo_findWithIn(t *testing.T) {
	monitor := event.CommandMonitor{
		Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
			fmt.Println(startedEvent.Command)
		},
		Succeeded: func(ctx context.Context, succeededEvent *event.CommandSucceededEvent) {

		},
		Failed: func(ctx context.Context, failedEvent *event.CommandFailedEvent) {

		},
	}
	ops := options.Client().
		ApplyURI("mongodb://root:root@localhost:27017").
		SetMonitor(&monitor)
	client, err := mongo.Connect(ops)
	if err != nil {
		panic(err)
	}
	wbDb := client.Database("webook")
	articleCol := wbDb.Collection("articles")

	cursor, err := articleCol.Find(context.Background(), bson.D{
		bson.E{
			Key: "id",
			Value: bson.D{
				bson.E{
					Key:   "$in",
					Value: []int{123, 1231},
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}
	var arts []Article
	err = cursor.All(context.Background(), &arts)
	fmt.Printf("%+v\n", arts)
}

func TestMongo_createIndex(t *testing.T) {
	monitor := event.CommandMonitor{
		Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
			fmt.Println(startedEvent.Command)
		},
		Succeeded: func(ctx context.Context, succeededEvent *event.CommandSucceededEvent) {

		},
		Failed: func(ctx context.Context, failedEvent *event.CommandFailedEvent) {

		},
	}
	ops := options.Client().
		ApplyURI("mongodb://root:root@localhost:27017").
		SetMonitor(&monitor)
	client, err := mongo.Connect(ops)
	if err != nil {
		panic(err)
	}
	wbDb := client.Database("webook")
	articleCol := wbDb.Collection("articles")
	one, err := articleCol.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "id", Value: 1},
			{Key: "title", Value: -1}, // id 升序，title 降序
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(one)

}

type Article struct {
	Id       int64  `bson:"id,omitempty"`
	Title    string `bson:"title,omitempty"`
	Content  string `bson:"content,omitempty"`
	AuthorId int64  `bson:"author_id,omitempty"`
	Status   int    `bson:"status,omitempty"`
	Ctime    int64  `bson:"ctime,omitempty"`
	Utime    int64  `bson:"utime,omitempty"`
}
