#!/bin/bash

# 创建beaver_network网络（如果不存在）
docker network ls | grep beaver_network || docker network create beaver_network

# 获取当前脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# 函数：按顺序启动服务
start_services() {
    local service_type=$1
    local pattern=$2
    
    echo "=== Starting $service_type services ==="
    
    # 查找匹配模式的目录
    for dir in "$SCRIPT_DIR"/*; do
        if [ -d "$dir" ] && [[ "$(basename "$dir")" == *"$pattern"* ]]; then
            if [ -f "$dir/docker-compose.yaml" ]; then
                echo "Processing $service_type service: $(basename "$dir")"
                cd "$dir" || { echo "Failed to navigate to directory: $dir"; continue; }

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
                
                # 短暂等待服务启动
                sleep 2
            else
                echo "docker-compose.yaml not found in $dir"
            fi
        fi
    done
}

# 按顺序启动服务：RPC -> API -> Admin
echo "Starting services in order: RPC -> API -> Admin"

# 1. 启动 RPC 服务
start_services "RPC" "_rpc"

# 2. 启动 API 服务
start_services "API" "_api"

# 3. 启动 Admin 服务
start_services "Admin" "_admin"

# 清理未被使用的镜像
echo "Pruning unused Docker images"
docker image prune -f

echo "All services started."