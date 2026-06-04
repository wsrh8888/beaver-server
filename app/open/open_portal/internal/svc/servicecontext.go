package svc

import (
	"beaver/app/auth/auth_rpc/auth"
	"beaver/app/group/group_rpc/group"
	"beaver/app/open/open_portal/internal/config"
	"beaver/app/open/open_portal/internal/middleware"
	"beaver/app/open/open_rpc/open"
	"beaver/app/user/user_rpc/user"
	"beaver/common/zrpc_interceptor"
	"beaver/core/coregorm"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config                     config.Config
	DB                         *gorm.DB
	UserRpc                    user.User
	AuthRpc                    auth.Auth
	GroupRpc                   group.Group
	OpenRpc                    open.Open
	DeveloperAuthMiddleware    rest.Middleware
	RequireDeveloperMiddleware rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := coregorm.InitGorm(c.Mysql.DataSource)
	rpcOpt := zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor)
	return &ServiceContext{
		Config:                     c,
		DB:                         db,
		UserRpc:                    user.NewUser(zrpc.MustNewClient(c.UserRpc, rpcOpt)),
		AuthRpc:                    auth.NewAuth(zrpc.MustNewClient(c.AuthRpc, rpcOpt)),
		GroupRpc:                   group.NewGroup(zrpc.MustNewClient(c.GroupRpc, rpcOpt)),
		OpenRpc:                    open.NewOpen(zrpc.MustNewClient(c.OpenRpc, rpcOpt)),
		DeveloperAuthMiddleware:    middleware.NewDeveloperAuthMiddleware(c.Auth.AccessSecret).Handle,
		RequireDeveloperMiddleware: middleware.NewRequireDeveloperMiddleware(db).Handle,
	}
}
