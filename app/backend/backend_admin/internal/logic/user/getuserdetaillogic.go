package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserDetailLogic {
	return &GetUserDetailLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

// GetUserDetail 管理后台：用户详情。
// admin 职责：校验 userId，复用 ListUsers 查单条并映射响应。
// RPC 职责：ListUsers(user_id) 精确查询，admin 不单独依赖 UserInfo RPC。
func (l *GetUserDetailLogic) GetUserDetail(req *types.GetUserDetailReq) (resp *types.GetUserDetailRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("用户ID不能为空")
	}

	rpcRes, err := l.svcCtx.UserRpc.ListUsers(l.ctx, &user_rpc.ListUsersReq{UserId: req.UserID})
	if err != nil {
		l.Errorf("获取用户详情失败: %v", err)
		return nil, err
	}
	if len(rpcRes.List) == 0 {
		return nil, errors.New("用户不存在")
	}

	u := rpcRes.List[0]
	return &types.GetUserDetailRes{
		Id: u.UserId, NickName: u.NickName, Avatar: u.Avatar, Email: u.Email,
		Abstract: u.Abstract, Status: int(u.Status), Source: int(u.Source), UserType: int(u.UserType),
		CreateTime: u.CreatedAt, UpdateTime: u.UpdatedAt,
	}, nil
}
