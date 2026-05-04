// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package contact

import (
	"context"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchGetUsersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量获取用户信息
func NewBatchGetUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchGetUsersLogic {
	return &BatchGetUsersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchGetUsersLogic) BatchGetUsers(req *types.BatchGetUsersReq) (resp *types.BatchGetUsersRes, err error) {
	// todo: add your logic here and delete this line

	return
}
