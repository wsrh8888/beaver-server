package logic

import (
	"context"
	"errors"
	"fmt"
	"time"

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/group/group_models"
	"beaver/app/open/open_rpc/types/open_rpc"
	"beaver/app/user/user_models"
	"beaver/app/user/user_rpc/types/user_rpc"
	"beaver/utils/logger"
	"beaver/utils/logger/model"

	"gorm.io/gorm"
)


type CreateBotLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logger.Logger
}

// 在群里创建通知机器人（群管理员操作，返回 Webhook URL + Secret）
func NewCreateBotLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateBotLogic {
	return &CreateBotLogic{
		ctx:    ctx,
		logger: logger.New("create_bot"),
		svcCtx: svcCtx,
	}
}

func (l *CreateBotLogic) CreateBot(req *types.CreateBotReq) (resp *types.CreateBotRes, err error) {
	// 1. 校验权限
	var member group_models.GroupMemberModel
	if err = l.svcCtx.DB.Take(&member, "group_id = ? AND user_id = ?", req.GroupID, req.UserID).Error; err != nil {
		return nil, errors.New("不是群成员")
	}
	if member.Role != 1 && member.Role != 2 {
		return nil, errors.New("无权限，仅群主或管理员可创建通知机器人")
	}

	// 2. 通过 user_rpc 创建机器人用户
	userRes, userErr := l.svcCtx.UserRpc.UserCreate(l.ctx, &user_rpc.UserCreateReq{
		NickName: req.Name,
		Source:   int32(user_models.SourceGroup), // 群内创建（机器人）
		UserType: int32(user_models.UserTypeBot), // 标记为机器人
	})
	if userErr != nil {
		return nil, fmt.Errorf("创建机器人用户失败: %v", userErr)
	}

	// 3. 调 open_rpc 创建推送机器人（关联刚创建的用户 ID）
	rpcRes, rpcErr := l.svcCtx.OpenRpc.CreateBot(l.ctx, &open_rpc.CreateBotReq{
		GroupId: req.GroupID,
		BotId:   userRes.UserID, // 使用 user_rpc 生成的用户 ID
	})
	if rpcErr != nil {
		return nil, errors.New("创建失败")
	}

	// 4. 添加机器人为群成员
	memberRow, err := l.prepareGroupMember(req.GroupID, userRes.UserID)
	if err != nil {
		return nil, err
	}

	txErr := l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		// 添加机器人为群成员
		if memberRow != nil {
			if err := tx.Create(memberRow).Error; err != nil {
				return errors.New("机器人入群失败")
			}
		}

		// 写入群内展示信息
		botType := req.Type
		if botType == "" {
			botType = "custom"
		}
		ref := &group_models.GroupBotModel{
			GroupID:   req.GroupID,
			BotID:     userRes.UserID, // 关联 users.user_id
			Type:      botType,
			Status:    1,
			CreatorID: req.UserID,
		}
		if err := tx.Create(ref).Error; err != nil {
			return errors.New("写入引用表失败")
		}
		return nil
	})

	if txErr != nil {
		// Saga 补偿：回滚 open 侧的机器人记录
		_, _ = l.svcCtx.OpenRpc.DeleteBot(l.ctx, &open_rpc.DeleteBotReq{Id: rpcRes.Id})
		return nil, txErr
	}

	// 5. 调用 open_rpc 获取 Secret（只在创建时返回一次）
	secretRes, err := l.svcCtx.OpenRpc.ResetBotSecret(l.ctx, &open_rpc.ResetBotSecretReq{
		Id: rpcRes.Id,
	})
	if err != nil {
		return nil, errors.New("获取密钥失败")
	}

	// 6. 拼接完整 Webhook URL（从配置获取基础 URL）
	fullWebhookURL := fmt.Sprintf("%s%s?token=%s", l.svcCtx.Config.Domain, rpcRes.WebhookUrl, rpcRes.Token)

	l.logger.Info(model.LogMsg{
		Text: "群机器人创建成功",
		Data: map[string]interface{}{
			"groupId": req.GroupID,
			"userId":  req.UserID,
			"botId":   userRes.UserID,
		},
	})

	return &types.CreateBotRes{
		BotID:      userRes.UserID, // 机器人用户 ID
		WebhookURL: fullWebhookURL,
		Secret:     secretRes.SignatureSecret,
	}, nil
}

// prepareGroupMember 准备群成员记录（已存在则跳过）
func (l *CreateBotLogic) prepareGroupMember(groupID, userID string) (memberRow *group_models.GroupMemberModel, err error) {
	// 检查机器人是否已在群内
	var existMember group_models.GroupMemberModel
	if l.svcCtx.DB.Where("group_id = ? AND user_id = ? AND status = 1", groupID, userID).First(&existMember).Error != nil {
		memberVersion := l.svcCtx.VersionGen.GetNextVersion("group_members", "group_id", groupID)
		if memberVersion != -1 {
			memberRow = &group_models.GroupMemberModel{
				GroupID:  groupID,
				UserID:   userID,
				Role:     3, // 普通成员
				Status:   1,
				JoinTime: time.Now(),
				Version:  memberVersion,
			}
		}
	}
	return memberRow, nil
}
