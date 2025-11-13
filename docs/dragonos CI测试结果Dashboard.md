# CI测试结果Dashboard系统实现计划

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
│   │   ├── middleware/     # 中间件（认证、日志等）
│   │   ├── config/         # 配置管理
│   │   └── storage/        # 文件存储处理
│   ├── migrations/         # 数据库迁移文件
│   ├── pkg/                # 公共包
│   ├── go.mod
│   └── Dockerfile
├── frontend/               # Vue3前端
│   ├── src/
│   │   ├── views/          # 页面组件
│   │   ├── components/     # 通用组件
│   │   ├── api/            # API调用
│   │   ├── router/         # 路由配置
│   │   ├── stores/         # 状态管理
│   │   └── utils/          # 工具函数
│   ├── package.json
│   └── Dockerfile
├── docker-compose.yml      # Docker编排配置
├── .env.example            # 环境变量示例
└── README.md               # 项目文档
```

## 数据库设计

### 核心表结构

1. **projects** - 项目表

   - id, name, description, created_at, updated_at

2. **test_runs** - 测试运行记录表

   - id, project_id, branch_name, commit_id, commit_short_id, test_type, status, started_at, completed_at, created_at
   - 索引：branch_name, commit_id, commit_short_id, created_at

3. **test_cases** - 测例详情表

   - id, test_run_id, name, status, duration_ms, error_log, debug_log, created_at
   - 索引：test_run_id, name, status

4. **test_output_files** - 测试输出文件表

   - id, test_run_id, filename, file_path, file_size, mime_type, created_at

5. **api_keys** - API密钥表

   - id, name, key_hash, project_id, created_at, last_used_at, expires_at

6. **users** - 用户表（后台管理）

   - id, username, password_hash, role, created_at, updated_at

## 后端实现要点

### API接口设计

**公开接口（无需认证）：**

- `GET /api/v1/test-runs` - 查询测试运行记录（支持多维度检索）
- `GET /api/v1/test-runs/:id` - 获取测试运行详情
- `GET /api/v1/test-runs/:id/test-cases` - 获取测例列表
- `GET /api/v1/test-runs/:id/output-files/:fileId` - 下载原始输出文件

**受保护接口（需要API Key）：**

- `POST /api/v1/test-runs` - 上传测试结果
- `POST /api/v1/test-runs/:id/output-files` - 上传原始输出文件

**管理接口（需要用户认证）：**

- `POST /api/v1/admin/login` - 管理员登录
- `GET /api/v1/admin/api-keys` - 查看API密钥列表
- `POST /api/v1/admin/api-keys` - 创建API密钥
- `DELETE /api/v1/admin/api-keys/:id` - 删除API密钥

### 关键技术选型

- **Web框架**: Gin
- **ORM**: GORM
- **数据库迁移**: golang-migrate/migrate
- **认证**: API Key (Bearer Token) + JWT (后台管理)
- **文件存储**: 本地文件系统（可配置路径）
- **配置管理**: viper + 环境变量

### 检索功能实现

支持以下查询参数：

- `branch` - 分支名（支持模糊匹配）
- `commit_id` - Commit ID（支持完整或短ID）
- `start_time` / `end_time` - 时间范围
- `status` - 测试状态（passed/failed/all）
- `test_case_name` - 测例名称（模糊匹配）
- `page` / `page_size` - 分页

## 前端实现要点

### 页面结构

1. **首页** (`/`) - 测试结果列表

   - 检索筛选栏（分支、Commit ID、时间范围、状态）
   - 测试运行记录表格（分页）
   - 卡片式展示，支持动画效果

2. **详情页** (`/test-runs/:id`) - 测试运行详情

   - 基本信息展示（分支、Commit、时间等）
   - 测例列表（通过/失败分组展示）
   - 失败测例的错误日志展示
   - 原始输出文件查看/下载

3. **后台管理** (`/admin/*`)

   - 登录页面
   - API密钥管理
   - 系统统计

### 技术栈

- **框架**: Vue3 + Composition API
- **UI库**: TDesign Vue Next
- **路由**: Vue Router
- **状态管理**: Pinia
- **HTTP客户端**: Axios
- **构建工具**: Vite

### UI设计要点

- 使用TDesign的现代组件和主题
- 添加平滑的过渡动画和加载效果
- 响应式设计，支持移动端
- 使用图表展示测试通过率趋势（可选）

## Docker部署

### 服务编排

- **backend**: Go服务容器
- **frontend**: Nginx容器（静态文件服务）
- **mysql**: MySQL 5.7容器
- **volumes**: 数据持久化（数据库数据、上传文件）

### 环境变量配置

- 数据库连接信息
- API密钥加密密钥
- 文件存储路径
- 服务端口配置

## 扩展性设计

1. **测试类型支持**: 

   - 当前仅实现 "gvisor" 系统调用测试
   - 使用`test_type`字段区分不同类型的测试，未来可扩展（如performance等）
   - 在服务层使用适配器模式处理不同测试类型的数据格式差异

2. **数据格式扩展**: 

   - 当前gvisor测试：标准化的test_cases结构（name, status, duration_ms, error_log, debug_log）
   - 未来新测试类型：可在test_runs表中添加JSON字段存储类型特定的元数据

3. **项目扩展**: 通过`project_id`关联，支持未来多项目

4. **插件化接口**: API设计时考虑未来不同测试类型的差异化需求，通过test_type参数路由到对应的处理逻辑

## 实施步骤

1. 初始化项目结构和依赖
2. 设计并实现数据库迁移文件
3. 实现后端核心功能（模型、服务、API）
4. 实现前端页面和交互
5. 配置Docker和部署文件
6. 编写文档和README
