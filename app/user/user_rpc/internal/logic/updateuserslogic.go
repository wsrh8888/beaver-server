package logic

import (
	"context"
	"errors"
	"fmt"
	"time"

	"beaver/app/user/user_models"
	"beaver/app/user/user_rpc/internal/svc"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type UpdateUsersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUsersLogic {
	return &UpdateUsersLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *UpdateUsersLogic) UpdateUsers(in *user_rpc.UpdateUsersReq) (*user_rpc.UpdateUsersRes, error) {
	var affected int64
	for _, uid := range in.UserIds {
		var user user_models.UserModel
		if err := l.svcCtx.DB.Where("user_id = ?", uid).First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				continue
			}
			return nil, err
		}

		if in.PatchEmail != nil && *in.PatchEmail != user.Email {
			var exist user_models.UserModel
			if err := l.svcCtx.DB.Where("email = ? AND user_id != ?", *in.PatchEmail, uid).First(&exist).Error; err == nil {
				return nil, status.Error(codes.AlreadyExists, "邮箱已存在")
			} else if !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}
		}

		updates := map[string]interface{}{}
		if in.PatchNickName != nil {
			updates["nick_name"] = *in.PatchNickName
		}
		if in.PatchEmail != nil {
			updates["email"] = *in.PatchEmail
		}
		if in.PatchAvatar != nil {
			updates["avatar"] = *in.PatchAvatar
		}
		if in.PatchAbstract != nil {
			updates["abstract"] = *in.PatchAbstract
		}
		if len(updates) == 0 {
			affected++
			continue
		}

		version := l.svcCtx.VersionGen.GetNextVersion("users", "user_id", uid)
		if version == -1 {
			return nil, errors.New("获取用户版本号失败")
		}
		updates["version"] = version

		if err := l.svcCtx.DB.Model(&user).Updates(updates).Error; err != nil {
			l.Errorf("更新用户失败: %v", err)
			return nil, err
		}
		l.recordUserChangeLog(uid, version, updates)
		affected++
	}
	return &user_rpc.UpdateUsersRes{AffectedCount: affected}, nil
}

func (l *UpdateUsersLogic) recordUserChangeLog(userID string, version int64, updateFields map[string]interface{}) {
	changeLogs := make([]user_models.UserChangeLogModel, 0, len(updateFields))
	now := time.Now().Unix()

	for field, newValue := range updateFields {
		if field == "version" {
			continue
		}
		changeType := field
		switch field {
		case "nick_name":
			changeType = "nickName"
		case "avatar":
			changeType = "avatar"
		case "abstract":
			changeType = "abstract"
		}
		changeLogs = append(changeLogs, user_models.UserChangeLogModel{
			UserID:     userID,
			ChangeType: changeType,
			NewValue:   fmt.Sprintf("%v", newValue),
			ChangeTime: now,
			Version:    version,
		})
	}

	if len(changeLogs) > 0 {
		if err := l.svcCtx.DB.Create(&changeLogs).Error; err != nil {
			l.Errorf("记录用户变更日志失败: userId=%s, error=%v", userID, err)
		}
	}
}
