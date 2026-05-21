package logic

import (
	"context"
	"errors"
	"time"

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/group/group_models"
	"beaver/app/user/user_models"
	"beaver/app/open/open_rpc/types/open_rpc"
	"beaver/utils/pwd"
	utils "beaver/utils/rand"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type CreateNotificationBotLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 在群里创建通知机器人（群管理员操作，返回 Webhook URL + Secret）
func NewCreateNotificationBotLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateNotificationBotLogic {
	return &CreateNotificationBotLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateNotificationBotLogic) CreateNotificationBot(req *types.CreateNotificationBotReq) (resp *types.CreateNotificationBotRes, err error) {
	var member group_models.GroupMemberModel
	if err = l.svcCtx.DB.Take(&member, "group_id = ? AND user_id = ?", req.GroupID, req.UserID).Error; err != nil {
		return nil, errors.New("不是群成员")
	}
	if member.Role != 1 && member.Role != 2 {
		return nil, errors.New("无权限，仅群主或管理员可创建通知机器人")
	}

	// 调 open_rpc 创建 Webhook 主记录（open 服务是 webhook 数据的 master）
	rpcRes, rpcErr := l.svcCtx.OpenRpc.CreateWebhook(l.ctx, &open_rpc.CreateWebhookReq{
		GroupId: req.GroupID,
	})
	if rpcErr != nil {
		return nil, errors.New("创建失败")
	}

	// 本地事务：建机器人用户 + 加群 + 写引用表
	botUser, memberRow := l.prepareBotResources(req.GroupID, rpcRes.BotUserId, req.Name)

	txErr := l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		if botUser != nil {
			if err := tx.Create(botUser).Error; err != nil {
				return errors.New("创建机器人用户失败")
			}
		}
		if memberRow != nil {
			if err := tx.Create(memberRow).Error; err != nil {
				return errors.New("机器人入群失败")
			}
		}
		botType := req.Type
		if botType == "" {
			botType = "custom"
		}
		ref := &group_models.GroupNotificationBotModel{
			GroupID:       req.GroupID,
			WebhookID:     uint(rpcRes.Id),
			BotUserID:     rpcRes.BotUserId,
			Name:          req.Name,
			Description:   req.Description,
			Avatar:        req.Avatar,
			WebhookURL:    rpcRes.WebhookUrl,
			Type:          botType,
			Status:        1,
			CreatorUserID: req.UserID,
		}
		if err := tx.Create(ref).Error; err != nil {
			return errors.New("写入引用表失败")
		}
		return nil
	})
	if txErr != nil {
		// Saga 补偿：回滚 open 侧的 webhook 记录
		_, _ = l.svcCtx.OpenRpc.DeleteWebhook(l.ctx, &open_rpc.DeleteWebhookReq{Id: rpcRes.Id})
		return nil, txErr
	}

	return &types.CreateNotificationBotRes{
		ID:         int64(rpcRes.Id),
		WebhookURL: rpcRes.WebhookUrl,
		Secret:     rpcRes.Secret,
	}, nil
}

// prepareBotResources 准备机器人用户和群成员记录（已存在则跳过）
func (l *CreateNotificationBotLogic) prepareBotResources(groupID, botUserID, botName string) (botUser *user_models.UserModel, memberRow *group_models.GroupMemberModel) {
	var existUser user_models.UserModel
	if l.svcCtx.DB.Where("user_id = ?", botUserID).First(&existUser).Error != nil {
		version := l.svcCtx.VersionGen.GetNextVersion("users", "user_id", botUserID)
		if version != -1 {
			randomPwd := utils.GenerateRandomString(32)
			botUser = &user_models.UserModel{
				UserID:   botUserID,
				NickName: botName,
				Password: pwd.HahPwd(randomPwd),
				Email:    botUserID + "@beaver.bot",
				Status:   1,
				IsBot:    1,
				BotAppID: "GROUP_NOTIFICATION",
				Version:  version,
			}
		}
	}

	var existMember group_models.GroupMemberModel
	if l.svcCtx.DB.Where("group_id = ? AND user_id = ? AND status = 1", groupID, botUserID).First(&existMember).Error != nil {
		memberVersion := l.svcCtx.VersionGen.GetNextVersion("group_members", "group_id", groupID)
		if memberVersion != -1 {
			memberRow = &group_models.GroupMemberModel{
				GroupID:  groupID,
				UserID:   botUserID,
				Role:     3,
				Status:   1,
				JoinTime: time.Now(),
				Version:  memberVersion,
			}
		}
	}
	return
}
