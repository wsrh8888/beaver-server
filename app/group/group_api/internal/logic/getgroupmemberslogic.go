package logic

import (
	"context"
	"errors"

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/group/group_models"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupMembersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetGroupMembersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupMembersLogic {
	return &GetGroupMembersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetGroupMembersLogic) GetGroupMembers(req *types.GroupMemberListReq) (resp *types.GroupMemberListRes, err error) {
	// 参数验证和默认值设置
	if req.GroupID == "" {
		return nil, errors.New("群组ID不能为空")
	}

	if req.Page <= 0 {
		req.Page = 1
	}

	if req.Limit <= 0 {
		req.Limit = 20
	}

	offset := (req.Page - 1) * req.Limit

	// 查询群成员列表
	var members []group_models.GroupMemberModel
	var count int64

	// 先检查群组是否存在
	var groupCount int64
	err = l.svcCtx.DB.Model(&group_models.GroupModel{}).Where("uuid = ?", req.GroupID).Count(&groupCount).Error
	if err != nil {
		l.Logger.Errorf("查询群组失败: %v", err)
		return nil, errors.New("查询群组失败")
	}

	if groupCount == 0 {
		return nil, errors.New("群组不存在")
	}

	// 查询成员总数
	err = l.svcCtx.DB.Model(&group_models.GroupMemberModel{}).
		Where("group_id = ?", req.GroupID).
		Count(&count).Error

	if err != nil {
		l.Logger.Errorf("统计群成员数量失败: %v", err)
		return nil, errors.New("获取群成员列表失败")
	}

	// 分页查询成员
	err = l.svcCtx.DB.Where("group_id = ?", req.GroupID).
		Limit(req.Limit).
		Offset(offset).
		Find(&members).Error

	if err != nil {
		l.Logger.Errorf("获取群成员列表失败: %v", err)
		return nil, errors.New("获取群成员列表失败")
	}

	// 构建响应
	resp = &types.GroupMemberListRes{
		List:  make([]types.GroupMember, 0, len(members)),
		Count: count,
	}

	// 如果没有成员，直接返回空列表
	if len(members) == 0 {
		return resp, nil
	}

	// 收集所有用户ID
	userIDs := make([]string, 0, len(members))
	for _, member := range members {
		userIDs = append(userIDs, member.UserID)
	}

	// 通过UserRpc批量获取用户信息
	userResp, err := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{
		UserIdList: userIDs,
	})

	if err != nil {
		l.Logger.Errorf("获取用户信息失败: %v", err)
		return nil, errors.New("获取用户信息失败")
	}

	// 直接使用返回的用户信息映射
	userMap := userResp.UserInfo
	// 打印userMap的值
	logx.Infof("userMap: %v", userIDs)
	// 组装最终结果
	for _, member := range members {
		user, exists := userMap[member.UserID]

		groupMember := types.GroupMember{
			UserID:   member.UserID,
			Role:     member.Role,
			JoinTime: member.CreatedAt.String(),
		}

		if exists {
			groupMember.Nickname = user.NickName
			groupMember.Avatar = user.Avatar
		} else {
			groupMember.Nickname = "未知用户"
			groupMember.Avatar = ""
		}

		resp.List = append(resp.List, groupMember)
	}

	return resp, nil
}
