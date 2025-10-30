package logic

import (
	"context"

	"beaver/app/group/group_models"
	"beaver/app/group/group_rpc/internal/svc"
	"beaver/app/group/group_rpc/types/group_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupVersionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupVersionLogic {
	return &GetGroupVersionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetGroupVersionLogic) GetGroupVersion(in *group_rpc.GetGroupVersionReq) (*group_rpc.GetGroupVersionRes, error) {
	var maxVersion int64
	err := l.svcCtx.DB.Model(&group_models.GroupModel{}).
		Joins("JOIN group_member_models ON group_models.group_id = group_member_models.group_id").
		Where("group_member_models.user_id = ?", in.UserId).
		Select("COALESCE(MAX(group_models.version), 0)").
		Scan(&maxVersion).Error

	if err != nil {
		l.Errorf("获取最新群组版本号失败: %v", err)
		return nil, err
	}

	return &group_rpc.GetGroupVersionRes{
		LatestVersion: maxVersion,
	}, nil
}
