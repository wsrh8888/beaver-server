package logic

import (
	"context"

	filecommon "beaver/app/backend/backend_admin/internal/handler/file/common"
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/file/file_rpc/types/file_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFileListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFileListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFileListLogic {
	return &GetFileListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetFileListLogic) GetFileList(req *types.GetFileListReq) (resp *types.GetFileListRes, err error) {
	rpcRes, err := l.svcCtx.FileRpc.ListFiles(l.ctx, &file_rpc.ListFilesReq{
		Page:     int32(req.Page),
		PageSize: int32(req.Limit),
		Type:     req.Type,
		Keywords: req.Keywords,
	})
	if err != nil {
		l.Errorf("获取文件列表失败: %v", err)
		return nil, err
	}

	list := make([]types.GetFileListItem, 0, len(rpcRes.List))
	for _, f := range rpcRes.List {
		list = append(list, types.GetFileListItem{
			Id:           uint(f.Id),
			FileName:     f.FileKey,
			OriginalName: f.OriginalName,
			Size:         f.Size,
			Path:         buildFileAccessURL(l.svcCtx, f),
			Md5:          f.Md5,
			Type:         f.Type,
			CreatedAt:    f.CreatedAt,
			UpdatedAt:    f.UpdatedAt,
		})
	}

	return &types.GetFileListRes{List: list, Total: rpcRes.Total}, nil
}

// buildFileAccessURL 按来源拼完整访问 URL（本地用 Domain，七牛用 Qiniu.Domain）
func buildFileAccessURL(svcCtx *svc.ServiceContext, f *file_rpc.FileItem) string {
	if f == nil {
		return ""
	}
	if f.Source == "qiniu" {
		return filecommon.BuildQiniuFileURL(svcCtx.Config.Qiniu.Domain, f.Path)
	}
	return filecommon.BuildLocalFileURL(svcCtx.Config.Domain, f.FileKey)
}
