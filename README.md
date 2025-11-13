# DragonOS CI Dashboard

CI测试结果Dashboard系统，用于展示和管理DragonOS项目的CI测试结果。

## 功能特性

- 📊 测试结果可视化展示
- 🔍 多维度检索（分支、Commit ID、时间范围、状态等）
- 📁 原始输出文件查看和下载
- 🔐 API Key认证机制
- 👤 后台管理系统
- 🐳 Docker一键部署

## 技术栈

### 后端
- Go 1.21+
- Gin Web框架
- GORM ORM
- MySQL 5.7+
- JWT认证

### 前端
- Vue 3 + Composition API
- TDesign Vue Next
- Pinia状态管理
- Vue Router
- Vite构建工具

## 快速开始

### 使用Docker Compose（推荐）

1. 复制环境变量文件：
```bash
cp .env.example .env
```

2. 修改`.env`文件中的配置（特别是数据库密码和JWT密钥）

3. 启动服务：
```bash
docker-compose up -d
```

4. 访问应用：
- 前端: http://localhost:3000
- 后端API: http://localhost:8080/api/v1

### 本地开发

#### 后端开发

1. 进入后端目录：
```bash
cd backend
```

2. 安装依赖：
```bash
go mod download
```

3. 配置环境变量（参考`.env.example`）

4. 运行数据库迁移：
```bash
migrate -path ./migrations -database "mysql://user:password@tcp(localhost:3306)/dragonos_ci" up
```

5. 启动服务：
```bash
go run cmd/server/main.go
```

#### 前端开发

1. 进入前端目录：
```bash
cd frontend
```

2. 安装依赖：
```bash
npm install
```

3. 启动开发服务器：
```bash
npm run dev
```

## API文档

### 公开接口

- `GET /api/v1/test-runs` - 查询测试运行记录
- `GET /api/v1/test-runs/:id` - 获取测试运行详情
- `GET /api/v1/test-runs/:id/test-cases` - 获取测例列表
- `GET /api/v1/test-runs/:id/output-files/:fileId` - 下载原始输出文件

### 受保护接口（需要API Key）

- `POST /api/v1/test-runs` - 上传测试结果
- `POST /api/v1/test-runs/:id/output-files` - 上传原始输出文件

### 管理接口（需要JWT认证）

- `POST /api/v1/admin/login` - 管理员登录
- `GET /api/v1/admin/api-keys` - 查看API密钥列表
- `POST /api/v1/admin/api-keys` - 创建API密钥
- `DELETE /api/v1/admin/api-keys/:id` - 删除API密钥

## 项目结构

```
dragonos-ci-dashboard/
├── backend/                 # Go后端服务
│   ├── cmd/
│   │   └── server/         # 主程序入口
│   ├── internal/
│   │   ├── api/            # API路由和处理器
│   │   ├── models/         # 数据模型
│   │   ├── services/       # 业务逻辑层
│   │   ├── middleware/     # 中间件
│   │   ├── config/         # 配置管理
│   │   └── storage/        # 文件存储处理
│   ├── migrations/         # 数据库迁移文件
│   └── pkg/                # 公共包
├── frontend/               # Vue3前端
│   └── src/
│       ├── views/          # 页面组件
│       ├── components/     # 通用组件
│       ├── api/            # API调用
│       ├── router/         # 路由配置
│       ├── stores/         # 状态管理
│       └── utils/          # 工具函数
└── docker-compose.yml      # Docker编排配置
```

## 开发规范

请参考项目根目录下的`.cursorrules`文件。

## 许可证

MIT License

