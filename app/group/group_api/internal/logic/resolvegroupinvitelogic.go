package logic

import (
	"context"
	"time"

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/group/group_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResolveGroupInviteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewResolveGroupInviteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResolveGroupInviteLogic {
	return &ResolveGroupInviteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResolveGroupInviteLogic) ResolveGroupInvite(req *types.ResolveGroupInviteReq) (*types.ResolveGroupInviteRes, error) {
	resp := &types.ResolveGroupInviteRes{Code: req.Code, Valid: false}
	if req.Code == "" {
		return resp, nil
	}

	var invite group_models.GroupInviteLinkModel
	if e := l.svcCtx.DB.Where("token = ?", req.Code).First(&invite).Error; e != nil {
		return resp, nil
	}
	if invite.Status == 2 || invite.Status == 3 {
		return resp, nil
	}
	if invite.MaxUses > 0 && invite.UsedCount >= invite.MaxUses {
		return resp, nil
	}
	if invite.ExpireAt > 0 && time.Now().Unix() >= invite.ExpireAt {
		return resp, nil
	}

	var group group_models.GroupModel
	if e := l.svcCtx.DB.Where("group_id = ? AND status = 1", invite.GroupID).First(&group).Error; e != nil {
		return resp, nil
	}

	var memberCount int64
	_ = l.svcCtx.DB.Model(&group_models.GroupMemberModel{}).
		Where("group_id = ? AND status = 1", invite.GroupID).Count(&memberCount).Error

	alreadyJoined := false
	var member group_models.GroupMemberModel
	if l.svcCtx.DB.Where("group_id = ? AND user_id = ? AND status = 1", invite.GroupID, req.UserID).First(&member).Error == nil {
		alreadyJoined = true
	}

	resp.Valid = true
	resp.GroupID = group.GroupID
	resp.Title = group.Title
	resp.Avatar = group.Avatar
	resp.Notice = group.Notice
	resp.MemberCount = memberCount
	resp.JoinType = group.JoinType
	resp.AlreadyJoined = alreadyJoined
	return resp, nil
}
