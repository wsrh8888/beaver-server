#!/bin/bash

# 定义模块列表
modules=(
    # ==================== RPC 服务 ====================
    "user_rpc"
    "group_rpc"
    "friend_rpc"
    "chat_rpc"
    "file_rpc"
    "dictionary_rpc"
    "datasync_rpc"
    
    # ==================== API 服务 ====================
    "auth_api"
    "chat_api"
    "dictionary_api"
    "emoji_api"
    "feedback_api"
    "file_api"
    "friend_api"
    "group_api"
    "moment_api"
    "track_api"
    "update_api"
    "user_api"
    "ws_api"
    "gateway_api"
    "datasync_api"
    
    
    
    # ==================== ADMIN 服务 ====================
    "auth_admin"
    "chat_admin"
    "emoji_admin"
    "gateway_admin"
    "group_admin"
    "moment_admin"
    "track_admin"
    "user_admin"
    "system_admin"
    "update_admin"
    "feedback_admin"
    "file_admin"
    "friend_admin"
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

echo "🎉 所有模块构建完成！"
