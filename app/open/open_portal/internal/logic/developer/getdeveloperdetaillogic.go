package developer

import (
	"context"
	"errors"
	"time"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetDeveloperDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDeveloperDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDeveloperDetailLogic {
	return &GetDeveloperDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDeveloperDetailLogic) GetDeveloperDetail(req *types.GetDeveloperDetailReq) (resp *types.GetDeveloperDetailRes, err error) {
	userID, ok := l.ctx.Value("userId").(string)
	if !ok || userID == "" {
		return nil, errors.New("未登录")
	}
	if req.ID == 0 {
		return nil, errors.New("id 不能为空")
	}

	var dev open_models.OpenDeveloper
	if err := l.svcCtx.DB.Where("id = ?", req.ID).First(&dev).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("开发者记录不存在")
		}
		return nil, err
	}

	return &types.GetDeveloperDetailRes{
		Developer: types.DeveloperInfo{
			ID:          dev.Id,
			UserID:      dev.UserID,
			RealName:    dev.RealName,
			CompanyName: dev.CompanyName,
			Phone:       dev.Phone,
			Email:       dev.Email,
			Description: dev.Description,
			Status:      dev.Status,
			AuditBy:     dev.AuditBy,
			AuditTime:   dev.AuditTime,
			AuditRemark: dev.AuditRemark,
			CreatedAt:   time.Time(dev.CreatedAt).Unix(),
		},
	}, nil
}