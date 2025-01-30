package logger

import (
	"context"
	"go.uber.org/zap"
	"my-go-basic-study/webook/internal/service/sms/tencent"
)

type Service struct {
	svc *tencent.Service
}

func (s *Service) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	zap.L().Debug("发送短信", zap.String("tplId", tplId), zap.Strings("args", args))
	err := s.svc.Send(ctx, tplId, numbers)
	if err != nil {
		zap.L().Debug("发送短信出现异常", zap.Error(err))
	}
	return err
}
