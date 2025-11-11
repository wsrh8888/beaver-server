package logic

import (
	"beaver/app/group/group_models"
	"beaver/app/group/group_rpc/internal/svc"
	"beaver/app/group/group_rpc/types/group_rpc"
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserGroupRequestVersionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserGroupRequestVersionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserGroupRequestVersionsLogic {
	return &GetUserGroupRequestVersionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取用户群组申请版本信息
func (l *GetUserGroupRequestVersionsLogic) GetUserGroupRequestVersions(in *group_rpc.GetUserGroupRequestVersionsReq) (*group_rpc.GetUserGroupRequestVersionsRes, error) {
	userID := in.UserID

	// 获取用户相关的所有群组申请版本
	versions, err := l.getUserGroupRequestVersions(userID, in.Since)
	if err != nil {
		l.Errorf("获取用户群组申请版本失败: %v", err)
		return nil, err
	}

	return &group_rpc.GetUserGroupRequestVersionsRes{
		Versions: versions,
	}, nil
}

// 获取用户相关的群组申请版本（包括我发出的和我管理的群）
func (l *GetUserGroupRequestVersionsLogic) getUserGroupRequestVersions(userID string, since int64) ([]*group_rpc.GroupRequestVersionItem, error) {
	// 1. 获取我申请加入的群的申请版本
	sentVersions, err := l.getSentRequestVersions(userID, since)
	if err != nil {
		return nil, err
	}

	// 2. 获取我管理的群收到的申请版本
	managedVersions, err := l.getManagedRequestVersions(userID, since)
	if err != nil {
		return nil, err
	}

	// 3. 合并版本信息（相同群组取最新版本）
	versionMap := make(map[string]int64)

	// 添加我发出的申请版本
	for _, version := range sentVersions {
		if currentVersion, exists := versionMap[version.GroupID]; !exists || version.Version > currentVersion {
			versionMap[version.GroupID] = version.Version
		}
	}

	// 添加我管理的申请版本
	for _, version := range managedVersions {
		if currentVersion, exists := versionMap[version.GroupID]; !exists || version.Version > currentVersion {
			versionMap[version.GroupID] = version.Version
		}
	}

	// 转换为响应格式
	var result []*group_rpc.GroupRequestVersionItem
	for groupID, version := range versionMap {
		result = append(result, &group_rpc.GroupRequestVersionItem{
			GroupID: groupID,
			Version: version,
		})
	}

	return result, nil
}

// 获取用户作为申请者发出的申请版本
func (l *GetUserGroupRequestVersionsLogic) getSentRequestVersions(userID string, since int64) ([]*group_rpc.GroupRequestVersionItem, error) {
	var requests []group_models.GroupJoinRequestModel
	query := l.svcCtx.DB.Where("applicant_user_id = ? AND status = 0", userID)
	if since > 0 {
		sinceTime := time.UnixMilli(since)
		query = query.Where("created_at > ?", sinceTime)
	}

	err := query.Find(&requests).Error
	if err != nil {
		return nil, err
	}

	// 按群组聚合版本
	versionMap := make(map[string]int64)
	for _, req := range requests {
		if currentVersion, exists := versionMap[req.GroupID]; !exists || req.Version > currentVersion {
			versionMap[req.GroupID] = req.Version
		}
	}

	var result []*group_rpc.GroupRequestVersionItem
	for groupID, version := range versionMap {
		result = append(result, &group_rpc.GroupRequestVersionItem{
			GroupID: groupID,
			Version: version,
		})
	}

	return result, nil
}

// 获取用户作为管理者收到的申请版本
func (l *GetUserGroupRequestVersionsLogic) getManagedRequestVersions(userID string, since int64) ([]*group_rpc.GroupRequestVersionItem, error) {
	// 1. 获取用户管理的群组
	var memberships []group_models.GroupMemberModel
	err := l.svcCtx.DB.Where("user_id = ? AND role IN (1, 2)", userID).Find(&memberships).Error
	if err != nil {
		return nil, err
	}

	if len(memberships) == 0 {
		return []*group_rpc.GroupRequestVersionItem{}, nil
	}

	var managedGroupIDs []string
	for _, membership := range memberships {
		managedGroupIDs = append(managedGroupIDs, membership.GroupID)
	}

	// 2. 获取这些群组的申请版本
	var requests []group_models.GroupJoinRequestModel
	query := l.svcCtx.DB.Where("group_id IN (?)", managedGroupIDs)
	if since > 0 {
		sinceTime := time.UnixMilli(since)
		query = query.Where("created_at > ?", sinceTime)
	}

	err = query.Find(&requests).Error
	if err != nil {
		return nil, err
	}

	// 按群组聚合版本
	versionMap := make(map[string]int64)
	for _, req := range requests {
		if currentVersion, exists := versionMap[req.GroupID]; !exists || req.Version > currentVersion {
			versionMap[req.GroupID] = req.Version
		}
	}

	var result []*group_rpc.GroupRequestVersionItem
	for groupID, version := range versionMap {
		result = append(result, &group_rpc.GroupRequestVersionItem{
			GroupID: groupID,
			Version: version,
		})
	}

	return result, nil
}
