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

// 入群申请同步
func NewGroupJoinRequestSyncLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupJoinRequestSyncLogic {
	return &GroupJoinRequestSyncLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupJoinRequestSyncLogic) GroupJoinRequestSync(req *types.GroupJoinRequestSyncReq) (resp *types.GroupJoinRequestSyncRes, err error) {
	resp = &types.GroupJoinRequestSyncRes{
		GroupJoinRequests: []types.GroupJoinRequestSyncItem{},
	}

	if len(req.Groups) == 0 {
		l.Infof("入群申请同步完成，用户ID: %s, 无需同步的群组", req.UserID)
		return resp, nil
	}

	// 为每个群组查询版本变化的数据
	for _, groupReq := range req.Groups {
		var requests []group_models.GroupJoinRequestModel
		err = l.svcCtx.DB.Where("group_id = ? AND version >= ?", groupReq.GroupID, groupReq.Version).
			Find(&requests).Error
		if err != nil {
			l.Errorf("查询入群申请数据失败，群组ID: %s, 错误: %v", groupReq.GroupID, err)
			continue
		}

		for _, request := range requests {
			resp.GroupJoinRequests = append(resp.GroupJoinRequests, types.GroupJoinRequestSyncItem{
				GroupID:         request.GroupID,
				ApplicantUserID: request.ApplicantUserID,
				Message:         request.Message,
				Status:          request.Status,
				HandledBy:       request.HandledBy,
				HandledAt: func() int64 {
					if request.HandledAt != nil {
						return request.HandledAt.Unix()
					}
					return 0
				}(),
				Version:  request.Version,
				CreateAt: time.Time(request.CreatedAt).Unix(),
				UpdateAt: time.Time(request.UpdatedAt).Unix(),
			})
		}
	}

	l.Infof("入群申请同步完成，用户ID: %s, 返回申请变化数: %d", req.UserID, len(resp.GroupJoinRequests))
	return resp, nil
}
