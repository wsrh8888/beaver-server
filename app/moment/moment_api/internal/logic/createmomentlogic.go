package logic

import (
	"context"
	"fmt"
	"strings"

	"beaver/app/moment/moment_api/internal/svc"
	"beaver/app/moment/moment_api/internal/types"
	"beaver/app/moment/moment_models"
	"beaver/app/user/user_rpc/types/user_rpc"
	"beaver/common/models/ctype"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateMomentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateMomentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateMomentLogic {
	return &CreateMomentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateMomentLogic) CreateMoment(req *types.CreateMomentReq) (resp *types.CreateMomentRes, err error) {
	// 构造MomentModel实例
	momentID := uuid.New().String()
	moment := moment_models.MomentModel{
		UserID:     req.UserID,
		Content:    req.Content,
		MomentID:   momentID,
		Files:      convertFiles(req.Files),
		Visibility: int8(req.Visibility),
		AllowList:  strings.Join(req.AllowList, ","),
		BlockList:  strings.Join(req.BlockList, ","),
	}

	// 插入数据库
	if err := l.svcCtx.DB.Create(&moment).Error; err != nil {
		return nil, fmt.Errorf("failed to create moment: %v", err)
	}

	userName := ""
	avatar := ""
	userResp, err := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{
		UserIdList: []string{req.UserID},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}
	if info := userResp.UserInfo[req.UserID]; info != nil {
		userName = info.NickName
		avatar = info.Avatar
	}

	// 构造响应数据
	resp = &types.CreateMomentRes{
		Id:           momentID,
		UserID:       req.UserID,
		UserName:     userName,
		Avatar:       avatar,
		Content:      moment.Content,
		Files:        convertToCreateMomentFileInfo(moment.Files),
		Comments:     []interface{}{}, // 创建时为空
		Likes:        []interface{}{}, // 创建时为空
		CommentCount: 0,               // 创建时为0
		LikeCount:    0,               // 创建时为0
		IsLiked:      false,           // 创建时为false
		CreatedAt:    moment.CreatedAt.String(),
	}

	return resp, nil
}

// 辅助函数：将请求中的文件信息转换为数据库模型所需的结构
func convertFiles(files []types.CreateFileInfo) *moment_models.Files {
	var result moment_models.Files
	for _, file := range files {
		result = append(result, moment_models.FileInfo{
			FileKey: file.FileKey,
			Type:    ctype.MsgType(file.Type),
		})
	}
	return &result
}

// 辅助函数：将数据库文件信息转换为响应结构
func convertToCreateMomentFileInfo(files *moment_models.Files) []types.CreateMomentFileInfo {
	if files == nil {
		return []types.CreateMomentFileInfo{}
	}

	var result []types.CreateMomentFileInfo
	for _, file := range *files {
		result = append(result, types.CreateMomentFileInfo{
			FileKey: file.FileKey,
			Type:    uint32(file.Type),
		})
	}
	return result
}
