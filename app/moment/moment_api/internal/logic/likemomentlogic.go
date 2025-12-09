package logic

import (
	"context"
	"encoding/json"
	"errors"

	"beaver/app/moment/moment_api/internal/svc"
	"beaver/app/moment/moment_api/internal/types"
	"beaver/app/moment/moment_models"
	"beaver/app/notification/notification_models"
	"beaver/app/notification/notification_rpc/types/notification_rpc"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type LikeMomentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLikeMomentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikeMomentLogic {
	return &LikeMomentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LikeMomentLogic) LikeMoment(req *types.LikeMomentReq) (resp *types.LikeMomentRes, err error) {
	var like moment_models.MomentLikeModel
	var likeId string

	// 查询动态作者
	var moment moment_models.MomentModel
	if err := l.svcCtx.DB.Where("moment_id = ? AND is_deleted = false", req.MomentID).
		First(&moment).Error; err != nil {
		return nil, errors.New("moment not found")
	}

	if req.Status {
		// 点赞操作
		like = moment_models.MomentLikeModel{
			LikeID:   uuid.New().String(),
			MomentID: req.MomentID,
			UserID:   req.UserID,
		}

		// 检查是否已点赞
		result := l.svcCtx.DB.Where("moment_id = ? AND user_id = ?", req.MomentID, req.UserID).First(&like)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// 未找到记录，进行点赞
			if err := l.svcCtx.DB.Create(&like).Error; err != nil {
				return nil, err
			}
			likeId = like.LikeID
		} else if result.Error != nil {
			// 其他错误
			return nil, result.Error
		} else {
			// 记录存在，如果被软删除了则恢复
			if like.IsDeleted {
				if err := l.svcCtx.DB.Model(&like).Update("is_deleted", false).Error; err != nil {
					return nil, err
				}
				likeId = like.LikeID
			}
			// 如果记录存在且未被删除，则无需操作
			if !like.IsDeleted {
				return &types.LikeMomentRes{}, nil
			}
		}
	} else {
		// 取消点赞操作（软删除）
		var existing moment_models.MomentLikeModel
		res := l.svcCtx.DB.Where("moment_id = ? AND user_id = ? AND is_deleted = false", req.MomentID, req.UserID).First(&existing)
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return &types.LikeMomentRes{}, nil
		}
		if res.Error != nil {
			return nil, res.Error
		}
		if err := l.svcCtx.DB.Model(&existing).Update("is_deleted", true).Error; err != nil {
			return nil, err
		}
		likeId = existing.LikeID
	}

	// 异步推送通知
	go func() {
		// 作者等于自己则不推
		if moment.UserID == "" || moment.UserID == req.UserID {
			return
		}
		payload, _ := json.Marshal(map[string]interface{}{
			"momentId": req.MomentID,
			"likeId":   likeId,
			"userId":   req.UserID,
			"status":   req.Status,
		})
		eventType := notification_models.EventTypeMomentLike
		if !req.Status {
			eventType = notification_models.EventTypeMomentUnlike
		}
		_, err := l.svcCtx.NotifyRpc.PushEvent(l.ctx, &notification_rpc.PushEventReq{
			EventType:   eventType,
			Category:    notification_models.CategoryMoment,
			FromUserId:  req.UserID,
			TargetId:    req.MomentID,
			TargetType:  notification_models.TargetTypeMoment,
			PayloadJson: string(payload),
			ToUserIds:   []string{moment.UserID},
			DedupHash:   likeId + "_" + eventType,
		})
		if err != nil {
			l.Logger.Errorf("推送点赞通知失败: %v", err)
		}
	}()

	return resp, nil
}
