package logic

import (
	"context"

	"beaver/app/emoji/emoji_api/internal/svc"
	"beaver/app/emoji/emoji_api/internal/types"
	"beaver/app/emoji/emoji_models"
	mqwsconst "beaver/common/const/mqwsconst"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"
	"beaver/utils/logger"
	"beaver/utils/logger/model"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type UpdateFavoriteEmojiPackageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logger.Logger
}

func NewUpdateFavoriteEmojiPackageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFavoriteEmojiPackageLogic {
	return &UpdateFavoriteEmojiPackageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: logger.New("update_favorite_emoji_package"),
	}
}

func (l *UpdateFavoriteEmojiPackageLogic) UpdateFavoriteEmojiPackage(req *types.UpdateFavoriteEmojiPackageReq) (*types.UpdateFavoriteEmojiPackageRes, error) {
	// 1. 检查表情包是否存在
	var emojiPackage emoji_models.EmojiPackage
	err := l.svcCtx.DB.Where("package_id = ?", req.PackageID).First(&emojiPackage).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "表情包不存在")
		}
		return nil, status.Error(codes.Internal, "获取表情包失败")
	}

	// 2. 检查表情包状态
	if emojiPackage.Status != 1 {
		return nil, status.Error(codes.PermissionDenied, "表情包已禁用")
	}

	// 3. 检查是否已收藏
	var collectRecord emoji_models.EmojiPackageCollect
	err = l.svcCtx.DB.Where("user_id = ? AND package_id = ?", req.UserID, req.PackageID).
		First(&collectRecord).Error

	// 4. 根据操作类型处理
	if req.Type == "favorite" {
		// 收藏
		if err == nil {
			return nil, status.Error(codes.AlreadyExists, "已经收藏过了")
		}
		// 生成收藏版本号（按用户ID分区）
		collectVersion := l.svcCtx.VersionGen.GetNextVersion("emoji_package_collect", "user_id", req.UserID)
		if collectVersion == -1 {
			logx.Error("生成表情包收藏版本号失败")
			return nil, status.Error(codes.Internal, "生成版本号失败")
		}

		collectRecord = emoji_models.EmojiPackageCollect{
			PackageCollectID: uuid.New().String(),
			UserID:           req.UserID,
			PackageID:        req.PackageID,
			Version:          collectVersion,
		}
		err = l.svcCtx.DB.Create(&collectRecord).Error
		if err != nil {
			return nil, status.Error(codes.Internal, "收藏失败")
		}

		// 异步通过 WS 通知用户其他客户端
		go func(etcdAddr string, userId string, packageCollectId string, version int64, packageId string, packageVersion int64) {
			// 查询表情包包含的表情数据
			var packageEmojis []emoji_models.EmojiPackageEmoji
			err := l.svcCtx.DB.Where("package_id = ?", packageId).Find(&packageEmojis).Error
			if err != nil {
				logx.Errorf("查询表情包内容失败: packageId=%s, error=%v", packageId, err)
				return
			}

			// 提取表情ID列表
			emojiIds := make([]string, 0, len(packageEmojis))
			for _, pe := range packageEmojis {
				emojiIds = append(emojiIds, pe.EmojiID)
			}

			// 查询表情版本信息（只需要ID和版本号）
			var emojis []struct {
				EmojiID string `gorm:"column:emoji_id"`
				Version int64  `gorm:"column:version"`
			}
			if len(emojiIds) > 0 {
				err = l.svcCtx.DB.Table("emojis").Where("emoji_id IN ?", emojiIds).
					Select("emoji_id, version").Find(&emojis).Error
				if err != nil {
					logx.Errorf("查询表情版本信息失败: emojiIds=%v, error=%v", emojiIds, err)
					return
				}
			}

			// 构建表更新数据
			var tableUpdates []map[string]interface{}

			// 1. 通知表情包收藏表更新
			collectUpdates := map[string]interface{}{
				"table":  "emoji_package_collect",
				"userId": userId,
				"data": []map[string]interface{}{
					{
						"version":          version,
						"packageCollectId": packageCollectId,
					},
				},
			}
			tableUpdates = append(tableUpdates, collectUpdates)

			// 2. 通知表情包表更新（只发送表情包ID和版本号）
			packageUpdates := map[string]interface{}{
				"table":  "emoji_package",
				"userId": userId,
				"data": []map[string]interface{}{
					{
						"packageId": packageId,
						"version":   packageVersion,
					},
				},
			}
			tableUpdates = append(tableUpdates, packageUpdates)

			// 3. 通知表情包内容表更新（发送表情包内容的ID和版本号列表）
			contentData := make([]map[string]interface{}, 0, len(packageEmojis))
			for _, pe := range packageEmojis {
				contentData = append(contentData, map[string]interface{}{
					"relationId": pe.RelationID,
					"version":    pe.Version,
				})
			}
			contentUpdates := map[string]interface{}{
				"table":  "emoji_package_emoji",
				"userId": userId,
				"data":   contentData,
			}
			tableUpdates = append(tableUpdates, contentUpdates)

			// 4. 通知表情表更新（发送表情的ID和版本号列表）
			emojiData := make([]map[string]interface{}, 0, len(emojis))
			for _, emoji := range emojis {
				emojiData = append(emojiData, map[string]interface{}{
					"emojiId": emoji.EmojiID,
					"version": emoji.Version,
				})
			}
			emojiUpdates := map[string]interface{}{
				"table":  "emoji",
				"userId": userId,
				"data":   emojiData,
			}
			tableUpdates = append(tableUpdates, emojiUpdates)

			// 通知给自己（用户ID作为接收者，空字符串作为发送者表示系统操作）
			payload := map[string]interface{}{
				"command":  wsCommandConst.EMOJI,
				"type":     wsTypeConst.EmojiReceive,
				"senderId": "",
				"targetId": userId,
				"body": map[string]interface{}{
					"tableUpdates": tableUpdates,
				},
				"conversationId": "",
			}
			l.svcCtx.RocketMQ.SendMessage(context.Background(), mqwsconst.MqTopicWs, payload)
		}(l.svcCtx.Config.Etcd, req.UserID, collectRecord.PackageCollectID, collectVersion, req.PackageID, emojiPackage.Version)

		l.logger.Info(model.LogMsg{
			Text: "收藏表情包成功",
			Data: map[string]interface{}{
				"userId":    req.UserID,
				"packageId": req.PackageID,
			},
		})
	} else if req.Type == "unfavorite" {
		// 取消收藏
		if err != nil {
			return nil, status.Error(codes.NotFound, "未收藏过")
		}

		// 软删除：设置IsDeleted为true并更新版本号（按用户ID分区）
		collectRecord.IsDeleted = true
		collectRecord.Version = l.svcCtx.VersionGen.GetNextVersion("emoji_package_collect", "user_id", req.UserID)
		if collectRecord.Version == -1 {
			logx.Error("生成版本号失败")
			return nil, status.Error(codes.Internal, "生成版本号失败")
		}

		err = l.svcCtx.DB.Save(&collectRecord).Error
		if err != nil {
			logx.Error("软删除收藏失败", err)
			return nil, status.Error(codes.Internal, "软删除收藏失败")
		}

		// 异步通过 WS 通知用户其他客户端
		go func(etcdAddr string, userId string, packageCollectId string, version int64, packageId string, packageVersion int64) {
			// 构建表更新数据
			var tableUpdates []map[string]interface{}

			// 1. 通知表情包收藏表更新
			collectUpdates := map[string]interface{}{
				"table":  "emoji_package_collect",
				"userId": userId,
				"data": []map[string]interface{}{
					{
						"version":          version,
						"packageCollectId": packageCollectId,
					},
				},
			}
			tableUpdates = append(tableUpdates, collectUpdates)

			// 2. 通知表情包表更新（取消收藏后用户失去访问权限）
			packageUpdates := map[string]interface{}{
				"table":  "emoji_package",
				"userId": userId,
				"data": []map[string]interface{}{
					{
						"packageId": packageId,
						"version":   packageVersion,
					},
				},
			}
			tableUpdates = append(tableUpdates, packageUpdates)

			// 3. 通知表情包内容表更新
			contentUpdates := map[string]interface{}{
				"table":  "emoji_package_emoji",
				"userId": userId,
				"data": []map[string]interface{}{
					{
						"packageId": packageId,
						"version":   emojiPackage.Version,
					},
				},
			}
			tableUpdates = append(tableUpdates, contentUpdates)

			// 4. 通知表情表更新（取消收藏表情包后可能失去对某些表情的访问权限）
			emojiUpdates := map[string]interface{}{
				"table":  "emoji",
				"userId": userId,
				"data": []map[string]interface{}{
					{
						"version": 0,
					},
				},
			}
			tableUpdates = append(tableUpdates, emojiUpdates)

			// 通知给自己（用户ID作为接收者，空字符串作为发送者表示系统操作）
			payload := map[string]interface{}{
				"command":  wsCommandConst.EMOJI,
				"type":     wsTypeConst.EmojiReceive,
				"senderId": "",
				"targetId": userId,
				"body": map[string]interface{}{
					"tableUpdates": tableUpdates,
				},
				"conversationId": "",
			}
			l.svcCtx.RocketMQ.SendMessage(context.Background(), mqwsconst.MqTopicWs, payload)
		}(l.svcCtx.Config.Etcd, req.UserID, collectRecord.PackageCollectID, collectRecord.Version, req.PackageID, emojiPackage.Version)

		l.logger.Info(model.LogMsg{
			Text: "取消收藏表情包成功",
			Data: map[string]interface{}{
				"userId":    req.UserID,
				"packageId": req.PackageID,
			},
		})
	} else {
		return nil, status.Error(codes.InvalidArgument, "无效的操作类型")
	}

	return &types.UpdateFavoriteEmojiPackageRes{}, nil
}
