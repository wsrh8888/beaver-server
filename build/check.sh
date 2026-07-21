#!/bin/bash

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🔍 开始检查服务一致性...${NC}"
echo "=================================="

# 获取app目录下的所有服务（数据源头）
echo -e "${YELLOW}📁 扫描 /app 目录下的服务（数据源头）...${NC}"
app_services=()

# 遍历 app 目录下的所有子目录
for dir in app/*/; do
    if [ -d "$dir" ]; then
        for subdir in "${dir}"*/; do
            if [ -d "$subdir" ]; then
                subname=$(basename "$subdir")
                if [[ "$subname" =~ _(api|rpc|admin|portal)$ ]]; then
                    app_services+=("$subname")
                fi
            fi
        done
    fi
done

echo -e "${GREEN}✅ 在 /app 目录下找到 ${#app_services[@]} 个服务:${NC}"
for service in "${app_services[@]}"; do
    echo "  - $service"
done
echo ""

# 检查Dockerfile中的服务
echo -e "${YELLOW}🐳 检查 Dockerfile 中的服务...${NC}"
dockerfile_services=()
if [ -f "Dockerfile" ]; then
    # 提取Dockerfile中的go build命令
    while IFS= read -r line; do
        if [[ $line =~ go\ build\ -o\ ([^/]+)/ ]]; then
            service_name="${BASH_REMATCH[1]}"
            dockerfile_services+=("$service_name")
        fi
    done < Dockerfile
fi

echo -e "${GREEN}✅ Dockerfile 中包含 ${#dockerfile_services[@]} 个服务:${NC}"
for service in "${dockerfile_services[@]}"; do
    echo "  - $service"
done
echo ""

