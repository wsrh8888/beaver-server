package post

import (
	"context"
	"fmt"
	"unicode/utf8"

	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/circle/circle_api/internal/svc"
	"beaver/app/circle/circle_api/internal/types"
	"beaver/app/circle/circle_models"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePostLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreatePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePostLogic {
	return &CreatePostLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreatePostLogic) CreatePost(req *types.CreatePostReq) (resp *types.CreatePostRes, err error) {
	var member circle_models.CircleMemberModel
	if err = l.svcCtx.DB.Where("circle_id = ? AND user_id = ?", req.CircleID, req.UserID).First(&member).Error; err != nil {
		return nil, fmt.Errorf("请先加入圈子再发帖")
	}

	postID := uuid.New().String()
	var files *circle_models.PostFiles
	if len(req.Files) > 0 {
		f := make(circle_models.PostFiles, 0, len(req.Files))
		for _, fi := range req.Files {
			f = append(f, circle_models.PostFileInfo{FileKey: fi.FileKey, Type: fi.Type})
		}
		files = &f
	}

	post := circle_models.CirclePostModel{
		PostID:   postID,
		CircleID: req.CircleID,
		UserID:   req.UserID,
		Content:  req.Content,
		Files:    files,
	}
	if err = l.svcCtx.DB.Create(&post).Error; err != nil {
		return nil, fmt.Errorf("发布帖子失败: %v", err)
	}

	// 确保所有成员的圈子会话存在，并把动态顶到会话列表
	var members []circle_models.CircleMemberModel
	l.svcCtx.DB.Where("circle_id = ?", req.CircleID).Find(&members)
	userIDs := make([]string, 0, len(members))
	for _, m := range members {
		userIDs = append(userIDs, m.UserID)
	}
	conversationID := fmt.Sprintf("circle_%s", req.CircleID)
	if len(userIDs) > 0 {
		_, _ = l.svcCtx.ChatRpc.InitializeConversation(l.ctx, &chat_rpc.InitializeConversationReq{
			ConversationId: conversationID,
			Type:           3,
			UserIds:        userIDs,
		})
	}

	preview := buildPostPreview(req.Content, len(req.Files) > 0)
	_, _ = l.svcCtx.ChatRpc.SendNotificationMessage(l.ctx, &chat_rpc.SendNotificationMessageReq{
		ConversationId: conversationID,
		MessageType:    1,
		Content:        preview,
		RelatedUserId:  req.UserID,
		ReadUserIds:    []string{req.UserID},
	})

	userName, avatar := "", ""
	userResp, _ := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{UserIdList: []string{req.UserID}})
	if userResp != nil {
		if info := userResp.UserInfo[req.UserID]; info != nil {
			userName = info.NickName
			avatar = info.Avatar
		}
	}

	return &types.CreatePostRes{
		PostID:    postID,
		CircleID:  req.CircleID,
		UserID:    req.UserID,
		UserName:  userName,
		Avatar:    avatar,
		Content:   req.Content,
		CreatedAt: post.CreatedAt.String(),
	}, nil
}

func buildPostPreview(content string, hasFiles bool) string {
	content = trimPostText(content)
	if content != "" {
		if utf8.RuneCountInString(content) > 40 {
			runes := []rune(content)
			return string(runes[:40]) + "..."
		}
		return content
	}
	if hasFiles {
		return "[图片]"
	}
	return "[新动态]"
}

func trimPostText(s string) string {
	start, end := 0, len(s)
	for start < end {
		r, size := utf8.DecodeRuneInString(s[start:])
		if r != ' ' && r != '\n' && r != '\t' && r != '\r' {
			break
		}
		start += size
	}
	for end > start {
		r, size := utf8.DecodeLastRuneInString(s[:end])
		if r != ' ' && r != '\n' && r != '\t' && r != '\r' {
			break
		}
		end -= size
	}
	return s[start:end]
}
