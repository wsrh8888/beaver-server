package logic

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"beaver/app/user/user_models"
	"beaver/app/user/user_rpc/internal/svc"
	"beaver/app/user/user_rpc/types/user_rpc"
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
	// todo: add your logic here and delete this line
	var user user_models.UserModel

	err := l.svcCtx.DB.Take(&user, "phone = ?", in.Phone).Error

	if err == nil {
		return nil, errors.New("用户已存在")
	}

	user = user_models.UserModel{
		UUID:     strings.Replace(uuid.New().String(), "-", "", -1),
		Password: in.Password,
		Phone:    in.Phone,
		Source:   in.Source,
		NickName: utils.GenerateRandomString(8),
		Abstract: "",
	}
	// 打印当前时间
	fmt.Println("Current time:", time.Now())
	err = l.svcCtx.DB.Create(&user).Error
	if err != nil {
		logx.Error(err)
		return nil, errors.New("创建用户失败")
	}

	return &user_rpc.UserCreateRes{
		UserID: user.UUID,
	}, nil
}
