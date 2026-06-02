package svc

import (
	"errors"

	models "beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/config"
	"beaver/app/open/open_portal/internal/middleware"
	"beaver/core/coregorm"

	"github.com/zeromicro/go-zero/rest"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config                  config.Config
	DB                      *gorm.DB
	DeveloperAuthMiddleware rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := coregorm.InitGorm(c.Mysql.DataSource)
	return &ServiceContext{
		Config:                  c,
		DB:                      db,
		DeveloperAuthMiddleware: middleware.NewDeveloperAuthMiddleware(c.Auth.AccessSecret).Handle,
	}
}

func (s *ServiceContext) RequireDeveloper(userID string) (*models.OpenDeveloper, error) {
	if userID == "" {
		return nil, errors.New("未登录")
	}

	var developer models.OpenDeveloper
	err := s.DB.Where("user_id = ? AND status = ?", userID, 1).First(&developer).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("您还不是认证开发者,请先申请开发者资质")
		}
		return nil, errors.New("服务内部异常")
	}

	return &developer, nil
}
