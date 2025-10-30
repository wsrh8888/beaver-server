package core

import (
	"github.com/BurntSushi/toml"      // 用于解析 toml 配置文件
	"github.com/pion/ion-sfu/pkg/sfu" // 导入 ion-sfu 包
)

func InitSFU(configFile string) *sfu.SFU {
	var cfg sfu.Config
	if _, err := toml.DecodeFile(configFile, &cfg); err != nil {
		panic(err) // 处理配置文件解析错误
	}
	return sfu.NewSFU(cfg)
}
