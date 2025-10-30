package logic

import (
	"context"
	"time"

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/group/group_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupJoinRequestSyncLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 群组申请数据同步
func NewGroupJoinRequestSyncLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupJoinRequestSyncLogic {
	return &GroupJoinRequestSyncLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupJoinRequestSyncLogic) GroupJoinRequestSync(req *types.GroupJoinRequestSyncReq) (resp *types.GroupJoinRequestSyncRes, err error) {
	var groupJoinRequests []group_models.GroupJoinRequestModel

	// 设置默认限制
	limit := req.Limit
	if limit <= 0 {
		limit = 100
	}

	// 查询用户发起的所有群组申请记录
	err = l.svcCtx.DB.Where("applicant_user_id = ? AND version > ? AND version <= ?",
		req.UserID, req.FromVersion, req.ToVersion).
		Order("version ASC").
		Limit(limit + 1).
		Find(&groupJoinRequests).Error
	if err != nil {
		l.Errorf("查询群组申请数据失败: %v", err)
		return nil, err
	}

	// 判断是否还有更多数据
	hasMore := len(groupJoinRequests) > limit
	if hasMore {
		groupJoinRequests = groupJoinRequests[:limit]
	}

	// 转换为响应格式
	var groupJoinRequestItems []types.GroupJoinRequestSyncItem
	var nextVersion int64 = req.FromVersion

	for _, request := range groupJoinRequests {
		handledAt := int64(0)
		if request.HandledAt != nil {
			handledAt = request.HandledAt.Unix()
		}

		groupJoinRequestItems = append(groupJoinRequestItems, types.GroupJoinRequestSyncItem{
			GroupID:         request.GroupID,
			ApplicantUserID: request.ApplicantUserID,
			Message:         request.Message,
			Status:          request.Status,
			HandledBy:       request.HandledBy,
			HandledAt:       handledAt,
			Version:         request.Version,
			CreateAt:        time.Time(request.CreatedAt).Unix(),
			UpdateAt:        time.Time(request.UpdatedAt).Unix(),
		})

		nextVersion = request.Version
	}

	// 如果没有更多数据，nextVersion应该是toVersion+1
	if !hasMore {
		nextVersion = req.ToVersion + 1
	}

	resp = &types.GroupJoinRequestSyncRes{
		GroupJoinRequests: groupJoinRequestItems,
		HasMore:           hasMore,
		NextVersion:       nextVersion,
	}

	l.Infof("群组申请数据同步完成，用户ID: %s, 返回群组申请记录数: %d, 还有更多: %v", req.UserID, len(groupJoinRequestItems), hasMore)
	return resp, nil
}
