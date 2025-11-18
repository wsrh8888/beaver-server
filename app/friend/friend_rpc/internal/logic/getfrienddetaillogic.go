package logic

import (
	"context"
	"time"

	"beaver/app/friend/friend_models"
	"beaver/app/friend/friend_rpc/internal/svc"
	"beaver/app/friend/friend_rpc/types/friend_rpc"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFriendDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendDetailLogic {
	return &GetFriendDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFriendDetailLogic) GetFriendDetail(in *friend_rpc.GetFriendDetailReq) (*friend_rpc.GetFriendDetailRes, error) {
	// 查询当前用户与指定好友的好友关系
	var friends []friend_models.FriendModel
	err := l.svcCtx.DB.Where("((send_user_id = ? AND rev_user_id IN (?)) OR (rev_user_id = ? AND send_user_id IN (?))) AND is_deleted = ?",
		in.UserId, in.FriendIds, in.UserId, in.FriendIds, false).Find(&friends).Error
	if err != nil {
		logx.Errorf("failed to query friends: %v", err)
		return nil, err
	}

	// 提取好友ID列表
	var friendIds []string
	for _, friend := range friends {
		if friend.SendUserID == in.UserId {
			friendIds = append(friendIds, friend.RevUserID)
		} else {
			friendIds = append(friendIds, friend.SendUserID)
		}
	}

	// 通过UserRpc批量获取用户信息
	var userInfoMap map[string]*user_rpc.UserInfo
	if len(friendIds) > 0 {
		userListRes, err := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{
			UserIdList: friendIds,
		})
		if err != nil {
			logx.Errorf("failed to get user info from UserRpc: %v", err)
			return nil, err
		}
		userInfoMap = userListRes.UserInfo
	}

	// 构建好友ID到好友关系的映射
	friendMap := make(map[string]friend_models.FriendModel)
	for _, friend := range friends {
		if friend.SendUserID == in.UserId {
			// 当前用户是发起方，对方是接收方
			friendMap[friend.RevUserID] = friend
		} else {
			// 当前用户是接收方，对方是发起方
			friendMap[friend.SendUserID] = friend
		}
	}

	// 构建响应
	var friendDetails []*friend_rpc.FriendDetailItem
	for _, friendId := range friendIds {
		userInfo, userExists := userInfoMap[friendId]
		friend, friendExists := friendMap[friendId]

		if !userExists || !friendExists {
			continue
		}

		// 获取备注信息
		var notice string
		if friend.SendUserID == in.UserId {
			// 当前用户是发起方，使用接收方的备注
			notice = friend.RevUserNotice
		} else {
			// 当前用户是接收方，使用发起方的备注
			notice = friend.SendUserNotice
		}

		// 获取成为好友的时间（取最早的创建时间）
		friendAt := time.Time(friend.CreatedAt).Unix()

		friendDetails = append(friendDetails, &friend_rpc.FriendDetailItem{
			UserId:   friendId, // userInfoMap的key就是用户ID
			Nickname: userInfo.NickName,
			Avatar:   userInfo.Avatar,
			Notice:   notice,
			FriendAt: friendAt,
		})
	}

	return &friend_rpc.GetFriendDetailRes{
		Friends: friendDetails,
	}, nil
}
