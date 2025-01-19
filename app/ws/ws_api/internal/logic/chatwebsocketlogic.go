package logic

import (
	"context"
	"fmt"
	"net/http"

	ws "beaver/app/ws/ws_api/internal/logic/websocket"
	websocket_utils "beaver/app/ws/ws_api/internal/logic/websocket/utils"
	"beaver/app/ws/ws_api/internal/svc"
	"beaver/app/ws/ws_api/internal/types"

	"github.com/gorilla/websocket"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChatWebsocketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatWebsocketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatWebsocketLogic {
	return &ChatWebsocketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChatWebsocketLogic) ChatWebsocket(req *types.WsReq, w http.ResponseWriter, r *http.Request) (resp *types.WsRes, err error) {
	conn, err := UpgradeToWebSocket(w, r)
	if err != nil {
		return
	}
	defer cleanupConnection(req.UserID, conn)
	fmt.Println("用户上线", req.UserID)
	userWsInfo, ok := websocket_utils.UserOnlineWsMap[req.UserID]
	addr := conn.RemoteAddr().String()

	if !ok {
		userWsInfo = &websocket_utils.UserWsInfo{
			WsClientMap: map[string]*websocket.Conn{
				addr: conn,
			},
		}
		websocket_utils.UserOnlineWsMap[req.UserID] = userWsInfo
	} else {
		// 则替换当前连接
		userWsInfo.WsClientMap[addr] = conn
	}

	ws.HandleWebSocketMessages(l.ctx, l.svcCtx, req, r, conn)
	return &types.WsRes{}, nil
}

func UpgradeToWebSocket(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	upGrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func cleanupConnection(userID string, conn *websocket.Conn) {
	conn.Close()
	addr := conn.RemoteAddr().String()

	userWsInfo, ok := websocket_utils.UserOnlineWsMap[userID]
	if ok {
		delete(userWsInfo.WsClientMap, addr)
	}
	if userWsInfo != nil && len(userWsInfo.WsClientMap) == 0 {
		delete(websocket_utils.UserOnlineWsMap, userID)
	}
}
