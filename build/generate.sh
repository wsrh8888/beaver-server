#!/bin/bash

# 自动编译脚本
# 自动遍历 app 目录下的所有模块，编译 .api 和 .proto 文件

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 从文件末尾往上查找包含 goctl 的行
extract_goctl_command() {
    local file="$1"
    
    if [ -f "$file" ]; then
        # 使用 sed 获取所有行，然后从后往前查找包含 goctl 的非空行
        local lines=()
        while IFS= read -r line; do
            lines+=("$line")
        done < "$file"
        
        # 从后往前遍历
        local i
        for ((i=${#lines[@]}-1; i>=0; i--)); do
            local trimmed_line=$(echo "${lines[i]}" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
            if [ -n "$trimmed_line" ] && [[ "$trimmed_line" =~ goctl ]]; then
                # 移除注释符号 // 或 //
                local command="${trimmed_line#// }"
                command="${command#//}"
                echo "$command"
                return 0
            fi
        done
    fi
    return 1
}

# 编译单个文件
compile_file() {
    local file_path="$1"
    local filename=$(basename "$file_path")
    
    log_info "编译文件: $filename"
    
    # 提取 goctl 命令
    local goctl_command=$(extract_goctl_command "$file_path")
    if [ $? -ne 0 ]; then
        log_warning "未找到 goctl 命令，跳过: $filename"
        return 1
    fi
    
    log_info "执行命令: $goctl_command"
    
    local dir=$(dirname "$file_path")
    local original_dir=$(pwd)
    
    cd "$dir" || {
        log_error "无法进入目录: $dir"
        return 1
    }
    
    # 执行 goctl 命令
    if eval "$goctl_command"; then
        log_success "编译成功: $filename"
        cd "$original_dir"
        return 0
    else
        log_error "编译失败: $filename"
        cd "$original_dir"
        return 1
    fi
}

# 处理单个模块
process_module() {
    local module_path="$1"
    local module_name="$2"
    
    log_info "处理模块: $module_name"
    
    local total_files=0
    local success_files=0
    
    # 遍历模块下的所有子目录
    for sub_dir_path in "$module_path"/*/; do
        if [ -d "$sub_dir_path" ]; then
            local sub_dir_name=$(basename "$sub_dir_path")
            log_info "  检查子目录: $sub_dir_name"
            
            # 查找 .api 文件
            while IFS= read -r api_file; do
                if [ -n "$api_file" ]; then
                    total_files=$((total_files + 1))
                    if compile_file "$api_file"; then
                        success_files=$((success_files + 1))
                    fi
                fi
            done < <(find "$sub_dir_path" -name "*.api" 2>/dev/null)
            
            # 查找 .proto 文件
            while IFS= read -r proto_file; do
                if [ -n "$proto_file" ]; then
                    total_files=$((total_files + 1))
                    if compile_file "$proto_file"; then
                        success_files=$((success_files + 1))
                    fi
                fi
            done < <(find "$sub_dir_path" -name "*.proto" 2>/dev/null)
        fi
    done
    
    if [ $total_files -gt 0 ]; then
        log_info "模块 $module_name 完成: $success_files/$total_files 成功"
    else
        log_info "模块 $module_name 没有找到需要编译的文件"
    fi
}

# 主函数
main() {
    log_info "开始自动编译 app 目录下的所有模块..."
    echo "========================================"
    
    # 获取 app 目录下的所有模块
    local app_path="app"
    if [ ! -d "$app_path" ]; then
        log_error "app 目录不存在"
        exit 1
    fi
    
    # 直接遍历 app 目录下的子目录
    for module_path in app/*/; do
        if [ -d "$module_path" ]; then
            local module_name=$(basename "$module_path")
            echo "----------------------------------------"
            process_module "$module_path" "$module_name"
        fi
    done
    
    echo "========================================"
    log_success "编译完成！"
}

# 脚本入口
main 