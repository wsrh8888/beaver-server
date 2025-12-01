package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/backend/backend_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateAuthorityLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建权限
func NewCreateAuthorityLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAuthorityLogic {
	return &CreateAuthorityLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateAuthorityLogic) CreateAuthority(req *types.CreateAuthorityReq) (resp *types.CreateAuthorityRes, err error) {
	// 创建权限数据
	authority := backend_models.AdminSystemAuthority{
		Name:        req.Name,
		Description: req.Description,
		Status:      1, // 默认启用
		Sort:        0, // 默认排序
	}

	// 创建权限
	err = l.svcCtx.DB.Create(&authority).Error
	if err != nil {
		logx.Errorf("创建权限失败: %v", err)
		return nil, err
	}

	logx.Infof("权限创建成功: ID=%d, Name=%s", authority.Id, authority.Name)
	return &types.CreateAuthorityRes{}, nil
}
