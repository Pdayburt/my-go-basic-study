package web

import (
	"bytes"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"my-go-basic-study/webook/internal/domain"
	"my-go-basic-study/webook/internal/service"
	svcmock "my-go-basic-study/webook/internal/service/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler_Signup(t *testing.T) {

	testCase := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) service.UserService
		reqBody  string
		wantCode int
		wantBody string
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) service.UserService {
				mockUserService := svcmock.NewMockUserService(ctrl)
				mockUserService.EXPECT().Signup(gomock.Any(), domain.User{
					Email:    "1123@qq.com",
					Password: "123456#qqq",
				}).Return(nil)
				return mockUserService
			},
			reqBody: `
{
   "email": "1123@qq.com",
   "password": "123456#qqq",
   "confirmPassword": "123456#qqq"
}
`,
			wantCode: http.StatusOK,
			wantBody: `signup success`,
		},
		{
			name: "参数不对，bind失败",
			mock: func(ctrl *gomock.Controller) service.UserService {
				mockUserService := svcmock.NewMockUserService(ctrl)
				//没有走到调用service中的signup这一步
				/*mockUserService.EXPECT().Signup(gomock.Any(), domain.User{
					Email:    "1123@qq.com",
					Password: "123456#qqq",
				}).Return(nil)*/
				return mockUserService
			},
			reqBody: `
{
   "email": "1123@qq.com",
   "password": "123456#qqq"
   "confirmPassword": "123456#qqq"
}
`,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "邮箱格式不对",
			mock: func(ctrl *gomock.Controller) service.UserService {
				mockUserService := svcmock.NewMockUserService(ctrl)
				//mockUserService.EXPECT().Signup(gomock.Any(), domain.User{
				//	Email:    "1123@qq.com",
				//	Password: "123456#qqq",
				//}).Return(nil)
				return mockUserService
			},
			reqBody: `
{
   "email": "1123qq.com",
   "password": "123456#qqq",
   "confirmPassword": "123456#qqq"
}
`,
			wantCode: http.StatusOK,
			wantBody: `邮件格式错误`,
		},
		{
			name: "邮箱已存在",
			mock: func(ctrl *gomock.Controller) service.UserService {
				mockUserService := svcmock.NewMockUserService(ctrl)
				mockUserService.EXPECT().Signup(gomock.Any(), domain.User{
					Email:    "1123@qq.com",
					Password: "123456#qqq",
				}).Return(service.ErrUserDuplicateEmail)
				return mockUserService
			},
			reqBody: `
{
   "email": "1123@qq.com",
   "password": "123456#qqq",
   "confirmPassword": "123456#qqq"
}
`,
			wantCode: http.StatusOK,
			wantBody: `邮箱已存在`,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			userHandler := NewUserHandler(tc.mock(ctrl))

			server := gin.Default()
			userHandler.RegisterRouters(server)
			request, err := http.NewRequest(http.MethodPost, "/users/signup",
				bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)
			request.Header.Set("Content-Type", "application/json")
			respWriter := httptest.NewRecorder()
			//http请求9进入gin框架的入口，也就是当调用这个方法是 Gin会处理这个请求同时将响应写回到resp中
			server.ServeHTTP(respWriter, request)
			assert.Equal(t, tc.wantCode, respWriter.Code)
			assert.Equal(t, tc.wantBody, respWriter.Body.String())
		})
	}

}

func TestMock(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockUserService := svcmock.NewMockUserService(controller)
	mockUserService.EXPECT().Signup(gomock.Any(), gomock.Any()).
		Return(errors.New("signup mock  error"))
	err := mockUserService.Signup(context.Background(), domain.User{
		Email:    "123@11.com",
		Password: "123@11.com",
	})
	t.Log(err)

}

func TestUserHttp_newReq(t *testing.T) {
	/*request, err := http.NewRequest(http.MethodGet, "www.baidu.com", nil)
	if err != nil {
		t.Fatal(err)
	}*/
}
