package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"my-go-basic-study/webook/internal/domain"
	"my-go-basic-study/webook/internal/repository/cache/redismock"
	"testing"
)

func TestRedisUserCache_Set(t *testing.T) {

	testCase := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) redis.Cmdable
		ctx      context.Context
		input    domain.User
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "test",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmdable := redismock.NewMockCmdable(ctrl)
				//Set(ctx, key, val, c.expiration)
				//cmdable.EXPECT().Set()
				return cmdable
			},
			ctx:      context.Background(),
			input:    domain.User{},
			wantUser: domain.User{},
			wantErr:  nil,
		},
	}

	for _, tc := range testCase {

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mock := tc.mock(ctrl)
			userCache := NewUserCache(mock)
			//Set(ctx context.Context, user domain.User) error
			err := userCache.Set(tc.ctx, tc.input)
			assert.Equal(t, tc.wantErr, err)

		})

	}

}
