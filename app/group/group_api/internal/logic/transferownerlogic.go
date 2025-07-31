package logic

import (
	"context"
	"errors"

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/group/group_models"
	"beaver/app/group/group_rpc/types/group_rpc"
	"beaver/common/ajax"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"

	"github.com/zeromicro/go-zero/core/logx"
)

type TransferOwnerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTransferOwnerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TransferOwnerLogic {
	return &TransferOwnerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TransferOwnerLogic) TransferOwner(req *types.TransferOwnerReq) (resp *types.TransferOwnerRes, err error) {
	// 检查当前用户是否为群主
	var currentOwner group_models.GroupMemberModel
	err = l.svcCtx.DB.Take(&currentOwner, "group_id = ? and user_id = ?", req.GroupID, req.UserID).Error
	if err != nil {
		return nil, errors.New("用户不是群组成员")
	}
	if currentOwner.Role != 1 {
		return nil, errors.New("只有群主可以转让群组")
	}

	// 检查新群主是否为群组成员
	var newOwner group_models.GroupMemberModel
	err = l.svcCtx.DB.Take(&newOwner, "group_id = ? and user_id = ?", req.GroupID, req.NewOwnerID).Error
	if err != nil {
		return nil, errors.New("新群主不是群组成员")
	}

	// 开始事务
	tx := l.svcCtx.DB.Begin()
	if tx.Error != nil {
		l.Logger.Errorf("开启事务失败: %v", tx.Error)
		return nil, errors.New("转让群组失败")
	}

	// 更新原群主角色为普通成员
	err = tx.Model(&group_models.GroupMemberModel{}).
		Where("group_id = ? and user_id = ?", req.GroupID, req.UserID).
		Update("role", 3).Error
	if err != nil {
		tx.Rollback()
		l.Logger.Errorf("更新原群主角色失败: %v", err)
		return nil, errors.New("转让群组失败")
	}

	// 更新新群主角色
	err = tx.Model(&group_models.GroupMemberModel{}).
		Where("group_id = ? and user_id = ?", req.GroupID, req.NewOwnerID).
		Update("role", 1).Error
	if err != nil {
		tx.Rollback()
		l.Logger.Errorf("更新新群主角色失败: %v", err)
		return nil, errors.New("转让群组失败")
	}

	// 提交事务
	if err = tx.Commit().Error; err != nil {
		l.Logger.Errorf("提交事务失败: %v", err)
		return nil, errors.New("转让群组失败")
	}

	// 异步通知群成员
	go func() {
		// 创建新的context，避免使用请求的context
		ctx := context.Background()

		// 获取群成员列表
		response, err := l.svcCtx.GroupRpc.GetGroupMembers(ctx, &group_rpc.GetGroupMembersReq{
			GroupID: req.GroupID,
		})
		if err != nil {
			l.Logger.Errorf("获取群成员列表失败: %v", err)
			return
		}

		// 通过ws推送给群成员
		for _, member := range response.Members {
			ajax.SendMessageToWs(l.svcCtx.Config.Etcd, wsCommandConst.GROUP_OPERATION, wsTypeConst.GroupUpdate, req.GroupID, member.UserID, map[string]interface{}{
				"groupId":    req.GroupID,
				"type":       "owner_transfer",
				"oldOwnerId": req.UserID,
				"newOwnerId": req.NewOwnerID,
			}, "")
		}
	}()

	return &types.TransferOwnerRes{}, nil
}
