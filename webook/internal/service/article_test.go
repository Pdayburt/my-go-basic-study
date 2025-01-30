package service

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"my-go-basic-study/webook/internal/domain"
	"my-go-basic-study/webook/internal/repository/article"
	artrepomock "my-go-basic-study/webook/internal/repository/article/mock"
	"testing"
)

func Test_articleService_Publish(t *testing.T) {

	testcases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) (article.ArticleAuthorRepository, article.ArticleReaderRepository)
		input    domain.Article
		wantErr  error
		wantCode int64
	}{
		{
			name: "发表成功～",
			mock: func(ctrl *gomock.Controller) (article.ArticleAuthorRepository, article.ArticleReaderRepository) {

				author := artrepomock.NewMockArticleAuthorRepository(ctrl)
				author.EXPECT().Create(gomock.Any(), domain.Article{
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(1), nil)
				reader := artrepomock.NewMockArticleReaderRepository(ctrl)
				reader.EXPECT().Create(gomock.Any(), domain.Article{
					Id:      1,
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(1), nil)

				return author, reader
			},
			input: domain.Article{
				Title:   "我的标题",
				Content: "我的内容",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantErr:  nil,
			wantCode: 1,
		},
		{
			name: "修改并发表成功～",
			mock: func(ctrl *gomock.Controller) (article.ArticleAuthorRepository, article.ArticleReaderRepository) {

				author := artrepomock.NewMockArticleAuthorRepository(ctrl)
				author.EXPECT().Create(gomock.Any(), domain.Article{
					Id:      2,
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(1), nil)
				reader := artrepomock.NewMockArticleReaderRepository(ctrl)
				reader.EXPECT().Create(gomock.Any(), domain.Article{
					Id:      1,
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(1), nil)

				return author, reader
			},
			input: domain.Article{
				Title:   "我的标题",
				Content: "我的内容",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantErr:  nil,
			wantCode: 1,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			author, reader := tc.mock(ctrl)
			service := NewArticleService(author, reader)
			id, err := service.PublishV1(context.Background(), tc.input)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantCode, id)

		})
	}

}
