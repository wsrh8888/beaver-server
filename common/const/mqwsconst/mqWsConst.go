package mqwsconst

// RocketMQ Topic 类型（用于 WebSocket 推送）
type TopicType string

// RocketMQ Topic 常量
const (
	// MqTopicWs WebSocket 推送专用 Topic
	MqTopicWs TopicType = "ws_push_topic"
)

// RocketMQ Consumer Group 类型
type GroupType string

// RocketMQ Consumer Group 常量
const (
	// MqGroupWs WS API 消费者组
	MqGroupWs GroupType = "ws_api_consumer_group"
)
