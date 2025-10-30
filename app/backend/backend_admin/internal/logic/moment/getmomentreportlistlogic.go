package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/moment/moment_models"
	"beaver/common/list_query"
	"beaver/common/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMomentReportListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取动态举报列表
func NewGetMomentReportListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMomentReportListLogic {
	return &GetMomentReportListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMomentReportListLogic) GetMomentReportList(req *types.GetMomentReportListReq) (resp *types.GetMomentReportListRes, err error) {
	// 构建查询条件
	whereClause := l.svcCtx.DB.Where("1 = 1")

	// 处理状态筛选
	if req.Status != 0 {
		whereClause = whereClause.Where("status = ?", req.Status)
	}

	// 动态ID筛选
	if req.MomentId != 0 {
		whereClause = whereClause.Where("moment_id = ?", req.MomentId)
	}

	// 分页查询举报
	reports, count, err := list_query.ListQuery(l.svcCtx.DB, moment_models.MomentReportModel{}, list_query.Option{
		PageInfo: models.PageInfo{
			Page:  req.Page,
			Limit: req.Limit,
			Sort:  "created_at desc",
		},
		Where: whereClause,
	})

	if err != nil {
		logx.Errorf("查询动态举报列表失败: %v", err)
		return nil, err
	}

	// 转换为响应格式
	var list []types.MomentReportInfo
	for _, report := range reports {
		var images []types.FileInfo
		if report.Images != nil {
			for _, img := range *report.Images {
				images = append(images, types.FileInfo{
					FileName: img.FileName,
				})
			}
		}

		list = append(list, types.MomentReportInfo{
			Id:        report.Id,
			UserId:    report.UserID,
			MomentId:  report.MomentID,
			Reason:    report.Reason,
			Images:    images,
			Status:    report.Status,
			CreatedAt: report.CreatedAt.String(),
			UpdatedAt: report.UpdatedAt.String(),
		})
	}

	return &types.GetMomentReportListRes{
		List:  list,
		Total: count,
	}, nil
}
