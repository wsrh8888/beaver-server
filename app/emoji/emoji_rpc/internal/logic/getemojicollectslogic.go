package logic

import (
	"context"

	"beaver/app/emoji/emoji_models"
	"beaver/app/emoji/emoji_rpc/internal/svc"
	"beaver/app/emoji/emoji_rpc/types/emoji_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEmojiCollectsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetEmojiCollectsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmojiCollectsLogic {
	return &GetEmojiCollectsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetEmojiCollectsLogic) GetEmojiCollects(in *emoji_rpc.GetEmojiCollectsReq) (*emoji_rpc.GetEmojiCollectsRes, error) {
	// 查询用户收藏的表情
	var collects []emoji_models.EmojiCollectEmoji
	query := l.svcCtx.DB.Where("user_id = ?", in.UserId)

	// 增量同步：只返回版本号大于since的记录
	if in.Since > 0 {
		query = query.Where("version > ?", in.Since)
	}

	err := query.Find(&collects).Error
	if err != nil {
		l.Errorf("查询用户收藏表情版本信息失败: userId=%s, since=%d, error=%v", in.UserId, in.Since, err)
		return nil, err
	}

	l.Infof("查询到用户 %s 的 %d 个收藏表情版本信息", in.UserId, len(collects))

	// 转换为响应格式
	var collectVersions []*emoji_rpc.GetEmojiCollectsRes_EmojiCollectVersion
	for _, collect := range collects {
		collectVersions = append(collectVersions, &emoji_rpc.GetEmojiCollectsRes_EmojiCollectVersion{
			Id:      uint32(collect.ID),
			Version: collect.Version,
		})
	}

	return &emoji_rpc.GetEmojiCollectsRes{CollectVersions: collectVersions}, nil
}
