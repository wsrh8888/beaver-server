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
	// 1. 校验权限
	var member group_models.GroupMemberModel
	if err = l.svcCtx.DB.Take(&member, "group_id = ? AND user_id = ?", req.GroupID, req.UserID).Error; err != nil {
		return nil, errors.New("不是群成员")
	}
	if member.Role != 1 && member.Role != 2 {
		return nil, errors.New("无权限，仅群主或管理员可创建通知机器人")
	}

	// 2. 调 open_rpc 创建推送机器人（open 服务负责生成 Token、BotUserID 等）
	rpcRes, rpcErr := l.svcCtx.OpenRpc.CreateBot(l.ctx, &open_rpc.CreateBotReq{
		GroupId: req.GroupID,
	})
	if rpcErr != nil {
		return nil, errors.New("创建失败")
	}

	// 3. 本地事务：建机器人用户 + 加群 + 写引用表
	botUser, memberRow := l.prepareBotResources(req.GroupID, rpcRes.BotUserId, req.Name)

	txErr := l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		// 创建机器人用户
		if botUser != nil {
			if err := tx.Create(botUser).Error; err != nil {
				return errors.New("创建机器人用户失败")
			}
		}

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
			GroupID:       req.GroupID,
			BotID:         uint(rpcRes.Id), // 关联 open_bots.id
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
		// Saga 补偿：回滚 open 侧的机器人记录
		_, _ = l.svcCtx.OpenRpc.DeleteBot(l.ctx, &open_rpc.DeleteBotReq{Id: rpcRes.Id})
		return nil, txErr
	}

	// 4. 调用 open_rpc 获取 Secret（只在创建时返回一次）
	secretRes, err := l.svcCtx.OpenRpc.ResetBotSecret(l.ctx, &open_rpc.ResetBotSecretReq{
		Id: rpcRes.Id,
	})
	if err != nil {
		return nil, errors.New("获取密钥失败")
	}

	return &types.CreateNotificationBotRes{
		ID:         int64(rpcRes.Id),
		WebhookURL: fmt.Sprintf("%s?token=%s", rpcRes.WebhookUrl, rpcRes.Token),
		Secret:     secretRes.SignatureSecret,
	}, nil
}

// prepareBotResources 准备机器人用户和群成员记录（已存在则跳过）
func (l *CreateNotificationBotLogic) prepareBotResources(groupID, botUserID, botName string) (botUser *user_models.UserModel, memberRow *group_models.GroupMemberModel) {
	var existUser user_models.UserModel
	if l.svcCtx.DB.Where("user_id = ?", botUserID).First(&existUser).Error != nil {
		version := l.svcCtx.VersionGen.GetNextVersion("users", "user_id", botUserID)
		if version != -1 {
			// 机器人不需要密码，由 auth 服务管理
			botUser = &user_models.UserModel{
				UserID:   botUserID,
				UserType: user_models.UserTypeBot, // 标记为推送机器人
				NickName: botName,
				Email:    botUserID + "@beaver.bot",
				Status:   1,
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
