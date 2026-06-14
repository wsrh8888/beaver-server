package logic

import (
	"context"
	"errors"

	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/datasync/datasync_api/internal/svc"
	"beaver/app/datasync/datasync_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSyncMessageMediasLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSyncMessageMediasLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSyncMessageMediasLogic {
	return &GetSyncMessageMediasLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSyncMessageMediasLogic) GetSyncMessageMedias(req *types.GetSyncMessageMediasReq) (*types.GetSyncMessageMediasRes, error) {
	if req.UserID == "" {
		return nil, errors.New("用户ID不能为空")
	}

	rpcResp, err := l.svcCtx.ChatRpc.GetSyncMessageMedias(l.ctx, &chat_rpc.GetSyncMessageMediasReq{
		UserId: req.UserID,
		Since:  req.Since,
	})
	if err != nil {
		l.Logger.Errorf("同步消息媒体状态失败: userId=%s, error=%v", req.UserID, err)
		return nil, errors.New("同步失败")
	}

	messageIDs := make([]string, 0)
	if rpcResp.MessageIds != nil {
		messageIDs = rpcResp.MessageIds
	}

	return &types.GetSyncMessageMediasRes{
		MessageIDs:      messageIDs,
		ServerTimestamp: rpcResp.ServerTimestamp,
	}, nil
}
