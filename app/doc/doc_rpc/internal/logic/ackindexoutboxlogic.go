package logic

import (
	"context"
	"strings"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"
	beaverlog "beaver/utils/beaverlog"

	"gorm.io/gorm"
)

type AckIndexOutboxLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewAckIndexOutboxLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AckIndexOutboxLogic {
	return &AckIndexOutboxLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("AckIndexOutboxLogic", ctx),
	}
}

func (l *AckIndexOutboxLogic) AckIndexOutbox(in *doc_rpc.AckIndexOutboxReq) (*doc_rpc.AckIndexOutboxRes, error) {
	status := trim(in.Status)
	if status != "done" && status != "failed" {
		return nil, errInvalidArg
	}
	updated := 0
	for _, key := range in.Keys {
		parts := strings.SplitN(trim(key), ":", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			continue
		}
		res := l.svcCtx.DB.Model(&doc_models.DocIndexOutbox{}).
			Where("doc_id = ? AND version = ?", parts[0], parts[1]).
			Updates(map[string]any{
				"status":   status,
				"attempts": gorm.Expr("attempts + 1"),
			})
		if res.Error == nil {
			updated += int(res.RowsAffected)
		}
	}
	return &doc_rpc.AckIndexOutboxRes{Updated: int32(updated)}, nil
}
