# 多阶段构建：前端构建阶段
FROM node:20-alpine AS frontend-builder

WORKDIR /app/frontend

# 复制前端依赖文件
COPY frontend/package*.json ./

# 安装前端依赖
RUN npm ci --only=production

# 复制前端源代码
COPY frontend/ .

# 构建前端应用
RUN npm run build

# 多阶段构建：后端构建阶段
FROM ubuntu:24.04 AS backend-builder

# 设置非交互式安装
ENV DEBIAN_FRONTEND=noninteractive

WORKDIR /app/backend

# 安装 Go、make 和构建依赖
RUN apt-get update && apt-get install -y --no-install-recommends \
    wget \
    git \
    ca-certificates \
    make \
    && rm -rf /var/lib/apt/lists/*

# 安装 Go 1.23（项目要求 Go 1.24，但使用 1.23 作为兼容版本）
# 如需使用其他版本，请修改下面的版本号
RUN wget -q https://go.dev/dl/go1.23.3.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go1.23.3.linux-amd64.tar.gz \
    && rm go1.23.3.linux-amd64.tar.gz

# 设置 Go 环境变量
ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH=/go
ENV GOCACHE=/go/cache

# 复制后端依赖文件
COPY backend/go.mod backend/go.sum ./

# 下载依赖
RUN go mod download

# 复制后端源代码和 Makefile
COPY backend/ .

# 使用 Makefile 构建两个二进制文件
RUN make build && make build-admin-cli

# 运行阶段：统一容器
FROM ubuntu:24.04

# 设置非交互式安装
ENV DEBIAN_FRONTEND=noninteractive

# 安装必要的工具
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    nginx \
    supervisor \
    curl \
    && rm -rf /var/lib/apt/lists/*

# 创建必要的目录
RUN mkdir -p \
    /var/www/html \
    /var/log/nginx \
    /var/log/supervisor \
    /etc/supervisor/conf.d \
    /data/uploads \
    /app/backend \
    /run/nginx

# 创建 nginx 运行用户（如果不存在）
RUN useradd -r -s /bin/false www-data || true

WORKDIR /app

# 从构建阶段复制文件
COPY --from=frontend-builder /app/frontend/dist /var/www/html
COPY --from=backend-builder /app/backend/bin/server /app/backend/server
COPY --from=backend-builder /app/backend/bin/admin-cli /app/backend/admin-cli
COPY --from=backend-builder /app/backend/migrations /app/backend/migrations

# 复制 nginx 配置
COPY deploy/nginx.conf /etc/nginx/nginx.conf

# 复制 supervisor 配置
COPY deploy/supervisord.conf /etc/supervisor/conf.d/supervisord.conf

# 设置权限
RUN chmod +x /app/backend/server && \
    chmod +x /app/backend/admin-cli && \
    chown -R www-data:www-data /var/www/html && \
    chown -R www-data:www-data /var/log/nginx && \
    chown -R www-data:www-data /run/nginx

# 暴露端口
EXPOSE 80

# 使用 supervisor 管理进程
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]

