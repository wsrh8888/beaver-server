package logic

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/group/group_models"
	"beaver/app/open/open_models"
	"beaver/app/user/user_models"
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

	tokenBytes := make([]byte, 32)
	if _, err = rand.Read(tokenBytes); err != nil {
		return nil, errors.New("生成 access_token 失败")
	}
	secretBytes := make([]byte, 32)
	if _, err = rand.Read(secretBytes); err != nil {
		return nil, errors.New("生成 secret 失败")
	}
	token := hex.EncodeToString(tokenBytes)
	secret := hex.EncodeToString(secretBytes)

	botUserID, botUser, memberRow, createErr := l.prepareBotUser(req.GroupID, req.Name, token)
	if createErr != nil {
		return nil, createErr
	}

	record := open_models.OpenIncomingWebhook{
		Token:         token,
		Secret:        secret,
		AppID:         "GROUP_NOTIFICATION",
		GroupID:       req.GroupID,
		BotUserID:     botUserID,
		Name:          req.Name,
		Description:   req.Description,
		Avatar:        req.Avatar,
		CreatorUserID: req.UserID,
		Status:        1,
	}

	txErr := l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&record).Error; err != nil {
			return errors.New("创建失败")
		}
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
		return nil
	})
	if txErr != nil {
		return nil, txErr
	}

	webhookURL := fmt.Sprintf("%s/api/open/v1/webhook/incoming?access_token=%s", l.svcCtx.Config.ApiBaseUrl, token)
	return &types.CreateNotificationBotRes{
		ID:         int64(record.ID),
		WebhookURL: webhookURL,
		Secret:     secret,
	}, nil
}

// prepareBotUser 为本 Webhook 创建独立机器人用户并准备入群记录（已存在则复用）
func (l *CreateNotificationBotLogic) prepareBotUser(groupID, botName, accessToken string) (botUserID string, botUser *user_models.UserModel, memberRow *group_models.GroupMemberModel, err error) {
	botUserID = "nbot_" + accessToken[:32]

	var existUser user_models.UserModel
	if e := l.svcCtx.DB.Where("user_id = ?", botUserID).First(&existUser).Error; e == nil {
		botUser = nil
	} else {
		version := l.svcCtx.VersionGen.GetNextVersion("users", "user_id", botUserID)
		if version == -1 {
			return "", nil, nil, errors.New("获取用户版本号失败")
		}
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

	var existMember group_models.GroupMemberModel
	if e := l.svcCtx.DB.Where("group_id = ? AND user_id = ? AND status = 1", groupID, botUserID).First(&existMember).Error; e == nil {
		memberRow = nil
		return botUserID, botUser, memberRow, nil
	}

	memberVersion := l.svcCtx.VersionGen.GetNextVersion("group_members", "group_id", groupID)
	if memberVersion == -1 {
		return "", nil, nil, errors.New("获取群成员版本号失败")
	}
	memberRow = &group_models.GroupMemberModel{
		GroupID:  groupID,
		UserID:   botUserID,
		Role:     3,
		Status:   1,
		JoinTime: time.Now(),
		Version:  memberVersion,
	}
	return botUserID, botUser, memberRow, nil
}
