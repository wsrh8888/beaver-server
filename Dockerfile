# 构建阶段
FROM golang:alpine AS builder
LABEL stage=gobuilder

ENV CGO_ENABLED=0
ENV GOPROXY=https://goproxy.cn,direct

# 安装必要的软件包
RUN apk add --no-cache tzdata

WORKDIR /build

COPY go.mod .
COPY go.sum .

COPY . .
RUN go mod tidy


# ==================== ACCOUNT 相关服务（user + auth 合并）====================
# account_api
RUN go build -ldflags="-s -w" -o account_api/account app/account/account_api/account.go
# account_rpc
RUN go build -ldflags="-s -w" -o account_rpc/accountrpc app/account/account_rpc/accountrpc.go

# ==================== ASSET 相关服务（file + emoji 合并）====================
# asset_api
RUN go build -ldflags="-s -w" -o asset_api/asset app/asset/asset_api/asset.go
# asset_rpc
RUN go build -ldflags="-s -w" -o asset_rpc/assetrpc app/asset/asset_rpc/assetrpc.go

# ==================== DATASYNC 相关服务 ====================
# datasync_api
RUN go build -ldflags="-s -w" -o datasync_api/datasync app/datasync/datasync_api/datasync.go

# ==================== GATEWAY 相关服务 ====================
# gateway_api
RUN go build -ldflags="-s -w" -o gateway_api/gateway app/gateway/gateway_api/gateway.go

# ==================== MESSAGE 相关服务（chat 更名）====================
# message_api
RUN go build -ldflags="-s -w" -o message_api/message app/message/message_api/message.go
# message_rpc
RUN go build -ldflags="-s -w" -o message_rpc/messagerpc app/message/message_rpc/messagerpc.go

# ==================== NOTIFICATION 相关服务 ====================
# notification_api
RUN go build -ldflags="-s -w" -o notification_api/notification app/notification/notification_api/notification.go
# notification_rpc
RUN go build -ldflags="-s -w" -o notification_rpc/notificationrpc app/notification/notification_rpc/notificationrpc.go

# ==================== PLATFORM 相关服务 ====================
# platform_api
RUN go build -ldflags="-s -w" -o platform_api/platform app/platform/platform_api/platform.go
# platform_rpc
RUN go build -ldflags="-s -w" -o platform_rpc/platformrpc app/platform/platform_rpc/platformrpc.go

# ==================== RELATION 相关服务（friend + group 合并）====================
# relation_api
RUN go build -ldflags="-s -w" -o relation_api/relation app/relation/relation_api/relation.go
# relation_rpc
RUN go build -ldflags="-s -w" -o relation_rpc/relationrpc app/relation/relation_rpc/relationrpc.go

# ==================== WS 相关服务 ====================
# ws_api
RUN go build -ldflags="-s -w" -o ws_api/ws app/ws/ws_api/ws.go
