package logic

import (
	"context"
	"errors"
	"fmt"
	"time"

	"beaver/app/friend/friend_rpc/types/friend_rpc"
	"beaver/app/user/user_api/internal/svc"
	"beaver/app/user/user_api/internal/types"
	"beaver/app/user/user_models"
	"beaver/common/ajax"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UpdateInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateInfoLogic {
	return &UpdateInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateInfoLogic) UpdateInfo(req *types.UpdateInfoReq) (resp *types.UpdateInfoRes, err error) {
	// 获取要更新的用户信息
	var user user_models.UserModel
	if err := l.svcCtx.DB.Where("uuid = ?", req.UserID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	// 准备更新的字段
	updateFields := make(map[string]interface{})
	if req.Nickname != nil {
		updateFields["nick_name"] = *req.Nickname
	}
	if req.Avatar != nil {
		updateFields["avatar"] = *req.Avatar
	}
	if req.Abstract != nil {
		updateFields["abstract"] = *req.Abstract
	}
	if req.Gender != nil {
		updateFields["gender"] = *req.Gender
	}

	var version int64 // 定义version变量在更外层作用域

	// 执行更新操作
	if len(updateFields) > 0 {
		// 获取新版本号（用户独立递增）
		version = l.svcCtx.VersionGen.GetNextVersion("users", "uuid", req.UserID)
		if version == -1 {
			l.Errorf("获取版本号失败")
			return nil, errors.New("获取版本号失败")
		}

		// 添加版本号到更新字段
		updateFields["version"] = version

		err = l.svcCtx.DB.Model(&user).Updates(updateFields).Error
		if err != nil {
			return nil, err
		}

		l.Infof("用户信息更新成功: userID=%s, version=%d", req.UserID, version)

		// 记录用户变更日志
		l.recordUserChangeLog(req.UserID, version, updateFields)
	}

	// 异步更新缓存和通知好友（只有在有更新的情况下才推送）
	if version > 0 {
		go func() {
			// 创建新的context，避免使用请求的context
			ctx := context.Background()
			// 拿到自己的好友列表
			response, err := l.svcCtx.FriendRpc.GetFriendIds(ctx, &friend_rpc.GetFriendIdsRequest{
				UserID: req.UserID,
			})
			if err != nil {
				logx.Errorf("failed to get friend ids: %v", err)
				return
			}
			logx.Infof("转发给好友的列表: %v", response.FriendIds)
			// 通过ws推送给自己和好友
			allRecipients := append(response.FriendIds, req.UserID) // 包含自己
			for _, recipientID := range allRecipients {
				ajax.SendMessageToWs(l.svcCtx.Config.Etcd, wsCommandConst.USER_PROFILE, wsTypeConst.ProfileChangeNotify, req.UserID, recipientID, map[string]interface{}{
					"table":    "users",        // 涉及的数据库表
					"version":  int32(version), // 最新版本号（转换为int32类型）
					"targetId": req.UserID,     // 变更的记录ID
				}, "")
			}
		}()
	}

	return &types.UpdateInfoRes{}, nil
}

/**
 * 记录用户变更日志
 */
func (l *UpdateInfoLogic) recordUserChangeLog(userID string, version int64, updateFields map[string]interface{}) {
	var changeLogs []user_models.UserChangeLogModel

	// 为每个变更的字段创建日志记录
	for field, newValue := range updateFields {
		if field == "version" {
			continue // 跳过版本字段
		}

		var changeType string
		switch field {
		case "nick_name":
			changeType = "nickname"
		case "avatar":
			changeType = "avatar"
		case "abstract":
			changeType = "abstract"
		case "gender":
			changeType = "gender"
		case "status":
			changeType = "status"
		default:
			changeType = field
		}

		changeLog := user_models.UserChangeLogModel{
			UserID:     userID,
			ChangeType: changeType,
			NewValue:   fmt.Sprintf("%v", newValue),
			ChangeTime: time.Now().Unix(),
			Version:    version,
		}

		changeLogs = append(changeLogs, changeLog)
	}

	// 批量插入变更日志
	if len(changeLogs) > 0 {
		if err := l.svcCtx.DB.Create(&changeLogs).Error; err != nil {
			l.Errorf("记录用户变更日志失败: userID=%s, error=%v", userID, err)
		} else {
			l.Infof("用户变更日志记录成功: userID=%s, 变更数=%d", userID, len(changeLogs))
		}
	}
}
