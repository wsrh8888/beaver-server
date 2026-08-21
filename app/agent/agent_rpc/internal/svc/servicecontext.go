package svc

import (
	"beaver/app/agent/agent_models"
	"beaver/app/agent/agent_rpc/internal/config"
	"beaver/core/coregorm"

	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := coregorm.InitGorm(c.Mysql.DataSource)
	_ = mysqlDb.AutoMigrate(&agent_models.Agent{}, &agent_models.AgentMessage{})

	return &ServiceContext{
		Config: c,
		DB:     mysqlDb,
	}
}
