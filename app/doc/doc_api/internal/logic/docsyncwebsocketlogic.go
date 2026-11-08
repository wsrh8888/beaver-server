package logic

import (
	"context"
	"net/http"
	"path"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"

	"github.com/zeromicro/go-zero/core/logx"
)

// goctl 会把 logx 写进 imports，占位避免未使用导入。
var _ = logx.WithContext

type DocSyncWebsocketLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

// 文档实时协同 WebSocket 网关（y-websocket sync + awareness）
func NewDocSyncWebsocketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DocSyncWebsocketLogic {
	return &DocSyncWebsocketLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("DocSyncWebsocket", ctx),
	}
}

// DocSyncWebsocket 把 HTTP 请求直接交给 ygo websocket.Server：
// 升级前鉴权与只读判定在 Server.Authorize 钩子内完成（见 ymsgw）。
// 这里只需把 URL path 改写为 "/"+docId，使 ygo 以 docId 作为房间名。
func (l *DocSyncWebsocketLogic) DocSyncWebsocket(req *types.WsReq, w http.ResponseWriter, r *http.Request) (resp *types.WsRes, err error) {
	if req.DocID == "" {
		l.logger.Error(model.LogMsg{Text: "文档协同缺少 docId", Data: map[string]any{"userId": req.UserID}})
		http.Error(w, "docId 不能为空", http.StatusBadRequest)
		return nil, nil
	}
	// 房间名 = docId：改写 path，令 ygo 的 path.Base 提取到 docId。
	r.URL.Path = path.Join("/", req.DocID)
	l.svcCtx.Gateway.ServeHTTP(w, r)
	return nil, nil
}
