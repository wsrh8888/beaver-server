package logic

import (
	"context"
	"errors"
	"fmt"
	"strings"
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
	var group group_models.GroupModel
	err = l.svcCtx.DB.Take(&group, "group_id = ?", req.GroupID).Error

	if err != nil {
		logx.Errorf("查询群组失败: %s", err.Error())
		return nil, errors.New("群组不存在")
	}

	var memberCount int64
	_ = l.svcCtx.DB.Model(&group_models.GroupMemberModel{}).
		Where("group_id = ?", req.GroupID).Count(&memberCount).Error

	resp = &types.GroupInfoRes{
		GroupID:        group.GroupID,
		Title:          group.Title,
		Avatar:         group.Avatar,
		ConversationID: group.GroupID,
		MemberCount:    int(memberCount),
		CreatorID:      group.CreatorID,
		Notice:         group.Notice,
		JoinType:       group.JoinType,
		Status:         group.Status,
		CreatedAt:      time.Time(group.CreatedAt).Unix(),
		UpdatedAt:      time.Time(group.UpdatedAt).Unix(),
		Version:        group.Version,
	}

	var member group_models.GroupMemberModel
	if l.svcCtx.DB.Where("group_id = ? AND user_id = ? AND status = 1", req.GroupID, req.UserID).First(&member).Error == nil {
		var invite group_models.GroupInviteLinkModel
		if e := l.svcCtx.DB.Where("group_id = ? AND status = 1", group.GroupID).Order("id asc").First(&invite).Error; e == nil {
			domain := strings.TrimRight(strings.TrimSpace(l.svcCtx.Config.Domain), "/")
			if domain != "" && invite.Token != "" {
				resp.InviteUrl = fmt.Sprintf("%s/api/group/v1/invite_code?code=%s", domain, invite.Token)
			}
		}
	}
	return resp, nil
}
