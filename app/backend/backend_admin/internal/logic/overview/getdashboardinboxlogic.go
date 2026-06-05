package overview

import (
	"context"
	"fmt"
	"sort"

	"beaver/app/backend/backend_models"
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/open/open_rpc/types/open_rpc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type inboxRow struct {
	item      types.DashboardInboxItem
	sortKey   string
}

type GetDashboardInboxLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDashboardInboxLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDashboardInboxLogic {
	return &GetDashboardInboxLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetDashboardInboxLogic) GetDashboardInbox(req *types.GetDashboardInboxReq) (resp *types.GetDashboardInboxRes, err error) {
	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 50 {
		limit = 50
	}

	rows := make([]inboxRow, 0, limit)
	perSource := int32(8)

	reportRes, err := l.svcCtx.PlatformRpc.ListContentReports(l.ctx, &platform_rpc.ListContentReportsReq{
		Status: 1, Page: 1, PageSize: perSource,
	})
	if err != nil {
		l.Errorf("收件箱拉取举报失败: %v", err)
	} else {
		for _, r := range reportRes.List {
			rows = append(rows, inboxRow{
				sortKey: r.CreatedAt,
				item: types.DashboardInboxItem{
					Category: "report", Title: "待处理举报",
					Summary: fmt.Sprintf("%s 举报 %s", r.ReporterUserId, r.TargetId),
					EntityID: fmt.Sprintf("%d", r.Id), CreatedAt: r.CreatedAt,
					Action: fmt.Sprintf("/safety/reports?reportId=%d", r.Id),
				},
			})
		}
	}

	var cases []backend_models.AdminModerationCase
	if err := l.svcCtx.DB.Where("status = ?", backend_models.CaseStatusPending).
		Order("id DESC").Limit(int(perSource)).Find(&cases).Error; err != nil {
		l.Errorf("收件箱拉取工单失败: %v", err)
	} else {
		for _, c := range cases {
			rows = append(rows, inboxRow{
				sortKey: c.CreatedAt.String(),
				item: types.DashboardInboxItem{
					Category: "case", Title: "待处理工单",
					Summary: fmt.Sprintf("%s · %s", c.CaseNo, c.Title),
					EntityID: fmt.Sprintf("%d", c.Id), CreatedAt: c.CreatedAt.String(),
					Action: fmt.Sprintf("/safety/cases?caseId=%d", c.Id),
				},
			})
		}
	}

	feedbackRes, err := l.svcCtx.PlatformRpc.ListFeedback(l.ctx, &platform_rpc.ListFeedbackReq{
		Status: 1, Page: 1, PageSize: perSource,
	})
	if err != nil {
		l.Errorf("收件箱拉取反馈失败: %v", err)
	} else {
		for _, f := range feedbackRes.List {
			rows = append(rows, inboxRow{
				sortKey: f.CreatedAt,
				item: types.DashboardInboxItem{
					Category: "feedback", Title: "待处理反馈",
					Summary: f.Content, EntityID: fmt.Sprintf("%d", f.Id),
					CreatedAt: f.CreatedAt, Action: "/service/feedback",
				},
			})
		}
	}

	devRes, err := l.svcCtx.OpenRpc.ListDevelopers(l.ctx, &open_rpc.ListDevelopersReq{
		Status: 0, Page: 1, PageSize: perSource,
	})
	if err != nil {
		l.Errorf("收件箱拉取开发者审核失败: %v", err)
	} else {
		for _, d := range devRes.List {
			createdAt := fmt.Sprintf("%d", d.CreatedAt)
			rows = append(rows, inboxRow{
				sortKey: createdAt,
				item: types.DashboardInboxItem{
					Category: "developer", Title: "待审开发者",
					Summary: fmt.Sprintf("%s (%s)", d.RealName, d.Email),
					EntityID: d.UserId, CreatedAt: createdAt, Action: "/open/developers",
				},
			})
		}
	}

	appRes, err := l.svcCtx.OpenRpc.ListOpenApps(l.ctx, &open_rpc.ListOpenAppsReq{
		AuditStatus: 0, Page: 1, PageSize: perSource,
	})
	if err != nil {
		l.Errorf("收件箱拉取应用审核失败: %v", err)
	} else {
		for _, a := range appRes.List {
			createdAt := fmt.Sprintf("%d", a.CreatedAt)
			rows = append(rows, inboxRow{
				sortKey: createdAt,
				item: types.DashboardInboxItem{
					Category: "app", Title: "待审应用",
					Summary: a.Name, EntityID: a.AppId, CreatedAt: createdAt, Action: "/open/apps",
				},
			})
		}
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i].sortKey > rows[j].sortKey
	})

	list := make([]types.DashboardInboxItem, 0, limit)
	for i := 0; i < len(rows) && i < limit; i++ {
		list = append(list, rows[i].item)
	}
	return &types.GetDashboardInboxRes{List: list, Total: int64(len(rows))}, nil
}
