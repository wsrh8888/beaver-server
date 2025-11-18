package logic

import (
	"context"
	"time"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/friend/friend_models"
	"beaver/app/user/user_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取好友关系列表
func NewGetFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendListLogic {
	return &GetFriendListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFriendListLogic) GetFriendList(req *types.GetFriendListReq) (resp *types.GetFriendListRes, err error) {
	// 设置默认分页参数
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// 构建查询条件
	query := l.svcCtx.DB.Model(&friend_models.FriendModel{})

	// 按用户ID筛选
	if req.UserID != "" {
		query = query.Where("send_user_id = ? OR rev_user_id = ?", req.UserID, req.UserID)
	}

	// 按好友ID筛选
	if req.FriendID != "" {
		query = query.Where("send_user_id = ? OR rev_user_id = ?", req.FriendID, req.FriendID)
	}

	// 按删除状态筛选
	if req.IsDeleted {
		query = query.Where("is_deleted = ?", req.IsDeleted)
	} else {
		query = query.Where("is_deleted = false")
	}

	// 时间范围筛选
	if req.StartTime != "" {
		query = query.Where("created_at >= ?", req.StartTime)
	}
	if req.EndTime != "" {
		query = query.Where("created_at <= ?", req.EndTime)
	}

	// 查询总数
	var total int64
	err = query.Count(&total).Error
	if err != nil {
		logx.Errorf("查询好友关系总数失败: %v", err)
		return nil, err
	}

	// 查询列表
	var friends []friend_models.FriendModel
	err = query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&friends).Error
	if err != nil {
		logx.Errorf("查询好友关系列表失败: %v", err)
		return nil, err
	}

	// 转换为响应格式
	list := make([]types.GetFriendListItem, len(friends))
	for i, friend := range friends {
		var sendUserName, revUserName string
		// 查询发送者信息
		if friend.SendUserID != "" {
			var sendUser user_models.UserModel
			if err := l.svcCtx.DB.Where("uuid = ?", friend.SendUserID).First(&sendUser).Error; err == nil {
				sendUserName = sendUser.NickName
			}
		}
		// 查询接收者信息
		if friend.RevUserID != "" {
			var revUser user_models.UserModel
			if err := l.svcCtx.DB.Where("uuid = ?", friend.RevUserID).First(&revUser).Error; err == nil {
				revUserName = revUser.NickName
			}
		}

		list[i] = types.GetFriendListItem{
			Id:             friend.UUID, // 使用UUID而不是数据库ID
			SendUserId:     friend.SendUserID,
			SendUserName:   sendUserName,
			RevUserId:      friend.RevUserID,
			RevUserName:    revUserName,
			SendUserNotice: friend.SendUserNotice,
			RevUserNotice:  friend.RevUserNotice,
			IsDeleted:      friend.IsDeleted,
			CreateTime:     time.Time(friend.CreatedAt).Format(time.RFC3339),
			UpdateTime:     time.Time(friend.UpdatedAt).Format(time.RFC3339),
		}
	}

	return &types.GetFriendListRes{
		List:  list,
		Total: total,
	}, nil
}
