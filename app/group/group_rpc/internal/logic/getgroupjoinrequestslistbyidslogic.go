package logic

import (
	"context"
	"fmt"
	"time"

	"beaver/app/group/group_models"
	"beaver/app/group/group_rpc/internal/svc"
	"beaver/app/group/group_rpc/types/group_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupJoinRequestsListByIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupJoinRequestsListByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupJoinRequestsListByIdsLogic {
	return &GetGroupJoinRequestsListByIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetGroupJoinRequestsListByIdsLogic) GetGroupJoinRequestsListByIds(in *group_rpc.GetGroupJoinRequestsListByIdsReq) (*group_rpc.GetGroupJoinRequestsListByIdsRes, error) {
	if len(in.GroupIDs) == 0 {
		return &group_rpc.GetGroupJoinRequestsListByIdsRes{Requests: []*group_rpc.GroupJoinRequestListById{}}, nil
	}

	// 查询指定群组ID列表中，自指定时间戳以来变更的入群申请
	var changedRequests []group_models.GroupJoinRequestModel
	query := l.svcCtx.DB.Where("group_id IN (?)", in.GroupIDs)
	if in.Since > 0 {
		query = query.Where("updated_at > ?", in.Since)
	}

	err := query.Find(&changedRequests).Error
	if err != nil {
		l.Errorf("查询变更的入群申请失败: %v", err)
		return nil, err
	}

	// 转换为响应格式
	var requests []*group_rpc.GroupJoinRequestListById
	for _, request := range changedRequests {
		requests = append(requests, &group_rpc.GroupJoinRequestListById{
			RequestID: fmt.Sprintf("%d", request.Id), // 使用主键Id作为RequestID
			GroupID:   request.GroupID,
			UserID:    request.ApplicantUserID,
			Message:   request.Message,
			Status:    int32(request.Status),
			AppliedAt: time.Time(request.CreatedAt).UnixMilli(), // 使用CreatedAt作为申请时间
			Version:   request.Version,
		})
	}

	return &group_rpc.GetGroupJoinRequestsListByIdsRes{Requests: requests}, nil
}
