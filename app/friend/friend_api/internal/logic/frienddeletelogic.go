package logic

import (
	"context"
	"errors"
	"fmt"

	"beaver/app/chat/chat_models"
	"beaver/app/friend/friend_api/internal/svc"
	"beaver/app/friend/friend_api/internal/types"
	"beaver/app/friend/friend_models"
	"beaver/common/ajax"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"
	"beaver/utils/conversation"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFriendDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendDeleteLogic {
	return &FriendDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendDeleteLogic) FriendDelete(req *types.FriendDeleteReq) (resp *types.FriendDeleteRes, err error) {
	// 参数验证
	if req.UserID == "" || req.FriendID == "" {
		return nil, errors.New("用户ID和好友ID不能为空")
	}

	// 不能删除自己
	if req.UserID == req.FriendID {
		return nil, errors.New("不能删除自己")
	}

	// 确认好友关系
	var friend friend_models.FriendModel
	if !friend.IsFriend(l.svcCtx.DB, req.UserID, req.FriendID) {
		l.Logger.Errorf("尝试删除非好友关系: userID=%s, friendID=%s", req.UserID, req.FriendID)
		return nil, errors.New("不是好友关系")
	}

	// 标记好友关系为已删除
	err = l.svcCtx.DB.Model(&friend_models.FriendModel{}).Where(
		"((send_user_id = ? AND rev_user_id = ?) OR (send_user_id = ? AND rev_user_id = ?)) AND is_deleted = 0",
		req.UserID, req.FriendID, req.FriendID, req.UserID).Update("is_deleted", 1).Error
	if err != nil {
		l.Logger.Errorf("标记好友关系删除失败: %v", err)
		return nil, errors.New("删除好友失败")
	}

	// 获取会话Id
	conversationID, err := conversation.GenerateConversation([]string{req.UserID, req.FriendID})
	if err != nil {
		l.Logger.Errorf("生成会话Id失败: %v", err)
		return nil, fmt.Errorf("生成会话Id失败: %v", err)
	}

	// 异步处理相关数据清理
	go func() {
		// 发送WebSocket通知
		ajax.SendMessageToWs(l.svcCtx.Config.Etcd, wsCommandConst.FRIEND_OPERATION, wsTypeConst.FriendDelete, req.UserID, req.FriendID, map[string]interface{}{
			"userId": req.FriendID,
		}, "")
		ajax.SendMessageToWs(l.svcCtx.Config.Etcd, wsCommandConst.FRIEND_OPERATION, wsTypeConst.FriendDelete, req.FriendID, req.UserID, map[string]interface{}{
			"userId": req.UserID,
		}, "")

		// 标记会话和聊天记录为已删除
		if err := l.markConversationAndChatsAsDeleted(req.UserID, conversationID); err != nil {
			l.Logger.Errorf("删除会话和聊天记录失败: userID=%s, conversationID=%s, error=%v", req.UserID, conversationID, err)
		}
	}()

	l.Logger.Infof("好友删除成功: userID=%s, friendID=%s", req.UserID, req.FriendID)
	return &types.FriendDeleteRes{}, nil
}

func (l *FriendDeleteLogic) markConversationAndChatsAsDeleted(userID, conversationID string) error {

	db := l.svcCtx.DB

	// 批量标记会话记录为已删除
	err := db.Model(&chat_models.ChatUserConversationModel{}).Where("user_id = ? AND conversation_id = ?", userID, conversationID).Update("is_deleted", true).Error
	if err != nil {
		return err
	}

	// 批量标记聊天记录为已删除
	err = db.Model(&chat_models.ChatModel{}).Where("conversation_id = ?", conversationID).Update("is_deleted", true).Error
	if err != nil {
		return err
	}

	return nil
}
