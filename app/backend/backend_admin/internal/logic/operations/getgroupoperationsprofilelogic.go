package operations

import (
	"context"
	"errors"
	"fmt"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/group/group_rpc/types/group_rpc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupOperationsProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetGroupOperationsProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupOperationsProfileLogic {
	return &GetGroupOperationsProfileLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetGroupOperationsProfileLogic) GetGroupOperationsProfile(req *types.GetGroupOperationsProfileReq) (resp *types.GetGroupOperationsProfileRes, err error) {
	if req.GroupID == "" {
		return nil, errors.New("群组ID不能为空")
	}

	groupRes, err := l.svcCtx.GroupRpc.ListGroups(l.ctx, &group_rpc.ListGroupsReq{
		GroupId: req.GroupID, Page: 1, PageSize: 1,
	})
	if err != nil {
		l.Errorf("查询群组失败: %v", err)
		return nil, err
	}
	if len(groupRes.List) == 0 {
		return nil, errors.New("群组不存在")
	}
	g := groupRes.List[0]

	resp = &types.GetGroupOperationsProfileRes{
		Profile: types.GroupOpsProfileInfo{
			GroupID: g.GroupId, Title: g.Title, Avatar: g.Avatar, CreatorID: g.CreatorId,
			Notice: g.Notice, Status: int(g.Status), MuteAll: g.MuteAll, CreatedAt: g.CreatedAt,
		},
		Members:  []types.GroupOpsMemberItem{},
		Messages: []types.GroupOpsMessageItem{},
		Reports:  []types.GroupOpsReportItem{},
	}

	memberRes, err := l.svcCtx.GroupRpc.ListGroupMembers(l.ctx, &group_rpc.ListGroupMembersReq{
		GroupId: req.GroupID, Page: 1, PageSize: opsPreviewLimit,
	})
	if err != nil {
		l.Errorf("查询群成员失败: %v", err)
	} else {
		resp.MemberTotal = memberRes.Total
		resp.Members = l.mapMembers(memberRes.List)
	}

	convIDs := []string{fmt.Sprintf("group_%s", req.GroupID), req.GroupID}
	for _, convID := range convIDs {
		msgRes, msgErr := l.svcCtx.ChatRpc.ListChatMessages(l.ctx, &chat_rpc.ListChatMessagesReq{
			ConversationId: convID, Page: 1, PageSize: opsPreviewLimit, Order: 2,
		})
		if msgErr != nil || len(msgRes.List) == 0 {
			continue
		}
		resp.MessageTotal = msgRes.Total
		resp.Messages = l.mapMessages(msgRes.List)
		break
	}

	reportRes, err := l.svcCtx.PlatformRpc.ListContentReports(l.ctx, &platform_rpc.ListContentReportsReq{
		TargetType: 4, TargetId: req.GroupID, Page: 1, PageSize: opsPreviewLimit,
	})
	if err != nil {
		l.Errorf("查询群举报失败: %v", err)
	} else {
		resp.ReportTotal = reportRes.Total
		resp.Reports = l.mapGroupReports(reportRes.List)
	}

	return resp, nil
}

func (l *GetGroupOperationsProfileLogic) mapMembers(list []*group_rpc.GroupMemberItem) []types.GroupOpsMemberItem {
	userIDs := make([]string, 0, len(list))
	seen := map[string]struct{}{}
	for _, m := range list {
		if m.UserId == "" {
			continue
		}
		if _, ok := seen[m.UserId]; ok {
			continue
		}
		seen[m.UserId] = struct{}{}
		userIDs = append(userIDs, m.UserId)
	}
	users := map[string]*user_rpc.UserInfo{}
	if len(userIDs) > 0 {
		if res, err := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{UserIdList: userIDs}); err == nil && res != nil {
			users = res.UserInfo
		}
	}
	out := make([]types.GroupOpsMemberItem, 0, len(list))
	for _, m := range list {
		name := m.UserId
		if u, ok := users[m.UserId]; ok && u != nil && u.NickName != "" {
			name = u.NickName
		}
		out = append(out, types.GroupOpsMemberItem{
			UserID: m.UserId, NickName: name, Role: int(m.Role),
			Status: int(m.Status), JoinTime: m.CreatedAt,
		})
	}
	return out
}

func (l *GetGroupOperationsProfileLogic) mapMessages(list []*chat_rpc.ChatMessageItem) []types.GroupOpsMessageItem {
	userIDs := make([]string, 0, len(list))
	seen := map[string]struct{}{}
	for _, m := range list {
		if m.SendUserId == "" {
			continue
		}
		if _, ok := seen[m.SendUserId]; ok {
			continue
		}
		seen[m.SendUserId] = struct{}{}
		userIDs = append(userIDs, m.SendUserId)
	}
	users := map[string]*user_rpc.UserInfo{}
	if len(userIDs) > 0 {
		if res, err := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{UserIdList: userIDs}); err == nil && res != nil {
			users = res.UserInfo
		}
	}
	const deletedStatus int32 = 4
	out := make([]types.GroupOpsMessageItem, 0, len(list))
	for _, m := range list {
		name := m.SendUserId
		if u, ok := users[m.SendUserId]; ok && u != nil {
			name = u.NickName
		}
		out = append(out, types.GroupOpsMessageItem{
			MessageID: m.MessageId, SendUserID: m.SendUserId, SendName: name,
			MsgPreview: m.MsgPreview, IsDeleted: m.Status == deletedStatus, CreateTime: m.CreatedAt,
		})
	}
	return out
}

func (l *GetGroupOperationsProfileLogic) mapGroupReports(list []*platform_rpc.ContentReportItem) []types.GroupOpsReportItem {
	reporterIDs := make([]string, 0, len(list))
	seen := map[string]struct{}{}
	for _, r := range list {
		if r.ReporterUserId == "" {
			continue
		}
		if _, ok := seen[r.ReporterUserId]; ok {
			continue
		}
		seen[r.ReporterUserId] = struct{}{}
		reporterIDs = append(reporterIDs, r.ReporterUserId)
	}
	users := map[string]*user_rpc.UserInfo{}
	if len(reporterIDs) > 0 {
		if res, err := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{UserIdList: reporterIDs}); err == nil && res != nil {
			users = res.UserInfo
		}
	}
	out := make([]types.GroupOpsReportItem, 0, len(list))
	for _, r := range list {
		name := ""
		if u, ok := users[r.ReporterUserId]; ok && u != nil {
			name = u.NickName
		}
		out = append(out, types.GroupOpsReportItem{
			ID: r.Id, ReporterUserID: r.ReporterUserId, ReporterName: name,
			ReasonType: int(r.ReasonType), Status: int(r.Status), CreatedAt: r.CreatedAt,
		})
	}
	return out
}
