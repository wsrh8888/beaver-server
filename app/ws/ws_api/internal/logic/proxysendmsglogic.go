package logic

import (
	"context"
	"encoding/json"
	"errors"

	ws_conn "beaver/app/ws/ws_api/internal/logic/websocket/conn"
	"beaver/app/ws/ws_api/internal/svc"
	"beaver/app/ws/ws_api/internal/types"
	type_struct "beaver/app/ws/ws_api/types"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProxySendMsgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProxySendMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProxySendMsgLogic {
	return &ProxySendMsgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProxySendMsgLogic) ProxySendMsg(req *types.ProxySendMsgReq) (resp *types.ProxySendMsgRes, err error) {
	// 内部接口鉴权：只允许内部服务调用
	if req.InternalSecret != l.svcCtx.Config.InternalSecret {
		return nil, errors.New("unauthorized")
	}

	bodyBytes, err := json.Marshal(req.Body)
	if err != nil {
		return nil, err
	}

	content := type_struct.WsContent{
		Data: type_struct.WsData{
			Type:           wsTypeConst.Type(req.Type),
			Body:           bodyBytes,
			ConversationID: req.ConversationId,
		},
	}

	ws_conn.SendMsgToUser(req.TargetID, wsCommandConst.Command(req.Command), content)

	return &types.ProxySendMsgRes{}, nil
}
