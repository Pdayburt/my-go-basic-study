package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"my-go-basic-study/webook/internal/service"
	svcmock "my-go-basic-study/webook/internal/service/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestArticleHandler_Publish(t *testing.T) {
	testCase := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) service.ArticleService
		reqBody  string
		wantCode int
		wantBody Result[int64]
	}{
		{
			name: "新建并发表",
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				articleService := svcmock.NewMockArticleService(ctrl)
				articleService.EXPECT().Publish(gomock.Any(),
					/*	domain.Article{
						Title:   "我的标题",
						Content: "我的内容",
						Author: domain.Author{
							Id: 1,
						},
					}*/
					gomock.Any(),
				).Return(int64(1), nil)

				return articleService
			},
			reqBody: `
{
"title":"我的标题",
"content":"我的内容"
}
`,
			wantCode: http.StatusOK,
			wantBody: Result[int64]{
				Data: 1,
				Msg:  "OK",
			},
		},
		{
			name: "发表失败",
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				articleService := svcmock.NewMockArticleService(ctrl)
				articleService.EXPECT().Publish(gomock.Any(),
					/*	domain.Article{
						Title:   "我的标题",
						Content: "我的内容",
						Author: domain.Author{
							Id: 1,
						},
					}*/
					gomock.Any(),
				).Return(int64(0), errors.New("publish fail"))

				return articleService
			},
			reqBody: `
{
"title":"我的标题",
"content":"我的内容"
}
`,
			wantCode: http.StatusOK,
			wantBody: Result[int64]{
				Code: 5,
				Msg:  "系统错误",
			},
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			server := gin.Default()
			articleHandler := NewArticleHandler(tc.mock(ctrl))
			articleHandler.RegisterRouters(server)
			request, err := http.NewRequest(http.MethodPost,
				"/articles/publish",
				bytes.NewBuffer([]byte(tc.reqBody)))
			if err != nil {
				t.Fatal(err)
			}
			request.Header.Set("Content-Type", "application/json")

			response := httptest.NewRecorder()
			server.ServeHTTP(response, request)
			assert.Equal(t, tc.wantCode, response.Code)

			if response.Code != http.StatusOK {
				t.Fatalf("response code error %d", response.Code)
			}

			var webRes Result[int64]
			err = json.NewDecoder(response.Body).Decode(&webRes)
			require.NoError(t, err)
			assert.Equal(t, tc.wantBody, webRes)
		})
	}
}
