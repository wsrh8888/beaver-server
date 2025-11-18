package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/file/file_rpc/types/file_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetQiniuUploadTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取七牛云上传token
func NewGetQiniuUploadTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetQiniuUploadTokenLogic {
	return &GetQiniuUploadTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetQiniuUploadTokenLogic) GetQiniuUploadToken(req *types.GetQiniuUploadTokenReq) (resp *types.GetQiniuUploadTokenRes, err error) {
	// 调用fileRpc服务获取七牛云上传token
	rpcResp, err := l.svcCtx.FileRpc.GetQiniuUploadToken(l.ctx, &file_rpc.GetQiniuUploadTokenReq{})
	if err != nil {
		l.Logger.Errorf("调用fileRpc获取七牛云token失败: %v", err)
		return nil, err
	}

	// 转换为HTTP响应格式
	resp = &types.GetQiniuUploadTokenRes{
		UploadToken: rpcResp.UploadToken,
		ExpiresIn:   rpcResp.ExpiresIn,
	}

	l.Logger.Infof("成功获取七牛云上传token，过期时间: %d秒", resp.ExpiresIn)
	return resp, nil
}
