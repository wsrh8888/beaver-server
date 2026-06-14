package user

import (
	"context"
	"errors"

	"beaver/app/open/constants"
	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserListLogic {
	return &GetUserListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserListLogic) GetUserList(req *types.GetUserListReq) (resp *types.GetUserListRes, err error) {
	tokenRecord, err := loadUserAccessToken(l.svcCtx.DB, req.Authorization, constants.ScopeUserProfileRead)
	if err != nil {
		return nil, err
	}
	if len(req.UserIDs) == 0 {
		return nil, errors.New("userIds 不能为空")
	}
	for _, uid := range req.UserIDs {
		if uid != tokenRecord.UserID {
			return nil, errors.New("无权查询该用户信息")
		}
	}

	listRes, err := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{
		UserIdList: req.UserIDs,
	})
	if err != nil {
		return nil, err
	}

	scopes := parseScopeList(tokenRecord.Scope)
	userList := make([]types.GetUserListUserItem, 0, len(req.UserIDs))
	for _, uid := range req.UserIDs {
		info, ok := listRes.UserInfo[uid]
		if !ok {
			continue
		}
		item := types.GetUserListUserItem{
			UserID:   info.UserId,
			Nickname: info.NickName,
		}
		if hasScope(scopes, constants.ScopeUserAvatarRead) {
			item.Avatar = info.Avatar
		}
		if hasScope(scopes, constants.ScopeUserEmailRead) {
			item.Email = info.Email
		}
		userList = append(userList, item)
	}

	return &types.GetUserListRes{Users: userList}, nil
}
