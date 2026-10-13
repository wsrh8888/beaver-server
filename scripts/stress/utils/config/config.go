/*
 * 配置：从 scripts/stress/.env 读取所有服务地址与各自的 DB/Redis 配置。
 *
 * 每个服务的 DB 和 Redis 是独立的（不同库、不同密码、甚至不同实例），
 * 所以这里按服务分别持有各自的配置。
 *
 * .env 已 gitignore，密码不进 git、不外传。
 */

package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// ServiceConfig 一个服务的连接配置
type ServiceConfig struct {
	APIBase       string // API 地址（host:port）
	RedisAddr     string // 该服务用的 Redis 地址
	RedisPassword string // 该服务用的 Redis 密码
	RedisDB       int    // 该服务用的 Redis DB 号
}

// Config 持有所有服务的配置
type Config struct {
	Auth   ServiceConfig
	User   ServiceConfig
	Friend ServiceConfig
	Chat   ServiceConfig
	Group  ServiceConfig
	Moment ServiceConfig
}

var (
	loaded     *Config
	stressRoot string // .env 所在目录，即 scripts/stress 根
)

// Load 从 scripts/stress/.env 读取配置，只在首次调用时读文件，之后缓存
func Load() (*Config, error) {
	if loaded != nil {
		return loaded, nil
	}

	// 找 .env：脚本可能在 setup/、等子目录下跑，统一往上找 stress/.env
	wd, _ := os.Getwd()
	dir := wd
	for i := 0; i < 5; i++ {
		p := filepath.Join(dir, ".env")
		if _, err := os.Stat(p); err == nil {
			return loadFromFile(p)
		}
		dir = filepath.Dir(dir)
	}
	return nil, fmt.Errorf("找不到 .env 文件，请在 scripts/stress/.env 填写配置（可复制 .env.example）")
}

func loadFromFile(path string) (*Config, error) {
	env, err := parseEnv(path)
	if err != nil {
		return nil, err
	}
	stressRoot = filepath.Dir(path)

	c := &Config{
		Auth:   readService(env, "AUTH"),
		User:   readService(env, "USER"),
		Friend: readService(env, "FRIEND"),
		Chat:   readService(env, "CHAT"),
		Group:  readService(env, "GROUP"),
		Moment: readService(env, "MOMENT"),
	}

	// 基本校验：至少 auth 要有
	if c.Auth.APIBase == "" {
		return nil, fmt.Errorf(".env 中 AUTH_API_BASE 不能为空")
	}
	loaded = c
	return c, nil
}

// AccountsPath 返回 accounts.json 的绝对路径（与 .env 同目录）。
// 必须先调用 Load()，否则 stressRoot 未初始化。
func AccountsPath() string {
	return dataPath("accounts.json")
}

// FriendsPath 返回 friends.json 的绝对路径（好友关系 + 会话ID）。
// 由 setup/friends 产出，供 chat 压测取 conversationId。
func FriendsPath() string {
	return dataPath("friends.json")
}

// dataPath 返回 stressRoot 下的数据文件路径
func dataPath(name string) string {
	return filepath.Join(stressRoot, name)
}

func parseEnv(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("打开 .env 失败: %w", err)
	}
	defer f.Close()

	env := map[string]string{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		val = strings.Trim(val, `"'`) // 容忍 .env 里值带引号
		env[key] = val
	}
	return env, nil
}

// readService 按 <PREFIX> 读取一个服务的 API_BASE / REDIS_ADDR / REDIS_PASSWORD / REDIS_DB
func readService(env map[string]string, prefix string) ServiceConfig {
	db := 0
	if v, ok := env[prefix+"_REDIS_DB"]; ok {
		db, _ = strconv.Atoi(v)
	}
	return ServiceConfig{
		APIBase:       env[prefix+"_API_BASE"],
		RedisAddr:     env[prefix+"_REDIS_ADDR"],
		RedisPassword: env[prefix+"_REDIS_PASSWORD"],
		RedisDB:       db,
	}
}
