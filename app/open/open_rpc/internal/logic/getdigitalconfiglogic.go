package logic

import (
	"context"
	"encoding/json"
	"errors"

	"beaver/app/open/open_models"
	"beaver/app/open/open_rpc/internal/svc"
	"beaver/app/open/open_rpc/types/open_rpc"

	beaverlog "beaver/utils/beaverlog"
)

type GetDigitalConfigLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewGetDigitalConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDigitalConfigLogic {
	return &GetDigitalConfigLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("GetDigitalConfigLogic", ctx),
	}
}

// GetDigitalConfig agent 运行时来取数字伙伴编排配置。
// 按 OpenDigital.Runtime 返回对应模式的配置 JSON，其余留空。
func (l *GetDigitalConfigLogic) GetDigitalConfig(in *open_rpc.GetDigitalConfigReq) (*open_rpc.GetDigitalConfigRes, error) {
	if in.DigitalId == "" {
		return nil, errors.New("digital_id 不能为空")
	}

	// 查数字伙伴主表取 Runtime
	var digital open_models.OpenDigital
	if err := l.svcCtx.DB.Where("digital_id = ?", in.DigitalId).First(&digital).Error; err != nil {
		return nil, errors.New("数字伙伴不存在")
	}

	res := &open_rpc.GetDigitalConfigRes{
		DigitalId: in.DigitalId,
		Runtime:   digital.Runtime,
	}

	// 查编排配置
	var cfg open_models.OpenDigitalConfig
	if err := l.svcCtx.DB.Where("digital_id = ?", in.DigitalId).First(&cfg).Error; err == nil {
		// 按 Runtime 序列化对应配置 JSON，其余留空
		switch digital.Runtime {
		case open_models.DigitalRuntimeBuiltin:
			if cfg.BuiltinConfig != nil {
				if b, err := json.Marshal(cfg.BuiltinConfig); err == nil {
					res.BuiltinConfig = string(b)
				}
			}
		case open_models.DigitalRuntimeExternal:
			if cfg.ExternalConfig != nil {
				if b, err := json.Marshal(cfg.ExternalConfig); err == nil {
					res.ExternalConfig = string(b)
				}
			}
		case open_models.DigitalRuntimeWorkflow:
			if cfg.WorkflowConfig != nil {
				if b, err := json.Marshal(cfg.WorkflowConfig); err == nil {
					res.WorkflowConfig = string(b)
				}
			}
		}
	}

	return res, nil
}
