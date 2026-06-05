package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFeedbackListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFeedbackListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFeedbackListLogic {
	return &GetFeedbackListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetFeedbackListLogic) GetFeedbackList(req *types.GetFeedbackListReq) (resp *types.GetFeedbackListRes, err error) {
	rpcRes, err := l.svcCtx.PlatformRpc.ListFeedback(l.ctx, &platform_rpc.ListFeedbackReq{
		Page:     int32(req.Page),
		PageSize: int32(req.Limit),
		Status:   int32(req.Status),
		Type:     int32(req.Type),
		UserId:   req.UserID,
		Keywords: req.Keywords,
	})
	if err != nil {
		l.Errorf("获取反馈列表失败: %v", err)
		return nil, err
	}

	list := make([]types.GetFeedbackListItem, 0, len(rpcRes.List))
	for _, item := range rpcRes.List {
		list = append(list, types.GetFeedbackListItem{
			Id:           uint(item.Id),
			UserId:       item.UserId,
			Content:      item.Content,
			Type:         int(item.Type),
			Status:       int(item.Status),
			FileNames:    item.FileNames,
			HandlerId:    item.HandlerId,
			HandleTime:   item.HandleTime,
			HandleResult: item.HandleResult,
			CreatedAt:    item.CreatedAt,
			UpdatedAt:    item.UpdatedAt,
		})
	}

	return &types.GetFeedbackListRes{List: list, Total: rpcRes.Total}, nil
}
