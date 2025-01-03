#!/bin/bash

# 定义包含docker-compose.yaml文件的目录列表
directories=(
    "/mnt/beaver/user_rpc"
    "/mnt/beaver/group_rpc"
    "/mnt/beaver/friend_rpc"
    "/mnt/beaver/chat_rpc"
    "/mnt/beaver/auth_api"
    "/mnt/beaver/chat_api"
    "/mnt/beaver/friend_api"
    "/mnt/beaver/gateway_api"
    "/mnt/beaver/group_api"
    "/mnt/beaver/user_api"
    "/mnt/beaver/ws_api"
    "/mnt/beaver/file_api"
)

# 循环遍历每个目录进行完整的停止、删除、拉取、启动流程
for dir in "${directories[@]}"; do
    if [ -d "$dir" ]; then
        echo "Navigating to directory: $dir"
        cd "$dir" || { echo "Failed to navigate to directory: $dir"; continue; }

        # 调试信息：列出当前目录内容
        echo "Listing contents of $dir:"
        ls -la

        if [ -f "docker-compose.yaml" ]; then
            echo "Stopping docker-compose services in $dir"
            docker-compose down
            echo "Services stopped in $dir"

            echo "Removing old Docker images in $dir"
            docker-compose rm -f

            echo "Pulling latest docker images in $dir"
            docker-compose pull
            
            echo "Building docker images in $dir (if required)"
            docker-compose build

            echo "Starting docker-compose in $dir"
            docker-compose up -d
            echo "docker-compose started in $dir"
        else
            echo "docker-compose.yaml not found in $dir"
        fi
    else
        echo "Directory not found: $dir"
    fi
done

# 清理未被使用的镜像
echo "Pruning unused Docker images"
docker image prune -f

echo "All services started."