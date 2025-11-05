package core

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

// VersionGenerator 版本号生成器（Redis缓存 + MySQL主存储）
type VersionGenerator struct {
	redisClient *redis.Client
	db          *gorm.DB
}

// NewVersionGenerator 创建版本号生成器
func NewVersionGenerator(redisClient *redis.Client, db *gorm.DB) *VersionGenerator {
	return &VersionGenerator{
		redisClient: redisClient,
		db:          db,
	}
}

// GetNextVersion 获取下一个版本号（Redis缓存 + MySQL主存储）
// 参数:
//   - table: 表名或业务标识符 (必传)
//   - field: 字段名，为空时表示全局版本
//   - value: 字段值，为空时表示全局版本
//   - defaultValue: 默认值，当Redis中没有时使用此值初始化，为nil时使用1
//
// 返回: 版本号，失败时返回-1
func (vg *VersionGenerator) GetNextVersion(table string, field string, value string, defaultValue *int64) int64 {
	var key string
	isConditional := field != "" && value != ""

	if isConditional {
		key = fmt.Sprintf("version:%s:%s:%s", table, field, value)
	} else {
		key = fmt.Sprintf("version:%s", table)
	}

	// 获取默认值
	initValue := int64(1) // 默认从1开始
	if defaultValue != nil && *defaultValue > 0 {
		initValue = *defaultValue
	}

	// 1. 先尝试从Redis获取
	version, err := vg.redisClient.Incr(key).Result()
	if err != nil {
		logx.Errorf("Redis获取版本号失败: table=%s, field=%s, value=%s, error=%v", table, field, value, err)
		// Redis失败，初始化Redis并返回初始值
		vg.redisClient.Set(key, initValue, 0) // 永不过期
		logx.Infof("初始化版本号: table=%s, field=%s, value=%s, initValue=%d", table, field, value, initValue)
		return initValue
	}

	// 2. 检查Redis版本号是否合理（防止Redis重启后从0开始）
	if version == 1 && !isConditional {
		// 只有全局版本才从MySQL同步，条件版本直接使用初始值
		mysqlVersion := vg.getMaxVersionFromMySQL(table, field, value)
		if mysqlVersion < 0 {
			logx.Errorf("MySQL获取版本号失败: table=%s", table)
			// MySQL查询失败，继续使用Redis的版本号1
			return version
		}

		if mysqlVersion > 0 {
			// MySQL有数据，更新Redis为MySQL版本号+1
			newVersion := mysqlVersion + 1
			vg.redisClient.Set(key, newVersion, 0) // 永不过期
			logx.Infof("从MySQL同步版本号: table=%s, mysqlVersion=%d, newVersion=%d", table, mysqlVersion, newVersion)
			return newVersion
		}

		// MySQL没有数据，使用初始值
		if initValue > 1 {
			vg.redisClient.Set(key, initValue, 0) // 永不过期
			logx.Infof("使用自定义初始值: table=%s, initValue=%d", table, initValue)
			return initValue
		}
		// MySQL没有数据，继续使用Redis的版本号1
		logx.Infof("MySQL没有数据，Redis从1开始: table=%s, version=%d", table, version)
	}

	logx.Infof("生成版本号: table=%s, field=%s, value=%s, version=%d", table, field, value, version)
	return version
}

// getVersionFromMySQL 从MySQL获取版本号（Redis失败时的备用方案）
// 对于条件版本，只返回1，不从MySQL获取，失败时返回-1
func (vg *VersionGenerator) getVersionFromMySQL(table string, field string, value string) int64 {
	isConditional := field != "" && value != ""
	// 检查是否为条件版本
	if isConditional {
		logx.Infof("条件版本从1开始: table=%s, field=%s, value=%s", table, field, value)
		return 1
	}

	var maxVersion int64
	err := vg.db.Raw("SELECT COALESCE(MAX(version), 0) FROM " + vg.getTableName(table)).Scan(&maxVersion).Error
	if err != nil {
		logx.Errorf("MySQL获取版本号失败: table=%s, error=%v", table, err)
		return -1
	}

	// 返回MySQL最大版本号+1，确保至少从1开始
	nextVersion := maxVersion + 1
	if nextVersion < 1 {
		nextVersion = 1
	}

	logx.Infof("从MySQL获取版本号: table=%s, maxVersion=%d, nextVersion=%d", table, maxVersion, nextVersion)
	return nextVersion
}

// getMaxVersionFromMySQL 从MySQL获取最大版本号
// 对于条件版本，返回0（表示没有历史版本），失败时返回-1
func (vg *VersionGenerator) getMaxVersionFromMySQL(table string, field string, value string) int64 {
	isConditional := field != "" && value != ""
	// 检查是否为条件版本
	if isConditional {
		logx.Infof("条件版本无历史数据: table=%s, field=%s, value=%s", table, field, value)
		return 0
	}

	var maxVersion int64
	err := vg.db.Raw("SELECT COALESCE(MAX(version), 0) FROM " + vg.getTableName(table)).Scan(&maxVersion).Error
	if err != nil {
		logx.Errorf("MySQL获取最大版本号失败: table=%s, error=%v", table, err)
		return -1
	}
	return maxVersion
}

// getTableName 根据数据类型获取表名
// 注意：条件版本不应该调用此方法，因为条件版本不从MySQL查询
func (vg *VersionGenerator) getTableName(dataType string) string {
	switch dataType {
	case "users":
		return "user_models"
	case "friends":
		return "friend_models"
	case "groups":
		return "group_models"
	case "group_members":
		return "group_member_models"
	case "group_member_logs":
		return "group_member_change_log_models"
	case "chats":
		return "chat_models"
	default:
		return dataType
	}
}
