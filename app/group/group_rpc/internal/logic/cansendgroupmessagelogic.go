package logic

import (
	"context"
	"strings"
	"time"

	"beaver/app/group/group_models"
	"beaver/app/group/group_rpc/internal/svc"
	"beaver/app/group/group_rpc/types/group_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type CanSendGroupMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCanSendGroupMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CanSendGroupMessageLogic {
	return &CanSendGroupMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CanSendGroupMessageLogic) CanSendGroupMessage(in *group_rpc.CanSendGroupMessageReq) (*group_rpc.CanSendGroupMessageRes, error) {
	var member group_models.GroupMemberModel
	if err := l.svcCtx.DB.Where("group_id = ? AND user_id = ? AND status = 1", in.GroupId, in.UserId).
		First(&member).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &group_rpc.CanSendGroupMessageRes{Allowed: false, Reason: "你不是该群成员"}, nil
		}
		return nil, err
	}

	if member.Role == 1 || member.Role == 2 || strings.HasPrefix(in.UserId, "nbot_") {
		return &group_rpc.CanSendGroupMessageRes{Allowed: true}, nil
	}

	var group group_models.GroupModel
	if err := l.svcCtx.DB.Where("group_id = ?", in.GroupId).First(&group).Error; err == nil {
		if group.IsMuteAll {
			return &group_rpc.CanSendGroupMessageRes{Allowed: false, Reason: "当前群已开启全员禁言"}, nil
		}
	}

	if member.MutedUntil != nil && member.MutedUntil.After(time.Now()) {
		return &group_rpc.CanSendGroupMessageRes{Allowed: false, Reason: "你已被禁言"}, nil
	}

	return &group_rpc.CanSendGroupMessageRes{Allowed: true}, nil
}
