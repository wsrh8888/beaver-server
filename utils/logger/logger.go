package logger

import (
	"sync"

	"beaver/utils/logger/model"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	sourceMu sync.RWMutex
	source   string
)

// Init 每个服务启动时调用一次，source 标识当前微服务（如 auth_api、chat_api）
func Init(serviceSource string) {
	sourceMu.Lock()
	source = serviceSource
	sourceMu.Unlock()
}

func currentSource() string {
	sourceMu.RLock()
	defer sourceMu.RUnlock()
	return source
}

// Logger module 为服务内子模块，source 由 Init 统一设置
type Logger struct {
	module string
}

func New(module string) *Logger {
	return &Logger{module: module}
}

func (l *Logger) Info(msg model.LogMsg) {
	l.send("info", msg)
}

func (l *Logger) Warn(msg model.LogMsg) {
	l.send("warn", msg)
}

func (l *Logger) Error(msg model.LogMsg) {
	l.send("error", msg)
}

func (l *Logger) send(level string, msg model.LogMsg) {
	fields := []logx.LogField{
		logx.Field("source", currentSource()),
	}
	if l.module != "" {
		fields = append(fields, logx.Field("module", l.module))
	}
	if msg.Data != nil {
		fields = append(fields, logx.Field("data", msg.Data))
	}

	switch level {
	case "warn":
		logx.Sloww(msg.Text, fields...)
	case "error":
		logx.Errorw(msg.Text, fields...)
	default:
		logx.Infow(msg.Text, fields...)
	}
}
