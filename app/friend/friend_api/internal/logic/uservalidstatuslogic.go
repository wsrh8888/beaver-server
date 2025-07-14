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
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"
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
	// 参数验证
	if req.UserID == "" || req.VerifyID == 0 {
		return nil, errors.New("用户ID和验证ID不能为空")
	}

	// 状态值验证
	if req.Status < 1 || req.Status > 4 {
		return nil, errors.New("无效的状态值")
	}

	var friendVerify friend_models.FriendVerifyModel
	var conversationID string

	// 查询好友验证记录，确保当前用户是接收方
	err = l.svcCtx.DB.Take(&friendVerify, "id = ? and rev_user_id = ?", req.VerifyID, req.UserID).Error
	if err != nil {
		l.Logger.Errorf("好友验证记录不存在: verifyID=%d, userID=%s, error=%v", req.VerifyID, req.UserID, err)
		return nil, errors.New("好友验证不存在")
	}

	// 检查验证状态是否已处理
	if friendVerify.RevStatus != 0 {
		l.Logger.Errorf("好友验证已处理: verifyID=%d, currentStatus=%d", req.VerifyID, friendVerify.RevStatus)
		return nil, errors.New("该验证已处理，无法重复操作")
	}

	// 处理不同的状态
	switch req.Status {
	case 1: // 同意
		friendVerify.RevStatus = 1

		// 创建好友关系
		err = l.svcCtx.DB.Create(&friend_models.FriendModel{
			SendUserID: friendVerify.SendUserID,
			RevUserID:  friendVerify.RevUserID,
		}).Error
		if err != nil {
			l.Logger.Errorf("创建好友关系失败: %v", err)
			return nil, errors.New("创建好友关系失败")
		}

		// 生成会话ID
		conversationID, err = conversation.GenerateConversation([]string{friendVerify.SendUserID, friendVerify.RevUserID})
		if err != nil {
			l.Logger.Errorf("生成会话ID失败: %v", err)
			return nil, fmt.Errorf("生成会话ID失败: %v", err)
		}

		// 发送默认欢迎消息
		_, err = l.svcCtx.ChatRpc.SendMsg(l.ctx, &chat_rpc.SendMsgReq{
			UserID:         friendVerify.RevUserID,
			ConversationId: conversationID,
			Msg: &chat_rpc.Msg{
				Type: 1,
				TextMsg: &chat_rpc.TextMsg{
					Content: "我们已经是好友了，开始聊天吧",
				},
			},
		})
		if err != nil {
			l.Logger.Errorf("发送欢迎消息失败: %v", err)
			// 不返回错误，因为好友关系已经创建成功
		}

	case 2: // 拒绝
		friendVerify.RevStatus = 2

	case 3: // 忽略
		friendVerify.RevStatus = 3

	case 4: // 删除
		// 直接删除验证记录
		err = l.svcCtx.DB.Delete(&friendVerify).Error
		if err != nil {
			l.Logger.Errorf("删除验证记录失败: %v", err)
			return nil, errors.New("删除验证记录失败")
		}

		l.Logger.Infof("删除好友验证记录成功: verifyID=%d, userID=%s", req.VerifyID, req.UserID)
		return &types.FriendValidStatusRes{}, nil
	}

	// 保存验证状态
	err = l.svcCtx.DB.Save(&friendVerify).Error
	if err != nil {
		l.Logger.Errorf("保存验证状态失败: %v", err)
		return nil, errors.New("保存验证状态失败")
	}

	// 发送WebSocket通知
	ajax.SendMessageToWs(l.svcCtx.Config.Etcd, wsCommandConst.FRIEND_OPERATION, wsTypeConst.FriendRequestReceive, friendVerify.SendUserID, friendVerify.RevUserID, map[string]interface{}{
		"userId": friendVerify.SendUserID,
		"status": friendVerify.RevStatus,
	}, conversationID)
	ajax.SendMessageToWs(l.svcCtx.Config.Etcd, wsCommandConst.FRIEND_OPERATION, wsTypeConst.FriendRequestReceive, friendVerify.RevUserID, friendVerify.SendUserID, map[string]interface{}{
		"userId": friendVerify.RevUserID,
		"status": friendVerify.RevStatus,
	}, conversationID)

	l.Logger.Infof("处理好友验证成功: verifyID=%d, userID=%s, status=%d", req.VerifyID, req.UserID, req.Status)
	return &types.FriendValidStatusRes{}, nil
}
