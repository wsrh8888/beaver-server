package etcd

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/netx"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// 获取服务地址
func GetServiceAddr(etcdAddr string, serviceName string) string {
	client, err := initEtcd(etcdAddr)
	if err != nil {
		logx.Errorf("初始化etcd客户端失败: %v", err)
		return ""
	}
	defer client.Close()

	// 拼接获取键的前缀
	keyPrefix := fmt.Sprintf("/services/%s", serviceName)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.Get(ctx, keyPrefix, clientv3.WithPrefix())
	if err != nil {
		logx.Errorf("获取etcd地址失败: %v", err)
		return ""
	}
	if len(resp.Kvs) == 0 {
		logx.Errorf("获取etcd地址失败，没有找到服务: %s", serviceName)
		return ""
	}

	// for _, kv := range resp.Kvs {
	// 	logx.Infof("键: %s, 值: %s", kv.Key, kv.Value)
	// }

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	selectedInstance := resp.Kvs[r.Intn(len(resp.Kvs))]
	return string(selectedInstance.Value)
}

// 初始化etcd客户端
func initEtcd(etcdAddr string) (*clientv3.Client, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{etcdAddr},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}

// 注册服务地址到etcd，并设置租约
func DeliveryAddress(etcdAddr string, serviceName string, addr string) {
	list := strings.Split(addr, ":")
	if len(list) != 2 {
		logx.Error("addr格式错误")
		return
	}
	ip := list[0]
	if ip == "0.0.0.0" {
		ip = netx.InternalIp()
		addr = strings.ReplaceAll(addr, "0.0.0.0", ip)
	}

	client, err := initEtcd(etcdAddr)
	if err != nil {
		logx.Error("初始化etcd客户端失败", err)
		return
	}
	defer client.Close()

	err = registerServiceWithRetry(client, etcdAddr, serviceName, addr)
	if err != nil {
		logx.Error("服务注册失败", err)
	}
}

// 注册服务，并在失败时自动重试
func registerServiceWithRetry(client *clientv3.Client, etcdAddr, serviceName, addr string) error {
	for {
		err := registerService(client, etcdAddr, serviceName, addr)
		if err == nil {
			return nil
		}
		logx.Error("注册服务失败，重试中...", err)
		time.Sleep(2 * time.Second)
	}
}

// 注册服务到etcd，并设置租约
func registerService(client *clientv3.Client, etcdAddr, serviceName, addr string) error {
	// 创建租约，TTL为10秒
	resp, err := client.Grant(context.Background(), 10)
	if err != nil {
		return fmt.Errorf("创建租约失败: %w", err)
	}

	// 使用服务名作为前缀，实例唯一标示作为键值
	key := fmt.Sprintf("/services/%s/%s", serviceName, addr)
	fmt.Println("key:", key)

	// 写入键值并附加租约
	_, err = client.Put(context.Background(), key, addr, clientv3.WithLease(resp.ID))
	if err != nil {
		return fmt.Errorf("上送etcd地址失败: %w", err)
	}
	logx.Info("上送etcd地址成功", addr)
	fmt.Println("resp.ID:", resp.ID)

	// 自动续约
	ch, kaerr := client.KeepAlive(context.Background(), resp.ID)
	if kaerr != nil {
		return fmt.Errorf("设置租约续约失败: %w", kaerr)
	}
	fmt.Println("续约通道:", ch)

	go keepAliveService(client, etcdAddr, serviceName, addr, ch)

	return nil
}

// 处理服务续约
func keepAliveService(client *clientv3.Client, etcdAddr, serviceName, addr string, ch <-chan *clientv3.LeaseKeepAliveResponse) {
	for {
		select {
		case ka, ok := <-ch:
			if !ok {
				logx.Error("续约频道关闭")
				reconnectEtcd(client, etcdAddr, serviceName, addr)
				return
			}
			if ka == nil {
				logx.Error("租约失效")
				reconnectEtcd(client, etcdAddr, serviceName, addr)
				return
			} else {
				// logx.Infof("续约成功: lease ID=%d", ka.ID)
			}
		}
	}
}

// 重新建立连接并注册服务
func reconnectEtcd(client *clientv3.Client, etcdAddr, serviceName, addr string) {
	logx.Info("尝试重新连接etcd并注册服务")
	client.Close()

	for {
		newClient, err := initEtcd(etcdAddr)
		if err == nil {
			client = newClient
			break
		}
		logx.Error("重新初始化etcd客户端失败", err)
		time.Sleep(2 * time.Second)
	}

	registerServiceWithRetry(client, etcdAddr, serviceName, addr)
}
