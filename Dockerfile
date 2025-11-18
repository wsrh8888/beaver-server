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
# auth_admin
RUN go build -ldflags="-s -w" -o auth_admin/auth app/auth/auth_admin/auth.go

# ==================== CHAT 相关服务 ====================
# chat_api
RUN go build -ldflags="-s -w" -o chat_api/chat app/chat/chat_api/chat.go
# chat_rpc
RUN go build -ldflags="-s -w" -o chat_rpc/chatrpc app/chat/chat_rpc/chatrpc.go
# chat_admin
RUN go build -ldflags="-s -w" -o chat_admin/chat app/chat/chat_admin/chat.go

# ==================== FRIEND 相关服务 ====================
# friend_api
RUN go build -ldflags="-s -w" -o friend_api/friend app/friend/friend_api/friend.go
# friend_rpc
RUN go build -ldflags="-s -w" -o friend_rpc/friendrpc app/friend/friend_rpc/friendrpc.go
# friend_admin
RUN go build -ldflags="-s -w" -o friend_admin/friend app/friend/friend_admin/friend.go  

# ==================== GATEWAY 相关服务 ====================
# gateway_api
RUN go build -ldflags="-s -w" -o gateway_api/gateway app/gateway/gateway_api/gateway.go 
# gateway_admin
RUN go build -ldflags="-s -w" -o gateway_admin/gateway app/gateway/gateway_admin/gateway.go

# ==================== GROUP 相关服务 ====================
# group_api
RUN go build -ldflags="-s -w" -o group_api/group app/group/group_api/group.go
# group_rpc
RUN go build -ldflags="-s -w" -o group_rpc/grouprpc app/group/group_rpc/grouprpc.go  
# group_admin
RUN go build -ldflags="-s -w" -o group_admin/group app/group/group_admin/group.go

# ==================== USER 相关服务 ====================
# user_api
RUN go build -ldflags="-s -w" -o user_api/user app/user/user_api/user.go 
# user_rpc
RUN go build -ldflags="-s -w" -o user_rpc/userrpc app/user/user_rpc/userrpc.go  
# user_admin
RUN go build -ldflags="-s -w" -o user_admin/user app/user/user_admin/user.go

# ==================== WS 相关服务 ====================
# ws_api
RUN go build -ldflags="-s -w" -o ws_api/ws app/ws/ws_api/ws.go

# ==================== FILE 相关服务 ====================
# file_api
RUN go build -ldflags="-s -w" -o file_api/file app/file/file_api/file.go
# file_rpc
RUN go build -ldflags="-s -w" -o file_rpc/filerpc app/file/file_rpc/filerpc.go
# file_admin
RUN go build -ldflags="-s -w" -o file_admin/file app/file/file_admin/file.go

# ==================== MOMENT 相关服务 ====================
# moment_api
RUN go build -ldflags="-s -w" -o moment_api/moment app/moment/moment_api/moment.go
# moment_admin
RUN go build -ldflags="-s -w" -o moment_admin/moment app/moment/moment_admin/moment.go

# ==================== EMOJI 相关服务 ====================
# emoji_api
RUN go build -ldflags="-s -w" -o emoji_api/emoji app/emoji/emoji_api/emoji.go
# emoji_admin
RUN go build -ldflags="-s -w" -o emoji_admin/emoji app/emoji/emoji_admin/emoji.go

# ==================== FEEDBACK 相关服务 ====================
# feedback_api
RUN go build -ldflags="-s -w" -o feedback_api/feedback app/feedback/feedback_api/feedback.go
# feedback_admin
RUN go build -ldflags="-s -w" -o feedback_admin/feedback app/feedback/feedback_admin/feedback.go

# ==================== TRACK 相关服务 ====================
# track_api
RUN go build -ldflags="-s -w" -o track_api/track app/track/track_api/track.go
# track_admin
RUN go build -ldflags="-s -w" -o track_admin/track app/track/track_admin/track.go

# ==================== UPDATE 相关服务 ====================
# update_api
RUN go build -ldflags="-s -w" -o update_api/update app/update/update_api/update.go
# update_admin
RUN go build -ldflags="-s -w" -o update_admin/update app/update/update_admin/update.go

# ==================== DICTIONARY 相关服务 ====================
# dictionary_api
RUN go build -ldflags="-s -w" -o dictionary_api/dictionary app/dictionary/dictionary_api/dictionary.go
# dictionary_rpc
RUN go build -ldflags="-s -w" -o dictionary_rpc/dictionaryrpc app/dictionary/dictionary_rpc/dictionaryrpc.go

# ==================== SYSTEM 相关服务 ====================
# system_admin
RUN go build -ldflags="-s -w" -o system_admin/system app/system/system_admin/system.go

