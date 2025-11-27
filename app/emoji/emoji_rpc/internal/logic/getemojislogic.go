package logic

import (
	"context"

	"beaver/app/emoji/emoji_models"
	"beaver/app/emoji/emoji_rpc/internal/svc"
	"beaver/app/emoji/emoji_rpc/types/emoji_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEmojisLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetEmojisLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmojisLogic {
	return &GetEmojisLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetEmojisLogic) GetEmojis(in *emoji_rpc.GetEmojisReq) (*emoji_rpc.GetEmojisRes, error) {
	// 查询用户相关的所有表情（官方表情 + 用户创建的表情 + 用户收藏的表情）
	// 这里主要返回用户创建的表情，官方表情通过表情包关联获取
	var emojis []emoji_models.Emoji
	query := l.svcCtx.DB.Where("(author_id = ? OR author_id = ?) AND status = ?",
		in.UserId, "official", 1) // 1=正常状态

	// 增量同步：只返回版本号大于since的记录
	if in.Since > 0 {
		query = query.Where("version > ?", in.Since)
	}

	err := query.Find(&emojis).Error
	if err != nil {
		l.Errorf("查询表情版本信息失败: userId=%s, since=%d, error=%v", in.UserId, in.Since, err)
		return nil, err
	}

	l.Infof("查询到用户 %s 的 %d 个表情版本信息", in.UserId, len(emojis))

	// 转换为响应格式
	var emojiVersions []*emoji_rpc.GetEmojisRes_EmojiVersion
	for _, emoji := range emojis {
		emojiVersions = append(emojiVersions, &emoji_rpc.GetEmojisRes_EmojiVersion{
			Id:      uint32(emoji.ID),
			Version: emoji.Version,
		})
	}

	return &emoji_rpc.GetEmojisRes{EmojiVersions: emojiVersions}, nil
}
