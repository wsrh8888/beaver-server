package logic

import (
	"context"
	"errors"

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
	// 查询群组信息，并预加载成员列表
	var group group_models.GroupModel
	err = l.svcCtx.DB.Preload("MemberList").Take(&group, "uuid = ?", req.GroupID).Error

	if err != nil {
		logx.Errorf("查询群组失败: %s", err.Error())
		return nil, errors.New("群组不存在")
	}

	return &types.GroupInfoRes{
		Title:          group.Title,
		Avatar:         group.Avatar,
		ConversationID: group.UUID,
		MemberCount:    len(group.MemberList), // 添加成员数量
	}, nil
}
