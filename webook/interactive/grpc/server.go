package grpc

import (
	"context"
	"my-go-basic-study/webook/api/proto/gen/interactive/v1"
	"my-go-basic-study/webook/interactive/domain"
	"my-go-basic-study/webook/interactive/service"
)

type InteractiveServiceServer struct {
	svc service.InteractiveService
	v1.UnimplementedInteractiveServiceServer
}

func (i *InteractiveServiceServer) IncrReadCnt(ctx context.Context, req *v1.IncrReadCntReq) (*v1.IncrReadCntResp, error) {
	err := i.svc.IncrReadCnt(ctx, req.GetBiz(), req.GetBizId())
	return &v1.IncrReadCntResp{}, err
}

func (i *InteractiveServiceServer) Like(ctx context.Context, req *v1.LikeReq) (*v1.LikeResp, error) {

	err := i.svc.Like(ctx, req.GetBiz(), req.GetId(), req.GetUid())
	return &v1.LikeResp{}, err
}

func (i *InteractiveServiceServer) CancelLike(ctx context.Context, req *v1.CancelLikeReq) (*v1.CancelLikeResp, error) {
	err := i.svc.CancelLike(ctx, req.GetBiz(), req.GetId(), req.GetUid())
	return &v1.CancelLikeResp{}, err
}

func (i *InteractiveServiceServer) Collect(ctx context.Context, req *v1.CollectReq) (*v1.CollectResp, error) {
	err := i.svc.Collect(ctx, req.GetBiz(), req.GetBizId(), req.GetCizId(), req.GetUid())
	return &v1.CollectResp{}, err
}

func (i *InteractiveServiceServer) Get(ctx context.Context, req *v1.GetReq) (*v1.GetResp, error) {
	interactive, err := i.svc.Get(ctx, req.GetBiz(), req.GetId(), req.GetUid())
	if err != nil {
		return nil, err
	}
	return &v1.GetResp{
		Interactive: i.toDTO(interactive),
	}, nil

}

func (i *InteractiveServiceServer) GetByIds(ctx context.Context, req *v1.GetByIdsReq) (*v1.GetByIdsResp, error) {

	interactiveMap, err := i.svc.GetByIds(ctx, req.GetBiz(), req.GetIds())
	if err != nil {
		return nil, err
	}
	m := make(map[int64]*v1.Interactive, len(interactiveMap))
	for k, v := range interactiveMap {
		m[k] = i.toDTO(v)
	}
	return &v1.GetByIdsResp{
		InteractiveMap: m,
	}, nil
}

func (i *InteractiveServiceServer) toDTO(interactive domain.Interactive) *v1.Interactive {
	return &v1.Interactive{
		Biz:        interactive.Biz,
		BizId:      interactive.BizId,
		ReadCnt:    interactive.ReadCnt,
		LikeCnt:    interactive.LikeCnt,
		CollectCnt: interactive.CollectCnt,
		Liked:      interactive.Liked,
		Collected:  interactive.Collected,
	}
}
