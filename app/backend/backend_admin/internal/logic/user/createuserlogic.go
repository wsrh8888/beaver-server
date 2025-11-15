package logic

import (
	"context"
	"errors"
	"strings"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/user/user_models"
	"beaver/utils/pwd"
	utils "beaver/utils/rand"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建用户
func NewCreateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserLogic {
	return &CreateUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateUserLogic) CreateUser(req *types.CreateUserReq) (resp *types.CreateUserRes, err error) {
	// 检查邮箱是否已存在
	var existUser user_models.UserModel
	err = l.svcCtx.DB.Where("email = ?", req.Email).First(&existUser).Error
	if err == nil {
		l.Logger.Errorf("邮箱已存在: %s", req.Email)
		return nil, errors.New("邮箱已存在")
	}

	// 生成用户UUID
	userUUID := strings.Replace(utils.GenerateUUId(), "-", "", -1)

	// 创建用户，设置默认值
	user := user_models.UserModel{
		UUID:     userUUID,
		NickName: req.Nickname,
		Password: pwd.HahPwd(req.Password),
		Email:    req.Email,
		Abstract: req.Abstract,
		Status:   int8(1),  // 默认状态：正常
		Source:   int32(1), // 默认来源：后台创建
	}

	err = l.svcCtx.DB.Create(&user).Error
	if err != nil {
		l.Logger.Errorf("创建用户失败: %v", err)
		return nil, errors.New("创建用户失败")
	}

	l.Logger.Infof("创建用户成功: userID=%s, email=%s", user.UUID, req.Email)
	return &types.CreateUserRes{
		Id: user.UUID,
	}, nil
}
