package operations

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/friend/friend_rpc/types/friend_rpc"
	"beaver/app/group/group_rpc/types/group_rpc"
	"beaver/app/moment/moment_rpc/types/moment_rpc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

const opsPreviewLimit int32 = 10

type GetUserOperationsProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserOperationsProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserOperationsProfileLogic {
	return &GetUserOperationsProfileLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetUserOperationsProfileLogic) GetUserOperationsProfile(req *types.GetUserOperationsProfileReq) (resp *types.GetUserOperationsProfileRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("用户ID不能为空")
	}

	userRes, err := l.svcCtx.UserRpc.ListUsers(l.ctx, &user_rpc.ListUsersReq{UserId: req.UserID})
	if err != nil {
		l.Errorf("查询用户失败: %v", err)
		return nil, err
	}
	if len(userRes.List) == 0 {
		return nil, errors.New("用户不存在")
	}
	u := userRes.List[0]

	resp = &types.GetUserOperationsProfileRes{
		Profile: types.UserOpsProfileInfo{
			UserID: u.UserId, NickName: u.NickName, Email: u.Email, Avatar: u.Avatar,
			Abstract: u.Abstract, Status: int(u.Status), Source: int(u.Source),
			CreateTime: u.CreatedAt,
		},
		Friends:  []types.UserOpsFriendItem{},
		Groups:   []types.UserOpsGroupItem{},
		Sessions: []types.UserOpsSessionItem{},
		Moments:  []types.UserOpsMomentItem{},
		Reports:  []types.UserOpsReportItem{},
		Blocks:   []types.UserOpsBlockItem{},
	}

	friendRes, err := l.svcCtx.FriendRpc.ListFriends(l.ctx, &friend_rpc.ListFriendsReq{
		UserId: req.UserID, Page: 1, PageSize: opsPreviewLimit,
	})
	if err != nil {
		l.Errorf("查询好友失败: %v", err)
	} else {
		resp.FriendTotal = friendRes.Total
		resp.Friends = l.mapFriends(friendRes.List, req.UserID)
	}

	groupIDsRes, err := l.svcCtx.GroupRpc.GetUserGroupIDs(l.ctx, &group_rpc.GetUserGroupIDsReq{UserID: req.UserID})
	if err != nil {
		l.Errorf("查询用户群组ID失败: %v", err)
	} else if groupIDsRes != nil {
		resp.GroupTotal = int64(len(groupIDsRes.GroupIDs))
		resp.Groups = l.mapGroups(groupIDsRes.GroupIDs)
	}

	sessionRes, err := l.svcCtx.ChatRpc.ListConversations(l.ctx, &chat_rpc.ListConversationsReq{
		UserId: req.UserID, Page: 1, PageSize: opsPreviewLimit,
	})
	if err != nil {
		l.Errorf("查询会话失败: %v", err)
	} else {
		resp.SessionTotal = sessionRes.Total
		resp.Sessions = l.mapSessions(sessionRes.List, req.UserID)
	}

	momentRes, err := l.svcCtx.MomentRpc.ListMoments(l.ctx, &moment_rpc.ListMomentsReq{
		UserId: req.UserID, Page: 1, PageSize: opsPreviewLimit,
	})
	if err != nil {
		l.Errorf("查询动态失败: %v", err)
	} else {
		resp.MomentTotal = momentRes.Total
		for _, m := range momentRes.List {
			resp.Moments = append(resp.Moments, types.UserOpsMomentItem{
				MomentID: m.MomentId, Content: m.Content, IsDeleted: m.IsDeleted, CreatedAt: m.CreatedAt,
			})
		}
	}

	reportRes, err := l.svcCtx.PlatformRpc.ListContentReports(l.ctx, &platform_rpc.ListContentReportsReq{
		ReporterUserId: req.UserID, Page: 1, PageSize: opsPreviewLimit,
	})
	if err != nil {
		l.Errorf("查询举报失败: %v", err)
	} else {
		resp.ReportTotal = reportRes.Total
		for _, r := range reportRes.List {
			resp.Reports = append(resp.Reports, types.UserOpsReportItem{
				ID: r.Id, TargetType: int(r.TargetType), TargetID: r.TargetId,
				ReasonType: int(r.ReasonType), Status: int(r.Status), CreatedAt: r.CreatedAt,
			})
		}
	}

	blockRes, err := l.svcCtx.FriendRpc.ListFriendBlocks(l.ctx, &friend_rpc.ListFriendBlocksReq{
		UserId: req.UserID, Page: 1, PageSize: opsPreviewLimit,
	})
	if err != nil {
		l.Errorf("查询黑名单失败: %v", err)
	} else {
		resp.BlockTotal = blockRes.Total
		resp.Blocks = l.mapBlocks(blockRes.List)
	}

	return resp, nil
}

