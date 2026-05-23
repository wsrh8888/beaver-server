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

type DeleteBotLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除机器人
func NewDeleteBotLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteBotLogic {
	return &DeleteBotLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteBotLogic) DeleteBot(req *types.DeleteBotReq) (resp *types.DeleteBotRes, err error) {
	// 1. 从本地引用表查（通过 bot_id）
	var ref group_models.GroupBotModel
	if err = l.svcCtx.DB.Where("bot_id = ?", req.BotID).First(&ref).Error; err != nil {
		return nil, errors.New("机器人不存在")
	}

	// 2. 校验权限
	var member group_models.GroupMemberModel
	if err = l.svcCtx.DB.Take(&member, "group_id = ? AND user_id = ?", ref.GroupID, req.UserID).Error; err != nil {
		return nil, errors.New("不是群成员")
	}
	if member.Role != 1 && member.Role != 2 {
		return nil, errors.New("无权限，仅群主或管理员可删除机器人")
	}

	// 3. 通过 open_rpc 获取 Bot ID，然后删除
	botInfoRes, err := l.svcCtx.OpenRpc.GetBotInfo(l.ctx, &open_rpc.GetBotInfoReq{
		BotId: ref.BotID,
	})
	if err != nil {
		return nil, errors.New("Open Bot 记录不存在")
	}

	if _, err = l.svcCtx.OpenRpc.DeleteBot(l.ctx, &open_rpc.DeleteBotReq{
		Id: botInfoRes.Id,
	}); err != nil {
		return nil, errors.New("删除失败")
	}

	// 4. 删本地引用表
	l.svcCtx.DB.Delete(&ref)

	// 5. 将机器人移出群（软删除）
	l.svcCtx.DB.Model(&group_models.GroupMemberModel{}).
		Where("group_id = ? AND user_id = ?", ref.GroupID, ref.BotID).
		Update("status", 0)

	return &types.DeleteBotRes{}, nil
}
