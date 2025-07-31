package logic

import (
	"context"
	"errors"
	"fmt"

	"beaver/app/friend/friend_api/internal/svc"
	"beaver/app/friend/friend_api/internal/types"
	"beaver/app/friend/friend_models"
	"beaver/app/user/user_models"
	"beaver/common/list_query"
	"beaver/common/models"
	"beaver/utils/conversation"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendListLogic {
	return &FriendListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendListLogic) FriendList(req *types.FriendListReq) (resp *types.FriendListRes, err error) {
	// 参数验证
	if req.UserID == "" {
		return nil, errors.New("用户ID不能为空")
	}

	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 20
	}

	// 查询好友列表
	friends, _, _ := list_query.ListQuery(l.svcCtx.DB, friend_models.FriendModel{}, list_query.Option{
		PageInfo: models.PageInfo{
			Page:  req.Page,
			Limit: req.Limit,
		},
		Where:   l.svcCtx.DB.Where("(send_user_id = ? OR rev_user_id = ?) AND is_deleted = ?", req.UserID, req.UserID, false),
		Preload: []string{"SendUserModel", "RevUserModel"},
	})

	var list []types.FriendInfoRes
	for _, friendUser := range friends {
		var info types.FriendInfoRes
		var targetUser user_models.UserModel
		var notice string

		// 确定目标用户和备注信息
		if friendUser.SendUserID == req.UserID {
			// 我是发起方，目标用户是接收方
			targetUser = friendUser.RevUserModel
			notice = friendUser.SendUserNotice
		} else if friendUser.RevUserID == req.UserID {
			// 我是接收方，目标用户是发起方
			targetUser = friendUser.SendUserModel
			notice = friendUser.RevUserNotice
		} else {
			// 这种情况理论上不应该发生，跳过
			continue
		}

		// 生成会话Id
		conversationID, err := conversation.GenerateConversation([]string{req.UserID, targetUser.UUID})
		if err != nil {
			l.Logger.Errorf("生成会话Id失败: userID=%s, targetID=%s, error=%v", req.UserID, targetUser.UUID, err)
			return nil, fmt.Errorf("生成会话Id失败: %v", err)
		}

		// 构造好友信息
		info = types.FriendInfoRes{
			UserID:         targetUser.UUID,
			Nickname:       targetUser.NickName,
			FileName:       targetUser.FileName,
			Abstract:       targetUser.Abstract,
			Notice:         notice,
			ConversationID: conversationID,
			Email:          targetUser.Email,
		}

		list = append(list, info)
	}

	l.Logger.Infof("获取好友列表成功: userID=%s, count=%d", req.UserID, len(list))
	return &types.FriendListRes{
		List: list,
	}, nil
}
