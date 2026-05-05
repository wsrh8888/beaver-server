package app

import (
	"context"
	"errors"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAppListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取应用列表
func NewGetAppListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAppListLogic {
	return &GetAppListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAppListLogic) GetAppList(req *types.GetAppListReq) (resp *types.GetAppListRes, err error) {
	// 1. 从 header 获取当前用户 ID（由中间件注入）
	userID := l.ctx.Value("userId")
	if userID == nil {
		return nil, errors.New("未登录")
	}

	// 2. 构建查询条件
	query := l.svcCtx.DB.Model(&open_models.OpenApp{}).Where("owner_user_id = ?", userID)

	// 3. 如果指定了状态，添加状态过滤
	if req.Status != 0 {
		query = query.Where("status = ?", req.Status)
	}

	// 4. 查询总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, errors.New("查询失败")
	}

	// 5. 分页查询
	var apps []open_models.OpenApp
	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&apps).Error; err != nil {
		return nil, errors.New("查询失败")
	}

	// 6. 转换为响应格式
	list := make([]types.AppInfo, 0, len(apps))
	for _, app := range apps {
		list = append(list, types.AppInfo{
			AppID:       app.AppID,
			Name:        app.Name,
			Description: app.Description,
			Status:      app.Status,
			WebhookURL:  app.WebhookURL,
			CreatedAt:   app.CreatedAt.Unix(),
		})
	}

	return &types.GetAppListRes{
		Total: total,
		List:  list,
	}, nil
}
