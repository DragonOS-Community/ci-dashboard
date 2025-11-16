# DragonOS CI Dashboard 生产环境部署指南

本目录包含生产环境的部署配置和脚本。

## 目录结构

```
deploy/
├── Dockerfile              # 统一构建文件（前端+后端+nginx）
├── nginx.conf              # Nginx 反向代理配置
├── supervisord.conf        # Supervisor 进程管理配置
├── docker-compose.yml      # 生产环境 Docker Compose 配置
├── build.sh                # 镜像构建脚本
├── .env.example            # 环境变量配置示例
└── README.md               # 本文件
```

## 快速开始

### 1. 准备环境变量

复制环境变量示例文件并修改配置：

```bash
cd deploy
cp .env.example .env
vim .env  # 修改数据库连接等配置
```

**重要配置项：**
- `DB_HOST`: 数据库地址（外部部署）
- `DB_USER`: 数据库用户名
- `DB_PASSWORD`: 数据库密码
- `JWT_SECRET`: JWT 密钥（必须修改为强随机字符串）
- `API_KEY_HASH_SALT`: API Key 哈希盐值（必须修改为强随机字符串）
- `CORS_ALLOW_ORIGINS`: 允许的跨域来源（生产环境建议指定具体域名）

### 2. 构建镜像

使用构建脚本构建 Docker 镜像：

```bash
# 构建默认标签的镜像
./build.sh

# 构建指定版本的镜像
./build.sh -t v1.0.0

# 构建并推送到镜像仓库
./build.sh -t v1.0.0 -r registry.example.com -p
```

### 3. 启动服务

使用 Docker Compose 启动服务：

```bash
docker-compose up -d
```

### 4. 查看日志

```bash
# 查看所有日志
docker-compose logs -f

# 查看应用日志
docker-compose logs -f app

# 查看后端日志
docker exec dragonos-ci-app tail -f /var/log/supervisor/backend.out.log

# 查看 Nginx 日志
docker exec dragonos-ci-app tail -f /var/log/nginx/access.log
```

### 5. 停止服务

```bash
docker-compose down
```

## 架构说明

### 容器架构

生产环境使用单一容器部署，包含以下组件：

1. **前端静态文件**: Vue 3 构建产物，由 Nginx 提供静态文件服务
2. **后端 API 服务**: Go 应用，监听 `127.0.0.1:8080`
3. **Nginx 反向代理**: 
   - 提供前端静态文件服务
   - 将 `/api/*` 请求代理到后端服务
4. **Supervisor**: 进程管理器，同时管理后端服务和 Nginx

### 网络架构

```
用户请求
   ↓
Nginx (端口 80)
   ├─ / → 前端静态文件 (/var/www/html)
   └─ /api → 代理到后端 (127.0.0.1:8080)
```

### 数据存储

- **上传文件**: 挂载到 `./data/uploads` 目录
- **日志文件**: 容器内 `/var/log/` 目录

## 配置说明

### 环境变量

所有配置通过环境变量传递，详见 `.env.example` 文件。

### Nginx 配置

Nginx 配置文件位于 `nginx.conf`，主要功能：

- 静态文件服务（前端构建产物）
- API 请求代理到后端服务
- Gzip 压缩
- 静态资源缓存
- 安全头设置

### Supervisor 配置

Supervisor 配置文件位于 `supervisord.conf`，管理：

- 后端服务进程
- Nginx 进程

## 健康检查

容器包含健康检查配置，检查 Nginx 的 `/health` 端点：

```bash
# 手动检查
curl http://localhost/health
```

## 数据备份

### 上传文件备份

```bash
# 备份上传文件
tar -czf uploads-backup-$(date +%Y%m%d).tar.gz ./data/uploads
```

### 数据库备份

数据库是外部部署的，请参考数据库的备份方案。

## 故障排查

### 查看容器状态

```bash
docker-compose ps
```

### 查看容器日志

```bash
# 查看所有日志
docker-compose logs

# 查看特定服务的日志
docker-compose logs app
```

### 进入容器调试

```bash
docker exec -it dragonos-ci-app sh
```

### 常见问题

1. **后端服务无法连接数据库**
   - 检查 `.env` 文件中的数据库配置
   - 确认数据库服务可访问
   - 检查网络连接

2. **前端页面无法访问**
   - 检查 Nginx 是否正常运行
   - 查看 Nginx 错误日志：`docker exec dragonos-ci-app cat /var/log/nginx/error.log`

3. **API 请求失败**
   - 检查后端服务是否正常运行
   - 查看后端日志：`docker exec dragonos-ci-app tail -f /var/log/supervisor/backend.out.log`
   - 检查 Nginx 代理配置

## 更新部署

### 更新应用

1. 构建新镜像：
   ```bash
   ./build.sh -t v1.0.1
   ```

2. 更新 docker-compose.yml 中的镜像标签（如果使用固定标签）

3. 重启服务：
   ```bash
   docker-compose up -d --force-recreate
   ```

### 零停机更新

使用滚动更新策略：

1. 构建新版本镜像
2. 启动新容器（使用不同端口）
3. 验证新版本正常
4. 切换流量到新容器
5. 停止旧容器

## 安全建议

1. **修改默认密钥**: 生产环境必须修改 `JWT_SECRET` 和 `API_KEY_HASH_SALT`
2. **限制 CORS**: 生产环境建议指定具体的 `CORS_ALLOW_ORIGINS`，避免使用 `*`
3. **使用 HTTPS**: 建议在 Nginx 前使用反向代理（如 Traefik、Nginx）提供 HTTPS
4. **定期更新**: 定期更新基础镜像和依赖包
5. **日志监控**: 配置日志收集和监控系统

## 性能优化

1. **静态资源缓存**: Nginx 已配置静态资源长期缓存
2. **Gzip 压缩**: 已启用 Gzip 压缩
3. **数据库连接池**: 后端已配置数据库连接池
4. **文件上传限制**: 通过 `MAX_FILE_SIZE` 限制上传文件大小

## 支持

如有问题，请查看项目主 README 或提交 Issue。

