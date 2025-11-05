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

// 获取用户管理的群组申请列表
func NewGroupJoinRequestListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupJoinRequestListLogic {
	return &GroupJoinRequestListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupJoinRequestListLogic) GroupJoinRequestList(req *types.GroupJoinRequestListReq) (resp *types.GroupJoinRequestListRes, err error) {
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

	// 先获取用户管理的群组ID列表（作为群主或管理员）
	var managedGroupIDs []string
	err = l.svcCtx.DB.Model(&group_models.GroupMemberModel{}).
		Where("user_id = ? AND role IN (1, 2)", req.UserID).
		Pluck("group_id", &managedGroupIDs).Error
	if err != nil {
		l.Errorf("获取用户管理的群组失败: %v", err)
		return nil, err
	}

	// 如果用户没有管理的群组，直接返回空结果
	if len(managedGroupIDs) == 0 {
		return &types.GroupJoinRequestListRes{
			List:  []types.GroupJoinRequestItem{},
			Count: 0,
		}, nil
	}

	var requests []group_models.GroupJoinRequestModel

	// 查询用户管理的所有群组的申请列表
	err = l.svcCtx.DB.Where("group_id IN (?)", managedGroupIDs).
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
		Where("group_id IN (?)", managedGroupIDs).
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
			Version:         request.Version,
		})
	}

	resp = &types.GroupJoinRequestListRes{
		List:  requestItems,
		Count: total,
	}

	l.Infof("获取群组申请列表完成，用户ID: %s, 返回申请数: %d", req.UserID, len(requestItems))
	return resp, nil
}
