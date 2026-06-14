package logic

import (
	"context"
	"errors"

	"beaver/app/user/user_models"
	"beaver/app/user/user_rpc/internal/svc"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UserUpdateDisplayLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserUpdateDisplayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserUpdateDisplayLogic {
	return &UserUpdateDisplayLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserUpdateDisplayLogic) UserUpdateDisplay(in *user_rpc.UserUpdateDisplayReq) (*user_rpc.UserUpdateDisplayRes, error) {
	if in.UserId == "" {
		return nil, errors.New("userId 不能为空")
	}
	if in.NickName == "" && in.Avatar == "" {
		return &user_rpc.UserUpdateDisplayRes{}, nil
	}

	var user user_models.UserModel
	if err := l.svcCtx.DB.Where("user_id = ?", in.UserId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	updates := map[string]interface{}{}
	if in.NickName != "" {
		updates["nick_name"] = in.NickName
	}
	if in.Avatar != "" {
		updates["avatar"] = in.Avatar
	}

	version := l.svcCtx.VersionGen.GetNextVersion("users", "user_id", in.UserId)
	if version == -1 {
		return nil, errors.New("获取用户版本号失败")
	}
	updates["version"] = version

	if err := l.svcCtx.DB.Model(&user).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &user_rpc.UserUpdateDisplayRes{Version: version}, nil
}
