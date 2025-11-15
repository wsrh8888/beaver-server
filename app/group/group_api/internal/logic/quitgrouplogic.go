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

type QuitGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQuitGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QuitGroupLogic {
	return &QuitGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QuitGroupLogic) QuitGroup(req *types.GroupQuitReq) (resp *types.GroupQuitRes, err error) {
	// 检查用户是否为群组成员
	var member group_models.GroupMemberModel
	err = l.svcCtx.DB.Take(&member, "group_id = ? and user_id = ?", req.GroupID, req.UserID).Error
	if err != nil {
		return nil, errors.New("用户不是群组成员")
	}

	// 检查用户是否为群主
	if member.Role == 1 {
		// 群主退出前需要先转让群主权限
		return nil, errors.New("群主不能直接退出，请先转让群主权限")
	}

	// 获取该群成员的版本号（按群独立递增）
	memberVersion := l.svcCtx.VersionGen.GetNextVersion("group_members", "group_id", req.GroupID)
	if memberVersion == -1 {
		l.Logger.Errorf("获取群成员版本号失败")
		return nil, errors.New("获取版本号失败")
	}

	// 更新成员状态为退出（Status = 2）
	err = l.svcCtx.DB.Model(&member).
		Updates(map[string]interface{}{
			"status":  2, // 2退出
			"version": memberVersion,
		}).Error
	if err != nil {
		l.Logger.Errorf("退出群组失败: %v", err)
		return nil, errors.New("退出群组失败")
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
			if member.UserID != req.UserID { // 不通知操作者自己
				ajax.SendMessageToWs(l.svcCtx.Config.Etcd, wsCommandConst.GROUP_OPERATION, wsTypeConst.GroupMemberUpdate, req.GroupID, member.UserID, map[string]interface{}{
					"groupId":  req.GroupID,
					"type":     "leave",
					"userId":   req.UserID,
					"username": member.Username,
				}, "")
			}
		}
	}()

	return &types.GroupQuitRes{
		Version: memberVersion,
	}, nil
}
