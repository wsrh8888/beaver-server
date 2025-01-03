#!/bin/bash

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

# 分割参数
IFS='_' read -r module sub_module <<<"$module_name"

# 检查分割后是否得到两个非空值
if [ -z "$module" ] || [ -z "$sub_module" ];then
    echo "Error: Invalid input, expected format 'module_submodule'"
    exit 1
fi


# 跳转到目标目录
cd deploy/"$module"/"$module_name" || {
    echo "Error: Directory not found"
    exit 1
}

# 打包镜像
docker build -t "$module_name" .

# 获取镜像Id
image_id=$(get_image_id "$module_name")

# 检查获取到的image_id是否为空
if [ -z "$image_id" ];then
    echo "Error: Failed to retrieve the image Id"
    exit 1
fi

# 打标签并推送镜像
docker tag "$image_id" registry.cn-hangzhou.aliyuncs.com/beaver_im/"$module_name":1.0.0
docker push registry.cn-hangzhou.aliyuncs.com/beaver_im/"$module_name":1.0.0

echo "Successfully built, tagged, and pushed the image: $module_name:1.0.0"
