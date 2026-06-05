package moderation

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/platform/platform_rpc/types/platform_rpc"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetContentReportListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetContentReportListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetContentReportListLogic {
	return &GetContentReportListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetContentReportListLogic) GetContentReportList(req *types.GetContentReportListReq) (resp *types.GetContentReportListRes, err error) {
	rpcRes, err := l.svcCtx.PlatformRpc.ListContentReports(l.ctx, &platform_rpc.ListContentReportsReq{
		Page:       int32(req.Page),
		PageSize:   int32(req.PageSize),
		Status:     int32(req.Status),
		TargetType: int32(req.TargetType),
		TargetId:   req.TargetID,
	})
	if err != nil {
		l.Errorf("获取内容举报列表失败: %v", err)
		return nil, err
	}

	users := map[string]*user_rpc.UserInfo{}
	seen := map[string]struct{}{}
	userIDs := make([]string, 0, len(rpcRes.List))
	for _, r := range rpcRes.List {
		if r.ReporterUserId == "" {
			continue
		}
		if _, ok := seen[r.ReporterUserId]; ok {
			continue
		}
		seen[r.ReporterUserId] = struct{}{}
		userIDs = append(userIDs, r.ReporterUserId)
	}
	if len(userIDs) > 0 {
		if res, err := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{UserIdList: userIDs}); err == nil && res != nil {
			users = res.UserInfo
		}
	}

	list := make([]types.ContentReportInfo, 0, len(rpcRes.List))
	for _, r := range rpcRes.List {
		reporterName := ""
		if u, ok := users[r.ReporterUserId]; ok && u != nil {
			reporterName = u.NickName
		}
		list = append(list, types.ContentReportInfo{
			ID:             r.Id,
			ReporterUserID: r.ReporterUserId,
			ReporterName:   reporterName,
			TargetType:     int(r.TargetType),
			TargetID:       r.TargetId,
			ReasonType:     int(r.ReasonType),
			Content:        r.Content,
			Status:         int(r.Status),
			CaseID:         r.CaseId,
			CreatedAt:      r.CreatedAt,
		})
	}
	return &types.GetContentReportListRes{List: list, Total: rpcRes.Total}, nil
}
