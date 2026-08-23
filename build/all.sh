#!/bin/bash

REGISTRY_URL="wsrh8888"
ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
VERSION="$(tr -d '[:space:]' < "$ROOT_DIR/VERSION")"
PUBLISH_PROBE_MODULE="gateway_api"

# 检查当前版本是否已发布到镜像仓库
assert_version_not_published() {
    local image="$REGISTRY_URL/$PUBLISH_PROBE_MODULE:$VERSION"
    echo "🔍 检查版本 $VERSION 是否已发布..."
    if docker manifest inspect "$image" >/dev/null 2>&1; then
        echo "❌ 版本 $VERSION 已发布（$image 已存在于镜像仓库）"
        echo "   请先更新 $ROOT_DIR/VERSION 后再执行构建发布"
        exit 1
    fi
    echo "✅ 版本 $VERSION 尚未发布，开始构建"
}

assert_version_not_published

# 定义模块列表
modules=(
    # ==================== RPC 服务 ====================
    "user_rpc"
    "auth_rpc"
    "group_rpc"
    "circle_rpc"
    "friend_rpc"
    "chat_rpc"
    "file_rpc"
    "emoji_rpc"
    "moment_rpc"
    "notification_rpc"
    "platform_rpc"
    "call_rpc"
    "open_rpc"
    
    # ==================== API 服务 ====================
    "auth_api"
    "chat_api"
    "datasync_api"
    "emoji_api"
    "file_api"
    "friend_api"
    "group_api"
    "circle_api"
    "moment_api"
    "notification_api"
    "platform_api"
    "user_api"
    "ws_api"
    "gateway_api"
    "call_api"
    "open_api"
    "open_portal"
    
    # ==================== ADMIN 服务 ====================
    "backend_admin"
    "gateway_admin"
)

# 在当前路径执行 docker build
docker build -t beaver_server .

# 并发构建函数
build_module() {
    local module_name="$1"
    echo "🚀 开始构建: $module_name"
    ./build/build.sh "$module_name"
    if [ $? -eq 0 ]; then
        echo "✅ 打包镜像成功: $module_name"
    else
        echo "❌ 打包镜像失败: $module_name"
    fi
    echo "----------------------------------------"
}

# 并发构建，限制最大并发数为4
max_jobs=4
current_jobs=0

for module_name in "${modules[@]}"; do
    # 如果当前并发数达到上限，等待一个任务完成
    while [ $current_jobs -ge $max_jobs ]; do
        wait -n
        current_jobs=$((current_jobs - 1))
    done
    
    # 启动后台任务
    build_module "$module_name" &
    current_jobs=$((current_jobs + 1))
done

# 等待所有后台任务完成
wait

echo "🎉 所有模块构建完成！"
