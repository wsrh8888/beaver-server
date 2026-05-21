package logic

import (
	"context"
	"errors"

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/group/group_models"
	"beaver/app/open/open_rpc/types/open_rpc"

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
	// 从本地引用表查（group 侧自有数据，无需跨服务）
	var ref group_models.GroupNotificationBotModel
	if err = l.svcCtx.DB.First(&ref, req.ID).Error; err != nil {
		return nil, errors.New("通知机器人不存在")
	}

	var member group_models.GroupMemberModel
	if err = l.svcCtx.DB.Take(&member, "group_id = ? AND user_id = ?", ref.GroupID, req.UserID).Error; err != nil {
		return nil, errors.New("不是群成员")
	}
	if member.Role != 1 && member.Role != 2 {
		return nil, errors.New("无权限，仅群主或管理员可删除通知机器人")
	}

	// 调 open_rpc 删除 master 记录
	if _, err = l.svcCtx.OpenRpc.DeleteWebhook(l.ctx, &open_rpc.DeleteWebhookReq{
		Id: uint32(ref.WebhookID),
	}); err != nil {
		return nil, errors.New("删除失败")
	}

	// 删本地引用表
	l.svcCtx.DB.Delete(&ref)

	// 将机器人移出群（软删除）
	l.svcCtx.DB.Model(&group_models.GroupMemberModel{}).
		Where("group_id = ? AND user_id = ?", ref.GroupID, ref.BotUserID).
		Update("status", 0)

	return &types.DeleteNotificationBotRes{Success: true}, nil
}
