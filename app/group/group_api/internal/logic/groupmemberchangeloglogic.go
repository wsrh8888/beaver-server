package logic

import (
	"context"

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/group/group_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupMemberChangeLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取群成员变更日志
func NewGroupMemberChangeLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupMemberChangeLogLogic {
	return &GroupMemberChangeLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupMemberChangeLogLogic) GroupMemberChangeLog(req *types.GroupMemberChangeLogReq) (resp *types.GroupMemberChangeLogRes, err error) {
	var changeLogs []group_models.GroupMemberChangeLogModel

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

	// 查询群成员变更日志
	err = l.svcCtx.DB.Where("group_id = ?", req.GroupID).
		Order("change_time DESC").
		Offset(offset).
		Limit(limit).
		Find(&changeLogs).Error
	if err != nil {
		l.Errorf("查询群成员变更日志失败: %v", err)
		return nil, err
	}

	// 获取总数
	var total int64
	err = l.svcCtx.DB.Model(&group_models.GroupMemberChangeLogModel{}).
		Where("group_id = ?", req.GroupID).
		Count(&total).Error
	if err != nil {
		l.Errorf("获取群成员变更日志总数失败: %v", err)
		return nil, err
	}

	// 转换为响应格式
	var logItems []types.GroupMemberChangeLogItem

	for _, log := range changeLogs {
		// 这里需要查询用户信息，但由于没有用户RPC，暂时使用默认值
		// 在实际项目中，应该通过用户RPC获取用户昵称
		logItems = append(logItems, types.GroupMemberChangeLogItem{
			UserID:       log.UserID,
			UserName:     "用户" + log.UserID, // 临时值，需要从用户服务获取
			ChangeType:   log.ChangeType,
			OperatedBy:   log.OperatedBy,
			OperatorName: "用户" + log.OperatedBy, // 临时值，需要从用户服务获取
			ChangeTime:   log.ChangeTime.Unix(),
		})
	}

	resp = &types.GroupMemberChangeLogRes{
		List:  logItems,
		Count: total,
	}

	l.Infof("获取群成员变更日志完成，群组ID: %s, 用户ID: %s, 返回日志数: %d", req.GroupID, req.UserID, len(logItems))
	return resp, nil
}
