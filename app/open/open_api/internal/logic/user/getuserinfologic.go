package user

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"beaver/app/open/constants"
	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/app/open/open_models"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserInfoLogic) GetUserInfo(req *types.GetUserInfoReq) (resp *types.GetUserInfoRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("userId 不能为空")
	}

	tokenRecord, err := loadUserAccessToken(l.svcCtx.DB, req.Authorization, constants.ScopeUserProfileRead)
	if err != nil {
		return nil, err
	}
	if req.UserID != tokenRecord.UserID {
		return nil, errors.New("无权查询该用户信息")
	}

	userRes, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &user_rpc.UserInfoReq{UserID: req.UserID})
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	scopes := parseScopeList(tokenRecord.Scope)
	item := types.GetUserInfoUserItem{
		UserID:   userRes.UserInfo.UserId,
		Nickname: userRes.UserInfo.NickName,
	}
	if hasScope(scopes, constants.ScopeUserAvatarRead) {
		item.Avatar = userRes.UserInfo.Avatar
	}
	if hasScope(scopes, constants.ScopeUserEmailRead) {
		item.Email = userRes.UserInfo.Email
	}

	return &types.GetUserInfoRes{User: item}, nil
}

func loadUserAccessToken(db *gorm.DB, authorization string, required constants.ScopeType) (*open_models.OpenOAuthToken, error) {
	token := parseBearerToken(authorization)
	if token == "" {
		return nil, errors.New("缺少访问令牌")
	}

	var record open_models.OpenOAuthToken
	if err := db.Where("token = ?", token).First(&record).Error; err != nil {
		return nil, errors.New("访问令牌无效")
	}
	if time.Now().Unix() > record.ExpiresAt {
		return nil, errors.New("访问令牌已过期")
	}
	if record.UserID == "" {
		return nil, errors.New("需要用户授权令牌")
	}
	if !hasScope(parseScopeList(record.Scope), required) {
		return nil, errors.New("权限不足: 缺少 " + string(required))
	}
	return &record, nil
}

func parseBearerToken(authorization string) string {
	if authorization == "" {
		return ""
	}
	parts := strings.SplitN(authorization, " ", 2)
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		return parts[1]
	}
	return authorization
}

func parseScopeList(scopeStr string) []string {
	scopeStr = strings.TrimSpace(scopeStr)
	if scopeStr == "" {
		return nil
	}
	if strings.HasPrefix(scopeStr, "[") {
		var scopes []string
		if err := json.Unmarshal([]byte(scopeStr), &scopes); err == nil {
			return scopes
		}
	}
	return strings.FieldsFunc(scopeStr, func(r rune) bool {
		return r == ' ' || r == ','
	})
}

func hasScope(granted []string, required constants.ScopeType) bool {
	target := string(required)
	for _, item := range granted {
		if item == target {
			return true
		}
	}
	return false
}
