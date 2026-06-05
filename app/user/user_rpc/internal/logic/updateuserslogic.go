package logic

import (
	"context"
	"errors"

	"beaver/app/user/user_models"
	"beaver/app/user/user_rpc/internal/svc"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

const (
	userActionPatch       int32 = 1 // 更新用户字段
	userActionSoftDelete  int32 = 2 // 软删除（status=3）
	userActionBatchStatus int32 = 3 // 批量修改状态
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
	switch in.Action {
	case userActionPatch:
		return l.patchUsers(in)
	case userActionSoftDelete:
		return l.softDelete(in.UserIds)
	case userActionBatchStatus:
		return l.batchStatus(in.UserIds, in.PatchStatus)
	default:
		return nil, errors.New("不支持的操作类型")
	}
}

func (l *UpdateUsersLogic) patchUsers(in *user_rpc.UpdateUsersReq) (*user_rpc.UpdateUsersRes, error) {
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
		if in.PatchStatus != nil {
			updates["status"] = int8(*in.PatchStatus)
		}
		if len(updates) == 0 {
			affected++
			continue
		}
		if err := l.svcCtx.DB.Model(&user).Updates(updates).Error; err != nil {
			l.Errorf("更新用户失败: %v", err)
			return nil, err
		}
		affected++
	}
	return &user_rpc.UpdateUsersRes{AffectedCount: affected}, nil
}

func (l *UpdateUsersLogic) softDelete(userIDs []string) (*user_rpc.UpdateUsersRes, error) {
	if len(userIDs) == 0 {
		return &user_rpc.UpdateUsersRes{}, nil
	}
	result := l.svcCtx.DB.Model(&user_models.UserModel{}).
		Where("user_id IN ?", userIDs).
		Update("status", 3)
	if result.Error != nil {
		l.Errorf("删除用户失败: %v", result.Error)
		return nil, result.Error
	}
	return &user_rpc.UpdateUsersRes{AffectedCount: result.RowsAffected}, nil
}

func (l *UpdateUsersLogic) batchStatus(userIDs []string, statusPtr *int32) (*user_rpc.UpdateUsersRes, error) {
	if statusPtr == nil {
		return &user_rpc.UpdateUsersRes{}, nil
	}
	status := *statusPtr
	if status < 1 || status > 3 {
		return nil, errors.New("无效的状态值")
	}
	if len(userIDs) == 0 {
		return &user_rpc.UpdateUsersRes{}, nil
	}
	result := l.svcCtx.DB.Model(&user_models.UserModel{}).
		Where("user_id IN ?", userIDs).
		Update("status", int8(status))
	if result.Error != nil {
		l.Errorf("批量更新用户状态失败: %v", result.Error)
		return nil, result.Error
	}
	return &user_rpc.UpdateUsersRes{AffectedCount: result.RowsAffected}, nil
}
