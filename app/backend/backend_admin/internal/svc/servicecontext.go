package svc

import (
	"beaver/app/backend/backend_models"
	"beaver/app/backend/backend_admin/internal/config"
	"beaver/app/auth/auth_rpc/auth"
	chatcli "beaver/app/chat/chat_rpc/chat"
	emojicli "beaver/app/emoji/emoji_rpc/emoji"
	"beaver/app/file/file_rpc/file"
	"beaver/app/file/file_rpc/types/file_rpc"
	friendcli "beaver/app/friend/friend_rpc/friend"
	groupcli "beaver/app/group/group_rpc/group"
	momentcli "beaver/app/moment/moment_rpc/moment"
	opencli "beaver/app/open/open_rpc/open"
	platformcli "beaver/app/platform/platform_rpc/platform"
	"beaver/app/user/user_rpc/types/user_rpc"
	"beaver/app/user/user_rpc/user"
	"beaver/common/zrpc_interceptor"
	"beaver/core/coregorm"
	"beaver/core/coreredis"
	versionPkg "beaver/core/version"

	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config      config.Config
	DB          *gorm.DB
	Redis       *redis.Client
	VersionGen  *versionPkg.VersionGenerator
	UserRpc     user_rpc.UserClient
	AuthRpc     auth.Auth
	FileRpc     file_rpc.FileClient
	PlatformRpc platformcli.Platform
	OpenRpc     opencli.Open
	FriendRpc   friendcli.Friend
	GroupRpc    groupcli.Group
	ChatRpc     chatcli.Chat
	MomentRpc   momentcli.Moment
	EmojiRpc    emojicli.Emoji
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := coregorm.InitGorm(c.Mysql.DataSource)
	client := coreredis.InitRedis(c.Redis.Addr, c.Redis.Password, c.Redis.Db)
	rpcOpt := zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor)

	return &ServiceContext{
		Config:      c,
		DB:          mysqlDb,
		Redis:       client,
		VersionGen:  versionPkg.NewVersionGenerator(client, mysqlDb),
		UserRpc:     user.NewUser(zrpc.MustNewClient(c.UserRpc, rpcOpt)),
		AuthRpc:     auth.NewAuth(zrpc.MustNewClient(c.AuthRpc, rpcOpt)),
		FileRpc:     file.NewFile(zrpc.MustNewClient(c.FileRpc, rpcOpt)),
		PlatformRpc: platformcli.NewPlatform(zrpc.MustNewClient(c.PlatformRpc, rpcOpt)),
		OpenRpc:     opencli.NewOpen(zrpc.MustNewClient(c.OpenRpc, rpcOpt)),
		FriendRpc:   friendcli.NewFriend(zrpc.MustNewClient(c.FriendRpc, rpcOpt)),
		GroupRpc:    groupcli.NewGroup(zrpc.MustNewClient(c.GroupRpc, rpcOpt)),
		ChatRpc:     chatcli.NewChat(zrpc.MustNewClient(c.ChatRpc, rpcOpt)),
		MomentRpc:   momentcli.NewMoment(zrpc.MustNewClient(c.MomentRpc, rpcOpt)),
		EmojiRpc:    emojicli.NewEmoji(zrpc.MustNewClient(c.EmojiRpc, rpcOpt)),
	}
}

// RecordOperation 写入管理员操作审计日志（后台域本地表）
func (s *ServiceContext) RecordOperation(operatorID, action, targetType, targetID string, caseID uint64, detail, result, errMsg string) {
	_ = s.DB.Create(&backend_models.AdminOperationLog{
		OperatorID:   operatorID,
		Action:       action,
		TargetType:   targetType,
		TargetID:     targetID,
		CaseID:       caseID,
		Detail:       detail,
		Result:       result,
		ErrorMessage: errMsg,
	}).Error
}
