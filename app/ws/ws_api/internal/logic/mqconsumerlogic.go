package logic

import (
	"context"
	"encoding/json"
	"fmt"

	ws_conn "beaver/app/ws/ws_api/internal/logic/websocket/conn"
	"beaver/app/ws/ws_api/internal/svc"
	type_struct "beaver/app/ws/ws_api/types"
	"beaver/common/const/mqwsconst"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"
	"beaver/core/corerocketmq"

	"github.com/zeromicro/go-zero/core/logx"
)

type MqConsumerLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMqConsumerLogic(svcCtx *svc.ServiceContext) *MqConsumerLogic {
	return &MqConsumerLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func payloadString(payload map[string]interface{}, key string) (string, error) {
	v, ok := payload[key]
	if !ok || v == nil {
		return "", fmt.Errorf("缺少字段 %s", key)
	}
	s, ok := v.(string)
	if !ok || s == "" {
		return "", fmt.Errorf("字段 %s 无效", key)
	}
	return s, nil
}

// StartConsumer 启动 RocketMQ 消费者
func (l *MqConsumerLogic) StartConsumer() error {
	mqClient := l.svcCtx.RocketMQ
	if mqClient == nil {
		logx.Error("RocketMQ 客户端未初始化")
		return nil
	}

	handler := func(msg *corerocketmq.Message) error {
		targetID, err := payloadString(msg.Payload, "targetId")
		if err != nil {
			logx.Errorf("MQ 消息格式错误: %v", err)
			return nil
		}
		command, err := payloadString(msg.Payload, "command")
		if err != nil {
			logx.Errorf("MQ 消息格式错误: %v", err)
			return nil
		}
		msgType, err := payloadString(msg.Payload, "type")
		if err != nil {
			logx.Errorf("MQ 消息格式错误: %v", err)
			return nil
		}
		conversationID, err := payloadString(msg.Payload, "conversationId")
		if err != nil {
			logx.Errorf("MQ 消息格式错误: %v", err)
			return nil
		}

		body, ok := msg.Payload["body"]
		if !ok || body == nil {
			logx.Error("MQ 消息缺少 body 字段")
			return nil
		}
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			logx.Errorf("序列化 body 失败: %v", err)
			return err
		}

		content := type_struct.WsContent{
			Data: type_struct.WsData{
				Type:           wsTypeConst.Type(msgType),
				Body:           bodyBytes,
				ConversationID: conversationID,
			},
		}

		ws_conn.SendMsgToUser(targetID, wsCommandConst.Command(command), content)
		return nil
	}

	err := mqClient.RegisterConsumer(
		mqwsconst.MqGroupWs,
		l.svcCtx.Config.RocketMQ.Addr,
		mqwsconst.MqTopicWs,
		true, // 广播模式：每个 WS 实例都消费，才能推送到本机在线连接
		handler,
	)

	if err != nil {
		logx.Errorf("启动 RocketMQ Consumer 失败: %v", err)
		return err
	}

	logx.Info("WS API RocketMQ Consumer 启动成功")
	return nil
}
