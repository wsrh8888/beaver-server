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

type SaveDigitalConfigLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

// 保存数字伙伴编排配置
func NewSaveDigitalConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveDigitalConfigLogic {
	return &SaveDigitalConfigLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("SaveDigitalConfig", ctx),
	}
}

func (l *SaveDigitalConfigLogic) SaveDigitalConfig(req *types.SaveDigitalConfigReq) (resp *types.SaveDigitalConfigRes, err error) {
	if req.DigitalID == "" {
		return nil, errors.New("digitalId 不能为空")
	}
	if _, err := loadOwnedDigital(l.svcCtx.DB, req.UserID, req.DigitalID); err != nil {
		return nil, err
	}

	// 查数字伙伴主表取 Runtime，按 Runtime 校验对应配置非空
	var digital open_models.OpenDigital
	if err := l.svcCtx.DB.Where("digital_id = ?", req.DigitalID).First(&digital).Error; err != nil {
		return nil, errors.New("数字伙伴不存在")
	}

	cfg := open_models.OpenDigitalConfig{
		DigitalID: req.DigitalID,
	}
	switch digital.Runtime {
	case open_models.DigitalRuntimeBuiltin:
		if req.BuiltinConfig == "" {
			return nil, errors.New("builtinConfig 不能为空")
		}
		var bc open_models.BuiltinConfig
		if err := json.Unmarshal([]byte(req.BuiltinConfig), &bc); err != nil {
			return nil, errors.New("builtinConfig JSON 格式错误")
		}
		cfg.BuiltinConfig = &bc
	case open_models.DigitalRuntimeExternal:
		if req.ExternalConfig == "" {
			return nil, errors.New("externalConfig 不能为空")
		}
		var ec open_models.ExternalConfig
		if err := json.Unmarshal([]byte(req.ExternalConfig), &ec); err != nil {
			return nil, errors.New("externalConfig JSON 格式错误")
		}
		cfg.ExternalConfig = &ec
	case open_models.DigitalRuntimeWorkflow:
		if req.WorkflowConfig == "" {
			return nil, errors.New("workflowConfig 不能为空")
		}
		var wc open_models.WorkflowConfig
		if err := json.Unmarshal([]byte(req.WorkflowConfig), &wc); err != nil {
			return nil, errors.New("workflowConfig JSON 格式错误")
		}
		cfg.WorkflowConfig = &wc
	default:
		return nil, errors.New("runtime 无效")
	}

	// upsert：有则更新，无则插入
	var existing open_models.OpenDigitalConfig
	if err := l.svcCtx.DB.Where("digital_id = ?", req.DigitalID).First(&existing).Error; err == nil {
		// 更新（只更新对应 Runtime 的配置列）
		updates := map[string]any{}
		switch digital.Runtime {
		case open_models.DigitalRuntimeBuiltin:
			updates["builtin_config"] = cfg.BuiltinConfig
		case open_models.DigitalRuntimeExternal:
			updates["external_config"] = cfg.ExternalConfig
		case open_models.DigitalRuntimeWorkflow:
			updates["workflow_config"] = cfg.WorkflowConfig
		}
		if err := l.svcCtx.DB.Model(&existing).Updates(updates).Error; err != nil {
			return nil, errors.New("保存配置失败")
		}
	} else {
		if err := l.svcCtx.DB.Create(&cfg).Error; err != nil {
			return nil, errors.New("保存配置失败")
		}
	}

	return &types.SaveDigitalConfigRes{}, nil
}
