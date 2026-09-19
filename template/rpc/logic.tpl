package {{.packageName}}

import (
	"context"

	{{.imports}}

	beaverlog "beaver/utils/beaverlog"
)

type {{.logicName}} struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func New{{.logicName}}(ctx context.Context, svcCtx *svc.ServiceContext) *{{.logicName}} {
	return &{{.logicName}}{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("{{.logicName}}", ctx),
	}
}
{{.functions}}
