package logic

import (
	"context"
	"errors"

	"beaver/app/auth/auth_api/internal/svc"
	"beaver/app/auth/auth_api/internal/types"
	"beaver/app/user/user_models"
	"beaver/app/user/user_rpc/types/user_rpc"
	"beaver/utils/pwd"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterRes, err error) {
	// todo: add your logic here and delete this line

	var user user_models.UserModel
	err = l.svcCtx.DB.Take(&user, "phone = ?", req.Phone).Error
	if err != nil {
		// 注册逻辑
		_, err := l.svcCtx.UserRpc.UserCreate(l.ctx, &user_rpc.UserCreateReq{
			Phone:    req.Phone,
			Password: pwd.HahPwd(req.Password),
			Source:   1,
		})
		if err != nil {
			logx.Error(err)
			return nil, errors.New("注册失败")
		}
	} else {
		return nil, errors.New("用户已存在")
	}

	return nil, nil
}
