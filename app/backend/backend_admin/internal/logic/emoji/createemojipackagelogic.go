package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/emoji/emoji_rpc/types/emoji_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateEmojiPackageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateEmojiPackageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateEmojiPackageLogic {
	return &CreateEmojiPackageLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

// CreateEmojiPackage 管理后台：创建表情包。
// admin 职责：校验运营必填项，从 header 取操作者/归属用户，映射创建请求。
// RPC 职责：SaveEmojiPackage 创建分支（同名检测、版本号、落库）。
func (l *CreateEmojiPackageLogic) CreateEmojiPackage(req *types.CreateEmojiPackageReq) (resp *types.CreateEmojiPackageRes, err error) {
	if req.Title == "" {
		return nil, errors.New("表情包名称不能为空")
	}
	if req.UserID == "" {
		return nil, errors.New("用户ID不能为空")
	}

	coverFile := ""
	if req.CoverFile != nil {
		coverFile = *req.CoverFile
	}

	rpcRes, err := l.svcCtx.EmojiRpc.SaveEmojiPackage(l.ctx, &emoji_rpc.SaveEmojiPackageReq{
		Title:       req.Title,
		CoverFile:   coverFile,
		UserId:      req.UserID,
		Description: req.Description,
		Type:        req.Type,
	})
	if err != nil {
		l.Errorf("创建表情包失败: %v", err)
		return nil, err
	}
	return &types.CreateEmojiPackageRes{PackageId: rpcRes.PackageId}, nil
}
