package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"
	beaverlog "beaver/utils/beaverlog"

	"gorm.io/gorm/clause"
)

type SetFavoriteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewSetFavoriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetFavoriteLogic {
	return &SetFavoriteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("SetFavoriteLogic", ctx),
	}
}

func (l *SetFavoriteLogic) SetFavorite(in *doc_rpc.SetFavoriteReq) (*doc_rpc.SetFavoriteRes, error) {
	doc, _, err := requireDoc(l.svcCtx.DB, in.DocId, in.UserId, roleViewer)
	if err != nil {
		return nil, err
	}
	uid := trim(in.UserId)
	if in.Favorite {
		row := doc_models.DocFavorite{UserID: uid, DocID: doc.DocID}
		if err := l.svcCtx.DB.Clauses(clause.OnConflict{DoNothing: true}).
			Create(&row).Error; err != nil {
			return nil, err
		}
	} else if err := l.svcCtx.DB.Where("user_id = ? AND doc_id = ?", uid, doc.DocID).
		Delete(&doc_models.DocFavorite{}).Error; err != nil {
		return nil, err
	}
	auditJSON(l.svcCtx.DB, uid, "doc.favorite", "doc", doc.DocID,
		map[string]any{"favorite": in.Favorite})
	return &doc_rpc.SetFavoriteRes{Favorite: in.Favorite}, nil
}
