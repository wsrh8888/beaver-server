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
	// 获取用户管理的群组ID列表（作为群主或管理员）
	var managedGroupIDs []string
	err := l.svcCtx.DB.Model(&group_models.GroupMemberModel{}).
		Where("user_id = ? AND role IN (1, 2) AND status = 1", in.UserId).
		Pluck("group_id", &managedGroupIDs).Error
	if err != nil {
		l.Errorf("获取用户管理的群组失败: %v", err)
		return nil, err
	}

	// 构建查询条件：用户发起的申请 或 发送给用户管理群组的申请
	var maxVersion int64
	query := l.svcCtx.DB.Model(&group_models.GroupJoinRequestModel{})

	if len(managedGroupIDs) > 0 {
		query = query.Where("(applicant_user_id = ? OR group_id IN (?))", in.UserId, managedGroupIDs)
	} else {
		query = query.Where("applicant_user_id = ?", in.UserId)
	}

	err = query.Select("COALESCE(MAX(version), 0)").Scan(&maxVersion).Error
	if err != nil {
		l.Errorf("获取最新群组申请版本号失败: %v", err)
		return nil, err
	}

	return &group_rpc.GetGroupJoinRequestVersionRes{
		LatestVersion: maxVersion,
	}, nil
}
