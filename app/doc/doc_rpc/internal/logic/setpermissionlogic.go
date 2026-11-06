package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"

	"beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"

	"gorm.io/gorm"
)

type SetPermissionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewSetPermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetPermissionLogic {
	return &SetPermissionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("SetPermissionLogic", ctx),
	}
}

func (l *SetPermissionLogic) SetPermission(in *doc_rpc.SetPermissionReq) (*doc_rpc.SetPermissionRes, error) {
	rt := trim(in.ResourceType)
	if rt != "space" && rt != "doc" {
		return nil, errInvalidArg
	}
	role := trim(in.Role)
	// owner 不通过此接口授予，避免越权提权。
	if role != roleManager && role != roleEditor && role != roleCommenter && role != roleViewer {
		return nil, errInvalidArg
	}

	// 操作者需对该资源有管理权。
	opRole, err := effectiveRoleForResource(l.svcCtx.DB, rt, in.ResourceId, in.OperatorUserId)
	if err != nil {
		return nil, err
	}
	if !canManage(opRole) {
		return nil, errNoPermission
	}

	pt := trim(in.PrincipalType)
	if pt == "" {
		pt = "user"
	}
	pi := trim(in.PrincipalId)
	if pi == "" {
		return nil, errInvalidArg
	}

	row := doc_models.DocPermission{
		ResourceType:  rt,
		ResourceID:    in.ResourceId,
		PrincipalType: pt,
		PrincipalID:   pi,
		Role:          role,
	}
	// upsert：存在则改角色，否则新建。
	var existing doc_models.DocPermission
	err = l.svcCtx.DB.Where("resource_type = ? AND resource_id = ? AND principal_type = ? AND principal_id = ?",
		rt, in.ResourceId, pt, pi).First(&existing).Error
	if err == gorm.ErrRecordNotFound {
		if cerr := l.svcCtx.DB.Create(&row).Error; cerr != nil {
			l.logger.Error(model.LogMsg{Text: "设置权限失败", Data: cerr})
			return nil, cerr
		}
	} else if err != nil {
		l.logger.Error(model.LogMsg{Text: "查询权限失败", Data: err})
		return nil, err
	} else {
		if existing.Role == roleOwner {
			return nil, errNoPermission // 不覆盖 owner 行
		}
		if uerr := l.svcCtx.DB.Model(&doc_models.DocPermission{}).Where("id = ?", existing.Id).Update("role", role).Error; uerr != nil {
			l.logger.Error(model.LogMsg{Text: "更新权限失败", Data: uerr})
			return nil, uerr
		}
	}

	return &doc_rpc.SetPermissionRes{}, nil
}

// effectiveRoleForResource 计算主体对某资源（space|doc）的有效角色，供权限接口复用。
func effectiveRoleForResource(db *gorm.DB, resourceType, resourceID, userID string) (string, error) {
	if resourceType == "space" {
		space, err := loadSpace(db, resourceID)
		if err != nil {
			return "", errNotFound
		}
		return resolveSpaceRole(db, space, userID)
	}
	doc, err := loadDoc(db, resourceID)
	if err != nil {
		return "", errNotFound
	}
	return resolveDocRole(db, doc.SpaceID, doc.DocID, userID)
}
