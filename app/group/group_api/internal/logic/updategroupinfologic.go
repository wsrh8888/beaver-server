package logic

import (
	"context"
	"errors"

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/group/group_models"
	"beaver/app/group/group_rpc/types/group_rpc"
	mqwsconst "beaver/common/const/mqwsconst"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"
	"beaver/utils/logger"
	"beaver/utils/logger/model"

	"github.com/zeromicro/go-zero/core/logx"
)


type UpdateGroupInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logger.Logger
}

func NewUpdateGroupInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateGroupInfoLogic {
	return &UpdateGroupInfoLogic{
		ctx:    ctx,
		logger: logger.New("update_group_info"),
		svcCtx: svcCtx,
	}
}

func (l *UpdateGroupInfoLogic) UpdateGroupInfo(req *types.UpdateGroupInfoReq) (resp *types.UpdateGroupInfoRes, err error) {
	var newVersion int64 // 定义newVersion变量在更外层作用域

	// 检查操作者权限
	var member group_models.GroupMemberModel
	err = l.svcCtx.DB.Take(&member, "group_id = ? and user_id = ?", req.GroupID, req.UserID).Error
	if err != nil {
		return nil, errors.New("操作者不是群组成员")
	}
	if !(member.Role == 1 || member.Role == 2) {
		return nil, errors.New("没有权限更新群组信息")
	}

	// 构建更新字段
	updateFields := make(map[string]interface{})
	if req.Title != "" {
		updateFields["title"] = req.Title
	}
	if req.Avatar != "" {
		updateFields["avatar"] = req.Avatar
	}
	if req.Notice != "" {
		updateFields["notice"] = req.Notice
	}

	// 执行更新
	if len(updateFields) > 0 {
		// 获取该群的新版本号（独立递增）
		newVersion = l.svcCtx.VersionGen.GetNextVersion("groups", "group_id", req.GroupID)
		if newVersion == -1 {
			logx.WithContext(l.ctx).Errorf("获取群组版本号失败")
			return nil, errors.New("获取版本号失败")
		}

		// 添加版本号到更新字段
		updateFields["version"] = newVersion

		err = l.svcCtx.DB.Model(&group_models.GroupModel{}).
			Where("group_id = ?", req.GroupID).
			Updates(updateFields).Error
		if err != nil {
			logx.WithContext(l.ctx).Errorf("更新群组信息失败: %v", err)
			return nil, errors.New("更新群组信息失败")
		}
	}

	// 异步通知群成员（只有在有更新的情况下才推送）
	if newVersion > 0 {
		go func() {
			// 创建新的context，避免使用请求的context
			ctx := context.Background()

			// 获取群成员列表
			response, err := l.svcCtx.GroupRpc.GetGroupMembers(ctx, &group_rpc.GetGroupMembersReq{
				GroupID: req.GroupID,
			})
			if err != nil {
				logx.WithContext(l.ctx).Errorf("获取群成员列表失败: %v", err)
				return
			}

			// 通过 RocketMQ 异步推送给自己和群成员 - 群组信息同步
			allRecipients := append(response.Members, &group_rpc.GroupMemberInfo{UserID: req.UserID}) // 包含自己
			for _, member := range allRecipients {
				payload := map[string]interface{}{
					"command":  wsCommandConst.GROUP_OPERATION,
					"type":     wsTypeConst.GroupReceive,
					"senderId": req.UserID,
					"targetId": member.UserID,
					"body": map[string]interface{}{
						"table": "groups",
						"data": []map[string]interface{}{
							{
								"version": newVersion,
								"groupId": req.GroupID,
							},
						},
					},
					"conversationId": req.GroupID,
				}
				l.svcCtx.RocketMQ.SendMessage(ctx, mqwsconst.MqTopicWs, payload)
			}
		}()
	}

	if newVersion > 0 {
		l.logger.Info(model.LogMsg{
			Text: "群信息更新成功",
			Data: map[string]interface{}{
				"groupId": req.GroupID,
				"userId":  req.UserID,
			},
		})
	}

	return &types.UpdateGroupInfoRes{
		Version: newVersion,
	}, nil
}
