package logic

import (
	"context"
	"errors"
	"time"

	"beaver/app/auth/auth_models"
	"beaver/app/auth/auth_rpc/internal/svc"
	"beaver/app/auth/auth_rpc/types/auth_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLoginRecordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateLoginRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLoginRecordLogic {
	return &UpdateLoginRecordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateLoginRecordLogic) UpdateLoginRecord(in *auth_rpc.UpdateLoginRecordReq) (*auth_rpc.UpdateLoginRecordRes, error) {
	// 验证必填字段
	if in.UserId == "" {
		return nil, errors.New("用户ID不能为空")
	}

	// 查询用户凭证
	var credential auth_models.AuthCredentialModel
	err := l.svcCtx.DB.Take(&credential, "user_id = ?", in.UserId).Error
	if err != nil {
		logx.Errorf("查询用户凭证失败: %v", err)
		return nil, errors.New("用户凭证不存在")
	}

	// 更新登录记录
	now := time.Now()
	credential.LastLoginAt = &now
	credential.LoginCount++

	err = l.svcCtx.DB.Save(&credential).Error
	if err != nil {
		logx.Errorf("更新登录记录失败: %v", err)
		return nil, errors.New("更新登录记录失败")
	}

	logx.Infof("登录记录更新成功: userID=%s, loginCount=%d", in.UserId, credential.LoginCount)

	return &auth_rpc.UpdateLoginRecordRes{
		Success: true,
	}, nil
}
