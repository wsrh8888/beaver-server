#!/bin/bash

# 获取当前脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# 函数：按顺序停止服务
stop_services() {
    local service_type=$1
    local pattern=$2
    
    echo "=== Stopping $service_type services ==="
    
    # 查找匹配模式的目录
    for dir in "$SCRIPT_DIR"/*; do
        if [ -d "$dir" ] && [[ "$(basename "$dir")" == *"$pattern"* ]]; then
            if [ -f "$dir/docker-compose.yaml" ]; then
                echo "Processing $service_type service: $(basename "$dir")"
                cd "$dir" || { echo "Failed to navigate to directory: $dir"; continue; }

                echo "Stopping docker-compose services in $dir"
                docker-compose down
                echo "Services stopped in $dir"
                
                # 短暂等待服务停止
                sleep 1
            else
                echo "docker-compose.yaml not found in $dir"
            fi
        fi
    done
}

# 按顺序停止服务：Admin -> API -> RPC
echo "Stopping services in order: Admin -> API -> RPC"

# 1. 停止 Admin 服务
stop_services "Admin" "_admin"

# 2. 停止 API 服务
stop_services "API" "_api"

# 3. 停止 RPC 服务
stop_services "RPC" "_rpc"

# 清理容器和网络
echo "Cleaning up containers and networks"
docker container prune -f
docker network prune -f


echo "All services stopped and cleaned up." 