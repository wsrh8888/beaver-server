package logic

import (
	"context"
	"errors"

	"beaver/app/auth/auth_rpc/types/auth_rpc"
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserLogic {
	return &CreateUserLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

// CreateUser 管理后台：创建用户。
// admin 职责：校验运营录入项，编排 UserRpc 创建用户 + AuthRpc 创建凭证。
// RPC 职责：UserCreate 落库用户资料，AuthRpc 管理密码凭证。
func (l *CreateUserLogic) CreateUser(req *types.CreateUserReq) (resp *types.CreateUserRes, err error) {
	if req.Email == "" {
		return nil, errors.New("邮箱不能为空")
	}
	if req.Password == "" {
		return nil, errors.New("密码不能为空")
	}

	createRes, err := l.svcCtx.UserRpc.UserCreate(l.ctx, &user_rpc.UserCreateReq{
		Email:    req.Email,
		NickName: req.NickName,
		Abstract: req.Abstract,
		Source:   2,
	})
	if err != nil {
		l.Errorf("创建用户失败: %v", err)
		return nil, err
	}

	credRes, err := l.svcCtx.AuthRpc.CreateCredential(l.ctx, &auth_rpc.CreateCredentialReq{
		UserId:   createRes.UserID,
		Password: req.Password,
	})
	if err != nil || !credRes.Success {
		l.Errorf("创建用户凭证失败: %v", err)
		return nil, errors.New("创建用户凭证失败")
	}

	return &types.CreateUserRes{Id: createRes.UserID}, nil
}
