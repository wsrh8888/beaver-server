package logic

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/friend/friend_api/internal/svc"
	"beaver/app/friend/friend_api/internal/types"
	"beaver/app/friend/friend_models"
	"beaver/common/ajax"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"

	"github.com/google/uuid"
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
	if req.UserID == "" || req.VerifyID == "" {
		return nil, errors.New("用户ID和验证ID不能为空")
	}

	// 状态值验证
	if req.Status < 1 || req.Status > 4 {
		return nil, errors.New("无效的状态值")
	}

	var friendVerify friend_models.FriendVerifyModel
	var conversationID string
	var friendNextVersion int64
	var friendUUID string

	// 查询好友验证记录，确保当前用户是接收方
	err = l.svcCtx.DB.Take(&friendVerify, "uuid = ? and rev_user_id = ?", req.VerifyID, req.UserID).Error
	if err != nil {
		l.Logger.Errorf("好友验证记录不存在: verifyID=%s, userID=%s, error=%v", req.VerifyID, req.UserID, err)
		return nil, errors.New("好友验证不存在")
	}

	// 检查验证状态是否已处理
	if friendVerify.RevStatus != 0 {
		l.Logger.Errorf("好友验证已处理: verifyID=%s, currentStatus=%d", req.VerifyID, friendVerify.RevStatus)
		return nil, errors.New("该验证已处理，无法重复操作")
	}

	// 处理不同的状态
	switch req.Status {
	case 1: // 同意
		friendVerify.RevStatus = 1

		// 获取下一个版本号
		friendNextVersion = l.svcCtx.VersionGen.GetNextVersion("friends", "", "")
		if friendNextVersion == -1 {
			l.Logger.Errorf("获取好友版本号失败")
			return nil, errors.New("系统错误")
		}

		// 生成UUID
		friendUUID = uuid.New().String()

		// 创建好友关系，同步来源信息
		err = l.svcCtx.DB.Create(&friend_models.FriendModel{
			UUID:       friendUUID, // 使用预生成的UUID
			SendUserID: friendVerify.SendUserID,
			RevUserID:  friendVerify.RevUserID,
			Source:     friendVerify.Source, // 同步来源字段
			Version:    friendNextVersion,   // 设置初始版本号
		}).Error
		if err != nil {
			l.Logger.Errorf("创建好友关系失败: %v", err)
			return nil, errors.New("创建好友关系失败")
		}

		// 生成私聊会话ID（微信风格：排序后拼接）
		userIds := []string{friendVerify.SendUserID, friendVerify.RevUserID}
		sort.Strings(userIds) // 确保ID顺序一致
		conversationId := fmt.Sprintf("private_%s_%s", userIds[0], userIds[1])

		// 调用Chat服务初始化私聊会话
		initResp, err := l.svcCtx.ChatRpc.InitializeConversation(context.Background(), &chat_rpc.InitializeConversationReq{
			ConversationId: conversationId,
			Type:           1, // 私聊
			UserIds:        []string{friendVerify.SendUserID, friendVerify.RevUserID},
		})
		if err != nil {
			l.Logger.Errorf("初始化私聊会话失败: %v", err)
			return nil, fmt.Errorf("初始化私聊会话失败: %v", err)
		}
		conversationID = initResp.ConversationId

		// 异步发送通知欢迎消息（通过专门的通知消息服务）
		go func() {
			defer func() {
				if r := recover(); r != nil {
					l.Logger.Errorf("异步发送欢迎消息时发生panic: %v", r)
				}
			}()

			// 调用Chat服务的通知消息发送接口
			_, err := l.svcCtx.ChatRpc.SendNotificationMessage(context.Background(), &chat_rpc.SendNotificationMessageReq{
				ConversationId: conversationID,
				MessageType:    1, // 好友添加成功欢迎消息
				Content:        "我们已经是好友了，开始聊天吧",
				RelatedUserId:  friendVerify.SendUserID, // 相关好友ID
			})
			if err != nil {
				l.Logger.Errorf("异步发送欢迎消息失败: %v", err)
			} else {
				l.Logger.Infof("异步发送欢迎消息成功: conversationID=%s", conversationID)
			}
		}()

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

		l.Logger.Infof("删除好友验证记录成功: verifyID=%s, userID=%s", req.VerifyID, req.UserID)
		return &types.FriendValidStatusRes{}, nil
	}

	// 获取下一个版本号并更新version字段
	nextVersion := l.svcCtx.VersionGen.GetNextVersion("friend_verify", "", "")
	if nextVersion == -1 {
		l.Logger.Errorf("获取版本号失败")
		return nil, errors.New("系统错误")
	}
	friendVerify.Version = nextVersion

	// 保存验证状态
	err = l.svcCtx.DB.Save(&friendVerify).Error
	if err != nil {
		l.Logger.Errorf("保存验证状态失败: %v", err)
		return nil, errors.New("保存验证状态失败")
	}

	// 异步发送WebSocket通知
	go func() {
		defer func() {
			if r := recover(); r != nil {
				l.Logger.Errorf("异步发送WebSocket消息时发生panic: %v", r)
			}
		}()

		// 构建表更新数据 - 包含版本号和UUID，让前端知道具体同步哪些数据
		if friendVerify.RevStatus == 1 { // 同意添加好友
			// 通知双方好友关系已建立 - 发送friends表版本更新
			friendUpdates := map[string]interface{}{
				"table": "friends",
				"data": []map[string]interface{}{
					{
						"version": friendNextVersion,
						"uuid":    friendUUID, // 使用预生成的UUID
					},
				},
			}

			ajax.SendMessageToWs(l.svcCtx.Config.Etcd, wsCommandConst.FRIEND_OPERATION, wsTypeConst.FriendReceive, friendVerify.SendUserID, friendVerify.RevUserID, map[string]interface{}{
				"tableUpdates": []map[string]interface{}{friendUpdates},
			}, conversationID)
			ajax.SendMessageToWs(l.svcCtx.Config.Etcd, wsCommandConst.FRIEND_OPERATION, wsTypeConst.FriendReceive, friendVerify.RevUserID, friendVerify.SendUserID, map[string]interface{}{
				"tableUpdates": []map[string]interface{}{friendUpdates},
			}, conversationID)
		} else {
			// 其他状态（拒绝、忽略）发送验证表版本更新
			verifyUpdates := map[string]interface{}{
				"table": "friend_verify",
				"data": []map[string]interface{}{
					{
						"version": nextVersion,
						"uuid":    friendVerify.UUID,
					},
				},
			}

			ajax.SendMessageToWs(l.svcCtx.Config.Etcd, wsCommandConst.FRIEND_OPERATION, wsTypeConst.FriendVerifyReceive, friendVerify.SendUserID, friendVerify.RevUserID, map[string]interface{}{
				"tableUpdates": []map[string]interface{}{verifyUpdates},
			}, conversationID)
			ajax.SendMessageToWs(l.svcCtx.Config.Etcd, wsCommandConst.FRIEND_OPERATION, wsTypeConst.FriendVerifyReceive, friendVerify.RevUserID, friendVerify.SendUserID, map[string]interface{}{
				"tableUpdates": []map[string]interface{}{verifyUpdates},
			}, conversationID)
		}

		l.Logger.Infof("异步发送WebSocket通知完成: verifyID=%s, status=%d", req.VerifyID, friendVerify.RevStatus)
	}()

	l.Logger.Infof("处理好友验证成功: verifyID=%s, userID=%s, status=%d, source=%s", req.VerifyID, req.UserID, req.Status, friendVerify.Source)
	return &types.FriendValidStatusRes{
		Version: nextVersion,
	}, nil
}
