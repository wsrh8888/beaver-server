package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/friend/friend_rpc/types/friend_rpc"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendListLogic {
	return &GetFriendListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

// GetFriendList 管理后台：好友关系列表。
// admin 职责：运营筛选条件映射、UserRpc 批量组装双方昵称。
// RPC 职责：ListFriends 领域查询，不与本 HTTP 接口 1:1。
func (l *GetFriendListLogic) GetFriendList(req *types.GetFriendListReq) (resp *types.GetFriendListRes, err error) {
	rpcReq := &friend_rpc.ListFriendsReq{
		Page:       int32(req.Page),
		PageSize:   int32(req.PageSize),
		UserId:     req.UserID,
		PeerUserId: req.FriendID,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
	}
	if req.IsDeleted {
		deleted := true
		rpcReq.IsDeleted = &deleted
	}

	rpcRes, err := l.svcCtx.FriendRpc.ListFriends(l.ctx, rpcReq)
	if err != nil {
		l.Errorf("获取好友列表失败: %v", err)
		return nil, err
	}

	users := map[string]*user_rpc.UserInfo{}
	seen := map[string]struct{}{}
	userIDs := make([]string, 0, len(rpcRes.List)*2)
	for _, f := range rpcRes.List {
		for _, id := range []string{f.SendUserId, f.RevUserId} {
			if id == "" {
				continue
			}
			if _, ok := seen[id]; ok {
				continue
			}
			seen[id] = struct{}{}
			userIDs = append(userIDs, id)
		}
	}
	if len(userIDs) > 0 {
		if res, err := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{UserIdList: userIDs}); err == nil && res != nil {
			users = res.UserInfo
		}
	}

	list := make([]types.GetFriendListItem, 0, len(rpcRes.List))
	for _, f := range rpcRes.List {
		sendName, revName := "", ""
		if u, ok := users[f.SendUserId]; ok && u != nil {
			sendName = u.NickName
		}
		if u, ok := users[f.RevUserId]; ok && u != nil {
			revName = u.NickName
		}
		list = append(list, types.GetFriendListItem{
			Id:             f.FriendId,
			SendUserId:     f.SendUserId,
			SendUserName:   sendName,
			RevUserId:      f.RevUserId,
			RevUserName:    revName,
			SendUserNotice: f.SendUserNotice,
			RevUserNotice:  f.RevUserNotice,
			IsDeleted:      f.IsDeleted,
			CreateTime:     f.CreatedAt,
			UpdateTime:     f.UpdatedAt,
		})
	}
	return &types.GetFriendListRes{List: list, Total: rpcRes.Total}, nil
}