func (l *GetUserOperationsProfileLogic) mapFriends(list []*friend_rpc.FriendItem, selfID string) []types.UserOpsFriendItem {
	users := l.loadUsers(list, selfID)
	out := make([]types.UserOpsFriendItem, 0, len(list))
	for _, f := range list {
		peerID := f.RevUserId
		if peerID == selfID {
			peerID = f.SendUserId
		}
		name := peerID
		if u, ok := users[peerID]; ok && u != nil && u.NickName != "" {
			name = u.NickName
		}
		out = append(out, types.UserOpsFriendItem{
			PeerUserID: peerID, PeerUserName: name, CreateTime: f.CreatedAt,
		})
	}
	return out
}

func (l *GetUserOperationsProfileLogic) loadUsers(list []*friend_rpc.FriendItem, selfID string) map[string]*user_rpc.UserInfo {
	seen := map[string]struct{}{}
	ids := make([]string, 0, len(list))
	for _, f := range list {
		for _, uid := range []string{f.SendUserId, f.RevUserId} {
			if uid == "" || uid == selfID {
				continue
			}
			if _, ok := seen[uid]; ok {
				continue
			}
			seen[uid] = struct{}{}
			ids = append(ids, uid)
		}
	}
	if len(ids) == 0 {
		return nil
	}
	res, err := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{UserIdList: ids})
	if err != nil || res == nil {
		return nil
	}
	return res.UserInfo
}

func (l *GetUserOperationsProfileLogic) mapGroups(groupIDs []string) []types.UserOpsGroupItem {
	if len(groupIDs) == 0 {
		return nil
	}
	limit := int(opsPreviewLimit)
	if len(groupIDs) < limit {
		limit = len(groupIDs)
	}
	ids := groupIDs[:limit]
	gRes, err := l.svcCtx.GroupRpc.GetGroupsListByIds(l.ctx, &group_rpc.GetGroupsListByIdsReq{GroupIDs: ids})
	if err != nil || gRes == nil {
		out := make([]types.UserOpsGroupItem, 0, len(ids))
		for _, gid := range ids {
			out = append(out, types.UserOpsGroupItem{GroupID: gid, Title: gid})
		}
		return out
	}
	out := make([]types.UserOpsGroupItem, 0, len(gRes.Groups))
	for _, g := range gRes.Groups {
		out = append(out, types.UserOpsGroupItem{
			GroupID: g.GroupID, Title: g.Name,
		})
	}
	return out
}

func (l *GetUserOperationsProfileLogic) mapSessions(list []*chat_rpc.ConversationDetailItem, selfID string) []types.UserOpsSessionItem {
	users := map[string]*user_rpc.UserInfo{}
	seen := map[string]struct{}{}
	for _, item := range list {
		for _, uid := range item.ParticipantIds {
			if uid == "" || uid == selfID {
				continue
			}
			if _, ok := seen[uid]; ok {
				continue
			}
			seen[uid] = struct{}{}
		}
	}
	userIDs := make([]string, 0, len(seen))
	for uid := range seen {
		userIDs = append(userIDs, uid)
	}
	if len(userIDs) > 0 {
		if res, err := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{UserIdList: userIDs}); err == nil && res != nil {
			users = res.UserInfo
		}
	}

	out := make([]types.UserOpsSessionItem, 0, len(list))
	for _, item := range list {
		title := item.ConversationId
		if item.Type == 1 {
			for _, uid := range item.ParticipantIds {
				if uid != selfID {
					if u, ok := users[uid]; ok && u != nil && u.NickName != "" {
						title = u.NickName
					} else {
						title = uid
					}
					break
				}
			}
		} else {
			title = item.ConversationId
		}
		out = append(out, types.UserOpsSessionItem{
			ConversationID: item.ConversationId, ConversationType: int(item.Type),
			Title: title, LastMessage: item.LastMessage, LastMessageTime: item.LastMessageTime,
			MessageCount: item.MessageCount,
		})
	}
	return out
}

func (l *GetUserOperationsProfileLogic) mapBlocks(list []*friend_rpc.FriendBlockItem) []types.UserOpsBlockItem {
	seen := map[string]struct{}{}
	ids := make([]string, 0, len(list))
	for _, b := range list {
		if b.BlockedUserId == "" {
			continue
		}
		if _, ok := seen[b.BlockedUserId]; ok {
			continue
		}
		seen[b.BlockedUserId] = struct{}{}
		ids = append(ids, b.BlockedUserId)
	}
	users := map[string]*user_rpc.UserInfo{}
	if len(ids) > 0 {
		if res, err := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{UserIdList: ids}); err == nil && res != nil {
			users = res.UserInfo
		}
	}
	out := make([]types.UserOpsBlockItem, 0, len(list))
	for _, b := range list {
		name := b.BlockedUserId
		if u, ok := users[b.BlockedUserId]; ok && u != nil {
			name = u.NickName
		}
		out = append(out, types.UserOpsBlockItem{
			ID: b.BlockId, BlockedUserID: b.BlockedUserId, BlockedUserName: name, CreateTime: b.CreatedAt,
		})
	}
	return out
}
