package developer

import (
	"context"
	"errors"

	models "beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApplyDeveloperLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApplyDeveloperLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyDeveloperLogic {
	return &ApplyDeveloperLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApplyDeveloperLogic) ApplyDeveloper(req *types.ApplyDeveloperReq) (resp *types.ApplyDeveloperRes, err error) {
	// 1. 从 context 获取用户ID (由中间件注入)
	userID, ok := l.ctx.Value("userId").(string)
	if !ok || userID == "" {
		return nil, errors.New("未登录或登录已过期")
	}

	// 2. 检查是否已经申请过
	var existingDeveloper models.OpenDeveloper
	err = l.svcCtx.DB.Where("user_id = ?", userID).First(&existingDeveloper).Error
	if err == nil {
		// 已存在记录
		if existingDeveloper.Status == 1 {
			return nil, errors.New("您已经是认证开发者,无需重复申请")
		} else if existingDeveloper.Status == 0 {
			return nil, errors.New("您的申请正在审核中,请耐心等待")
		} else if existingDeveloper.Status == 2 {
			// 之前被拒绝,可以重新申请,更新信息
			existingDeveloper.RealName = req.RealName
			existingDeveloper.CompanyName = req.CompanyName
			existingDeveloper.Phone = req.Phone
			existingDeveloper.Email = req.Email
			existingDeveloper.Description = req.Description
			existingDeveloper.Status = 0 // 重置为待审核
			existingDeveloper.AuditBy = ""
			existingDeveloper.AuditTime = 0
			existingDeveloper.AuditRemark = ""

			if err := l.svcCtx.DB.Save(&existingDeveloper).Error; err != nil {
				logx.Errorf("更新开发者申请失败: %v", err)
				return nil, errors.New("申请失败,请稍后重试")
			}

			logx.Infof("重新申请开发者: user_id=%s", userID)
			return &types.ApplyDeveloperRes{}, nil
		}
	}

	// 3. 创建新的开发者申请
	newDeveloper := models.OpenDeveloper{
		UserID:      userID,
		RealName:    req.RealName,
		CompanyName: req.CompanyName,
		Phone:       req.Phone,
		Email:       req.Email,
		Description: req.Description,
		Status:      0, // 待审核
	}

	if err := l.svcCtx.DB.Create(&newDeveloper).Error; err != nil {
		logx.Errorf("创建开发者申请失败: %v", err)
		return nil, errors.New("申请失败,请稍后重试")
	}

	logx.Infof("开发者申请成功: user_id=%s, real_name=%s", userID, req.RealName)

	return &types.ApplyDeveloperRes{}, nil
}
