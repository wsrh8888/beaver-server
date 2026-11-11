package ymsgw

import (
	"context"
	"encoding/json"
	"sort"
	"strings"

	doccli "beaver/app/doc/doc_rpc/doc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"

	"github.com/reearth/ygo/crdt"
	ygws "github.com/reearth/ygo/provider/websocket"
)

// textTypeName 是 TipTap/Yjs 协同正文使用的根类型名（与客户端 ydoc 对齐）。
// 现代 TipTap Collaboration 使用 Y.XmlFragment；旧纯文本协同使用 YText。
// 两者同名，运行时按 CRDT 实际根类型择一读取（见 deriveContent）。
const textTypeName = "prosemirror"

// persistence 是 ygo websocket.Server 的持久化适配器，把房间读写落到 doc_rpc 权威存储。
//   - LoadDoc    -> GetDocSyncState：拉快照 + 增量，合并为一份完整 V1 update 供 ygo 应用。
//   - StoreUpdate-> AppendDocUpdate：把每条客户端增量写进增量日志（doc_yjs_updates）。
//   - Compact    -> CommitDocSnapshot：压实全量快照，并尽力解码正文更新权威文档 + 生成版本。
//
// 房间名统一为 docId@epoch：版本回滚会递增 epoch，旧房间的增量/快照一律隔离，
// 不会把回滚掉的正文重新写回。实现 ygws.PersistenceAdapter + ygws.CompactableAdapter。
// srv 反向引用在 New 后绑定，Compact 通过它读取房间内的实时 crdt.Doc 以解码正文。
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
	docID, epoch := ParseRoom(room)
	res, err := p.doc.GetDocSyncState(context.Background(), &doc_rpc.GetDocSyncStateReq{
		DocId: docID,
		Epoch: epoch,
	})
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

// StoreUpdate 追加一条增量 update 到权威增量日志（限定在当前纪元）。
func (p *persistence) StoreUpdate(room string, update []byte) error {
	if len(update) == 0 {
		return nil
	}
	docID, epoch := ParseRoom(room)
	_, err := p.doc.AppendDocUpdate(context.Background(), &doc_rpc.AppendDocUpdateReq{
		DocId:  docID,
		Update: update,
		Epoch:  epoch,
	})
	return err
}

