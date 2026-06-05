package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/friend/friend_rpc/types/friend_rpc"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendVerifyListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFriendVerifyListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendVerifyListLogic {
	return &GetFriendVerifyListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

// GetFriendVerifyList 管理后台：好友验证列表。
// admin 职责：运营筛选条件映射、UserRpc 批量组装双方昵称。
// RPC 职责：ListFriendVerifies 领域查询。
func (l *GetFriendVerifyListLogic) GetFriendVerifyList(req *types.GetFriendVerifyListReq) (resp *types.GetFriendVerifyListRes, err error) {
	rpcRes, err := l.svcCtx.FriendRpc.ListFriendVerifies(l.ctx, &friend_rpc.ListFriendVerifiesReq{
		Page:       int32(req.Page),
		PageSize:   int32(req.PageSize),
		SendUserId: req.SendUserId,
		RevUserId:  req.RevUserId,
		SendStatus: int32(req.SendStatus),
		RevStatus:  int32(req.RevStatus),
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
	})
	if err != nil {
		l.Errorf("获取好友验证列表失败: %v", err)
		return nil, err
	}

	users := map[string]*user_rpc.UserInfo{}
	seen := map[string]struct{}{}
	userIDs := make([]string, 0, len(rpcRes.List)*2)
	for _, v := range rpcRes.List {
		for _, id := range []string{v.SendUserId, v.RevUserId} {
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

	list := make([]types.GetFriendVerifyListItem, 0, len(rpcRes.List))
	for _, v := range rpcRes.List {
		sendName, revName := "", ""
		if u, ok := users[v.SendUserId]; ok && u != nil {
			sendName = u.NickName
		}
		if u, ok := users[v.RevUserId]; ok && u != nil {
			revName = u.NickName
		}
		list = append(list, types.GetFriendVerifyListItem{
			Id:           v.VerifyId,
			SendUserId:   v.SendUserId,
			SendUserName: sendName,
			RevUserId:    v.RevUserId,
			RevUserName:  revName,
			SendStatus:   int(v.SendStatus),
			RevStatus:    int(v.RevStatus),
			Message:      v.Message,
			CreateTime:   v.CreatedAt,
			UpdateTime:   v.UpdatedAt,
		})
	}
	return &types.GetFriendVerifyListRes{List: list, Total: rpcRes.Total}, nil
}
