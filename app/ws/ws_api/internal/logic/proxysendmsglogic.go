package logic

import (
	"context"
	"encoding/json"
	"fmt"

	websocket_utils "beaver/app/ws/ws_api/internal/logic/websocket/utils"
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
	fmt.Println("收到ws转发的消息")
	// 将map转换为json.RawMessage
	bodyBytes, err := json.Marshal(req.Body)
	if err != nil {
		return nil, err
	}

	content := type_struct.WsContent{
		Timestamp: 0,
		Data: type_struct.WsData{
			Type:           wsTypeConst.Type(req.Type),
			Body:           bodyBytes,
			ConversationID: req.ConversationId,
		},
	}
	// todo: add your logic here and delete this line
	websocket_utils.SendMsgByUser(req.TargetID, req.UserID, wsCommandConst.Command(req.Command), content)
	return
}
