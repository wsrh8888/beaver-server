package corerocketmq

import (
	"beaver/common/const/mqwsconst"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

// Message RocketMQ 通用消息体
type Message struct {
	MessageID string                 `json:"messageId"` // 消息唯一ID
	Timestamp int64                  `json:"timestamp"` // 时间戳
	Payload   map[string]interface{} `json:"payload"`   // 业务数据
}

// Client RocketMQ 客户端
type Client struct {
	Producer rocketmq.Producer
	Consumer rocketmq.PushConsumer
}

// NewClient 创建 RocketMQ 客户端
func InitRocketMQ(addr string) *Client {
	// 创建生产者
	p, err := rocketmq.NewProducer(
		producer.WithNameServer([]string{addr}),
		producer.WithRetry(3),
		producer.WithSendMsgTimeout(10*time.Second),
	)
	if err != nil {
		logx.Errorf("创建 Producer 失败: %v", err)
		return nil
	}

	err = p.Start()
	if err != nil {
		logx.Errorf("启动 Producer 失败: %v", err)
		return nil
	}

	client := &Client{
		Producer: p,
	}

	logx.Infof("RocketMQ Producer 启动成功, Addr: %s", addr)
	return client
}

// generateMessageID 生成消息唯一ID
func generateMessageID() string {
	return uuid.New().String()
}

// SendMessage 发送消息到 MQ
func (c *Client) SendMessage(ctx context.Context, topic mqwsconst.TopicType, payload map[string]interface{}) error {
	msg := &Message{
		MessageID: generateMessageID(),
		Timestamp: time.Now().UnixMilli(),
		Payload:   payload,
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("序列化消息失败: %v", err)
	}

	result, err := c.Producer.SendSync(ctx, primitive.NewMessage(string(topic), body))
	if err != nil {
		return fmt.Errorf("发送消息失败: %v", err)
	}

	if result.Status != primitive.SendOK {
		return fmt.Errorf("消息发送状态异常: %v", result.Status)
	}

	logx.Debugf("消息发送成功, Topic: %s, MsgID: %s", topic, result.MsgID)
	return nil
}

// RegisterConsumer 注册消费者
func (c *Client) RegisterConsumer(group mqwsconst.GroupType, addr string, topic mqwsconst.TopicType, handler func(msg *Message) error) error {
	c2, err := rocketmq.NewPushConsumer(
		consumer.WithGroupName(string(group)),
		consumer.WithNameServer([]string{addr}),
		consumer.WithConsumerModel(consumer.Clustering), // 集群模式，负载均衡
	)
	if err != nil {
		return fmt.Errorf("创建 Consumer 失败: %v", err)
	}

	err = c2.Subscribe(string(topic), consumer.MessageSelector{}, func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for _, msg := range msgs {
			var mqMsg Message
			if err := json.Unmarshal(msg.Body, &mqMsg); err != nil {
				logx.Errorf("反序列化消息失败: %v", err)
				continue
			}

			if err := handler(&mqMsg); err != nil {
				logx.Errorf("处理消息失败: %v", err)
				return consumer.ConsumeRetryLater, nil
			}
		}
		return consumer.ConsumeSuccess, nil
	})

	if err != nil {
		return fmt.Errorf("订阅 Topic 失败: %v", err)
	}

	err = c2.Start()
	if err != nil {
		return fmt.Errorf("启动 Consumer 失败: %v", err)
	}

	c.Consumer = c2
	logx.Infof("RocketMQ Consumer 启动成功, Group: %s, Topic: %s", group, topic)
	return nil
}

// Shutdown 关闭客户端
func (c *Client) Shutdown() {
	if c.Producer != nil {
		c.Producer.Shutdown()
	}
	if c.Consumer != nil {
		c.Consumer.Shutdown()
	}
	logx.Info("RocketMQ 客户端已关闭")
}
