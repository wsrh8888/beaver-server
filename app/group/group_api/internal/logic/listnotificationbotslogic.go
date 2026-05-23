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

type ListNotificationBotsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取群内所有通知机器人列表
func NewListNotificationBotsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListNotificationBotsLogic {
	return &ListNotificationBotsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListNotificationBotsLogic) ListNotificationBots(req *types.ListNotificationBotsReq) (resp *types.ListNotificationBotsRes, err error) {
	var member group_models.GroupMemberModel
	if err = l.svcCtx.DB.Take(&member, "group_id = ? AND user_id = ?", req.GroupID, req.UserID).Error; err != nil {
		return nil, errors.New("不是群成员")
	}
	if member.Role != 1 && member.Role != 2 {
		return nil, errors.New("无权限，仅群主或管理员可查看通知机器人")
	}

	// 直接读本地引用表，无需跨服务调用
	var bots []group_models.GroupBotModel
	if err = l.svcCtx.DB.Where("group_id = ?", req.GroupID).
		Order("id DESC").Find(&bots).Error; err != nil {
		return nil, err
	}

	items := make([]types.ListNotificationBotsItem, 0, len(bots))
	for _, b := range bots {
		items = append(items, types.ListNotificationBotsItem{
			ID:          int64(b.Id),
			Name:        b.Name,
			Description: b.Description,
			Avatar:      b.Avatar,
			Type:        b.Type,
			Status:      b.Status,
			CreatedAt:   time.Time(b.CreatedAt).Unix(),
		})
	}

	return &types.ListNotificationBotsRes{List: items}, nil
}
