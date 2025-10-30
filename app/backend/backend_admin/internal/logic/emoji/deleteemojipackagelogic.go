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

type DeleteEmojiPackageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除表情包集合
func NewDeleteEmojiPackageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteEmojiPackageLogic {
	return &DeleteEmojiPackageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteEmojiPackageLogic) DeleteEmojiPackage(req *types.DeleteEmojiPackageReq) (resp *types.DeleteEmojiPackageRes, err error) {
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

	// 使用逻辑删除
	err = l.svcCtx.DB.Delete(&pkg).Error
	if err != nil {
		logx.Errorf("删除表情包失败: %v", err)
		return nil, errors.New("删除表情包失败")
	}

	return &types.DeleteEmojiPackageRes{}, nil
}
