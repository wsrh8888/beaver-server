package logic

import (
	"context"
	"errors"

	"beaver/app/auth/auth_models"
	"beaver/app/auth/auth_rpc/internal/svc"
	"beaver/app/auth/auth_rpc/types/auth_rpc"
	"beaver/utils/pwd"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCredentialLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateCredentialLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCredentialLogic {
	return &CreateCredentialLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateCredentialLogic) CreateCredential(in *auth_rpc.CreateCredentialReq) (*auth_rpc.CreateCredentialRes, error) {
	// 验证必填字段
	if in.UserId == "" {
		return nil, errors.New("用户ID不能为空")
	}
	if in.Password == "" {
		return nil, errors.New("密码不能为空")
	}

	// 检查凭证是否已存在
	var credential auth_models.AuthCredentialModel
	err := l.svcCtx.DB.Take(&credential, "user_id = ?", in.UserId).Error
	if err == nil {
		return nil, errors.New("用户凭证已存在")
	}

	// 加密密码
	hashedPassword := pwd.HahPwd(in.Password)

	// 创建用户凭证
	credential = auth_models.AuthCredentialModel{
		UserID:   in.UserId,
		Password: hashedPassword,
	}

	err = l.svcCtx.DB.Create(&credential).Error
	if err != nil {
		logx.Errorf("创建用户凭证失败: %v", err)
		return nil, errors.New("创建用户凭证失败")
	}

	logx.Infof("用户凭证创建成功: userID=%s", in.UserId)

	return &auth_rpc.CreateCredentialRes{
		Success: true,
	}, nil
}
