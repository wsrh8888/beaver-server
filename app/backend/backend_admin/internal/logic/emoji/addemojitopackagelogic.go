package logic

import (
	"context"
	"errors"
	"strconv"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/emoji/emoji_models"

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
	if req.FileName == "" {
		return nil, errors.New("文件ID不能为空")
	}
	if req.Title == "" {
		return nil, errors.New("表情名称不能为空")
	}
	if req.AuthorID == "" {
		return nil, errors.New("创建者ID不能为空")
	}

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
	err = l.svcCtx.DB.Model(&emoji_models.Emoji{}).Where("title = ? AND author_id = ?", req.Title, req.AuthorID).Count(&count).Error
	if err != nil {
		logx.Errorf("检查表情名称失败: %v", err)
		return nil, errors.New("检查表情名称失败")
	}
	if count > 0 {
		return nil, errors.New("该创建者已存在同名表情")
	}

	// 创建表情
	emoji := emoji_models.Emoji{
		FileName: req.FileName,
		Title:    req.Title,
		AuthorID: req.AuthorID,
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
		PackageID: uint(packageID),
		EmojiID:   emoji.Id,
		SortOrder: maxSortOrder + 1, // 添加到末尾
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
