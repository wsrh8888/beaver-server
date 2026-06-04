package logic

import (
	"context"
	"encoding/json"
	"errors"

	"beaver/app/friend/friend_api/internal/svc"
	"beaver/app/friend/friend_api/internal/types"
	"beaver/app/friend/friend_models"
	"beaver/app/open/openevent"
	"beaver/app/open/open_rpc/types/open_rpc"
	mqwsconst "beaver/common/const/mqwsconst"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteFriendLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除好友
func NewDeleteFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFriendLogic {
	return &DeleteFriendLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteFriendLogic) DeleteFriend(req *types.DeleteFriendReq) (resp *types.DeleteFriendRes, err error) {
	// 查找好友关系
	var friendModel friend_models.FriendModel
	if err := l.svcCtx.DB.Where(
		"((send_user_id = ? AND rev_user_id = ?) OR (send_user_id = ? AND rev_user_id = ?)) AND is_deleted = false",
		req.UserID, req.FriendID, req.FriendID, req.UserID,
	).First(&friendModel).Error; err != nil {
		return nil, errors.New("好友关系不存在")
	}

	// 获取下一个版本号并软删除
	nextVersion := l.svcCtx.VersionGen.GetNextVersion("friend", "", "")
	if nextVersion == -1 {
		return nil, errors.New("系统错误")
	}

	if err := l.svcCtx.DB.Model(&friendModel).Updates(map[string]interface{}{
		"is_deleted": true,
		"version":    nextVersion,
	}).Error; err != nil {
		return nil, errors.New("删除好友失败")
	}

	// 异步通知双方同步好友数据
	go func() {
		tableUpdates := []map[string]interface{}{
			{
				"table": "friends",
				"data": []map[string]interface{}{
					{"friendId": friendModel.FriendID, "version": nextVersion},
				},
			},
		}
		payload1 := map[string]interface{}{
			"command":        wsCommandConst.FRIEND_OPERATION,
			"type":           wsTypeConst.FriendReceive,
			"senderId":       req.UserID,
			"targetId":       req.FriendID,
			"body":           map[string]interface{}{"tableUpdates": tableUpdates},
			"conversationId": "",
		}
		l.svcCtx.RocketMQ.SendMessage(context.Background(), mqwsconst.MqTopicWs, payload1)

		payload2 := map[string]interface{}{
			"command":        wsCommandConst.FRIEND_OPERATION,
			"type":           wsTypeConst.FriendReceive,
			"senderId":       req.FriendID,
			"targetId":       req.UserID,
			"body":           map[string]interface{}{"tableUpdates": tableUpdates},
			"conversationId": "",
		}
		l.svcCtx.RocketMQ.SendMessage(context.Background(), mqwsconst.MqTopicWs, payload2)
	}()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				l.Logger.Errorf("Robot 好友事件推送 panic: %v", r)
			}
		}()
		ctx := context.Background()
		res, err := l.svcCtx.OpenRpc.GetRobotByUserID(ctx, &open_rpc.GetRobotByUserIDReq{RobotUserId: req.FriendID})
		if err == nil && res != nil && res.Found {
			body, _ := json.Marshal(map[string]interface{}{
				"robot_id": req.FriendID,
				"user_id":  req.UserID,
			})
			_, _ = l.svcCtx.OpenRpc.DispatchPlatformEvent(ctx, &open_rpc.DispatchPlatformEventReq{
				AppId:     res.AppId,
				EventType: openevent.EventIMBotUnfollowed,
				EventJson: string(body),
			})
		}
		res, err = l.svcCtx.OpenRpc.GetRobotByUserID(ctx, &open_rpc.GetRobotByUserIDReq{RobotUserId: req.UserID})
		if err == nil && res != nil && res.Found {
			body, _ := json.Marshal(map[string]interface{}{
				"robot_id": req.UserID,
				"user_id":  req.FriendID,
			})
			_, _ = l.svcCtx.OpenRpc.DispatchPlatformEvent(ctx, &open_rpc.DispatchPlatformEventReq{
				AppId:     res.AppId,
				EventType: openevent.EventIMBotUnfollowed,
				EventJson: string(body),
			})
		}
	}()

	return &types.DeleteFriendRes{}, nil
}
