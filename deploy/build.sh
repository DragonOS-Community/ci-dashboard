#!/bin/bash

# DragonOS CI Dashboard 生产环境构建脚本
# 使用方法: ./build.sh [选项]
# 选项:
#   -t, --tag TAG        指定镜像标签 (默认: dragonos-ci-dashboard:latest)
#   -p, --push           构建后推送到镜像仓库
#   -r, --registry REG   指定镜像仓库地址
#   -h, --help           显示帮助信息

set -e

# 默认配置
IMAGE_TAG="dragonos-ci-dashboard:latest"
PUSH_IMAGE=false
REGISTRY=""

# 解析命令行参数
while [[ $# -gt 0 ]]; do
    case $1 in
        -t|--tag)
            IMAGE_TAG="$2"
            shift 2
            ;;
        -p|--push)
            PUSH_IMAGE=true
            shift
            ;;
        -r|--registry)
            REGISTRY="$2"
            shift 2
            ;;
        -h|--help)
            echo "DragonOS CI Dashboard 生产环境构建脚本"
            echo ""
            echo "使用方法: $0 [选项]"
            echo ""
            echo "选项:"
            echo "  -t, --tag TAG        指定镜像标签 (默认: dragonos-ci-dashboard:latest)"
            echo "  -p, --push           构建后推送到镜像仓库"
            echo "  -r, --registry REG   指定镜像仓库地址"
            echo "  -h, --help           显示帮助信息"
            echo ""
            echo "示例:"
            echo "  $0                                    # 构建默认标签的镜像"
            echo "  $0 -t v1.0.0                         # 构建指定版本的镜像"
            echo "  $0 -t v1.0.0 -r registry.example.com # 构建并指定仓库地址"
            echo "  $0 -t v1.0.0 -r registry.example.com -p # 构建并推送"
            exit 0
            ;;
        *)
            echo "未知选项: $1"
            echo "使用 -h 或 --help 查看帮助信息"
            exit 1
            ;;
    esac
done

# 获取脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

echo "=========================================="
echo "DragonOS CI Dashboard 生产环境构建"
echo "=========================================="
echo "项目根目录: $PROJECT_ROOT"
echo "镜像标签: $IMAGE_TAG"
echo ""

# 检查必要的文件
echo "检查必要文件..."
if [ ! -f "$PROJECT_ROOT/backend/go.mod" ]; then
    echo "错误: 未找到 backend/go.mod"
    exit 1
fi

if [ ! -f "$PROJECT_ROOT/frontend/package.json" ]; then
    echo "错误: 未找到 frontend/package.json"
    exit 1
fi

if [ ! -f "$SCRIPT_DIR/Dockerfile" ]; then
    echo "错误: 未找到 deploy/Dockerfile"
    exit 1
fi

echo "✓ 文件检查通过"
echo ""

# 构建镜像
echo "开始构建 Docker 镜像..."
FULL_IMAGE_TAG="$IMAGE_TAG"
if [ -n "$REGISTRY" ]; then
    # 移除 registry 末尾的斜杠
    REGISTRY="${REGISTRY%/}"
    FULL_IMAGE_TAG="$REGISTRY/$IMAGE_TAG"
fi

cd "$PROJECT_ROOT"
docker build -f deploy/Dockerfile -t "$FULL_IMAGE_TAG" .

if [ $? -eq 0 ]; then
    echo "✓ 镜像构建成功: $FULL_IMAGE_TAG"
else
    echo "✗ 镜像构建失败"
    exit 1
fi

# 推送镜像
if [ "$PUSH_IMAGE" = true ]; then
    if [ -z "$REGISTRY" ]; then
        echo "错误: 推送镜像需要指定镜像仓库地址 (-r 或 --registry)"
        exit 1
    fi
    
    echo ""
    echo "推送镜像到 $REGISTRY..."
    docker push "$FULL_IMAGE_TAG"
    
    if [ $? -eq 0 ]; then
        echo "✓ 镜像推送成功: $FULL_IMAGE_TAG"
    else
        echo "✗ 镜像推送失败"
        exit 1
    fi
fi

echo ""
echo "=========================================="
echo "构建完成!"
echo "=========================================="
echo "镜像名称: $FULL_IMAGE_TAG"
echo ""
echo "使用以下命令运行容器:"
echo "  cd deploy"
echo "  docker-compose up -d"
echo ""

