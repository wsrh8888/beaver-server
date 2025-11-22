package logic

import (
	"beaver/app/datasync/datasync_api/internal/svc"
	"beaver/app/datasync/datasync_api/internal/types"
	"beaver/app/group/group_rpc/types/group_rpc"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSyncGroupInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取所有需要更新的群组信息版本
func NewGetSyncGroupInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSyncGroupInfoLogic {
	return &GetSyncGroupInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSyncGroupInfoLogic) GetSyncGroupInfo(req *types.GetSyncGroupInfoReq) (resp *types.GetSyncGroupInfoRes, err error) {
	userId := req.UserID
	if userId == "" {
		l.Errorf("用户ID为空")
		return nil, errors.New("用户ID不能为空")
	}

	// 1. 获取用户加入的群组ID列表
	groupIDsResp, err := l.svcCtx.GroupRpc.GetUserGroupIDs(l.ctx, &group_rpc.GetUserGroupIDsReq{
		UserID: userId,
	})
	if err != nil {
		l.Errorf("获取用户群组ID列表失败: %v", err)
		return nil, err
	}

	groupIDs := groupIDsResp.GroupIDs
	if len(groupIDs) == 0 {
		return &types.GetSyncGroupInfoRes{
			GroupVersions:   []types.GroupInfoVersionItem{},
			ServerTimestamp: time.Now().UnixMilli(),
		}, nil
	}

	fmt.Println("111111111111111111")
	fmt.Println(groupIDs)
	// 2. 获取变更的群组资料
	serverTimestamp := time.Now().UnixMilli()

	groupResp, err := l.svcCtx.GroupRpc.GetGroupsListByIds(l.ctx, &group_rpc.GetGroupsListByIdsReq{
		GroupIDs: groupIDs,
		Since:    req.Since,
	})
	if err != nil {
		l.Errorf("获取变更的群组资料失败: %v", err)
		return nil, err
	}

	// 3. 转换为响应格式，确保返回空数组而不是null
	groupVersions := make([]types.GroupInfoVersionItem, 0)
	if groupResp.Groups != nil {
		for _, group := range groupResp.Groups {
			groupVersions = append(groupVersions, types.GroupInfoVersionItem{
				GroupID: group.GroupID,
				Version: group.Version,
			})
		}
	}

	return &types.GetSyncGroupInfoRes{
		GroupVersions:   groupVersions,
		ServerTimestamp: serverTimestamp,
	}, nil
}
