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

# ==================== CHAT 相关服务 ====================
# chat_api
RUN go build -ldflags="-s -w" -o chat_api/chat app/chat/chat_api/chat.go
# chat_rpc
RUN go build -ldflags="-s -w" -o chat_rpc/chatrpc app/chat/chat_rpc/chatrpc.go

# ==================== FRIEND 相关服务 ====================
# friend_api
RUN go build -ldflags="-s -w" -o friend_api/friend app/friend/friend_api/friend.go
# friend_rpc
RUN go build -ldflags="-s -w" -o friend_rpc/friendrpc app/friend/friend_rpc/friendrpc.go

# ==================== GATEWAY 相关服务 ====================
# gateway_api
RUN go build -ldflags="-s -w" -o gateway_api/gateway app/gateway/gateway_api/gateway.go 

# ==================== GROUP 相关服务 ====================
# group_api
RUN go build -ldflags="-s -w" -o group_api/group app/group/group_api/group.go
# group_rpc
RUN go build -ldflags="-s -w" -o group_rpc/grouprpc app/group/group_rpc/grouprpc.go  

# ==================== USER 相关服务 ====================
# user_api
RUN go build -ldflags="-s -w" -o user_api/user app/user/user_api/user.go 
# user_rpc
RUN go build -ldflags="-s -w" -o user_rpc/userrpc app/user/user_rpc/userrpc.go  

# ==================== WS 相关服务 ====================
# ws_api
RUN go build -ldflags="-s -w" -o ws_api/ws app/ws/ws_api/ws.go

# ==================== FILE 相关服务 ====================
# file_api
RUN go build -ldflags="-s -w" -o file_api/file app/file/file_api/file.go
# file_rpc
RUN go build -ldflags="-s -w" -o file_rpc/filerpc app/file/file_rpc/filerpc.go

# ==================== FEEDBACK 相关服务 ====================
# feedback_api
RUN go build -ldflags="-s -w" -o feedback_api/feedback app/feedback/feedback_api/feedback.go

# ==================== DICTIONARY 相关服务 ====================
# dictionary_api
RUN go build -ldflags="-s -w" -o dictionary_api/dictionary app/dictionary/dictionary_api/dictionary.go
# dictionary_rpc
RUN go build -ldflags="-s -w" -o dictionary_rpc/dictionaryrpc app/dictionary/dictionary_rpc/dictionaryrpc.go


# ==================== MOMENT 相关服务 ====================
# moment_api
RUN go build -ldflags="-s -w" -o moment_api/moment app/moment/moment_api/moment.go
