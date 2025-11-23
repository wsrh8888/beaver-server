package logic

import (
	"context"
	"errors"

	"beaver/app/friend/friend_api/internal/svc"
	"beaver/app/friend/friend_api/internal/types"
	"beaver/app/friend/friend_models"
	"beaver/app/user/user_rpc/types/user_rpc"
	"beaver/common/ajax"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddFriendLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddFriendLogic {
	return &AddFriendLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddFriendLogic) AddFriend(req *types.AddFriendReq) (resp *types.AddFriendRes, err error) {
	var friend friend_models.FriendModel

	// 检查是否已经是好友
	if friend.IsFriend(l.svcCtx.DB, req.UserID, req.FriendID) {
		return nil, errors.New("已经是好友了")
	}

	// 检查目标用户是否存在（通过RPC）
	_, err = l.svcCtx.UserRpc.UserInfo(l.ctx, &user_rpc.UserInfoReq{
		UserID: req.FriendID,
	})
	if err != nil {
		l.Logger.Errorf("目标用户不存在: friendID=%s, error=%v", req.FriendID, err)
		return nil, errors.New("用户不存在")
	}

	// 检查是否已经有待处理的好友请求
	var existingVerify friend_models.FriendVerifyModel
	err = l.svcCtx.DB.Take(&existingVerify,
		"(send_user_id = ? AND rev_user_id = ? AND rev_status = 0) OR (send_user_id = ? AND rev_user_id = ? AND rev_status = 0)",
		req.UserID, req.FriendID, req.FriendID, req.UserID).Error

	if err == nil {
		l.Logger.Infof("已存在待处理的好友请求: userID=%s, friendID=%s", req.UserID, req.FriendID)
		return &types.AddFriendRes{}, nil
	}

	// 获取下一个版本号
	nextVersion := l.svcCtx.VersionGen.GetNextVersion("friend_verify", "", "")
	if nextVersion == -1 {
		l.Logger.Errorf("获取版本号失败")
		return nil, errors.New("系统错误")
	}

	// 创建好友验证请求
	verifyModel := friend_models.FriendVerifyModel{
		SendUserID: req.UserID,
		RevUserID:  req.FriendID,
		Message:    req.Verify,
		Source:     req.Source, // 添加来源字段
		Version:    nextVersion,
		UUID:       uuid.New().String(),
	}

	err = l.svcCtx.DB.Create(&verifyModel).Error
	if err != nil {
		l.Logger.Errorf("创建好友验证请求失败: %v", err)
		return nil, errors.New("添加好友请求失败")
	}

	// 异步发送WebSocket通知给发送方和接收方
	go func() {
		defer func() {
			if r := recover(); r != nil {
				l.Logger.Errorf("异步发送WebSocket消息时发生panic: %v", r)
			}
		}()

		// 获取发送者和接收者的用户信息
		senderInfo, senderErr := l.svcCtx.UserRpc.UserInfo(l.ctx, &user_rpc.UserInfoReq{
			UserID: req.UserID,
		})
		if senderErr != nil {
			l.Logger.Errorf("获取发送者用户信息失败: %v", senderErr)
		}

		receiverInfo, receiverErr := l.svcCtx.UserRpc.UserInfo(l.ctx, &user_rpc.UserInfoReq{
			UserID: req.FriendID,
		})
		if receiverErr != nil {
			l.Logger.Errorf("获取接收者用户信息失败: %v", receiverErr)
		}

		// 构建好友验证表更新数据
		verifyUpdates := map[string]interface{}{
			"table": "friend_verify",
			"data": []map[string]interface{}{
				{
					"version": nextVersion,
					"uuid":    verifyModel.UUID,
				},
			},
		}

		// 构建用户表更新数据数组
		var userUpdates []map[string]interface{}
		if senderInfo != nil {
			userUpdates = append(userUpdates, map[string]interface{}{
				"table": "users",
				"data": []map[string]interface{}{
					{
						"userId":  senderInfo.UserInfo.UserId,
						"version": senderInfo.UserInfo.Version,
					},
				},
			})
		}
		if receiverInfo != nil {
			userUpdates = append(userUpdates, map[string]interface{}{
				"table": "users",
				"data": []map[string]interface{}{
					{
						"userId":  senderInfo.UserInfo.UserId,
						"version": senderInfo.UserInfo.Version,
					},
				},
			})
		}

		// 合并所有表更新
		tableUpdates := append([]map[string]interface{}{verifyUpdates}, userUpdates...)

		// 通知接收方
		ajax.SendMessageToWs(l.svcCtx.Config.Etcd, wsCommandConst.FRIEND_OPERATION, wsTypeConst.FriendVerifyReceive, req.UserID, req.FriendID, map[string]interface{}{
			"tableUpdates": tableUpdates,
		}, "")

		// 通知发送方
		ajax.SendMessageToWs(l.svcCtx.Config.Etcd, wsCommandConst.FRIEND_OPERATION, wsTypeConst.FriendVerifyReceive, req.FriendID, req.UserID, map[string]interface{}{
			"tableUpdates": tableUpdates,
		}, "")

		l.Logger.Infof("异步发送好友验证请求通知完成: sender=%s, receiver=%s, version=%d, uuid=%s", req.UserID, req.FriendID, nextVersion, verifyModel.UUID)
	}()

	l.Logger.Infof("好友请求发送成功: userID=%s, friendID=%s, source=%s", req.UserID, req.FriendID, req.Source)
	return &types.AddFriendRes{
		Version: nextVersion,
	}, nil
}
