package version

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

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
	// 验证应用是否存在
	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ?", req.AppID).First(&app).Error; err != nil {
		return nil, err
	}

	// 验证 UserID 是否为应用所有者
	if app.OwnerUserID != req.UserID {
		return nil, errors.New("无权查看此应用")
	}

	// 设置默认分页
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	// 查询总数
	var total int64
	if err := l.svcCtx.DB.Model(&open_models.OpenAppVersion{}).Where("app_id = ?", req.AppID).Count(&total).Error; err != nil {
		return nil, err
	}

	// 查询列表
	var versions []open_models.OpenAppVersion
	offset := (page - 1) * pageSize
	if err := l.svcCtx.DB.Where("app_id = ?", req.AppID).Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&versions).Error; err != nil {
		return nil, err
	}

	// 转换为响应格式
	list := make([]types.VersionInfo, 0, len(versions))
	for _, v := range versions {
		// 解析 capabilities
		var capabilities []string
		if v.Capabilities != "" {
			json.Unmarshal([]byte(v.Capabilities), &capabilities)
		}

		list = append(list, types.VersionInfo{
			ID:           fmt.Sprintf("%d", v.ID),
			Version:      v.Version,
			Description:  v.Description,
			Visibility:   v.Visibility,
			Status:       v.Status,
			Capabilities: capabilities,
			CreatedAt:    int64(v.CreatedAt.Unix()),
		})
	}

	return &types.GetVersionListRes{
		Total: total,
		List:  list,
	}, nil
}
