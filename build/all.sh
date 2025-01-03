#!/bin/bash

# 定义模块列表
modules=(
    "user_rpc"
    "group_rpc"
    "friend_rpc"
    "chat_rpc"
    "auth_api"
    "chat_api"
    "friend_api"
    "gateway_api"
    "group_api"
    "user_api"
    "ws_api"
    "file_api"
)

# 在当前路径执行 docker build
docker build -t beaver_server .


# 循环处理每个模块名
for module_name in "${modules[@]}"; do
    ./build/build.sh "$module_name"
    echo "打包镜像成功: $module_name:1.0.0"
    echo "            "
    echo "            "
    echo "            "
    echo "            "
    echo "            "
    sleep 2s

done
