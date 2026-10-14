package digital

import (
	"context"
	"encoding/json"
	"errors"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"

	beaverlog "beaver/utils/beaverlog"
)

// goctl 会把 logx 写进 imports，占位避免未使用导入。
var _ = logx.WithContext

type GetDigitalConfigLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

// 获取数字伙伴编排配置
func NewGetDigitalConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDigitalConfigLogic {
	return &GetDigitalConfigLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("GetDigitalConfig", ctx),
	}
}

func (l *GetDigitalConfigLogic) GetDigitalConfig(req *types.GetDigitalConfigReq) (resp *types.GetDigitalConfigRes, err error) {
	if req.DigitalID == "" {
		return nil, errors.New("digitalId 不能为空")
	}
	if _, err := loadOwnedDigital(l.svcCtx.DB, req.UserID, req.DigitalID); err != nil {
		return nil, err
	}

	// 查数字伙伴主表取 Runtime
	var digital open_models.OpenDigital
	if err := l.svcCtx.DB.Where("digital_id = ?", req.DigitalID).First(&digital).Error; err != nil {
		return nil, errors.New("数字伙伴不存在")
	}

	resp = &types.GetDigitalConfigRes{
		DigitalID: req.DigitalID,
		Runtime:   digital.Runtime,
	}

	// 查编排配置，按 Runtime 填对应字段
	var cfg open_models.OpenDigitalConfig
	if err := l.svcCtx.DB.Where("digital_id = ?", req.DigitalID).First(&cfg).Error; err == nil {
		switch digital.Runtime {
		case open_models.DigitalRuntimeBuiltin:
			if cfg.BuiltinConfig != nil {
				if b, mErr := json.Marshal(cfg.BuiltinConfig); mErr == nil {
					resp.BuiltinConfig = string(b)
				}
			}
		case open_models.DigitalRuntimeExternal:
			if cfg.ExternalConfig != nil {
				if b, mErr := json.Marshal(cfg.ExternalConfig); mErr == nil {
					resp.ExternalConfig = string(b)
				}
			}
		case open_models.DigitalRuntimeWorkflow:
			if cfg.WorkflowConfig != nil {
				if b, mErr := json.Marshal(cfg.WorkflowConfig); mErr == nil {
					resp.WorkflowConfig = string(b)
				}
			}
		}
	}

	return resp, nil
}
