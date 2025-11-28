package logic

import (
	"context"
	"time"

	"beaver/app/emoji/emoji_models"
	"beaver/app/emoji/emoji_rpc/internal/svc"
	"beaver/app/emoji/emoji_rpc/types/emoji_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEmojiPackagesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetEmojiPackagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmojiPackagesLogic {
	return &GetEmojiPackagesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetEmojiPackagesLogic) GetEmojiPackages(in *emoji_rpc.GetEmojiPackagesReq) (*emoji_rpc.GetEmojiPackagesRes, error) {
	// 查询用户相关的所有表情包（官方表情包 + 用户创建的表情包）
	var packages []emoji_models.EmojiPackage
	query := l.svcCtx.DB.Where("(type = ? OR user_id = ?) AND status = ?",
		"official", in.UserId, 1) // 1=正常状态

	// 增量同步：只返回更新时间大于since的记录
	if in.Since > 0 {
		sinceTime := time.UnixMilli(in.Since)
		query = query.Where("updated_at > ?", sinceTime)
	}

	err := query.Find(&packages).Error
	if err != nil {
		l.Errorf("查询表情包版本信息失败: userId=%s, since=%d, error=%v", in.UserId, in.Since, err)
		return nil, err
	}

	l.Infof("查询到 %d 个表情包版本变更", len(packages))

	// 转换为版本摘要格式
	var packageVersions []*emoji_rpc.EmojiPackageVersionItem
	for _, pkg := range packages {
		packageVersions = append(packageVersions, &emoji_rpc.EmojiPackageVersionItem{
			Id:      pkg.UUID,
			Version: pkg.Version,
		})
	}

	return &emoji_rpc.GetEmojiPackagesRes{
		EmojiPackageVersions: packageVersions,
		ServerTimestamp:      time.Now().UnixMilli(),
	}, nil
}
