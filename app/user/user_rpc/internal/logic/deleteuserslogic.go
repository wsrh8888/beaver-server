package logic

import (
	"context"

	"beaver/app/user/user_models"
	"beaver/app/user/user_rpc/internal/svc"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteUsersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUsersLogic {
	return &DeleteUsersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteUsersLogic) DeleteUsers(in *user_rpc.DeleteUsersReq) (*user_rpc.DeleteUsersRes, error) {
	if len(in.UserIds) == 0 {
		return &user_rpc.DeleteUsersRes{}, nil
	}
	result := l.svcCtx.DB.Model(&user_models.UserModel{}).
		Where("user_id IN ?", in.UserIds).
		Update("status", 3)
	if result.Error != nil {
		l.Errorf("删除用户失败: %v", result.Error)
		return nil, result.Error
	}
	return &user_rpc.DeleteUsersRes{AffectedCount: result.RowsAffected}, nil
}
