package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	Mysql struct {
		DataSource string
	}
	// KnowledgeRpc：把文档正文单向投影为 RAG 索引（Create/Update/Delete Document + Reindex）。
	KnowledgeRpc zrpc.RpcClientConf
	// RAG 投影调参。Enabled=false 时 doc 保存不触发 knowledge 同步（纯文档，不检索）。
	Rag struct {
		Enabled bool `json:",optional,default=true"`
		// ChunkIndex 与 knowledge 域共享的分段索引名，仅记录用途（实际投影走 knowledge_rpc）。
		CollectionKind string `json:",optional,default=docs"`
	}
	// Import 文件导入解析参数。
	Import struct {
		// PdfBinPath pdftotext（poppler）可执行文件路径；缺失时 pdf 导入降级为仅登记附件不解析。
		PdfBinPath string `json:",optional,default=pdftotext"`
		// MaxBytes 单文件导入上限（字节），超过拒绝。
		MaxBytes int64 `json:",optional,default=20971520"`
	}
}
