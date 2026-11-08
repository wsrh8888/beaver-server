package ymsgw

import (
	"context"
	"encoding/json"

	doccli "beaver/app/doc/doc_rpc/doc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"

	"github.com/reearth/ygo/crdt"
	ygws "github.com/reearth/ygo/provider/websocket"
)

// textTypeName TipTap/Yjs 协同正文使用的 YText 根类型名（与客户端 ydoc.getText 对齐）。
const textTypeName = "prosemirror"

// persistence 是 ygo websocket.Server 的持久化适配器，把房间读写落到 doc_rpc 权威存储。
//   - LoadDoc    -> GetDocSyncState：拉快照 + 增量，合并为一份完整 V1 update 供 ygo 应用。
//   - StoreUpdate-> AppendDocUpdate：把每条客户端增量写进增量日志（doc_yjs_updates）。
//   - Compact    -> CommitDocSnapshot：压实全量快照，并尽力解码正文更新权威文档 + 生成版本。
//
// 实现 ygws.PersistenceAdapter + ygws.CompactableAdapter。srv 反向引用在 New 后绑定，
// Compact 通过它读取房间内的实时 crdt.Doc 以解码正文。
type persistence struct {
	doc doccli.Doc
	srv *ygws.Server
}

var (
	_ ygws.PersistenceAdapter = (*persistence)(nil)
	_ ygws.CompactableAdapter = (*persistence)(nil)
)

// LoadDoc 返回房间当前完整状态的 V1 update；无历史状态时返回 (nil, nil)。
func (p *persistence) LoadDoc(room string) ([]byte, error) {
	res, err := p.doc.GetDocSyncState(context.Background(), &doc_rpc.GetDocSyncStateReq{DocId: room})
	if err != nil {
		return nil, err
	}
	var parts [][]byte
	if len(res.GetSnapshot()) > 0 {
		parts = append(parts, res.GetSnapshot())
	}
	for _, u := range res.GetUpdates() {
		if len(u) > 0 {
			parts = append(parts, u)
		}
	}
	switch len(parts) {
	case 0:
		return nil, nil
	case 1:
		return parts[0], nil
	default:
		return crdt.MergeUpdatesV1(parts...)
	}
}

// StoreUpdate 追加一条增量 update 到权威增量日志。
func (p *persistence) StoreUpdate(room string, update []byte) error {
	if len(update) == 0 {
		return nil
	}
	_, err := p.doc.AppendDocUpdate(context.Background(), &doc_rpc.AppendDocUpdateReq{
		DocId:  room,
		Update: update,
	})
	return err
}

// StoreUpdateContext 是 StoreUpdate 的 context 感知变体，配合 Server.Shutdown 中止在途写。
func (p *persistence) StoreUpdateContext(ctx context.Context, room string, update []byte) error {
	if len(update) == 0 {
		return nil
	}
	_, err := p.doc.AppendDocUpdate(ctx, &doc_rpc.AppendDocUpdateReq{
		DocId:  room,
		Update: update,
	})
	return err
}

// Compact 压实房间状态：写全量快照，并尽力解码正文更新权威文档。
// 无正文（content_json 为空）时 doc_rpc 只做快照压实，不改动 content/version。
func (p *persistence) Compact(ctx context.Context, room string) error {
	doc := p.srv.GetDoc(room)
	if doc == nil {
		return nil
	}
	state := doc.EncodeStateAsUpdate()
	contentJSON, contentText := deriveContent(doc)
	_, err := p.doc.CommitDocSnapshot(ctx, &doc_rpc.CommitDocSnapshotReq{
		DocId:       room,
		Snapshot:    state,
		ContentJson: contentJSON,
		ContentText: contentText,
	})
	return err
}

// deriveContent 从房间 crdt.Doc 派生 ProseMirror block 快照与纯文本。
// 现阶段取根 YText 文本，包装成单段落 doc；正文为空时返回空串（快照-only）。
// TipTap 富 block 结构的完整还原随 M4 客户端编辑器对齐后细化。
func deriveContent(doc *crdt.Doc) (contentJSON, contentText string) {
	txt := doc.GetText(textTypeName)
	contentText = txt.ToString()
	if trim(contentText) == "" {
		return "", ""
	}
	type pmText struct {
		Type string `json:"type"`
		Text string `json:"text"`
	}
	type pmPara struct {
		Type    string   `json:"type"`
		Content []pmText `json:"content"`
	}
	type pmDoc struct {
		Type    string   `json:"type"`
		Content []pmPara `json:"content"`
	}
	// 按换行切段，还原为多段落 block 结构。
	var paras []pmPara
	for _, line := range splitLines(contentText) {
		para := pmPara{Type: "paragraph"}
		if line != "" {
			para.Content = []pmText{{Type: "text", Text: line}}
		}
		paras = append(paras, para)
	}
	payload := pmDoc{Type: "doc", Content: paras}
	buf, err := json.Marshal(payload)
	if err != nil {
		return "", contentText
	}
	return string(buf), contentText
}
