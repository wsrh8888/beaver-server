package logic

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/friend/friend_models"
	"beaver/app/user/user_models"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetFriendDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取好友关系详情
func NewGetFriendDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendDetailLogic {
	return &GetFriendDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFriendDetailLogic) GetFriendDetail(req *types.GetFriendDetailReq) (resp *types.GetFriendDetailRes, err error) {
	// 转换ID
	friendID, err := strconv.ParseUint(req.FriendID, 10, 32)
	if err != nil {
		logx.Errorf("无效的好友关系ID: %s", req.FriendID)
		return nil, errors.New("无效的好友关系ID")
	}

	var friend friend_models.FriendModel
	err = l.svcCtx.DB.Where("id = ?", uint(friendID)).
		First(&friend).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logx.Errorf("好友关系不存在, Id: %s", req.FriendID)
			return nil, errors.New("好友关系不存在")
		}
		logx.Errorf("查询好友关系详情失败: %v", err)
		return nil, err
	}

	var sendUserName, sendUserFileName, revUserName, revUserFileName string
	// 查询发送者信息
	if friend.SendUserID != "" {
		var sendUser user_models.UserModel
		if err := l.svcCtx.DB.Where("uuid = ?", friend.SendUserID).First(&sendUser).Error; err == nil {
			sendUserName = sendUser.NickName
			sendUserFileName = sendUser.Avatar
		}
	}
	// 查询接收者信息
	if friend.RevUserID != "" {
		var revUser user_models.UserModel
		if err := l.svcCtx.DB.Where("uuid = ?", friend.RevUserID).First(&revUser).Error; err == nil {
			revUserName = revUser.NickName
			revUserFileName = revUser.Avatar
		}
	}

	return &types.GetFriendDetailRes{
		Id:               fmt.Sprintf("%d", friend.Id),
		SendUserId:       friend.SendUserID,
		SendUserName:     sendUserName,
		SendUserFileName: sendUserFileName,
		RevUserId:        friend.RevUserID,
		RevUserName:      revUserName,
		RevUserFileName:  revUserFileName,
		SendUserNotice:   friend.SendUserNotice,
		RevUserNotice:    friend.RevUserNotice,
		IsDeleted:        friend.IsDeleted,
		CreateTime:       time.Time(friend.CreatedAt).Format(time.RFC3339),
		UpdateTime:       time.Time(friend.UpdatedAt).Format(time.RFC3339),
	}, nil
}
