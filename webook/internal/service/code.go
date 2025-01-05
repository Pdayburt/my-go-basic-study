package service

import "context"

type CodeService struct {
}

func (svc *CodeService) Send(ctx context.Context,
	//区分业务场景
	biz string,
	phone string) error {
	return nil
}

func (svc *CodeService) Verify(ctx context.Context,
	biz string, phone string, inputCode string) (bool, error) {

	return true, nil
}
