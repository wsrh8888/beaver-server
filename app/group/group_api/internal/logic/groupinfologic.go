package logic

import (
	"context"
	"errors"
	"time"

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/group/group_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupInfoLogic {
	return &GroupInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupInfoLogic) GroupInfo(req *types.GroupInfoReq) (resp *types.GroupInfoRes, err error) {
	// 查询群组信息
	var group group_models.GroupModel
	err = l.svcCtx.DB.Take(&group, "uuid = ?", req.GroupID).Error

	if err != nil {
		logx.Errorf("查询群组失败: %s", err.Error())
		return nil, errors.New("群组不存在")
	}

	// 统计成员数量
	var memberCount int64
	_ = l.svcCtx.DB.Model(&group_models.GroupMemberModel{}).
		Where("group_id = ?", req.GroupID).Count(&memberCount).Error

	return &types.GroupInfoRes{
		GroupID:        group.GroupID,
		Title:          group.Title,
		Avatar:         group.Avatar,
		ConversationID: group.GroupID,
		MemberCount:    int(memberCount),
		CreatorID:      group.CreatorID,
		Notice:         group.Notice,
		JoinType:       group.JoinType,
		Status:         group.Status,
		CreateAt:       time.Time(group.CreatedAt).Unix(),
		UpdateAt:       time.Time(group.UpdatedAt).Unix(),
		Version:        group.Version,
	}, nil
}
