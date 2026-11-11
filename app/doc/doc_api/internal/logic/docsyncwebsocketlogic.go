package logic

import (
	"context"
	"net/http"
	"path"
	"strings"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	"beaver/app/doc/doc_api/internal/ymsgw"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
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

// DocSyncWebsocket 把 HTTP 请求交给 ygo websocket.Server。
//
// 房间名由服务端按 docId@epoch 权威生成：epoch 取自 doc_rpc 的文档当前纪元，
// 客户端携带的 epoch 仅用于"是否过期"判定，不能自行指定房间。版本回滚会递增
// epoch，旧纪元连接一律 409，迫使客户端丢弃旧 Y.Doc 重建 provider，避免把回滚
// 掉的 CRDT 状态重新写回。升级前鉴权与只读判定在 Server.Authorize 内完成（见 ymsgw）。
func (l *DocSyncWebsocketLogic) DocSyncWebsocket(req *types.WsReq, w http.ResponseWriter, r *http.Request) (resp *types.WsRes, err error) {
	docID := strings.TrimSpace(req.DocID)
	if docID == "" {
		l.logger.Error(model.LogMsg{Text: "文档协同缺少 docId", Data: map[string]any{"userId": req.UserID}})
		http.Error(w, "docId 不能为空", http.StatusBadRequest)
		return nil, nil
	}

	// y-websocket 客户端可能额外携带 room 参数；不允许指向另一个文档，防串房间。
	if req.Room != "" {
		roomID, _ := ymsgw.ParseRoom(path.Base(req.Room))
		if roomID != docID {
			l.logger.Error(model.LogMsg{Text: "协同房间与 docId 不一致", Data: map[string]any{"docId": docID, "room": req.Room}})
			http.Error(w, "协同房间与 docId 不一致", http.StatusBadRequest)
			return nil, nil
		}
	}

	// 以 doc_rpc 当前纪元为准。Epoch=0 表示客户端未缓存纪元，直接返回当前值。
	state, serr := l.svcCtx.Doc.ResolveRoomEpoch(context.Background(), &doc_rpc.ResolveRoomEpochReq{
		DocId: docID,
		Epoch: int32(req.Epoch),
	})
	if serr != nil {
		l.logger.Error(model.LogMsg{Text: "解析协同房间纪元失败", Data: map[string]any{"docId": docID, "error": serr}})
		http.Error(w, "协同服务暂不可用", http.StatusBadGateway)
		return nil, nil
	}
	if !state.GetValid() {
		// 409 是客户端重建 provider 的信号，不是鉴权失败。
		http.Error(w, "文档协同纪元已过期，请刷新后重新连接", http.StatusConflict)
		return nil, nil
	}
	epoch := state.GetEpoch()
	if epoch <= 0 {
		epoch = 1
	}

	room := ymsgw.BuildRoom(docID, epoch)
	// ygo 用 path.Base(r.URL.Path) 取房间名，改写 path 为权威 docId@epoch。
	r.URL.Path = path.Join("/", room)
	// 同步把权威纪元/房间写回 query，便于下游日志与客户端诊断对齐。
	q := r.URL.Query()
	q.Set("docId", docID)
	q.Set("epoch", ymsgw.EpochString(epoch))
	q.Set("room", room)
	r.URL.RawQuery = q.Encode()

	l.svcCtx.Gateway.ServeHTTP(w, r)
	return nil, nil
}
