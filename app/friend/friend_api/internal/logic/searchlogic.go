package logic

import (
	"context"
	"errors"
	"fmt"

	"beaver/app/friend/friend_api/internal/svc"
	"beaver/app/friend/friend_api/internal/types"
	"beaver/app/friend/friend_models"
	"beaver/app/user/user_models"
	"beaver/utils/conversation"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchLogic {
	return &SearchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchLogic) Search(req *types.SearchReq) (resp *types.SearchRes, err error) {
	// todo: add your logic here and delete this line
	var user user_models.UserModel

	// 根据手机号查询用户信息
	err = l.svcCtx.DB.Take(&user, "phone = ?", req.Phone).Error
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 获取好友关系
	var friend friend_models.FriendModel
	isFriend := friend.IsFriend(l.svcCtx.DB, req.UserId, user.UserId)

	// 生成会话Id
	conversationId, err := conversation.GenerateConversation([]string{req.UserId, user.UserId})
	if err != nil {
		return nil, fmt.Errorf("生成会话Id失败: %v", err)
	}

	// 构造返回值
	resp = &types.SearchRes{
		UserId:         user.UserId,
		Nickname:       user.NickName,
		Avatar:         user.Avatar,
		Abstract:       user.Abstract,
		IsFriend:       isFriend,
		ConversationId: conversationId,
	}

	return resp, nil

}
