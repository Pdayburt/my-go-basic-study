package intrgration

import (
	"testing"
)

func TestArticleHandler_Edit(t *testing.T) {

	/*testCases := []struct {
		name     string
		before   func(t *testing.T) *gin.Engine
		after    func(t *testing.T)
		art      web.ArticleReq    //输入
		wantCode int               //http 响应码
		wantRes  web.Result[int64] //http的响应，带上帖子的ID
	}{
		{
			name: "新建帖子--保存成功～",
			before: func(t *testing.T) *gin.Engine {
				articleHandler := InitArticleHandler() // 获取 handler
				server := gin.Default()
				articleHandler.RegisterRouters(server)
				return server
			},
			after: func(t *testing.T) {
				//验证数据库
			},
			art: web.ArticleReq{
				Title:   "handler",
				Content: "handler content~",
			},
			wantCode: http.StatusOK,
			wantRes: web.Result[int64]{
				Code: 0,
				Data: 1,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//构造请求 执行 验证
			artBody, err := json.Marshal(tc.art)
			require.NoError(t, err)
			request, err := http.NewRequest("POST", "/articles/edit",
				bytes.NewBuffer(artBody))
			require.NoError(t, err)
			request.Header.Set("Content-Type", "application/json")
			svc := tc.before(t)
			recorder := httptest.NewRecorder()
			svc.ServeHTTP(recorder, request)
			assert.Equal(t, tc.wantCode, recorder.Code)
			var res web.Result[int64]
			err = json.Unmarshal(recorder.Body.Bytes(), &res)
			assert.Equal(t, tc.wantRes, res)
			tc.after(t)
		})
	}*/
}
