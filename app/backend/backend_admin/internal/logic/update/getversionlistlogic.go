package logic

import (
	"beaver/app/update/update_models"
	"beaver/common/list_query"
	"beaver/common/models"
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetVersionListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取版本列表
func NewGetVersionListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVersionListLogic {
	return &GetVersionListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetVersionListLogic) GetVersionList(req *types.GetVersionListReq) (resp *types.GetVersionListRes, err error) {
	// 构建查询条件
	query := l.svcCtx.DB

	// 添加查询条件
	if req.ArchitectureID > 0 {
		query = query.Where("architecture_id = ?", req.ArchitectureID)
	}

	// 使用 list_query 进行查询
	list, total, err := list_query.ListQuery(l.svcCtx.DB.Preload("Architecture"), update_models.UpdateVersion{}, list_query.Option{
		PageInfo: models.PageInfo{
			Page:  req.Page,
			Limit: req.PageSize,
		},
		Where: query,
	})

	if err != nil {
		logx.Errorf("Failed to get versions: %v", err)
		return nil, err
	}

	// 构建响应
	versionList := make([]types.VersionInfo, 0, len(list))
	for _, ver := range list {
		versionList = append(versionList, types.VersionInfo{
			VersionID:      uint(ver.Id),
			ArchitectureID: ver.ArchitectureID,
			Version:        ver.Version,
			FileName:       ver.FileName,
			Description:    ver.Description,
			ReleaseNotes:   ver.ReleaseNotes,
			CreatedAt:      ver.CreatedAt.String(),
			UpdatedAt:      ver.UpdatedAt.String(),
		})
	}

	return &types.GetVersionListRes{
		Total:    total,
		Versions: versionList,
	}, nil
}
