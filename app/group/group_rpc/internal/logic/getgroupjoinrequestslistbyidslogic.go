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
		l.Errorf("群组ID列表为空")
		return &group_rpc.GetGroupJoinRequestsListByIdsRes{Requests: []*group_rpc.GroupJoinRequestListById{}}, nil
	}

	// 查询指定群组ID列表中的入群申请
	var requestsData []group_models.GroupJoinRequestModel
	query := l.svcCtx.DB.Where("group_id IN (?)", in.GroupIDs)

	// 注意：Since在这里表示客户端已知的最新版本号，用于增量同步
	if in.Since > 0 {
		query = query.Where("version > ?", in.Since)
	}

	err := query.Find(&requestsData).Error
	if err != nil {
		l.Errorf("查询入群申请失败: groupIDs=%v, since=%d, error=%v", in.GroupIDs, in.Since, err)
		return nil, err
	}

	l.Infof("查询到 %d 个入群申请", len(requestsData))

	// 转换为响应格式
	var requests []*group_rpc.GroupJoinRequestListById
	for _, request := range requestsData {
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
