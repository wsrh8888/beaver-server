package logic

import (
	"context"

	"beaver/app/emoji/emoji_models"
	"beaver/app/emoji/emoji_rpc/internal/svc"
	"beaver/app/emoji/emoji_rpc/types/emoji_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEmojiCollectsListByIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetEmojiCollectsListByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmojiCollectsListByIdsLogic {
	return &GetEmojiCollectsListByIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetEmojiCollectsListByIdsLogic) GetEmojiCollectsListByIds(in *emoji_rpc.GetEmojiCollectsListByIdsReq) (*emoji_rpc.GetEmojiCollectsListByIdsRes, error) {
	if len(in.Ids) == 0 {
		return &emoji_rpc.GetEmojiCollectsListByIdsRes{Collects: []*emoji_rpc.EmojiCollectListById{}}, nil
	}

	var collects []emoji_models.EmojiCollectEmoji
	query := l.svcCtx.DB.Where("id IN (?)", in.Ids)

	// 时间戳过滤：只返回更新时间大于since的记录
	if in.Since > 0 {
		query = query.Where("updated_at > ?", in.Since)
	}

	err := query.Find(&collects).Error
	if err != nil {
		l.Errorf("查询收藏表情列表失败: ids=%v, since=%d, error=%v", in.Ids, in.Since, err)
		return nil, err
	}

	l.Infof("查询到 %d 个收藏表情详情", len(collects))

	// 转换为响应格式
	var collectList []*emoji_rpc.EmojiCollectListById
	for _, collect := range collects {
		collectList = append(collectList, &emoji_rpc.EmojiCollectListById{
			Id:      uint32(collect.ID),
			UserId:  collect.UserID,
			EmojiId: uint32(collect.EmojiID),
			Version: collect.Version,
			CreateAt: collect.CreatedAt.UnixMilli(),
			UpdateAt: collect.UpdatedAt.UnixMilli(),
		})
	}

	return &emoji_rpc.GetEmojiCollectsListByIdsRes{Collects: collectList}, nil
}
