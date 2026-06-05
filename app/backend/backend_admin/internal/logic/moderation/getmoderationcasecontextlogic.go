package moderation

import (
	"context"
	"errors"

	"beaver/app/backend/backend_models"
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/group/group_rpc/types/group_rpc"
	"beaver/app/moment/moment_rpc/types/moment_rpc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetModerationCaseContextLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetModerationCaseContextLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetModerationCaseContextLogic {
	return &GetModerationCaseContextLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetModerationCaseContextLogic) GetModerationCaseContext(req *types.GetModerationCaseContextReq) (resp *types.GetModerationCaseContextRes, err error) {
	if req.CaseID == 0 {
		return nil, errors.New("工单ID不能为空")
	}

	var c backend_models.AdminModerationCase
	if err = l.svcCtx.DB.Where("id = ?", req.CaseID).First(&c).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("工单不存在")
		}
		l.Errorf("查询工单失败: %v", err)
		return nil, err
	}

	res := &types.GetModerationCaseContextRes{
		Case:           formatCaseInfo(c),
		RelatedReports: []types.ContentReportInfo{},
		RecentMessages: []types.CaseContextMessage{},
	}

	switch c.TargetType {
	case backend_models.CaseTargetUser:
		res.TargetUser = l.loadUser(c.TargetID)
	case backend_models.CaseTargetMessage:
		msg := l.loadMessage(c.TargetID)
		res.TargetMessage = msg
		if msg != nil && msg.ConversationID != "" {
			res.RecentMessages = l.loadRecentMessages(msg.ConversationID)
		}
	case backend_models.CaseTargetMoment:
		res.TargetMoment = l.loadMoment(c.TargetID)
	case backend_models.CaseTargetGroup:
		res.TargetGroup = l.loadGroup(c.TargetID)
	}

	reportRes, err := l.svcCtx.PlatformRpc.ListContentReports(l.ctx, &platform_rpc.ListContentReportsReq{
		Page:       1,
		PageSize:   20,
		TargetType: int32(c.TargetType),
		TargetId:   c.TargetID,
	})
	if err == nil && reportRes != nil {
		userIDs := make([]string, 0)
		seen := map[string]struct{}{}
		for _, r := range reportRes.List {
			if r.ReporterUserId == "" {
				continue
			}
			if _, ok := seen[r.ReporterUserId]; ok {
				continue
			}
			seen[r.ReporterUserId] = struct{}{}
			userIDs = append(userIDs, r.ReporterUserId)
		}
		users := map[string]*user_rpc.UserInfo{}
		if len(userIDs) > 0 {
			if uRes, uErr := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{UserIdList: userIDs}); uErr == nil && uRes != nil {
				users = uRes.UserInfo
			}
		}
		for _, r := range reportRes.List {
			name := ""
			if u, ok := users[r.ReporterUserId]; ok && u != nil {
				name = u.NickName
			}
			res.RelatedReports = append(res.RelatedReports, types.ContentReportInfo{
				ID:             r.Id,
				ReporterUserID: r.ReporterUserId,
				ReporterName:   name,
				TargetType:     int(r.TargetType),
				TargetID:       r.TargetId,
				ReasonType:     int(r.ReasonType),
				Content:        r.Content,
				Status:         int(r.Status),
				CaseID:         r.CaseId,
				CreatedAt:      r.CreatedAt,
			})
		}
	}

	return res, nil
}

func (l *GetModerationCaseContextLogic) loadUser(userID string) *types.CaseContextUser {
	if userID == "" {
		return nil
	}
	res, err := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{UserIdList: []string{userID}})
	if err != nil || res == nil {
		return nil
	}
	u, ok := res.UserInfo[userID]
	if !ok || u == nil {
		return nil
	}
	return &types.CaseContextUser{
		UserID:   u.UserId,
		NickName: u.NickName,
		Email:    u.Email,
		Status:   int(u.Status),
	}
}

func (l *GetModerationCaseContextLogic) loadMessage(messageID string) *types.CaseContextMessage {
	if messageID == "" {
		return nil
	}
	msgRes, err := l.svcCtx.ChatRpc.ListChatMessages(l.ctx, &chat_rpc.ListChatMessagesReq{
		MessageId:   messageID,
		WithContent: true,
		Page:        1,
		PageSize:    1,
	})
	if err != nil || len(msgRes.List) == 0 {
		return nil
	}
	m := msgRes.List[0]
	sendName := ""
	if m.SendUserId != "" {
		if uRes, uErr := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{UserIdList: []string{m.SendUserId}}); uErr == nil && uRes != nil {
			if u, ok := uRes.UserInfo[m.SendUserId]; ok && u != nil {
				sendName = u.NickName
			}
		}
	}
	return &types.CaseContextMessage{
		MessageID:      m.MessageId,
		ConversationID: m.ConversationId,
		SendUserID:     m.SendUserId,
		SendUserName:   sendName,
		MsgPreview:     m.MsgPreview,
		IsDeleted:      m.Status == chatMessageStatusDeleted,
		CreateTime:     m.CreatedAt,
	}
}

func (l *GetModerationCaseContextLogic) loadRecentMessages(conversationID string) []types.CaseContextMessage {
	msgRes, err := l.svcCtx.ChatRpc.ListChatMessages(l.ctx, &chat_rpc.ListChatMessagesReq{
		ConversationId: conversationID,
		WithContent:    true,
		Order:          2,
		Page:           1,
		PageSize:       30,
	})
	if err != nil || len(msgRes.List) == 0 {
		return nil
	}

	userIDs := make([]string, 0)
	seen := map[string]struct{}{}
	for _, m := range msgRes.List {
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
		if uRes, uErr := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{UserIdList: userIDs}); uErr == nil && uRes != nil {
			users = uRes.UserInfo
		}
	}

	list := make([]types.CaseContextMessage, 0, len(msgRes.List))
	for _, m := range msgRes.List {
		sendName := ""
		if u, ok := users[m.SendUserId]; ok && u != nil {
			sendName = u.NickName
		}
		list = append(list, types.CaseContextMessage{
			MessageID:      m.MessageId,
			ConversationID: m.ConversationId,
			SendUserID:     m.SendUserId,
			SendUserName:   sendName,
			MsgPreview:     m.MsgPreview,
			IsDeleted:      m.Status == chatMessageStatusDeleted,
			CreateTime:     m.CreatedAt,
		})
	}
	return list
}

func (l *GetModerationCaseContextLogic) loadMoment(momentID string) *types.CaseContextMoment {
	if momentID == "" {
		return nil
	}
	res, err := l.svcCtx.MomentRpc.ListMoments(l.ctx, &moment_rpc.ListMomentsReq{
		MomentId: momentID,
		Page:     1,
		PageSize: 1,
	})
	if err != nil || len(res.List) == 0 {
		return nil
	}
	m := res.List[0]
	return &types.CaseContextMoment{
		MomentID: m.MomentId,
		UserID:   m.UserId,
		Content:  m.Content,
	}
}

func (l *GetModerationCaseContextLogic) loadGroup(groupID string) *types.CaseContextGroup {
	if groupID == "" {
		return nil
	}
	res, err := l.svcCtx.GroupRpc.ListGroups(l.ctx, &group_rpc.ListGroupsReq{
		GroupId:  groupID,
		Page:     1,
		PageSize: 1,
	})
	if err != nil || len(res.List) == 0 {
		return nil
	}
	g := res.List[0]
	return &types.CaseContextGroup{
		GroupID: g.GroupId,
		Title:   g.Title,
		Status:  int(g.Status),
	}
}
