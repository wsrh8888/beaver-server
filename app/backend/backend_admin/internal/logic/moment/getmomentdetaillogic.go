package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/moment/moment_rpc/types/moment_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMomentDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMomentDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMomentDetailLogic {
	return &GetMomentDetailLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetMomentDetailLogic) GetMomentDetail(req *types.GetMomentDetailReq) (resp *types.GetMomentDetailRes, err error) {
	if req.MomentId == "" {
		return nil, errors.New("动态ID不能为空")
	}

	rpcRes, err := l.svcCtx.MomentRpc.ListMoments(l.ctx, &moment_rpc.ListMomentsReq{
		MomentId: req.MomentId,
		Page:     1,
		PageSize: 1,
	})
	if err != nil {
		l.Errorf("获取动态详情失败: %v", err)
		return nil, err
	}
	if len(rpcRes.List) == 0 {
		return nil, errors.New("动态不存在")
	}

	m := rpcRes.List[0]
	files := make([]types.GetMomentDetailFileInfo, 0, len(m.Files))
	for _, f := range m.Files {
		files = append(files, types.GetMomentDetailFileInfo{FileName: f.FileKey})
	}
	return &types.GetMomentDetailRes{
		MomentId:  m.MomentId,
		UserId:    m.UserId,
		Content:   m.Content,
		Files:     files,
		IsDeleted: m.IsDeleted,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}, nil
}
