package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/friend/friend_rpc/types/friend_rpc"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFriendDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendDetailLogic {
	return &GetFriendDetailLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

// GetFriendDetail 管理后台：好友关系详情。
// admin 职责：校验 relationId，复用 ListFriends 查单条，UserRpc 组装双方昵称与头像。
// RPC 职责：ListFriends(relation_id) 精确查询，不单独暴露 GetFriend RPC。
func (l *GetFriendDetailLogic) GetFriendDetail(req *types.GetFriendDetailReq) (resp *types.GetFriendDetailRes, err error) {
	if req.FriendID == "" {
		return nil, errors.New("好友关系ID不能为空")
	}

	rpcRes, err := l.svcCtx.FriendRpc.ListFriends(l.ctx, &friend_rpc.ListFriendsReq{
		RelationId: req.FriendID,
	})
	if err != nil {
		l.Errorf("获取好友详情失败: %v", err)
		return nil, err
	}
	if len(rpcRes.List) == 0 {
		return nil, errors.New("好友关系不存在")
	}

	f := rpcRes.List[0]
	users := map[string]*user_rpc.UserInfo{}
	if res, err := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{UserIdList: []string{f.SendUserId, f.RevUserId}}); err == nil && res != nil {
		users = res.UserInfo
	}

	sendName, sendAvatar, revName, revAvatar := "", "", "", ""
	if u, ok := users[f.SendUserId]; ok && u != nil {
		sendName, sendAvatar = u.NickName, u.Avatar
	}
	if u, ok := users[f.RevUserId]; ok && u != nil {
		revName, revAvatar = u.NickName, u.Avatar
	}

	return &types.GetFriendDetailRes{
		Id:               f.FriendId,
		SendUserId:       f.SendUserId,
		SendUserName:     sendName,
		SendUserFileName: sendAvatar,
		RevUserId:        f.RevUserId,
		RevUserName:      revName,
		RevUserFileName:  revAvatar,
		SendUserNotice:   f.SendUserNotice,
		RevUserNotice:    f.RevUserNotice,
		IsDeleted:        f.IsDeleted,
		CreateTime:       f.CreatedAt,
		UpdateTime:       f.UpdatedAt,
	}, nil
}
