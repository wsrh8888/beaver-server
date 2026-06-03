package svc

import (
	"errors"

	models "beaver/app/open/open_models"
	"beaver/app/auth/auth_rpc/auth"
	"beaver/app/open/open_portal/internal/config"
	"beaver/app/open/open_portal/internal/middleware"
	"beaver/app/user/user_rpc/user"
	"beaver/common/zrpc_interceptor"
	"beaver/core/coregorm"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config                  config.Config
	DB                      *gorm.DB
	UserRpc                 user.User
	AuthRpc                 auth.Auth
	DeveloperAuthMiddleware rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := coregorm.InitGorm(c.Mysql.DataSource)
	rpcOpt := zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor)
	return &ServiceContext{
		Config:                  c,
		DB:                      db,
		UserRpc:                 user.NewUser(zrpc.MustNewClient(c.UserRpc, rpcOpt)),
		AuthRpc:                 auth.NewAuth(zrpc.MustNewClient(c.AuthRpc, rpcOpt)),
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
