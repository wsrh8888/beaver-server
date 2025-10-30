package logic

import (
	"context"
	"time"

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/group/group_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupJoinRequestListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取群组申请列表
func NewGroupJoinRequestListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupJoinRequestListLogic {
	return &GroupJoinRequestListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupJoinRequestListLogic) GroupJoinRequestList(req *types.GroupJoinRequestListReq) (resp *types.GroupJoinRequestListRes, err error) {
	var requests []group_models.GroupJoinRequestModel

	// 设置默认分页参数
	page := req.Page
	if page <= 0 {
		page = 1
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}
	offset := (page - 1) * limit

	// 查询群组申请列表
	err = l.svcCtx.DB.Where("group_id = ?", req.GroupID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&requests).Error
	if err != nil {
		l.Errorf("查询群组申请列表失败: %v", err)
		return nil, err
	}

	// 获取总数
	var total int64
	err = l.svcCtx.DB.Model(&group_models.GroupJoinRequestModel{}).
		Where("group_id = ?", req.GroupID).
		Count(&total).Error
	if err != nil {
		l.Errorf("获取群组申请总数失败: %v", err)
		return nil, err
	}

	// 转换为响应格式
	var requestItems []types.GroupJoinRequestItem

	for _, request := range requests {
		// 这里需要查询用户信息，但由于没有用户RPC，暂时使用默认值
		// 在实际项目中，应该通过用户RPC获取用户昵称和头像
		requestItems = append(requestItems, types.GroupJoinRequestItem{
			RequestID:       request.Id,
			GroupID:         request.GroupID,
			ApplicantID:     request.ApplicantUserID,
			ApplicantName:   "用户" + request.ApplicantUserID, // 临时值，需要从用户服务获取
			ApplicantAvatar: "",                             // 临时值，需要从用户服务获取
			Message:         request.Message,
			Status:          request.Status,
			CreateAt:        time.Time(request.CreatedAt).Unix(),
		})
	}

	resp = &types.GroupJoinRequestListRes{
		List:  requestItems,
		Count: total,
	}

	l.Infof("获取群组申请列表完成，群组ID: %s, 用户ID: %s, 返回申请数: %d", req.GroupID, req.UserID, len(requestItems))
	return resp, nil
}
