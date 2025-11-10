package logic

import (
	"context"

	"beaver/app/chat/chat_api/internal/svc"
	"beaver/app/chat/chat_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetConversationsListByIdsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量获取会话数据
func NewGetConversationsListByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationsListByIdsLogic {
	return &GetConversationsListByIdsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetConversationsListByIdsLogic) GetConversationsListByIds(req *types.GetConversationsListByIdsReq) (resp *types.GetConversationsListByIdsRes, err error) {
	// todo: add your logic here and delete this line

	return
}
