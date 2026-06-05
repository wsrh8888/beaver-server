package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetChatMessageDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetChatMessageDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatMessageDetailLogic {
	return &GetChatMessageDetailLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetChatMessageDetailLogic) GetChatMessageDetail(req *types.GetChatMessageDetailReq) (resp *types.GetChatMessageDetailRes, err error) {
	if req.MessageID == "" {
		return nil, errors.New("消息ID不能为空")
	}

	rpcRes, err := l.svcCtx.ChatRpc.ListChatMessages(l.ctx, &chat_rpc.ListChatMessagesReq{
		MessageId:   req.MessageID,
		WithContent: true,
		Page:        1,
		PageSize:    1,
	})
	if err != nil {
		l.Errorf("获取聊天消息详情失败: %v", err)
		return nil, err
	}
	if len(rpcRes.List) == 0 {
		return nil, errors.New("聊天消息不存在")
	}

	m := rpcRes.List[0]
	sendName, sendAvatar := "", ""
	if m.SendUserId != "" {
		if res, err := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{UserIdList: []string{m.SendUserId}}); err == nil && res != nil {
			if u, ok := res.UserInfo[m.SendUserId]; ok && u != nil {
				sendName, sendAvatar = u.NickName, u.Avatar
			}
		}
	}

	return &types.GetChatMessageDetailRes{
		Id:               m.MessageId,
		MessageID:        m.MessageId,
		ConversationID:   m.ConversationId,
		SendUserID:       m.SendUserId,
		SendUserName:     sendName,
		SendUserFileName: sendAvatar,
		MsgType:          int(m.MsgType),
		MsgPreview:       m.MsgPreview,
		MsgContent:       m.MsgContent,
		IsDeleted:        m.Status == chatMessageStatusDeleted,
		CreateTime:       m.CreatedAt,
		UpdateTime:       m.UpdatedAt,
	}, nil
}
