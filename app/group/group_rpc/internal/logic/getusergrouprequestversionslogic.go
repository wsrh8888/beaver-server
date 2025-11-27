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

// 获取用户作为申请者发出的申请版本 - 优化查询性能
func (l *GetUserGroupRequestVersionsLogic) getSentRequestVersions(userID string, since int64) ([]*group_rpc.GroupRequestVersionItem, error) {
	var results []struct {
		GroupID string
		Version int64
	}

	query := l.svcCtx.DB.Model(&group_models.GroupJoinRequestModel{}).
		Select("group_id, MAX(version) as version").
		Where("applicant_user_id = ? AND status = 0", userID)

	if since > 0 {
		sinceTime := time.UnixMilli(since)
		query = query.Where("created_at > ?", sinceTime)
	}

	query = query.Group("group_id")

	err := query.Scan(&results).Error
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	var response []*group_rpc.GroupRequestVersionItem
	for _, item := range results {
		response = append(response, &group_rpc.GroupRequestVersionItem{
			GroupID: item.GroupID,
			Version: item.Version,
		})
	}

	return response, nil
}

// 获取用户作为管理者收到的申请版本 - 优化查询性能
func (l *GetUserGroupRequestVersionsLogic) getManagedRequestVersions(userID string, since int64) ([]*group_rpc.GroupRequestVersionItem, error) {
	var results []struct {
		GroupID string
		Version int64
	}

	// 使用JOIN查询避免IN子查询，提高性能
	query := l.svcCtx.DB.Table("group_join_request_models gjr").
		Select("gjr.group_id, MAX(gjr.version) as version").
		Joins("JOIN group_member_models gm ON gjr.group_id = gm.group_id").
		Where("gm.user_id = ? AND gm.role IN (1, 2)", userID)

	if since > 0 {
		sinceTime := time.UnixMilli(since)
		query = query.Where("gjr.created_at > ?", sinceTime)
	}

	query = query.Group("gjr.group_id")

	err := query.Scan(&results).Error
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	var response []*group_rpc.GroupRequestVersionItem
	for _, item := range results {
		response = append(response, &group_rpc.GroupRequestVersionItem{
			GroupID: item.GroupID,
			Version: item.Version,
		})
	}

	return response, nil
}
