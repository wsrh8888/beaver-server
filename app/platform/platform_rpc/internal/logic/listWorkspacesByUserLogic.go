package logic

import (
	"context"

	ai_models "beaver/app/platform/platform_models/ai"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListWorkspacesByUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListWorkspacesByUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListWorkspacesByUserLogic {
	return &ListWorkspacesByUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListWorkspacesByUserLogic) ListWorkspacesByUser(in *platform_rpc.ListWorkspacesByUserReq) (*platform_rpc.ListWorkspacesByUserRes, error) {
	var rows []ai_models.AgentWorkspace
	if err := l.svcCtx.DB.Where("user_id = ?", in.UserId).
		Order("version DESC").Find(&rows).Error; err != nil {
		return &platform_rpc.ListWorkspacesByUserRes{}, nil
	}
	out := make([]*platform_rpc.GetLatestWorkspaceRes, 0, len(rows))
	for _, row := range rows {
		out = append(out, &platform_rpc.GetLatestWorkspaceRes{
			Id:       row.ID,
			UserId:   row.UserID,
			DeviceId: row.DeviceID,
			Name:     row.Name,
			Path:     row.Path,
			Version:  row.Version,
		})
	}
	return &platform_rpc.ListWorkspacesByUserRes{Workspaces: out}, nil
}
