package logic

import (
	"context"

	"beaver/app/group/group_models"
	"beaver/app/group/group_rpc/internal/svc"
	"beaver/app/group/group_rpc/types/group_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupJoinRequestVersionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupJoinRequestVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupJoinRequestVersionLogic {
	return &GetGroupJoinRequestVersionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetGroupJoinRequestVersionLogic) GetGroupJoinRequestVersion(in *group_rpc.GetGroupJoinRequestVersionReq) (*group_rpc.GetGroupJoinRequestVersionRes, error) {
	var maxVersion int64
	err := l.svcCtx.DB.Model(&group_models.GroupJoinRequestModel{}).
		Where("applicant_user_id = ?", in.UserId).
		Select("COALESCE(MAX(version), 0)").
		Scan(&maxVersion).Error

	if err != nil {
		l.Errorf("获取最新群组申请版本号失败: %v", err)
		return nil, err
	}

	return &group_rpc.GetGroupJoinRequestVersionRes{
		LatestVersion: maxVersion,
	}, nil
}
