package ymsgw

import (
	"encoding/hex"
	"testing"

	"github.com/reearth/ygo/crdt"
)

// jsUpdateFixture 是一段由真实客户端 yjs@13.x（lib0 编码）产出的 Yjs update，
// 对 YText("prosemirror") 插入文本 "协同编辑 Hello — 海狸 🦫"（含中文/emoji/em-dash）。
// 由 spike 生成并固化，作为服务端解码正文一致性的回归断言（无需测试期依赖 node）。
const jsUpdateFixture = "010188bbb3e90c0004010b70726f73656d6972726f7222e58d8fe5908ce7bc96e8be912048656c6c6f20e2809420e6b5b7e78bb820f09fa6ab00"

const jsUpdateWantText = "协同编辑 Hello — 海狸 🦫"

// TestInterop_DecodeClientUpdate 断言：客户端产出的 update 被 Go 服务端解码后正文完全一致。
func TestInterop_DecodeClientUpdate(t *testing.T) {
	blob, err := hex.DecodeString(jsUpdateFixture)
	if err != nil {
		t.Fatalf("fixture 解码失败: %v", err)
	}
	doc := crdt.New()
	if err := doc.ApplyUpdate(blob); err != nil {
		t.Fatalf("ApplyUpdate 客户端 update 失败: %v", err)
	}
	if got := doc.GetText(textTypeName).ToString(); got != jsUpdateWantText {
		t.Fatalf("解码正文不一致:\n got=%q\nwant=%q", got, jsUpdateWantText)
	}
}

// TestInterop_GoUpdateRoundTrip 断言：Go 产出的 update 可被再次应用并收敛到同一正文。
func TestInterop_GoUpdateRoundTrip(t *testing.T) {
	src := crdt.New()
	txt := src.GetText(textTypeName)
	src.Transact(func(txn *crdt.Transaction) {
		txt.Insert(txn, 0, jsUpdateWantText, nil)
	})
	update := src.EncodeStateAsUpdate()

	dst := crdt.New()
	if err := dst.ApplyUpdate(update); err != nil {
		t.Fatalf("ApplyUpdate Go update 失败: %v", err)
	}
	if got := dst.GetText(textTypeName).ToString(); got != jsUpdateWantText {
		t.Fatalf("往返正文不一致:\n got=%q\nwant=%q", got, jsUpdateWantText)
	}
}
