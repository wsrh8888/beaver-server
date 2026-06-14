package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/emoji/emoji_rpc/types/emoji_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteEmojiPackageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteEmojiPackageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteEmojiPackageLogic {
	return &DeleteEmojiPackageLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

// DeleteEmojiPackage 管理后台：删除表情包。
// admin 职责：校验 packageId，映射为 SaveEmojiPackage.delete=true。
// RPC 职责：领域删除，不与 HTTP 路由 1:1。
func (l *DeleteEmojiPackageLogic) DeleteEmojiPackage(req *types.DeleteEmojiPackageReq) (resp *types.DeleteEmojiPackageRes, err error) {
	if req.PackageId == "" {
		return nil, errors.New("表情包ID不能为空")
	}

	del := true
	_, err = l.svcCtx.EmojiRpc.SaveEmojiPackage(l.ctx, &emoji_rpc.SaveEmojiPackageReq{
		PackageId: req.PackageId,
		Delete:    &del,
	})
	if err != nil {
		l.Errorf("删除表情包失败: %v", err)
		return nil, err
	}
	return &types.DeleteEmojiPackageRes{}, nil
}
