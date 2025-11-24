package logic

import (
	"context"
	"errors"

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/group/group_models"
	"beaver/common/ajax"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupDeleteLogic {
	return &GroupDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupDeleteLogic) GroupDelete(req *types.GroupDeleteReq) (resp *types.GroupDeleteRes, err error) {
	// todo: add your logic here and delete this line
	var groupMember group_models.GroupMemberModel
	err = l.svcCtx.DB.Take(&groupMember, "group_id = ? and user_id = ?", req.GroupID, req.UserID).Error
	if err != nil {
		return nil, errors.New("用户不是群组成员")
	}
	if groupMember.Role != 1 {
		return nil, errors.New("只有群主才可以删除群组")
	}

	// 获取该群的版本号（独立递增）
	groupVersion := l.svcCtx.VersionGen.GetNextVersion("groups", "group_id", req.GroupID)
	if groupVersion == -1 {
		l.Logger.Errorf("获取群组版本号失败")
		return nil, errors.New("获取版本号失败")
	}

	// 获取群成员列表，用于推送通知
	var memberList []group_models.GroupMemberModel
	err = l.svcCtx.DB.Find(&memberList, "group_id = ?", req.GroupID).Error
	if err != nil {
		return nil, errors.New("获取群成员失败")
	}

	// 将群组状态改为解散（逻辑删除），并更新版本号
	err = l.svcCtx.DB.Model(&group_models.GroupModel{}).
		Where("group_id = ?", req.GroupID).
		Updates(map[string]interface{}{
			"status":  3, // 3=解散
			"version": groupVersion,
		}).Error
	if err != nil {
		return nil, errors.New("解散群组失败")
	}

	// 异步通知所有成员群组已被解散
	go func() {
		// 推送给所有成员 - 群组信息同步（标记为删除状态）
		for _, member := range memberList {
			ajax.SendMessageToWs(l.svcCtx.Config.Etcd, wsCommandConst.GROUP_OPERATION, wsTypeConst.GroupMemberReceive, req.UserID, member.UserID, map[string]interface{}{
				"tables": []map[string]interface{}{
					{
						"table": "groups",
						"data": []map[string]interface{}{
							{
								"version": groupVersion,
								"groupId": req.GroupID,
							},
						},
					},
				},
			}, "")
		}
	}()

	return &types.GroupDeleteRes{
		Version: groupVersion,
	}, nil
}
