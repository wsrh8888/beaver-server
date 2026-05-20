package logic

import (
	"context"
	"errors"

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/group/group_models"
	"beaver/app/open/open_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteNotificationBotLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除通知机器人
func NewDeleteNotificationBotLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteNotificationBotLogic {
	return &DeleteNotificationBotLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteNotificationBotLogic) DeleteNotificationBot(req *types.DeleteNotificationBotReq) (resp *types.DeleteNotificationBotRes, err error) {
	var record open_models.OpenIncomingWebhook
	if err = l.svcCtx.DB.Where("id = ? AND app_id = ?", req.ID, "GROUP_NOTIFICATION").First(&record).Error; err != nil {
		return nil, errors.New("通知机器人不存在")
	}

	var member group_models.GroupMemberModel
	if err = l.svcCtx.DB.Take(&member, "group_id = ? AND user_id = ?", record.GroupID, req.UserID).Error; err != nil {
		return nil, errors.New("不是群成员")
	}
	if member.Role != 1 && member.Role != 2 {
		return nil, errors.New("无权限，仅群主或管理员可删除通知机器人")
	}

	if err = l.svcCtx.DB.Delete(&record).Error; err != nil {
		return nil, errors.New("删除失败")
	}

	return &types.DeleteNotificationBotRes{Success: true}, nil
}
