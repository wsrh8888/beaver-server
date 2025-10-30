package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/moment/moment_models"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type HandleMomentReportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 处理动态举报
func NewHandleMomentReportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HandleMomentReportLogic {
	return &HandleMomentReportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HandleMomentReportLogic) HandleMomentReport(req *types.HandleMomentReportReq) (resp *types.HandleMomentReportRes, err error) {
	// 检查举报是否存在
	var report moment_models.MomentReportModel
	err = l.svcCtx.DB.Where("id = ?", req.Id).First(&report).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logx.Errorf("举报不存在: %d", req.Id)
			return nil, errors.New("举报不存在")
		}
		logx.Errorf("查询举报失败: %v", err)
		return nil, errors.New("查询举报失败")
	}

	// 更新举报状态
	err = l.svcCtx.DB.Model(&report).Update("status", req.Status).Error
	if err != nil {
		logx.Errorf("更新举报状态失败: %v", err)
		return nil, errors.New("更新举报状态失败")
	}

	// 如果是处理状态（已处理），可能需要对被举报的动态进行相应处理
	if req.Status == 1 {
		// 可以根据举报类型和严重程度决定是否删除动态或给用户发警告
		// 这里示例是逻辑删除被举报的动态
		err = l.svcCtx.DB.Model(&moment_models.MomentModel{}).
			Where("id = ?", report.MomentID).
			Update("is_deleted", true).Error
		if err != nil {
			logx.Errorf("处理被举报动态失败: %v", err)
			// 这里不返回错误，因为举报状态已经更新成功
		}
	}

	return &types.HandleMomentReportRes{}, nil
}
