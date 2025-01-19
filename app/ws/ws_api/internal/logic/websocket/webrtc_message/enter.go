package webrtc_message

import (
	"context"
	"fmt"
	"net/http"

	"beaver/app/ws/ws_api/internal/svc"
	"beaver/app/ws/ws_api/internal/types"
	type_struct "beaver/app/ws/ws_api/types"

	"github.com/gorilla/websocket"
)

func HandleWebRTCAnswer(ctx context.Context, svcCtx *svc.ServiceContext, req *types.WsReq, r *http.Request, conn *websocket.Conn, content type_struct.WsContent) {
	switch content.Data.Type {
	default:
		fmt.Println("未支持的 WebRTC 消息类型", content.Data.Type)
	}
}
