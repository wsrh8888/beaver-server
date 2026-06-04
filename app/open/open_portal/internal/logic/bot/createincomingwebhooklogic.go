package bot

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"beaver/app/group/group_rpc/types/group_rpc"
	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"
	"beaver/app/open/open_rpc/types/open_rpc"
	"beaver/app/user/user_models"
	"beaver/app/user/user_rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateIncomingWebhookLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateIncomingWebhookLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateIncomingWebhookLogic {
	return &CreateIncomingWebhookLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateIncomingWebhookLogic) CreateIncomingWebhook(req *types.CreateIncomingWebhookReq) (resp *types.CreateIncomingWebhookRes, err error) {
	if req.AppID == "" || req.GroupID == "" {
		return nil, errors.New("appId 和 groupId 不能为空")
	}

	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ? AND owner_user_id = ?", req.AppID, req.UserID).First(&app).Error; err != nil {
		return nil, errors.New("应用不存在或无权限操作")
	}

	groupRes, err := l.svcCtx.GroupRpc.GetGroupsListByIds(l.ctx, &group_rpc.GetGroupsListByIdsReq{
		GroupIDs: []string{req.GroupID},
	})
	if err != nil || len(groupRes.Groups) == 0 {
		return nil, errors.New("群组不存在")
	}

	name := req.Name
	if name == "" {
		name = "通知机器人"
	}

	userRes, err := l.svcCtx.UserRpc.UserCreate(l.ctx, &user.UserCreateReq{
		NickName: name,
		UserType: int32(user_models.UserTypeBot),
		Source:   int32(user_models.SourceGroup),
	})
	if err != nil {
		return nil, errors.New("创建推送 Bot 用户失败")
	}

	rpcRes, err := l.svcCtx.OpenRpc.CreateBot(l.ctx, &open_rpc.CreateBotReq{
		GroupId: req.GroupID,
		BotId:   userRes.UserID,
	})
	if err != nil {
		return nil, errors.New("创建推送 Bot 失败")
	}

	if err := l.svcCtx.DB.Model(&open_models.OpenBotModel{}).Where("id = ?", rpcRes.Id).Updates(map[string]interface{}{
		"app_id": req.AppID,
		"name":   name,
	}).Error; err != nil {
		_, _ = l.svcCtx.OpenRpc.DeleteBot(l.ctx, &open_rpc.DeleteBotReq{Id: rpcRes.Id})
		return nil, errors.New("保存 Bot 元数据失败")
	}

	_, err = l.svcCtx.GroupRpc.AddGroupMember(l.ctx, &group_rpc.AddGroupMemberReq{
		GroupId:    req.GroupID,
		UserId:     userRes.UserID,
		OperatedBy: req.UserID,
	})
	if err != nil {
		_, _ = l.svcCtx.OpenRpc.DeleteBot(l.ctx, &open_rpc.DeleteBotReq{Id: rpcRes.Id})
		logx.Errorf("推送 Bot 入群失败: group=%s bot=%s err=%v", req.GroupID, userRes.UserID, err)
		return nil, errors.New("推送 Bot 入群失败")
	}

	var bot open_models.OpenBotModel
	if err := l.svcCtx.DB.Where("id = ?", rpcRes.Id).First(&bot).Error; err != nil {
		return nil, errors.New("查询 Bot 失败")
	}

	return &types.CreateIncomingWebhookRes{
		Webhook: toIncomingWebhookInfo(&bot, l.svcCtx.Config.ApiBaseUrl, true),
	}, nil
}

func toIncomingWebhookInfo(bot *open_models.OpenBotModel, apiBase string, withSecret bool) types.IncomingWebhookInfo {
	info := types.IncomingWebhookInfo{
		ID:         fmt.Sprintf("%d", bot.ID),
		Token:      bot.Token,
		AppID:      bot.AppID,
		GroupID:    bot.GroupID,
		BotID:      bot.BotID,
		Name:       bot.Name,
		WebhookURL: buildBotWebhookURL(apiBase, bot.Token),
		Status:     bot.Status,
		CreatedAt:  bot.CreatedAt.Unix(),
	}
	if withSecret && bot.Security.SignatureEnabled {
		info.Secret = bot.Security.SignatureSecret
	}
	return info
}

func buildBotWebhookURL(apiBase, token string) string {
	base := strings.TrimSuffix(apiBase, "/")
	return fmt.Sprintf("%s/api/open/bot_public/v1/send?token=%s", base, token)
}