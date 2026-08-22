package circle

import (
	"context"
	"time"

	"beaver/app/circle/circle_api/internal/svc"
	"beaver/app/circle/circle_api/internal/types"
	"beaver/app/circle/circle_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResolveCircleInviteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewResolveCircleInviteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResolveCircleInviteLogic {
	return &ResolveCircleInviteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResolveCircleInviteLogic) ResolveCircleInvite(req *types.ResolveCircleInviteReq) (resp *types.ResolveCircleInviteRes, err error) {
	resp = &types.ResolveCircleInviteRes{Code: req.Code, Valid: false}
	if req.Code == "" {
		return resp, nil
	}

	var invite circle_models.CircleInviteModel
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

	var circle circle_models.CircleModel
	if e := l.svcCtx.DB.Where("circle_id = ? AND is_deleted = false", invite.CircleID).First(&circle).Error; e != nil {
		return resp, nil
	}

	alreadyJoined := false
	var member circle_models.CircleMemberModel
	if l.svcCtx.DB.Where("circle_id = ? AND user_id = ?", invite.CircleID, req.UserID).First(&member).Error == nil {
		alreadyJoined = true
	}

	resp.Valid = true
	resp.CircleID = circle.CircleID
	resp.Name = circle.Name
	resp.Avatar = circle.Avatar
	resp.Description = circle.Description
	resp.MemberCount = countMembers(l.svcCtx.DB, circle.CircleID)
	resp.JoinType = circle.JoinType
	resp.AlreadyJoined = alreadyJoined
	return resp, nil
}
