FROM golang:alpine AS builder
LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.cn,direct



# 安装必要的软件包
RUN apk add --no-cache tzdata

WORKDIR /build

COPY go.mod .
COPY go.sum .

COPY . .
RUN go mod tidy


# auth_api
RUN go build -o auth_api/auth app/auth/auth_api/auth.go 

# chat_api
RUN go build -o chat_api/chat app/chat/chat_api/chat.go
# chat_rpc
RUN go build -o chat_rpc/chatrpc app/chat/chat_rpc/chatrpc.go

# friend_api
RUN go build -o friend_api/friend app/friend/friend_api/friend.go
# friend_rpc
RUN go build -o friend_rpc/friendrpc app/friend/friend_rpc/friendrpc.go  


# gateway_api
RUN go build -o gateway_api/gateway app/gateway/gateway_api/gateway.go 

# group_api
RUN go build -o group_api/group app/group/group_api/group.go
# group_rpc
RUN go build -o group_rpc/grouprpc app/group/group_rpc/grouprpc.go  


# user_api
RUN go build -o user_api/user app/user/user_api/user.go 
# user_rpc
RUN go build -o user_rpc/userrpc app/user/user_rpc/userrpc.go  

# ws_api
RUN go build -o ws_api/ws app/ws/ws_api/ws.go

# file_api
RUN go build -o file_api/file app/file/file_api/file.go


# feedback_api
RUN go build -o feedback_api/feedback app/feedback/feedback_api/feedback.go
