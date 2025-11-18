#!/bin/bash

# 检查是否传入了版本号参数
if [ -z "$1" ]; then
    echo "Usage: $0 <version>"
    echo "Example: $0 20241201-143052"
    exit 1
fi

VERSION="$1"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "Updating all docker-compose.yaml files to version: $VERSION"
echo "Working directory: $SCRIPT_DIR"

# 计数器
updated_count=0
total_count=0

# 遍历当前目录下的所有子目录
for dir in "$SCRIPT_DIR"/*; do
    if [ -d "$dir" ]; then
        docker_compose_file="$dir/docker-compose.yaml"
        if [ -f "$docker_compose_file" ]; then
            total_count=$((total_count + 1))
            service_name=$(basename "$dir")
            echo "Processing: $service_name"
            
            # 直接替换版本号，不使用备份
            sed -i "s|:latest|:$VERSION|g" "$docker_compose_file"
            sed -i "s|:[0-9]\+\.[0-9]\+\.[0-9]\+|:$VERSION|g" "$docker_compose_file"
            sed -i "s|:[0-9]\{8\}-[0-9]\{6\}|:$VERSION|g" "$docker_compose_file"
            
            # 检查是否成功更新
            if grep -q ":$VERSION" "$docker_compose_file"; then
                echo "  ✅ Updated: $service_name"
                updated_count=$((updated_count + 1))
            else
                echo "  ❌ Failed to update: $service_name"
            fi
        fi
    fi
done

echo ""
echo "Summary:"
echo "  Total docker-compose.yaml files found: $total_count"
echo "  Successfully updated: $updated_count"
echo "  Version set to: $VERSION"

echo ""
echo "All docker-compose.yaml files have been updated to use version: $VERSION" 