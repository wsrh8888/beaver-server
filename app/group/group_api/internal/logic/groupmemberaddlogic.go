package logic

import (
	"context"
	"errors"

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/group/group_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupMemberAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupMemberAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupMemberAddLogic {
	return &GroupMemberAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupMemberAddLogic) GroupMemberAdd(req *types.GroupMemberAddReq) (resp *types.GroupMemberAddRes, err error) {
	// 群成员邀请好友，isInvite 为true
	var member group_models.GroupMemberModel
	err = l.svcCtx.DB.Take(&member, "group_id = ? and user_id = ?", req.GroupID, req.UserID).Error

	if err != nil {
		return nil, errors.New("用户不是群组成员")
	}

	// 去查一下哪些用户已经进群了。
	var memberList []group_models.GroupMemberModel

	l.svcCtx.DB.Find(&memberList, "group_id = ? and user_id in ?", req.GroupID, req.MemberIdList)

	if len(memberList) > 0 {
		return nil, errors.New("已经有用户已经是群成员")
	}

	for _, memberID := range req.MemberIdList {
		memberList = append(memberList, group_models.GroupMemberModel{
			GroupID: req.GroupID,
			UserID:  memberID,
			Role:    3,
		})

	}
	err = l.svcCtx.DB.Create(&memberList).Error

	if err != nil {
		return nil, errors.New("添加失败")
	}

	return
}
