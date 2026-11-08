package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	Mysql struct {
		DataSource string
	}
	// Import 文件导入解析参数。
	Import struct {
		// PdfBinPath pdftotext（poppler）可执行文件路径；缺失时 pdf 导入降级为仅登记附件不解析。
		PdfBinPath string `json:",optional,default=pdftotext"`
		// MaxBytes 单文件导入上限（字节），超过拒绝。
		MaxBytes int64 `json:",optional,default=20971520"`
	}
}
