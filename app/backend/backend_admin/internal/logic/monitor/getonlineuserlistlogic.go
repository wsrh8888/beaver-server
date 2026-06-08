package monitor

import (
	"context"
	"strings"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/user/user_rpc/types/user_rpc"
	"beaver/core/coreonline"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOnlineUserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 在线用户列表
func NewGetOnlineUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOnlineUserListLogic {
	return &GetOnlineUserListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOnlineUserListLogic) GetOnlineUserList(req *types.GetOnlineUserListReq) (resp *types.GetOnlineUserListRes, err error) {
	online, err := coreonline.List(l.svcCtx.Redis)
	if err != nil {
		l.Errorf("获取在线用户列表失败: %v", err)
		return nil, err
	}

	userIDs := make([]string, 0, len(online))
	for _, user := range online {
		userIDs = append(userIDs, user.UserID)
	}

	userMap := map[string]*user_rpc.UserInfo{}
	if len(userIDs) > 0 {
		res, rpcErr := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{UserIdList: userIDs})
		if rpcErr != nil {
			l.Errorf("批量获取用户信息失败: %v", rpcErr)
		} else if res != nil {
			userMap = res.UserInfo
		}
	}

	keyword := strings.TrimSpace(strings.ToLower(req.Keyword))
	filtered := make([]coreonline.User, 0, len(online))
	for _, user := range online {
		if keyword == "" {
			filtered = append(filtered, user)
			continue
		}

		if strings.Contains(strings.ToLower(user.UserID), keyword) {
			filtered = append(filtered, user)
			continue
		}

		info := userMap[user.UserID]
		if info != nil {
			if strings.Contains(strings.ToLower(info.NickName), keyword) ||
				strings.Contains(strings.ToLower(info.Email), keyword) {
				filtered = append(filtered, user)
			}
		}
	}

	total := int64(len(filtered))
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	start := (page - 1) * pageSize
	if start >= len(filtered) {
		return &types.GetOnlineUserListRes{List: []types.OnlineUserItem{}, Total: total}, nil
	}
	end := start + pageSize
	if end > len(filtered) {
		end = len(filtered)
	}

	list := make([]types.OnlineUserItem, 0, end-start)
	for _, user := range filtered[start:end] {
		info := userMap[user.UserID]
		item := types.OnlineUserItem{
			UserID: user.UserID,
			Slots:  make([]types.OnlineUserSlotItem, 0, len(user.Slots)),
		}
		if info != nil {
			item.NickName = info.NickName
			item.Email = info.Email
			item.Avatar = info.Avatar
		}
		for _, slot := range user.Slots {
			item.Slots = append(item.Slots, types.OnlineUserSlotItem{
				InstanceID: slot.InstanceID,
				Slot:       slot.Slot,
			})
		}
		list = append(list, item)
	}

	return &types.GetOnlineUserListRes{
		List:  list,
		Total: total,
	}, nil
}
