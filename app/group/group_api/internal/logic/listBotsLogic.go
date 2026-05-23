package logic

import (
	"context"
	"errors"
	"time"

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/group/group_models"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListBotsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取群内所有机器人列表
func NewListBotsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListBotsLogic {
	return &ListBotsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListBotsLogic) ListBots(req *types.ListBotsReq) (resp *types.ListBotsRes, err error) {
	var member group_models.GroupMemberModel
	if err = l.svcCtx.DB.Take(&member, "group_id = ? AND user_id = ?", req.GroupID, req.UserID).Error; err != nil {
		return nil, errors.New("不是群成员")
	}
	if member.Role != 1 && member.Role != 2 {
		return nil, errors.New("无权限，仅群主或管理员可查看机器人")
	}

	// 查询群内机器人
	var bots []group_models.GroupBotModel
	if err = l.svcCtx.DB.Where("group_id = ?", req.GroupID).
		Order("id DESC").Find(&bots).Error; err != nil {
		return nil, err
	}

	// 批量获取用户信息（通过 user_rpc）
	botIDs := make([]string, 0, len(bots))
	for _, b := range bots {
		botIDs = append(botIDs, b.BotID)
	}

	userRes, err := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{
		UserIdList: botIDs,
	})
	if err != nil {
		return nil, err
	}

	items := make([]types.ListBotsItem, 0, len(bots))
	for _, b := range bots {
		userInfo := userRes.UserInfo[b.BotID]
		if userInfo == nil {
			continue
		}

		items = append(items, types.ListBotsItem{
			BotID:       b.BotID,
			Name:        userInfo.NickName,
			Description: "", // TODO: 从 open_rpc 获取
			Avatar:      userInfo.Avatar,
			Type:        b.Type,
			Status:      b.Status,
			CreatedAt:   time.Time(b.CreatedAt).Unix(),
		})
	}

	return &types.ListBotsRes{List: items}, nil
}
