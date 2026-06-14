package logic

import (
	"context"
	"errors"

	"beaver/app/user/user_models"
	"beaver/app/user/user_rpc/internal/svc"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUsersStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUsersStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUsersStatusLogic {
	return &UpdateUsersStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUsersStatusLogic) UpdateUsersStatus(in *user_rpc.UpdateUsersStatusReq) (*user_rpc.UpdateUsersStatusRes, error) {
	if in.Status < 1 || in.Status > 3 {
		return nil, errors.New("无效的状态值")
	}
	if len(in.UserIds) == 0 {
		return &user_rpc.UpdateUsersStatusRes{}, nil
	}
	result := l.svcCtx.DB.Model(&user_models.UserModel{}).
		Where("user_id IN ?", in.UserIds).
		Update("status", int8(in.Status))
	if result.Error != nil {
		l.Errorf("批量更新用户状态失败: %v", result.Error)
		return nil, result.Error
	}
	return &user_rpc.UpdateUsersStatusRes{AffectedCount: result.RowsAffected}, nil
}
