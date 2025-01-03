package logic

import (
	"context"
	"errors"
	"fmt"

	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/friend/friend_api/internal/svc"
	"beaver/app/friend/friend_api/internal/types"
	"beaver/app/friend/friend_models"
	"beaver/common/ajax"
	"beaver/utils/conversation"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserValidStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserValidStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserValidStatusLogic {
	return &UserValidStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserValidStatusLogic) UserValidStatus(req *types.FriendValidStatusReq) (resp *types.FriendValidStatusRes, err error) {
	// todo: add your logic here and delete this line
	var friendVerify friend_models.FriendVerifyModel
	// 我要操作状态，我自己是接受方
	err = l.svcCtx.DB.Take(&friendVerify, "id = ? and rev_user_id = ?", req.VerifyId, req.UserId).Error
	if err != nil {
		return nil, errors.New("好友验证不存在")
	}
	if friendVerify.RevStatus != 0 {
		return nil, errors.New("操作异常")
	}
	switch req.Status {
	case 1: // 同意
		friendVerify.RevStatus = 1
		// 往好友表里插入
		l.svcCtx.DB.Create(&friend_models.FriendModel{
			SendUserId: friendVerify.SendUserId,
			RevUserId:  friendVerify.RevUserId,
		})

		fmt.Println("发送消息")
		conversationId, _ := conversation.GenerateConversation([]string{friendVerify.SendUserId, friendVerify.RevUserId})
		// 默认发送一条消息
		l.svcCtx.ChatRpc.SendMsg(l.ctx, &chat_rpc.SendMsgReq{
			UserId:         friendVerify.RevUserId,
			ConversationId: conversationId,
			Msg: &chat_rpc.Msg{
				Type: 1,
				TextMsg: &chat_rpc.TextMsg{
					Content: "我们已经是好友了，开始聊天吧",
				},
			},
		})

	case 2: // 拒绝
		friendVerify.RevStatus = 2

	case 3: //忽略
		friendVerify.RevStatus = 3
	case 4: //删除
		// 一条验证记录是两个人看的
		l.svcCtx.DB.Delete(&friendVerify)
		return nil, nil
	}
	err = l.svcCtx.DB.Save(&friendVerify).Error

	ajax.SendMessageToWs(l.svcCtx.Config.Etcd, "user_valid_type_update", friendVerify.SendUserId, friendVerify.RevUserId, map[string]interface{}{
		"userId": friendVerify.SendUserId,
		"status": friendVerify.RevStatus,
	})
	ajax.SendMessageToWs(l.svcCtx.Config.Etcd, "user_valid_type_update", friendVerify.RevUserId, friendVerify.SendUserId, map[string]interface{}{
		"userId": friendVerify.RevUserId,
		"status": friendVerify.RevStatus,
	})
	return
}
