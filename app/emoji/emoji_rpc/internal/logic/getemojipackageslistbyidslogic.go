package logic

import (
	"context"

	"beaver/app/emoji/emoji_models"
	"beaver/app/emoji/emoji_rpc/internal/svc"
	"beaver/app/emoji/emoji_rpc/types/emoji_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEmojiPackagesListByIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetEmojiPackagesListByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmojiPackagesListByIdsLogic {
	return &GetEmojiPackagesListByIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetEmojiPackagesListByIdsLogic) GetEmojiPackagesListByIds(in *emoji_rpc.GetEmojiPackagesListByIdsReq) (*emoji_rpc.GetEmojiPackagesListByIdsRes, error) {
	if len(in.Ids) == 0 {
		return &emoji_rpc.GetEmojiPackagesListByIdsRes{Packages: []*emoji_rpc.EmojiPackageListById{}}, nil
	}

	var packages []emoji_models.EmojiPackage
	query := l.svcCtx.DB.Where("id IN (?)", in.Ids)

	// 时间戳过滤：只返回更新时间大于since的记录
	if in.Since > 0 {
		query = query.Where("updated_at > ?", in.Since)
	}

	err := query.Find(&packages).Error
	if err != nil {
		l.Errorf("查询表情包列表失败: ids=%v, since=%d, error=%v", in.Ids, in.Since, err)
		return nil, err
	}

	l.Infof("查询到 %d 个表情包详情", len(packages))

	// 转换为响应格式
	var packageList []*emoji_rpc.EmojiPackageListById
	for _, pkg := range packages {
		packageList = append(packageList, &emoji_rpc.EmojiPackageListById{
			Id:          uint32(pkg.ID),
			Title:       pkg.Title,
			CoverFile:   pkg.CoverFile,
			UserId:      pkg.UserID,
			Description: pkg.Description,
			Type:        pkg.Type,
			Status:      int32(pkg.Status),
			Version:     pkg.Version,
			CreateAt:    pkg.CreatedAt.UnixMilli(),
			UpdateAt:    pkg.UpdatedAt.UnixMilli(),
		})
	}

	return &emoji_rpc.GetEmojiPackagesListByIdsRes{Packages: packageList}, nil
}
