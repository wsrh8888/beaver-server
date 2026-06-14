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


# ==================== AUTH 相关服务 ====================
# auth_api
RUN go build -ldflags="-s -w" -o auth_api/auth app/auth/auth_api/auth.go

# ==================== BACKEND 相关服务 ====================
# backend_admin
RUN go build -ldflags="-s -w" -o backend_admin/backend_admin app/backend/backend_admin/backend.go

# ==================== CHAT 相关服务 ====================
# chat_api
RUN go build -ldflags="-s -w" -o chat_api/chat app/chat/chat_api/chat.go
# chat_rpc
RUN go build -ldflags="-s -w" -o chat_rpc/chatrpc app/chat/chat_rpc/chatrpc.go

# ==================== DATASYNC 相关服务 ====================
# datasync_api
RUN go build -ldflags="-s -w" -o datasync_api/datasync app/datasync/datasync_api/datasync.go

# ==================== PLATFORM 相关服务 ====================
# platform_api
RUN go build -ldflags="-s -w" -o platform_api/platform app/platform/platform_api/platform.go

# ==================== EMOJI 相关服务 ====================
# emoji_api
RUN go build -ldflags="-s -w" -o emoji_api/emoji app/emoji/emoji_api/emoji.go
# emoji_rpc
RUN go build -ldflags="-s -w" -o emoji_rpc/emojirpc app/emoji/emoji_rpc/emojirpc.go

# ==================== FILE 相关服务 ====================
# file_api
RUN go build -ldflags="-s -w" -o file_api/file app/file/file_api/file.go
# file_rpc
RUN go build -ldflags="-s -w" -o file_rpc/filerpc app/file/file_rpc/filerpc.go

# ==================== FRIEND 相关服务 ====================
# friend_api
RUN go build -ldflags="-s -w" -o friend_api/friend app/friend/friend_api/friend.go
# friend_rpc
RUN go build -ldflags="-s -w" -o friend_rpc/friendrpc app/friend/friend_rpc/friendrpc.go

# ==================== GATEWAY 相关服务 ====================
# gateway_api
RUN go build -ldflags="-s -w" -o gateway_api/gateway app/gateway/gateway_api/gateway.go
# gateway_admin
RUN go build -ldflags="-s -w" -o gateway_admin/gateway_admin app/gateway/gateway_admin/gateway.go

# ==================== GROUP 相关服务 ====================
# group_api
RUN go build -ldflags="-s -w" -o group_api/group app/group/group_api/group.go
# group_rpc
RUN go build -ldflags="-s -w" -o group_rpc/grouprpc app/group/group_rpc/grouprpc.go

# ==================== MOMENT 相关服务 ====================
# moment_api
RUN go build -ldflags="-s -w" -o moment_api/moment app/moment/moment_api/moment.go

# ==================== DOCUMENT 相关服务 ====================
# document_api
RUN go build -ldflags="-s -w" -o document_api/document app/document/document_api/document.go

# ==================== NOTIFICATION 相关服务 ====================
# notification_api
RUN go build -ldflags="-s -w" -o notification_api/notification app/notification/notification_api/notification.go
# notification_rpc
RUN go build -ldflags="-s -w" -o notification_rpc/notificationrpc app/notification/notification_rpc/notificationrpc.go


# ==================== USER 相关服务 ====================
# user_api
RUN go build -ldflags="-s -w" -o user_api/user app/user/user_api/user.go
# user_rpc
RUN go build -ldflags="-s -w" -o user_rpc/userrpc app/user/user_rpc/userrpc.go


# ==================== CALL 相关服务 ====================
# call_api
RUN go build -ldflags="-s -w" -o call_api/call app/call/call_api/call.go
# call_rpc
RUN go build -ldflags="-s -w" -o call_rpc/callrpc app/call/call_rpc/callrpc.go


# ==================== WS 相关服务 ====================
# ws_api
RUN go build -ldflags="-s -w" -o ws_api/ws app/ws/ws_api/ws.go

