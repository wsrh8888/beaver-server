package logic

import (
	"context"

	"beaver/app/group/group_models"
	"beaver/app/group/group_rpc/internal/svc"
	"beaver/app/group/group_rpc/types/group_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupMemberVersionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupMemberVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupMemberVersionLogic {
	return &GetGroupMemberVersionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetGroupMemberVersionLogic) GetGroupMemberVersion(in *group_rpc.GetGroupMemberVersionReq) (*group_rpc.GetGroupMemberVersionRes, error) {
	var maxVersion int64
	err := l.svcCtx.DB.Model(&group_models.GroupMemberModel{}).
		Where("user_id = ?", in.UserId).
		Select("COALESCE(MAX(version), 0)").
		Scan(&maxVersion).Error

	if err != nil {
		l.Errorf("获取最新群成员版本号失败: %v", err)
		return nil, err
	}

	return &group_rpc.GetGroupMemberVersionRes{
		LatestVersion: maxVersion,
	}, nil
}
