package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFeedbackDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFeedbackDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFeedbackDetailLogic {
	return &GetFeedbackDetailLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetFeedbackDetailLogic) GetFeedbackDetail(req *types.GetFeedbackDetailReq) (resp *types.GetFeedbackDetailRes, err error) {
	rpcRes, err := l.svcCtx.PlatformRpc.GetFeedback(l.ctx, &platform_rpc.GetFeedbackReq{Id: uint64(req.Id)})
	if err != nil {
		l.Errorf("获取反馈详情失败: %v", err)
		return nil, err
	}

	f := rpcRes.Feedback
	return &types.GetFeedbackDetailRes{
		Id:           uint(f.Id),
		UserId:       f.UserId,
		Content:      f.Content,
		Type:         int(f.Type),
		Status:       int(f.Status),
		FileNames:    f.FileNames,
		HandlerId:    f.HandlerId,
		HandleTime:   f.HandleTime,
		HandleResult: f.HandleResult,
		CreatedAt:    f.CreatedAt,
		UpdatedAt:    f.UpdatedAt,
	}, nil
}
