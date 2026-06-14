package bot

import (
	"context"
	"errors"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListIncomingWebhooksLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListIncomingWebhooksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListIncomingWebhooksLogic {
	return &ListIncomingWebhooksLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListIncomingWebhooksLogic) ListIncomingWebhooks(req *types.ListIncomingWebhooksReq) (resp *types.ListIncomingWebhooksRes, err error) {
	if req.AppID == "" {
		return nil, errors.New("appId 不能为空")
	}

	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ? AND owner_user_id = ?", req.AppID, req.UserID).First(&app).Error; err != nil {
		return nil, errors.New("应用不存在或无权限操作")
	}

	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	query := l.svcCtx.DB.Model(&open_models.OpenBotModel{}).Where("app_id = ?", req.AppID)
	if req.GroupID != "" {
		query = query.Where("group_id = ?", req.GroupID)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, errors.New("查询失败")
	}

	var bots []open_models.OpenBotModel
	if err := query.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&bots).Error; err != nil {
		return nil, errors.New("查询失败")
	}

	list := make([]types.IncomingWebhookInfo, 0, len(bots))
	for i := range bots {
		list = append(list, toIncomingWebhookInfo(&bots[i], l.svcCtx.Config.Domain, false))
	}

	return &types.ListIncomingWebhooksRes{
		Total: total,
		List:  list,
	}, nil
}