# 检查deploy目录中的服务
echo -e "${YELLOW}📦 检查 /deploy 目录中的服务...${NC}"
deploy_services=()
if [ -d "deploy" ]; then
    for dir in deploy/*/; do
        if [ -d "$dir" ]; then
            service_name=$(basename "$dir")
            deploy_services+=("$service_name")
        fi
    done
fi

echo -e "${GREEN}✅ /deploy 目录中包含 ${#deploy_services[@]} 个服务:${NC}"
for service in "${deploy_services[@]}"; do
    echo "  - $service"
done
echo ""

# 检查all.sh中的服务
echo -e "${YELLOW}📜 检查 all.sh 脚本中的服务...${NC}"
allsh_services=()
if [ -f "build/all.sh" ]; then
    # 提取all.sh中的模块列表，只提取引号中的服务名
    while IFS= read -r line; do
        if [[ $line =~ \"([a-zA-Z_]+)\" ]]; then
            service_name="${BASH_REMATCH[1]}"
            # 过滤掉变量名和其他非服务名
            if [[ ! "$service_name" =~ ^\$|^[A-Z] ]]; then
                allsh_services+=("$service_name")
            fi
        fi
    done < build/all.sh
fi

echo -e "${GREEN}✅ all.sh 脚本中包含 ${#allsh_services[@]} 个服务:${NC}"
for service in "${allsh_services[@]}"; do
    echo "  - $service"
done
echo ""

# 比较和报告
echo -e "${BLUE}🔍 开始比较服务一致性...${NC}"
echo "=================================="

# 检查Dockerfile中缺失的服务
echo -e "${YELLOW}📋 检查 Dockerfile 中缺失的服务...${NC}"
missing_in_dockerfile=()
for service in "${app_services[@]}"; do
    found=false
    for dockerfile_service in "${dockerfile_services[@]}"; do
        if [ "$service" = "$dockerfile_service" ]; then
            found=true
            break
        fi
    done
    if [ "$found" = false ]; then
        missing_in_dockerfile+=("$service")
    fi
done

if [ ${#missing_in_dockerfile[@]} -eq 0 ]; then
    echo -e "${GREEN}✅ Dockerfile 包含所有 app 目录中的服务${NC}"
else
    echo -e "${RED}❌ Dockerfile 中缺失以下服务:${NC}"
    for service in "${missing_in_dockerfile[@]}"; do
        echo "  - $service"
    done
fi

# 检查Dockerfile中多余的服务
echo -e "${YELLOW}🔍 检查 Dockerfile 中多余的服务...${NC}"
extra_in_dockerfile=()
for dockerfile_service in "${dockerfile_services[@]}"; do
    found=false
    for service in "${app_services[@]}"; do
        if [ "$dockerfile_service" = "$service" ]; then
            found=true
            break
        fi
    done
    if [ "$found" = false ]; then
        extra_in_dockerfile+=("$dockerfile_service")
    fi
done

if [ ${#extra_in_dockerfile[@]} -eq 0 ]; then
    echo -e "${GREEN}✅ Dockerfile 中没有多余的服务${NC}"
else
    echo -e "${RED}❌ Dockerfile 中有多余的服务:${NC}"
    for service in "${extra_in_dockerfile[@]}"; do
        echo "  - $service"
    done
fi

# 检查all.sh中缺失的服务
echo ""
echo -e "${YELLOW}📋 检查 all.sh 脚本中缺失的服务...${NC}"
missing_in_allsh=()
for service in "${app_services[@]}"; do
    found=false
    for allsh_service in "${allsh_services[@]}"; do
        if [ "$service" = "$allsh_service" ]; then
            found=true
            break
        fi
    done
    if [ "$found" = false ]; then
        missing_in_allsh+=("$service")
    fi
done

if [ ${#missing_in_allsh[@]} -eq 0 ]; then
    echo -e "${GREEN}✅ all.sh 脚本包含所有 app 目录中的服务${NC}"
else
    echo -e "${RED}❌ all.sh 脚本中缺失以下服务:${NC}"
    for service in "${missing_in_allsh[@]}"; do
        echo "  - $service"
    done
fi

# 检查all.sh中多余的服务
echo -e "${YELLOW}🔍 检查 all.sh 脚本中多余的服务...${NC}"
extra_in_allsh=()
for allsh_service in "${allsh_services[@]}"; do
    found=false
    for service in "${app_services[@]}"; do
        if [ "$allsh_service" = "$service" ]; then
            found=true
            break
        fi
    done
    if [ "$found" = false ]; then
        extra_in_allsh+=("$allsh_service")
    fi
done

if [ ${#extra_in_allsh[@]} -eq 0 ]; then
    echo -e "${GREEN}✅ all.sh 脚本中没有多余的服务${NC}"
else
    echo -e "${RED}❌ all.sh 脚本中有多余的服务:${NC}"
    for service in "${extra_in_allsh[@]}"; do
        echo "  - $service"
    done
fi

# 检查deploy目录中的服务（只提示，不报错）
echo ""
echo -e "${YELLOW}📦 deploy 目录检查（仅供参考）:${NC}"
missing_in_deploy=()
for service in "${app_services[@]}"; do
    found=false
    for deploy_service in "${deploy_services[@]}"; do
        if [ "$service" = "$deploy_service" ]; then
            found=true
            break
        fi
    done
    if [ "$found" = false ]; then
        missing_in_deploy+=("$service")
    fi
done

if [ ${#missing_in_deploy[@]} -eq 0 ]; then
    echo -e "${GREEN}✅ deploy 目录包含所有 app 目录中的服务${NC}"
else
    echo -e "${YELLOW}⚠️  deploy 目录中可能缺失以下服务（需要手动添加）:${NC}"
    for service in "${missing_in_deploy[@]}"; do
        echo "  - $service"
    done
fi

# 检查deploy目录中多余的服务
extra_in_deploy=()
for deploy_service in "${deploy_services[@]}"; do
    found=false
    for service in "${app_services[@]}"; do
        if [ "$deploy_service" = "$service" ]; then
            found=true
            break
        fi
    done
    if [ "$found" = false ]; then
        extra_in_deploy+=("$deploy_service")
    fi
done

if [ ${#extra_in_deploy[@]} -eq 0 ]; then
    echo -e "${GREEN}✅ deploy 目录中没有多余的服务${NC}"
else
    echo -e "${YELLOW}⚠️  deploy 目录中可能有多余的服务:${NC}"
    for service in "${extra_in_deploy[@]}"; do
        echo "  - $service"
    done
fi

# 总结
echo ""
echo -e "${BLUE}📊 检查总结:${NC}"
echo "=================================="
echo -e "App 目录中的服务数量（数据源头）: ${#app_services[@]}"
echo -e "Dockerfile 中的服务数量: ${#dockerfile_services[@]}"
echo -e "Deploy 目录中的服务数量: ${#deploy_services[@]}"
echo -e "All.sh 脚本中的服务数量: ${#allsh_services[@]}"

# 计算差异
dockerfile_diff=$(( ${#dockerfile_services[@]} - ${#app_services[@]} ))
allsh_diff=$(( ${#allsh_services[@]} - ${#app_services[@]} ))
deploy_diff=$(( ${#deploy_services[@]} - ${#app_services[@]} ))

echo ""
echo -e "${BLUE}📈 差异分析:${NC}"
if [ $dockerfile_diff -eq 0 ]; then
    echo -e "${GREEN}✅ Dockerfile 服务数量正确${NC}"
else
    if [ $dockerfile_diff -gt 0 ]; then
        echo -e "${RED}❌ Dockerfile 服务数量差异: +$dockerfile_diff (多余)${NC}"
    else
        echo -e "${RED}❌ Dockerfile 服务数量差异: $dockerfile_diff (缺失)${NC}"
    fi
fi

if [ $allsh_diff -eq 0 ]; then
    echo -e "${GREEN}✅ all.sh 脚本服务数量正确${NC}"
else
    if [ $allsh_diff -gt 0 ]; then
        echo -e "${RED}❌ all.sh 脚本服务数量差异: +$allsh_diff (多余)${NC}"
    else
        echo -e "${RED}❌ all.sh 脚本服务数量差异: $allsh_diff (缺失)${NC}"
    fi
fi

if [ $deploy_diff -eq 0 ]; then
    echo -e "${GREEN}✅ deploy 目录服务数量正确${NC}"
else
    if [ $deploy_diff -gt 0 ]; then
        echo -e "${YELLOW}⚠️  deploy 目录服务数量差异: +$deploy_diff (多余)${NC}"
        if [ ${#extra_in_deploy[@]} -gt 0 ]; then
            echo -e "${YELLOW}   多余的服务: ${extra_in_deploy[*]}${NC}"
        fi
    else
        echo -e "${YELLOW}⚠️  deploy 目录服务数量差异: $deploy_diff (缺失)${NC}"
        if [ ${#missing_in_deploy[@]} -gt 0 ]; then
            echo -e "${YELLOW}   缺失的服务: ${missing_in_deploy[*]}${NC}"
        fi
    fi
fi

if [ ${#missing_in_dockerfile[@]} -eq 0 ] && [ ${#missing_in_allsh[@]} -eq 0 ] && [ ${#extra_in_dockerfile[@]} -eq 0 ] && [ ${#extra_in_allsh[@]} -eq 0 ]; then
    echo ""
    echo -e "${GREEN}🎉 所有核心服务配置完全一致！${NC}"
else
    echo ""
    echo -e "${RED}⚠️  发现不一致，请根据上述报告进行修复${NC}"
fi

echo ""
echo -e "${BLUE}💡 修复建议:${NC}"
if [ ${#missing_in_dockerfile[@]} -gt 0 ]; then
    echo "- 在 Dockerfile 中添加缺失的 go build 命令"
fi
if [ ${#extra_in_dockerfile[@]} -gt 0 ]; then
    echo "- 从 Dockerfile 中删除多余的服务"
fi
if [ ${#missing_in_allsh[@]} -gt 0 ]; then
    echo "- 在 all.sh 脚本中添加缺失的服务到模块列表"
fi
if [ ${#extra_in_allsh[@]} -gt 0 ]; then
    echo "- 从 all.sh 脚本中删除多余的服务"
fi
echo "- 如果 deploy 目录缺失服务，需要手动创建对应的部署配置"
echo "- 建议定期运行此脚本确保服务配置的一致性" 