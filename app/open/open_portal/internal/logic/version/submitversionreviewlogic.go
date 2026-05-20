package version

import (
	"context"
	"errors"
	"fmt"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SubmitVersionReviewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 提交版本审核
func NewSubmitVersionReviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubmitVersionReviewLogic {
	return &SubmitVersionReviewLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SubmitVersionReviewLogic) SubmitVersionReview(req *types.SubmitVersionReviewReq) (resp *types.SubmitVersionReviewRes, err error) {
	// 解析 versionId
	var versionID uint
	fmt.Sscanf(req.VersionID, "%d", &versionID)

	// 查询版本
	var version open_models.OpenAppVersion
	if err := l.svcCtx.DB.First(&version, versionID).Error; err != nil {
		return nil, err
	}

	// 验证应用所有权
	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ?", version.AppID).First(&app).Error; err != nil {
		return nil, err
	}

	// 验证 UserID 是否为应用所有者
	if app.OwnerUserID != req.UserID {
		return nil, errors.New("无权操作此应用")
	}

	// 只有草稿状态才能提交审核
	if version.Status != "draft" {
		return nil, fmt.Errorf("当前状态(%s)不允许提交审核", version.Status)
	}

	// 更新状态为审核中
	if err := l.svcCtx.DB.Model(&version).Update("status", "reviewing").Error; err != nil {
		return nil, err
	}

	return &types.SubmitVersionReviewRes{}, nil
}
