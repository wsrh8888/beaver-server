package logic

import (
	"context"
	"strconv"

	"beaver/app/ai/ai_models"
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteAIBotLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除AI机器人
func NewDeleteAIBotLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAIBotLogic {
	return &DeleteAIBotLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteAIBotLogic) DeleteAIBot(req *types.DeleteAIBotReq) (resp *types.DeleteAIBotRes, err error) {
	// 解析机器人ID
	botID, err := strconv.ParseUint(req.BotID, 10, 64)
	if err != nil {
		return nil, err
	}

	// 查询机器人是否存在
	var bot ai_models.AIModelConfig
	err = l.svcCtx.DB.First(&bot, botID).Error
	if err != nil {
		return nil, err
	}

	// 删除机器人
	err = l.svcCtx.DB.Delete(&bot).Error
	if err != nil {
		return nil, err
	}

	return &types.DeleteAIBotRes{}, nil
}
