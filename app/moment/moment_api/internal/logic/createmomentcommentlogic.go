package logic

import (
	"context"
	"errors"
	"strings"

	"beaver/app/moment/moment_api/internal/svc"
	"beaver/app/moment/moment_api/internal/types"
	"beaver/app/moment/moment_models"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type CreateMomentCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateMomentCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateMomentCommentLogic {
	return &CreateMomentCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CreateMomentComment 发表评论
func (l *CreateMomentCommentLogic) CreateMomentComment(req *types.CreateMomentCommentReq) (resp *types.CreateMomentCommentRes, err error) {
	content := strings.TrimSpace(req.Content)
	if content == "" {
		return nil, errors.New("comment content is empty")
	}

	parentId := strings.TrimSpace(req.ParentId)
	replyToCommentId := strings.TrimSpace(req.ReplyToCommentId)

	// 规则：存储只保留两层
	// - parentId 始终指向顶层评论（根）
	// - replyToCommentId 指向具体被回复的评论（可为顶层或子评论）
	var targetComment moment_models.MomentCommentModel
	if replyToCommentId != "" {
		if err := l.svcCtx.DB.Where("uuid = ? AND is_deleted = false", replyToCommentId).
			First(&targetComment).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("reply comment not found")
			}
			return nil, err
		}
		if targetComment.MomentID != req.MomentID {
			return nil, errors.New("comment not belong to moment")
		}
		// 顶层评论自身就是根；子评论的父就是根
		if targetComment.ParentID == "" {
			parentId = targetComment.UUID
		} else {
			parentId = targetComment.ParentID
		}
	} else if parentId != "" {
		// 直接指定 parentId，要求它必须是顶层，防止出现第三层
		if err := l.svcCtx.DB.Where("uuid = ? AND is_deleted = false", parentId).
			First(&targetComment).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("parent comment not found")
			}
			return nil, err
		}
		if targetComment.MomentID != req.MomentID {
			return nil, errors.New("comment not belong to moment")
		}
		if targetComment.ParentID != "" {
			return nil, errors.New("only support two-level comments")
		}
		// replyToCommentId 为空时，视为回复顶层
		replyToCommentId = targetComment.UUID
	}

	comment := moment_models.MomentCommentModel{
		UUID:             uuid.New().String(),
		MomentID:         req.MomentID,
		UserID:           req.UserID,
		Content:          content,
		ParentID:         parentId,
		ReplyToCommentID: replyToCommentId,
	}

	if err := l.svcCtx.DB.Create(&comment).Error; err != nil {
		return nil, err
	}

	// 查询用户信息用于展示昵称和头像
	var userName, avatar string
	replyToUserName := ""
	userIdList := []string{req.UserID}
	if targetComment.UserID != "" {
		userIdList = append(userIdList, targetComment.UserID)
	}

	if len(userIdList) > 0 {
		userResp, err := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{
			UserIdList: userIdList,
		})
		if err == nil && userResp.UserInfo != nil {
			if info := userResp.UserInfo[req.UserID]; info != nil {
				userName = info.NickName
				avatar = info.Avatar
			}
			if targetComment.UserID != "" {
				if info := userResp.UserInfo[targetComment.UserID]; info != nil {
					replyToUserName = info.NickName
				}
			}
		}
	}

	resp = &types.CreateMomentCommentRes{
		Id:               comment.UUID,
		UserID:           comment.UserID,
		UserName:         userName,
		Avatar:           avatar,
		Content:          comment.Content,
		ParentId:         comment.ParentID,
		ReplyToCommentId: comment.ReplyToCommentID,
		ReplyToUserName:  replyToUserName,
		CreatedAt:        comment.CreatedAt.String(),
	}

	return resp, nil
}
