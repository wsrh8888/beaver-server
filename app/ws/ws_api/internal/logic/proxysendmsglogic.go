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

	// 打印内容

	fmt.Println("消息内容：", string(bodyBytes))
	fmt.Println("发送者ID：", req.UserID, "，目标ID：", req.TargetID)

	fmt.Println("命令类型：", req.Command, "，消息类型：", req.Type, "，会话ID：", req.ConversationId)

	websocket_utils.SendMsgToUser(req.TargetID, wsCommandConst.Command(req.Command), content)

	return
}
