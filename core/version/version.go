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
func (vg *VersionGenerator) GetNextVersion(dataType string) (int64, error) {
	key := fmt.Sprintf("version:%s", dataType)

	// 1. 先尝试从Redis获取
	version, err := vg.redisClient.Incr(key).Result()
	if err != nil {
		logx.Errorf("Redis获取版本号失败: dataType=%s, error=%v", dataType, err)
		// Redis失败，从MySQL获取
		return vg.getVersionFromMySQL(dataType)
	}

	// 2. 检查Redis版本号是否合理（防止Redis重启后从0开始）
	if version == 1 {
		// 可能是Redis重启，从MySQL同步最新版本号
		mysqlVersion, err := vg.getMaxVersionFromMySQL(dataType)
		if err != nil {
			logx.Errorf("MySQL获取版本号失败: dataType=%s, error=%v", dataType, err)
			// MySQL查询失败，继续使用Redis的版本号1
			return version, nil
		}

		if mysqlVersion > 0 {
			// MySQL有数据，更新Redis为MySQL版本号+1
			newVersion := mysqlVersion + 1
			vg.redisClient.Set(key, newVersion, 0) // 永不过期
			logx.Infof("从MySQL同步版本号: dataType=%s, mysqlVersion=%d, newVersion=%d", dataType, mysqlVersion, newVersion)
			return newVersion, nil
		}

		// MySQL没有数据，继续使用Redis的版本号1
		logx.Infof("MySQL没有数据，Redis从1开始: dataType=%s, version=%d", dataType, version)
	}

	logx.Infof("生成版本号: dataType=%s, version=%d", dataType, version)
	return version, nil
}

// getVersionFromMySQL 从MySQL获取版本号（Redis失败时的备用方案）
func (vg *VersionGenerator) getVersionFromMySQL(dataType string) (int64, error) {
	var maxVersion int64
	err := vg.db.Raw("SELECT COALESCE(MAX(version), 0) FROM " + vg.getTableName(dataType)).Scan(&maxVersion).Error
	if err != nil {
		logx.Errorf("MySQL获取版本号失败: dataType=%s, error=%v", dataType, err)
		return 1, fmt.Errorf("无法获取版本号: %v", err)
	}

	// 返回MySQL最大版本号+1，确保至少从1开始
	nextVersion := maxVersion + 1
	if nextVersion < 1 {
		nextVersion = 1
	}

	logx.Infof("从MySQL获取版本号: dataType=%s, maxVersion=%d, nextVersion=%d", dataType, maxVersion, nextVersion)
	return nextVersion, nil
}

// getMaxVersionFromMySQL 从MySQL获取最大版本号
func (vg *VersionGenerator) getMaxVersionFromMySQL(dataType string) (int64, error) {
	var maxVersion int64
	err := vg.db.Raw("SELECT COALESCE(MAX(version), 0) FROM " + vg.getTableName(dataType)).Scan(&maxVersion).Error
	if err != nil {
		return 0, err
	}
	return maxVersion, nil
}

// getTableName 根据数据类型获取表名
func (vg *VersionGenerator) getTableName(dataType string) string {
	switch dataType {
	case "users":
		return "user_models"
	case "friends":
		return "friend_models"
	case "groups":
		return "group_models"
	case "chats":
		return "chat_models"
	default:
		return dataType
	}
}

// GetCurrentVersion 获取当前版本号（不分配新的）
func (vg *VersionGenerator) GetCurrentVersion(dataType string) (int64, error) {
	key := fmt.Sprintf("version:%s", dataType)

	result, err := vg.redisClient.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			// Redis没有，从MySQL获取
			maxVersion, err := vg.getMaxVersionFromMySQL(dataType)
			if err != nil {
				logx.Errorf("获取当前版本号失败: dataType=%s, error=%v", dataType, err)
				return 0, err
			}
			logx.Infof("Redis无数据，从MySQL获取当前版本号: dataType=%s, currentVersion=%d", dataType, maxVersion)
			return maxVersion, nil
		}
		return 0, err
	}

	// 解析版本号
	var version int64
	fmt.Sscanf(result, "%d", &version)
	logx.Infof("从Redis获取当前版本号: dataType=%s, currentVersion=%d", dataType, version)
	return version, nil
}
