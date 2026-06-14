package logic

import (
	"context"
	"errors"
	"time"

	"beaver/app/group/group_models"
	"beaver/app/group/group_rpc/internal/svc"
	"beaver/app/group/group_rpc/types/group_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UpdateGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateGroupLogic {
	return &UpdateGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateGroupLogic) UpdateGroup(in *group_rpc.UpdateGroupReq) (*group_rpc.UpdateGroupRes, error) {
	var group group_models.GroupModel
	if err := l.svcCtx.DB.Where("id = ?", in.Id).First(&group).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("群组不存在")
		}
		l.Errorf("查询群组失败: %v", err)
		return nil, err
	}

	updates := map[string]interface{}{}
	if in.Title != "" {
		updates["title"] = in.Title
	}
	if in.Avatar != "" {
		updates["avatar"] = in.Avatar
	}
	if in.Notice != "" {
		updates["notice"] = in.Notice
	}
	if in.Status != 0 {
		updates["status"] = in.Status
		if in.Status == 3 {
			updates["deleted_at"] = time.Now()
		}
	}
	if in.MuteAll != nil {
		updates["is_mute_all"] = *in.MuteAll
		if *in.MuteAll {
			now := time.Now()
			updates["mute_all_at"] = now
		} else {
			updates["mute_all_at"] = nil
		}
	}
	if len(updates) == 0 {
		return &group_rpc.UpdateGroupRes{}, nil
	}

	groupVersion := l.svcCtx.VersionGen.GetNextVersion("groups", "group_id", group.GroupID)
	if groupVersion == -1 {
		return nil, errors.New("获取群组版本号失败")
	}
	updates["version"] = groupVersion

	if err := l.svcCtx.DB.Model(&group).Updates(updates).Error; err != nil {
		l.Errorf("更新群组失败: %v", err)
		return nil, err
	}
	return &group_rpc.UpdateGroupRes{}, nil
}
