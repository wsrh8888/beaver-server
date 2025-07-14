package logic

import (
	"context"
	"errors"
	"strings"

	"beaver/app/user/user_models"
	"beaver/app/user/user_rpc/internal/svc"
	"beaver/app/user/user_rpc/types/user_rpc"
	"beaver/utils/pwd"
	utils "beaver/utils/rand"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type UserCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCreateLogic {
	return &UserCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserCreateLogic) UserCreate(in *user_rpc.UserCreateReq) (*user_rpc.UserCreateRes, error) {
	// 验证必填字段
	if in.Password == "" {
		return nil, errors.New("密码不能为空")
	}
	if in.Phone == "" && in.Email == "" {
		return nil, errors.New("手机号或邮箱至少需要提供一个")
	}

	// 检查用户是否已存在
	var user user_models.UserModel
	var err error

	if in.Phone != "" && in.Email != "" {
		// 手机号和邮箱都提供，检查是否已存在
		err = l.svcCtx.DB.Take(&user, "phone = ? OR email = ?", in.Phone, in.Email).Error
	} else if in.Phone != "" {
		// 只提供手机号
		err = l.svcCtx.DB.Take(&user, "phone = ?", in.Phone).Error
	} else {
		// 只提供邮箱
		err = l.svcCtx.DB.Take(&user, "email = ?", in.Email).Error
	}

	if err == nil {
		return nil, errors.New("用户已存在")
	}

	// 加密密码
	hashedPassword := pwd.HahPwd(in.Password)

	// 生成随机昵称
	nickname := utils.GenerateRandomString(8)
	if in.NickName != "" {
		nickname = in.NickName
	}

	user = user_models.UserModel{
		UUID:     strings.Replace(uuid.New().String(), "-", "", -1),
		Password: hashedPassword,
		Email:    in.Email,
		Phone:    in.Phone,
		Source:   in.Source,
		NickName: nickname,
		Abstract: "",
	}

	err = l.svcCtx.DB.Create(&user).Error
	if err != nil {
		logx.Errorf("创建用户失败: %v", err)
		return nil, errors.New("创建用户失败")
	}

	logx.Infof("用户创建成功: %s", user.UUID)

	return &user_rpc.UserCreateRes{
		UserID: user.UUID,
	}, nil
}
