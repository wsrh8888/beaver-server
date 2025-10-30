#!/bin/bash

# 从配置文件批量替换地址
# 使用方法: ./replace_from_config.sh
# 需要先修改 config.txt 文件中的目标地址

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CONFIG_FILE="$SCRIPT_DIR/config.txt"

# 检查配置文件是否存在
if [ ! -f "$CONFIG_FILE" ]; then
    echo "❌ Error: config.txt not found in $SCRIPT_DIR"
    exit 1
fi

echo "📖 Reading configuration from $CONFIG_FILE..."
echo ""

# 读取替换规则
REPLACEMENTS=()
while IFS= read -r line; do
    # 跳过注释行和空行
    if echo "$line" | grep -q "^[[:space:]]*#" || [ -z "$(echo "$line" | tr -d ' ')" ]; then
        continue
    fi
    
    # 解析替换规则 (格式: 当前地址 -> 目标地址)
    if echo "$line" | grep -q " -> "; then
        current=$(echo "$line" | sed 's/^[[:space:]]*\([^[:space:]]*\)[[:space:]]*->[[:space:]]*\([^[:space:]]*\).*/\1/')
        target=$(echo "$line" | sed 's/^[[:space:]]*\([^[:space:]]*\)[[:space:]]*->[[:space:]]*\([^[:space:]]*\).*/\2/')
        REPLACEMENTS+=("$current|$target")
        echo "🔄 $current -> $target"
    fi
done < "$CONFIG_FILE"

if [ ${#REPLACEMENTS[@]} -eq 0 ]; then
    echo "❌ No valid replacement rules found in config.txt"
    exit 1
fi

echo ""
# 确认是否继续
read -p "Do you want to continue? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "❌ Operation cancelled."
    exit 0
fi

# 计数器
updated_files=0
total_files=0

echo "🚀 Starting replacement..."

# 遍历所有子目录
for dir in "$SCRIPT_DIR"/*; do
    if [ -d "$dir" ]; then
        service_name=$(basename "$dir")
        
        # 查找该目录下的所有yaml文件
        for yaml_file in "$dir"/*.yaml; do
            if [ -f "$yaml_file" ]; then
                total_files=$((total_files + 1))
                filename=$(basename "$yaml_file")
                echo "📝 Processing: $service_name/$filename"
                
                # 不创建备份文件，直接修改
                
                # 应用所有替换规则
                file_updated=false
                for replacement in "${REPLACEMENTS[@]}"; do
                    IFS='|' read -r current target <<< "$replacement"
                    if grep -q "$current" "$yaml_file"; then
                        sed -i "s|$current|$target|g" "$yaml_file"
                        file_updated=true
                    fi
                done
                
                if [ "$file_updated" = true ]; then
                    echo "  ✅ Updated: $service_name/$filename"
                    updated_files=$((updated_files + 1))
                else
                    echo "  ⚠️  No changes needed: $service_name/$filename"
                fi
            fi
        done
    fi
done

echo ""
echo "📊 Summary:"
echo "  Total yaml files processed: $total_files"
echo "  Files updated: $updated_files"
echo ""
echo "✅ Replacement completed!" 