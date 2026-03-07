package logic

import (
	"context"
	"encoding/json"

	"beaver/app/chat/chat_api/internal/svc"
	"beaver/app/chat/chat_api/internal/types"
	"beaver/app/chat/chat_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetForwardDetailsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetForwardDetailsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetForwardDetailsLogic {
	return &GetForwardDetailsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetForwardDetailsLogic) GetForwardDetails(req *types.GetForwardDetailsReq) (resp *types.GetForwardDetailsRes, err error) {
	var detail chat_models.ChatForward
	err = l.svcCtx.DB.Where("record_id = ?", req.RecordID).First(&detail).Error
	if err != nil {
		l.Logger.Errorf("获取合并转发详情失败: %v", err)
		return nil, err
	}

	// 因为使用了自定义类型 ForwardContent，GORM 的 Scan 已经帮我们把 JSON 转成了结构体
	// 我们需要将其转换为 API 定义的消息类型
	var list []types.Message
	for _, m := range detail.Content {
		// 这里需要将 chat_models.ChatMessage 转换为 types.Message
		// 简单起见，这里假设你可以通过 JSON 转换或者手动构建
		msgJSON, _ := json.Marshal(m)
		var tMsg types.Message
		json.Unmarshal(msgJSON, &tMsg)
		list = append(list, tMsg)
	}

	return &types.GetForwardDetailsRes{
		List: list,
	}, nil
}
