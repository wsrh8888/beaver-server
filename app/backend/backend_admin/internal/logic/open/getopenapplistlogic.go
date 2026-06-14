package open

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/open/open_rpc/types/open_rpc"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOpenAppListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOpenAppListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOpenAppListLogic {
	return &GetOpenAppListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetOpenAppListLogic) GetOpenAppList(req *types.GetOpenAppListReq) (resp *types.GetOpenAppListRes, err error) {
	rpcRes, err := l.svcCtx.OpenRpc.ListOpenApps(l.ctx, &open_rpc.ListOpenAppsReq{
		Page:            int32(req.Page),
		PageSize:        int32(req.PageSize),
		Keyword:         req.Keyword,
		OwnerUserId:     req.OwnerUserID,
		AppId:           req.AppID,
		Status:          int32(req.Status),
		AuditStatus:     int32(req.AuditStatus),
		CapabilityType:  int32(req.CapabilityType),
	})
	if err != nil {
		l.Errorf("获取应用列表失败: %v", err)
		return nil, err
	}

	users := map[string]*user_rpc.UserInfo{}
	seen := map[string]struct{}{}
	userIDs := make([]string, 0, len(rpcRes.List))
	for _, app := range rpcRes.List {
		if app.OwnerUserId == "" {
			continue
		}
		if _, ok := seen[app.OwnerUserId]; ok {
			continue
		}
		seen[app.OwnerUserId] = struct{}{}
		userIDs = append(userIDs, app.OwnerUserId)
	}
	if len(userIDs) > 0 {
		if res, err := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{UserIdList: userIDs}); err == nil && res != nil {
			users = res.UserInfo
		}
	}

	list := make([]types.OpenAppInfo, 0, len(rpcRes.List))
	for _, app := range rpcRes.List {
		ownerName := ""
		if u, ok := users[app.OwnerUserId]; ok && u != nil {
			ownerName = u.NickName
		}
		list = append(list, types.OpenAppInfo{
			AppID:         app.AppId,
			Name:          app.Name,
			Description:   app.Description,
			Icon:          app.Icon,
			OwnerUserID:   app.OwnerUserId,
			OwnerUserName: ownerName,
			AppType:       int(app.AppType),
			Category:      app.Category,
			Status:        int(app.Status),
			AuditStatus:   int(app.AuditStatus),
			AuditedBy:     app.AuditedBy,
			AuditedAt:     app.AuditedAt,
			EnableRobot:   int(app.EnableRobot),
			EnableOAuth:   int(app.EnableOauth),
			EnableWebhook: int(app.EnableWebhook),
			CreatedAt:     app.CreatedAt,
		})
	}

	return &types.GetOpenAppListRes{Total: rpcRes.Total, List: list}, nil
}
