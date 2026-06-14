package overview

import (
	"context"
	"strings"
	"time"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/backend/backend_models"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDashboardTrendsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDashboardTrendsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDashboardTrendsLogic {
	return &GetDashboardTrendsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func buildTrendDays(days int) []string {
	if days <= 0 {
		days = 7
	}
	if days > 30 {
		days = 30
	}
	now := time.Now()
	list := make([]string, days)
	for i := days - 1; i >= 0; i-- {
		list[days-1-i] = now.AddDate(0, 0, -i).Format("2006-01-02")
	}
	return list
}

func emptyCounts(days []string) map[string]int64 {
	m := make(map[string]int64, len(days))
	for _, d := range days {
		m[d] = 0
	}
	return m
}

func toSeries(key, label string, days []string, counts map[string]int64) types.DashboardTrendSeries {
	values := make([]int64, len(days))
	for i, d := range days {
		values[i] = counts[d]
	}
	return types.DashboardTrendSeries{Key: key, Label: label, Values: values}
}

func parseDateKey(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	if len(raw) >= 10 {
		return raw[:10]
	}
	if t, err := time.Parse("2006-01-02 15:04:05", raw); err == nil {
		return t.Format("2006-01-02")
	}
	if t, err := time.Parse(time.RFC3339, raw); err == nil {
		return t.Format("2006-01-02")
	}
	return ""
}

func (l *GetDashboardTrendsLogic) countCasesByDay(days []string, since time.Time) map[string]int64 {
	counts := emptyCounts(days)
	type row struct {
		Day   string
		Count int64
	}
	var rows []row
	if err := l.svcCtx.DB.Model(&backend_models.AdminModerationCase{}).
		Select("DATE(created_at) as day, COUNT(*) as count").
		Where("created_at >= ?", since).
		Group("DATE(created_at)").
		Scan(&rows).Error; err != nil {
		l.Errorf("统计工单趋势失败: %v", err)
		return counts
	}
	for _, r := range rows {
		if _, ok := counts[r.Day]; ok {
			counts[r.Day] = r.Count
		}
	}
	return counts
}

func (l *GetDashboardTrendsLogic) countAuditOpsByDay(days []string, since time.Time) map[string]int64 {
	counts := emptyCounts(days)
	type row struct {
		Day   string
		Count int64
	}
	var rows []row
	if err := l.svcCtx.DB.Model(&backend_models.AdminOperationLog{}).
		Select("DATE(created_at) as day, COUNT(*) as count").
		Where("created_at >= ?", since).
		Group("DATE(created_at)").
		Scan(&rows).Error; err != nil {
		l.Errorf("统计审计操作趋势失败: %v", err)
		return counts
	}
	for _, r := range rows {
		if _, ok := counts[r.Day]; ok {
			counts[r.Day] = r.Count
		}
	}
	return counts
}

func (l *GetDashboardTrendsLogic) countReportsByDay(days []string, since time.Time) map[string]int64 {
	counts := emptyCounts(days)
	res, err := l.svcCtx.PlatformRpc.ListContentReports(l.ctx, &platform_rpc.ListContentReportsReq{
		Page: 1, PageSize: 500,
	})
	if err != nil {
		l.Errorf("统计举报趋势失败: %v", err)
		return counts
	}
	sinceDate := since.Format("2006-01-02")
	for _, r := range res.List {
		day := parseDateKey(r.CreatedAt)
		if day == "" || day < sinceDate {
			continue
		}
		if _, ok := counts[day]; ok {
			counts[day]++
		}
	}
	return counts
}

func (l *GetDashboardTrendsLogic) countFeedbackByDay(days []string, since time.Time) map[string]int64 {
	counts := emptyCounts(days)
	res, err := l.svcCtx.PlatformRpc.ListFeedback(l.ctx, &platform_rpc.ListFeedbackReq{
		Page: 1, PageSize: 500,
	})
	if err != nil {
		l.Errorf("统计反馈趋势失败: %v", err)
		return counts
	}
	sinceDate := since.Format("2006-01-02")
	for _, f := range res.List {
		day := parseDateKey(f.CreatedAt)
		if day == "" || day < sinceDate {
			continue
		}
		if _, ok := counts[day]; ok {
			counts[day]++
		}
	}
	return counts
}

func (l *GetDashboardTrendsLogic) GetDashboardTrends(req *types.GetDashboardTrendsReq) (resp *types.GetDashboardTrendsRes, err error) {
	days := req.Days
	if days <= 0 {
		days = 7
	}
	dayLabels := buildTrendDays(days)
	since := time.Now().AddDate(0, 0, -(days - 1))
	since = time.Date(since.Year(), since.Month(), since.Day(), 0, 0, 0, 0, since.Location())

	resp = &types.GetDashboardTrendsRes{
		Days: dayLabels,
		Series: []types.DashboardTrendSeries{
			toSeries("cases", "新建工单", dayLabels, l.countCasesByDay(dayLabels, since)),
			toSeries("reports", "新增举报", dayLabels, l.countReportsByDay(dayLabels, since)),
			toSeries("feedback", "用户反馈", dayLabels, l.countFeedbackByDay(dayLabels, since)),
			toSeries("auditOps", "安全操作", dayLabels, l.countAuditOpsByDay(dayLabels, since)),
		},
	}
	return resp, nil
}
