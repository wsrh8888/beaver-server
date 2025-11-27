package logic

import (
	"context"

	"beaver/app/emoji/emoji_models"
	"beaver/app/emoji/emoji_rpc/internal/svc"
	"beaver/app/emoji/emoji_rpc/types/emoji_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserEmojiCollectsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserEmojiCollectsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserEmojiCollectsLogic {
	return &GetUserEmojiCollectsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取用户收藏的表情完整数据（用于datasync调用）
func (l *GetUserEmojiCollectsLogic) GetUserEmojiCollects(in *emoji_rpc.GetUserEmojiCollectsReq) (*emoji_rpc.GetUserEmojiCollectsRes, error) {
	var collects []emoji_models.EmojiCollectEmoji
	query := l.svcCtx.DB.Where("user_id = ?", in.UserId)

	// 时间戳过滤：只返回更新时间大于since的记录
	if in.Since > 0 {
		query = query.Where("updated_at > ?", in.Since)
	}

	err := query.Find(&collects).Error
	if err != nil {
		l.Errorf("查询用户收藏表情失败: userId=%s, since=%d, error=%v", in.UserId, in.Since, err)
		return nil, err
	}

	l.Infof("查询到用户 %s 的 %d 个收藏表情", in.UserId, len(collects))

	// 转换为响应格式
	var collectList []*emoji_rpc.EmojiCollectListById
	for _, collect := range collects {
		collectList = append(collectList, &emoji_rpc.EmojiCollectListById{
			Uuid:    collect.UUID,
			UserId:  collect.UserID,
			EmojiId: collect.EmojiID,
			Version: collect.Version,
			CreateAt: collect.CreatedAt.UnixMilli(),
			UpdateAt: collect.UpdatedAt.UnixMilli(),
		})
	}

	return &emoji_rpc.GetUserEmojiCollectsRes{Collects: collectList}, nil
}