// StoreUpdateContext 是 StoreUpdate 的 context 感知变体，配合 Server.Shutdown 中止在途写。
func (p *persistence) StoreUpdateContext(ctx context.Context, room string, update []byte) error {
	if len(update) == 0 {
		return nil
	}
	docID, epoch := ParseRoom(room)
	_, err := p.doc.AppendDocUpdate(ctx, &doc_rpc.AppendDocUpdateReq{
		DocId:  docID,
		Update: update,
		Epoch:  epoch,
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
	docID, epoch := ParseRoom(room)
	_, err := p.doc.CommitDocSnapshot(ctx, &doc_rpc.CommitDocSnapshotReq{
		DocId:       docID,
		Snapshot:    state,
		ContentJson: contentJSON,
		ContentText: contentText,
		Epoch:       epoch,
	})
	return err
}

// ── ProseMirror 正文派生 ──────────────────────────────────────────────────────

type pmNode struct {
	Type    string         `json:"type"`
	Attrs   map[string]any `json:"attrs,omitempty"`
	Content []*pmNode      `json:"content,omitempty"`
	Text    string         `json:"text,omitempty"`
	Marks   []pmMark       `json:"marks,omitempty"`
}

type pmMark struct {
	Type  string         `json:"type"`
	Attrs map[string]any `json:"attrs,omitempty"`
}

// xmlMarkAlias 兼容 HTML 风格标签；y-prosemirror 默认直接用 PM mark.type.name 作标签名。
var xmlMarkAlias = map[string]string{
	"strong": "bold",
	"b":      "bold",
	"em":     "italic",
	"i":      "italic",
	"del":    "strike",
	"s":      "strike",
	"u":      "underline",
	"a":      "link",
	"mark":   "highlight",
}

// xmlMarkType 判断一个 YXmlElement 是否是内联 mark 包装节点。
// 块级/叶子节点（paragraph、heading、image、hardBreak 等）走普通 node 分支。
func xmlMarkType(name string) (string, bool) {
	if alias, ok := xmlMarkAlias[name]; ok {
		return alias, true
	}
	switch name {
	case "bold", "italic", "strike", "underline", "link", "code", "highlight",
		"textStyle", "subscript", "superscript", "comment", "suggestion":
		return name, true
	}
	return "", false
}

// xmlFragmentRoot 安全读取 Y.XmlFragment 根类型。
// ygo 的根类型名互斥：同名 root 若是 YText，GetXmlFragment 会 panic，
// 这里用 recover 转成"不可用"，让 deriveContent 回退到纯文本派生。
func xmlFragmentRoot(doc *crdt.Doc) (frag *crdt.YXmlFragment) {
	defer func() { _ = recover() }()
	return doc.GetXmlFragment(textTypeName)
}

// textRoot 安全读取 YText 根类型（旧纯文本协同）。
func textRoot(doc *crdt.Doc) (txt *crdt.YText) {
	defer func() { _ = recover() }()
	return doc.GetText(textTypeName)
}

// deriveContent 从房间 crdt.Doc 派生 ProseMirror JSON 快照与纯文本。
// 优先读取 Y.XmlFragment('prosemirror')（TipTap Collaboration 的实际结构），
// 完整保留块类型、块属性、内联 mark 与富文本格式；fragment 为空时回退到
// 根 YText，兼容旧纯文本协同数据。两者皆空则返回空串（快照-only，不动正文/版本）。
func deriveContent(doc *crdt.Doc) (contentJSON, contentText string) {
	if frag := xmlFragmentRoot(doc); frag != nil {
		root := &pmNode{Type: "doc", Content: normalizeTopLevel(pmNodesFrom(any(frag), nil))}
		if len(root.Content) > 0 {
			contentText = pmPlainText(root)
			buf, err := json.Marshal(root)
			if err != nil {
				return "", contentText
			}
			return string(buf), contentText
		}
	}
	// 回退：根 YText 纯文本，按换行还原成多段落 block 结构。
	if txt := textRoot(doc); txt != nil {
		contentText = txt.ToString()
		if trim(contentText) != "" {
			root := pmFromPlainText(contentText)
			buf, err := json.Marshal(root)
			if err != nil {
				return "", contentText
			}
			return string(buf), contentText
		}
	}
	return "", ""
}

// pmNodesFrom 把 XML 节点展开成 ProseMirror 节点序列。
// inherited 是外层 mark 元素累积下来的 marks；一个 YXmlText 可能按格式切成多个 text 节点。
// 注意：ygo 的 Children() 返回未导出接口切片，这里通过 any 逐个分发，避免依赖内部类型。
func pmNodesFrom(n any, inherited []pmMark) []*pmNode {
	switch v := n.(type) {
	case *crdt.YXmlFragment:
		var out []*pmNode
		for _, child := range v.Children() {
			out = append(out, pmNodesFrom(any(child), inherited)...)
		}
		return out

	case *crdt.YXmlElement:
		attrs := v.GetAttributeValues()
		if markType, isMark := xmlMarkType(v.NodeName); isMark {
			marks := append(append([]pmMark{}, inherited...), pmMark{Type: markType, Attrs: attrs})
			var out []*pmNode
			for _, child := range v.Children() {
				out = append(out, pmNodesFrom(any(child), marks)...)
			}
			return out
		}
		node := &pmNode{Type: v.NodeName, Attrs: attrs}
		for _, child := range v.Children() {
			node.Content = append(node.Content, pmNodesFrom(any(child), inherited)...)
		}
		return []*pmNode{node}

	case *crdt.YXmlText:
		var out []*pmNode
		for _, d := range v.ToDelta() {
			s, ok := d.Insert.(string)
			if d.Op != crdt.DeltaOpInsert || !ok || s == "" {
				continue
			}
			marks := append([]pmMark{}, inherited...)
			marks = append(marks, marksFromAttributes(d.Attributes)...)
			node := &pmNode{Type: "text", Text: s}
			if len(marks) > 0 {
				node.Marks = marks
			}
			out = append(out, node)
		}
		return out
	}
	return nil
}

// marksFromAttributes 把 YText 富文本格式属性转成 PM marks。
// 值 true 表示开关型 mark（无 attrs），其余值按同名属性保留。
// 按 key 排序保证同一 CRDT 状态派生出的 JSON 稳定，避免无意义的版本抖动。
func marksFromAttributes(attrs crdt.Attributes) []pmMark {
	if len(attrs) == 0 {
		return nil
	}
	keys := make([]string, 0, len(attrs))
	for k := range attrs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var out []pmMark
	for _, k := range keys {
		v := attrs[k]
		if v == nil {
			continue
		}
		m := pmMark{Type: k}
		if alias, ok := xmlMarkAlias[k]; ok {
			m.Type = alias
		}
		if b, ok := v.(bool); !ok || !b {
			m.Attrs = map[string]any{k: v}
		}
		out = append(out, m)
	}
	return out
}

// normalizeTopLevel 保证 doc 根下是 block 节点：裸 text/hardBreak 序列包成 paragraph。
func normalizeTopLevel(nodes []*pmNode) []*pmNode {
	if len(nodes) == 0 {
		return nodes
	}
	for _, n := range nodes {
		switch n.Type {
		case "text", "hardBreak", "image", "mention", "emoji":
			continue
		default:
			return nodes
		}
	}
	return []*pmNode{{Type: "paragraph", Content: nodes}}
}

// pmPlainText 从 PM 树拼出搜索/摘要/知识库索引用的纯文本，块之间以换行分隔。
func pmPlainText(root *pmNode) string {
	var sb strings.Builder
	for i, n := range root.Content {
		appendPMText(n, &sb)
		if i < len(root.Content)-1 && !strings.HasSuffix(sb.String(), "\n") {
			sb.WriteString("\n")
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}

func appendPMText(n *pmNode, sb *strings.Builder) {
	if n.Text != "" {
		sb.WriteString(n.Text)
		return
	}
	if n.Type == "hardBreak" {
		sb.WriteString("\n")
		return
	}
	for _, c := range n.Content {
		appendPMText(c, sb)
	}
}

// pmFromPlainText 把纯文本按换行还原成多段落 doc（旧 YText 协同兼容路径）。
func pmFromPlainText(text string) *pmNode {
	root := &pmNode{Type: "doc"}
	for _, line := range splitLines(text) {
		para := &pmNode{Type: "paragraph"}
		if line != "" {
			para.Content = []*pmNode{{Type: "text", Text: line}}
		}
		root.Content = append(root.Content, para)
	}
	return root
}
