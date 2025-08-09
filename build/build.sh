#!/bin/bash

# 配置变量
REGISTRY_URL="wsrh8888/beaver_im"
VERSION="1.0.0"

# 函数：获取Docker镜像的Id
get_image_id() {
    docker images -q $1 | head -n1
}

# 检查是否传入了参数
if [ -z "$1" ];then
    echo "Usage: $0 module_name"
    exit 1
fi

# 获取参数
module_name="$1"

# 跳转到目标目录
cd deploy/"$module_name" || {
    echo "Error: Directory deploy/$module_name not found"
    exit 1
}

# 构建镜像
echo "Building image: $module_name"
docker build -t "$module_name" .
if [ $? -ne 0 ]; then
    echo "Error: Failed to build the image"
    exit 1
fi

# 获取镜像Id
image_id=$(get_image_id "$module_name")

# 检查获取到的image_id是否为空
if [ -z "$image_id" ];then
    echo "Error: Failed to retrieve the image Id"
    exit 1
fi

# 打标签并推送镜像
docker tag "$image_id" "$REGISTRY_URL/$module_name:$VERSION"
if [ $? -ne 0 ]; then
    echo "Error: Failed to tag the image"
    exit 1
fi

# 推送镜像，带重试机制
max_retries=3
retry_count=0

while [ $retry_count -lt $max_retries ]; do
    echo "Pushing image to registry (attempt $((retry_count + 1))/$max_retries)..."
    docker push "$REGISTRY_URL/$module_name:$VERSION"
    
    if [ $? -eq 0 ]; then
        echo "Successfully pushed the image"
        break
    else
        retry_count=$((retry_count + 1))
        if [ $retry_count -lt $max_retries ]; then
            echo "Push failed, retrying in 5 seconds..."
            sleep 5
        else
            echo "Error: Failed to push the image to registry after $max_retries attempts"
            exit 1
        fi
    fi
done

echo "Successfully built, tagged, and pushed the image: $module_name:$VERSION"
