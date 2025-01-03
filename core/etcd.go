package core

import (
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/zeromicro/go-zero/core/logx"
)

/**
 * @description: 初始化etcd
 */
func InitEtcd(add string) *clientv3.Client {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{add},
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		logx.Error("etcd连接失败", err)
		panic(err)
	}
	return cli
}
