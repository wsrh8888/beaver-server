package logic

import (
	"context"
	"errors"
	"strconv"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/emoji/emoji_models"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type AddEmojiToPackageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 向表情包集合中添加表情图片
func NewAddEmojiToPackageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddEmojiToPackageLogic {
	return &AddEmojiToPackageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddEmojiToPackageLogic) AddEmojiToPackage(req *types.AddEmojiToPackageReq) (resp *types.AddEmojiToPackageRes, err error) {
	// 验证必填字段
	if req.FileKey == "" {
		return nil, errors.New("文件ID不能为空")
	}
	if req.Title == "" {
		return nil, errors.New("表情名称不能为空")
	}
	// 创建者ID验证暂时移除，由上层中间件处理

	// 转换PackageID为uint
	packageID, err := strconv.ParseUint(req.PackageID, 10, 32)
	if err != nil {
		logx.Errorf("表情包ID格式错误: %s", req.PackageID)
		return nil, errors.New("表情包ID格式错误")
	}

	// 检查表情包是否存在
	var pkg emoji_models.EmojiPackage
	err = l.svcCtx.DB.Where("id = ?", packageID).First(&pkg).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logx.Errorf("表情包不存在: %s", req.PackageID)
			return nil, errors.New("表情包不存在")
		}
		logx.Errorf("查询表情包失败: %v", err)
		return nil, errors.New("查询表情包失败")
	}

	// 检查表情名称是否已存在（同创建者下）
	var count int64
	// 检查表情标题是否已存在（暂时不检查作者，允许同名表情）
	count = 0
	if err != nil {
		logx.Errorf("检查表情名称失败: %v", err)
		return nil, errors.New("检查表情名称失败")
	}
	if count > 0 {
		return nil, errors.New("该创建者已存在同名表情")
	}

	// 创建表情
	emoji := emoji_models.Emoji{
		UUID:    uuid.New().String(),
		FileKey: req.FileKey,
		Title:   req.Title,
		Status:  1, // 默认状态为正常
		Version: 0, // 暂时设为0
	}

	err = l.svcCtx.DB.Create(&emoji).Error
	if err != nil {
		logx.Errorf("创建表情失败: %v", err)
		return nil, errors.New("创建表情失败")
	}

	// 获取当前表情包中表情的最大排序值
	var maxSortOrder int
	var emojiPackageEmojis []emoji_models.EmojiPackageEmoji
	err = l.svcCtx.DB.Where("package_id = ?", packageID).Find(&emojiPackageEmojis).Error
	if err != nil {
		logx.Errorf("查询表情包关联失败: %v", err)
		return nil, errors.New("查询表情包关联失败")
	}

	// 手动计算最大排序值
	maxSortOrder = 0
	for _, emojiPackageEmoji := range emojiPackageEmojis {
		if emojiPackageEmoji.SortOrder > maxSortOrder {
			maxSortOrder = emojiPackageEmoji.SortOrder
		}
	}

	// 创建表情包与表情的关联
	emojiPackageEmoji := emoji_models.EmojiPackageEmoji{
		UUID:      uuid.New().String(),
		PackageID: strconv.FormatUint(packageID, 10),
		EmojiID:   emoji.UUID,
		SortOrder: maxSortOrder + 1, // 添加到末尾
		Version:   0,                // 暂时设为0，后续需要实现版本控制
	}

	err = l.svcCtx.DB.Create(&emojiPackageEmoji).Error
	if err != nil {
		logx.Errorf("添加表情到表情包失败: %v", err)
		return nil, errors.New("添加表情到表情包失败")
	}

	return &types.AddEmojiToPackageRes{
		Id: strconv.Itoa(int(emoji.Id)),
	}, nil
}
