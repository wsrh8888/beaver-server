package logic

import (
	"context"

	"beaver/app/chat/chat_api/internal/svc"
	"beaver/app/chat/chat_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserConversationSettingsListByIdsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量获取用户会话设置数据
func NewGetUserConversationSettingsListByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserConversationSettingsListByIdsLogic {
	return &GetUserConversationSettingsListByIdsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserConversationSettingsListByIdsLogic) GetUserConversationSettingsListByIds(req *types.GetUserConversationSettingsListByIdsReq) (resp *types.GetUserConversationSettingsListByIdsRes, err error) {
	// todo: add your logic here and delete this line

	return
}
