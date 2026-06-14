package logic

import (
	"context"
	"errors"

	"beaver/app/auth/auth_models"
	"beaver/app/auth/auth_rpc/internal/svc"
	"beaver/app/auth/auth_rpc/types/auth_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLastLoginAtLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetLastLoginAtLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLastLoginAtLogic {
	return &GetLastLoginAtLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetLastLoginAtLogic) GetLastLoginAt(in *auth_rpc.GetLastLoginAtReq) (*auth_rpc.GetLastLoginAtRes, error) {
	// 验证必填字段
	if in.UserId == "" {
		return nil, errors.New("用户ID不能为空")
	}

	// 查询用户凭证
	var credential auth_models.AuthCredentialModel
	err := l.svcCtx.DB.Take(&credential, "user_id = ?", in.UserId).Error
	if err != nil {
		logx.Errorf("查询用户凭证失败: %v", err)
		return &auth_rpc.GetLastLoginAtRes{
			LastLoginAt: 0,
		}, nil
	}

	// 返回最后登录时间
	var lastLoginAt int64
	if credential.LastLoginAt != nil {
		lastLoginAt = credential.LastLoginAt.Unix()
	}

	return &auth_rpc.GetLastLoginAtRes{
		LastLoginAt: lastLoginAt,
	}, nil
}
