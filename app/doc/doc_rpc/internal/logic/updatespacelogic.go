package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"

	"beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type UpdateSpaceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewUpdateSpaceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSpaceLogic {
	return &UpdateSpaceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("UpdateSpaceLogic", ctx),
	}
}

func (l *UpdateSpaceLogic) UpdateSpace(in *doc_rpc.UpdateSpaceReq) (*doc_rpc.UpdateSpaceRes, error) {
	// 权限校验在 doc_api 层完成（本 RPC 无 operator 字段，为数据权威写入端）。
	space, err := loadSpace(l.svcCtx.DB, in.SpaceId)
	if err != nil {
		return nil, errNotFound
	}

	upd := map[string]any{}
	if name := trim(in.Name); name != "" {
		upd["name"] = name
	}
	if icon := trim(in.Icon); icon != "" {
		upd["icon"] = icon
	}
	if desc := trim(in.Description); desc != "" {
		upd["description"] = desc
	}
	if in.ParentId != space.ParentID {
		upd["parent_id"] = trim(in.ParentId)
	}
	if in.Status > 0 {
		upd["status"] = int(in.Status)
	}
	if int(in.Sort) != space.Sort {
		upd["sort"] = int(in.Sort)
	}
	if len(upd) > 0 {
		if err := l.svcCtx.DB.Model(&doc_models.DocSpace{}).Where("space_id = ?", space.SpaceID).Updates(upd).Error; err != nil {
			l.logger.Error(model.LogMsg{Text: "更新空间失败", Data: err})
			return nil, err
		}
	}

	return &doc_rpc.UpdateSpaceRes{}, nil
}
