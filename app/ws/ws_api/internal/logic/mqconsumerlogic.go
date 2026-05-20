package logic

import (
	"context"
	"encoding/json"

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

// StartConsumer 启动 RocketMQ 消费者
func (l *MqConsumerLogic) StartConsumer() error {
	mqClient := l.svcCtx.RocketMQ
	if mqClient == nil {
		logx.Error("RocketMQ 客户端未初始化")
		return nil
	}

	handler := func(msg *corerocketmq.Message) error {
		bodyBytes, err := json.Marshal(msg.Payload)
		if err != nil {
			logx.Errorf("序列化消息体失败: %v", err)
			return err
		}

		content := type_struct.WsContent{
			Data: type_struct.WsData{
				Type:           wsTypeConst.Type(msg.Payload["type"].(string)),
				Body:           bodyBytes,
				ConversationID: msg.Payload["conversationId"].(string),
			},
		}

		ws_conn.SendMsgToUser(msg.Payload["targetId"].(string), wsCommandConst.Command(msg.Payload["command"].(string)), content)
		return nil
	}

	err := mqClient.RegisterConsumer(
		"ws_api_consumer_group",
		l.svcCtx.Config.RocketMQ.Addr,
		mqwsconst.MqTopicWs,
		handler,
	)

	if err != nil {
		logx.Errorf("启动 RocketMQ Consumer 失败: %v", err)
		return err
	}

	logx.Info("WS API RocketMQ Consumer 启动成功")
	return nil
}
