package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"

	"beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type CreateSpaceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewCreateSpaceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSpaceLogic {
	return &CreateSpaceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("CreateSpaceLogic", ctx),
	}
}

func (l *CreateSpaceLogic) CreateSpace(in *doc_rpc.CreateSpaceReq) (*doc_rpc.CreateSpaceRes, error) {
	name := trim(in.Name)
	if name == "" {
		return nil, errInvalidArg
	}
	kind := trim(in.Kind)
	if kind == "" {
		kind = "owner"
	}

	spaceID := genID("space")
	row := doc_models.DocSpace{
		SpaceID:     spaceID,
		Name:        name,
		Kind:        kind,
		OwnerUserID: trim(in.OwnerUserId),
		ParentID:    trim(in.ParentId),
		Icon:        trim(in.Icon),
		Description: trim(in.Description),
		Status:      1,
		Sort:        int(in.Sort),
	}
	if err := l.svcCtx.DB.Create(&row).Error; err != nil {
		l.logger.Error(model.LogMsg{Text: "创建空间失败", Data: err})
		return nil, err
	}

	// 归属人写入 owner 授权行，供权限继承与列表可见性使用。
	if owner := trim(in.OwnerUserId); owner != "" {
		perm := doc_models.DocPermission{
			ResourceType:  "space",
			ResourceID:    spaceID,
			PrincipalType: "user",
			PrincipalID:   owner,
			Role:          roleOwner,
		}
		if err := l.svcCtx.DB.Create(&perm).Error; err != nil {
			l.logger.Error(model.LogMsg{Text: "写入空间owner授权失败", Data: err})
		}
	}

	return &doc_rpc.CreateSpaceRes{SpaceId: spaceID}, nil
}
