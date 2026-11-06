package doc_models

import "beaver/common/models"

// DocYjsUpdate 客户端 Yjs 增量 update 日志。表 doc_yjs_updates。
// 实时协同网关把每条 update 追加到此，seq 文档内递增；断线重连/回放按 snapshot + 其后的 update 还原。
// update 为 Yjs 二进制（lib0 编码），跨端必须 pin 同一 Yjs 主版本保证互操作。
type DocYjsUpdate struct {
	models.Model
	DocID          string `gorm:"column:doc_id;type:varchar(64);not null;uniqueIndex:idx_doc_seq,priority:1;index;comment:所属文档" json:"docId"`
	Seq            int64  `gorm:"column:seq;type:bigint;not null;uniqueIndex:idx_doc_seq,priority:2;comment:文档内递增序号" json:"seq"`
	Update         []byte `gorm:"column:update;type:blob;not null;comment:Yjs update 二进制"`
	OriginUserID   string `gorm:"column:origin_user_id;type:varchar(64);not null;default:'';comment:产生该 update 的用户" json:"originUserId"`
}

func (DocYjsUpdate) TableName() string { return "doc_yjs_updates" }

// DocYjsSnapshot 压实后的 Yjs 全量 state。表 doc_yjs_snapshots。每文档一行（doc_id 唯一）。
// update_seq_up_to 记录已并入快照的最大 update seq；GetDocSyncState 返回 state + seq 大于该值的 update。
// content_json/content_text 为解码派生的权威正文快照，与正文保存同源。
type DocYjsSnapshot struct {
	models.Model
	DocID        string `gorm:"column:doc_id;type:varchar(64);uniqueIndex;not null;comment:所属文档" json:"docId"`
	State        []byte `gorm:"column:state;type:longblob;not null;comment:压实后的 Yjs 全量 state"`
	UpdateSeqUpTo int64 `gorm:"column:update_seq_up_to;type:bigint;not null;default:0;comment:已并入快照的最大 update seq" json:"updateSeqUpTo"`
	Version      int    `gorm:"type:int;not null;default:0;comment:快照对应文档版本" json:"version"`
}

func (DocYjsSnapshot) TableName() string { return "doc_yjs_snapshots" }
