package contact

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	user_models "beaver/app/user/user_models"

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
	// 1. 生成用户 ID
	userIDBytes := make([]byte, 16)
	_, _ = rand.Read(userIDBytes)
	userID := "u_" + hex.EncodeToString(userIDBytes)

	// 2. 创建用户
	user := user_models.UserModel{
		UserID:   userID,
		NickName: req.Nickname,
		Phone:    req.Phone,
		Email:    req.Email,
		Avatar:   req.Avatar,
		Status:   int8(req.Status),
	}
	if user.Status == 0 {
		user.Status = 1 // 默认启用
	}

	if err := l.svcCtx.DB.Create(&user).Error; err != nil {
		logx.Errorf("创建用户失败: %v", err)
		return nil, errors.New("创建用户失败")
	}

	return &types.CreateUserRes{
		UserID: userID,
	}, nil
}
