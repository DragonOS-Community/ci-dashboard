.PHONY: help build run dev test clean migrate-up migrate-down install deps docker-up docker-down docker-restart

# 变量定义
BINARY_NAME=server
CMD_DIR=cmd/server
BUILD_DIR=bin
MIGRATE_CMD=migrate
MIGRATE_PATH=./migrations
DB_DSN?=mysql://root:rootpassword@tcp(localhost:3306)/dragonos_ci

# 默认目标
help: ## 显示帮助信息
	@echo "可用的命令:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

# 安装依赖
deps: ## 安装Go依赖
	@echo "安装Go依赖..."
	go mod download
	@echo "依赖安装完成"

# 整理依赖
tidy: ## 整理Go模块依赖
	@echo "整理Go模块依赖..."
	go mod tidy
	@echo "依赖整理完成"

# 构建后端
build: ## 构建后端二进制文件
	@echo "构建后端..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)
	@echo "构建完成: $(BUILD_DIR)/$(BINARY_NAME)"

# 构建admin-cli工具
build-admin-cli: ## 构建admin-cli工具
	@echo "构建admin-cli工具..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/admin-cli ./cmd/admin-cli
	@echo "构建完成: $(BUILD_DIR)/admin-cli"

# 运行后端（开发模式）
run: ## 运行后端服务
	@echo "启动后端服务..."
	go run ./$(CMD_DIR)/main.go

# 开发模式（带热重载，需要安装air）
dev: ## 开发模式运行（需要安装 air: go install github.com/cosmtrek/air@latest）
	@if ! command -v air > /dev/null; then \
		echo "错误: 未安装 air，请运行: go install github.com/cosmtrek/air@latest"; \
		exit 1; \
	fi
	air

# 测试
test: ## 运行测试
	@echo "运行测试..."
	go test -v ./...

# 测试覆盖率
test-coverage: ## 运行测试并生成覆盖率报告
	@echo "运行测试并生成覆盖率报告..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"

# 数据库迁移 - 升级
migrate-up: ## 执行数据库迁移（升级）
	@if ! command -v $(MIGRATE_CMD) > /dev/null; then \
		echo "错误: 未安装 migrate，请运行: go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"; \
		exit 1; \
	fi
	@echo "执行数据库迁移（升级）..."
	$(MIGRATE_CMD) -path $(MIGRATE_PATH) -database "$(DB_DSN)" up

# 数据库迁移 - 回滚
migrate-down: ## 回滚数据库迁移
	@if ! command -v $(MIGRATE_CMD) > /dev/null; then \
		echo "错误: 未安装 migrate，请运行: go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"; \
		exit 1; \
	fi
	@echo "回滚数据库迁移..."
	$(MIGRATE_CMD) -path $(MIGRATE_PATH) -database "$(DB_DSN)" down 1

# 创建新的迁移文件
migrate-create: ## 创建新的迁移文件（用法: make migrate-create NAME=create_users_table）
	@if ! command -v $(MIGRATE_CMD) > /dev/null; then \
		echo "错误: 未安装 migrate，请运行: go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"; \
		exit 1; \
	fi
	@if [ -z "$(NAME)" ]; then \
		echo "错误: 请提供迁移文件名，例如: make migrate-create NAME=create_users_table"; \
		exit 1; \
	fi
	@echo "创建迁移文件: $(NAME)"
	$(MIGRATE_CMD) create -ext sql -dir $(MIGRATE_PATH) -seq $(NAME)

# 格式化代码
fmt: ## 格式化Go代码
	@echo "格式化代码..."
	go fmt ./...

# 代码检查
lint: ## 运行代码检查（需要安装 golangci-lint）
	@if ! command -v golangci-lint > /dev/null; then \
		echo "错误: 未安装 golangci-lint"; \
		exit 1; \
	fi
	@echo "运行代码检查..."
	golangci-lint run

# 清理构建文件
clean: ## 清理构建文件和临时文件
	@echo "清理构建文件..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	@echo "清理完成"

# Docker相关命令（需要在项目根目录执行）
docker-up: ## 启动Docker服务（需要在项目根目录执行）
	@echo "启动Docker服务..."
	@cd .. && docker-compose up -d
	@echo "等待服务启动..."
	@sleep 5
	@echo "服务已启动"

docker-down: ## 停止Docker服务（需要在项目根目录执行）
	@echo "停止Docker服务..."
	@cd .. && docker-compose down

docker-restart: ## 重启Docker服务（需要在项目根目录执行）
	@echo "重启Docker服务..."
	@cd .. && docker-compose restart

docker-logs: ## 查看Docker日志（需要在项目根目录执行）
	@cd .. && docker-compose logs -f

docker-logs-backend: ## 查看后端日志（需要在项目根目录执行）
	@cd .. && docker-compose logs -f backend

docker-logs-frontend: ## 查看前端日志（需要在项目根目录执行）
	@cd .. && docker-compose logs -f frontend

docker-logs-mysql: ## 查看MySQL日志（需要在项目根目录执行）
	@cd .. && docker-compose logs -f mysql

# 一键启动（安装依赖 + 运行迁移 + 启动服务）
start: deps migrate-up run ## 一键启动：安装依赖、运行迁移、启动服务

# 完整开发环境启动
dev-start: deps migrate-up dev ## 开发环境启动（包含热重载）

# 检查环境
check-env: ## 检查开发环境
	@echo "检查开发环境..."
	@echo "Go版本:"
	@go version || echo "错误: 未安装Go"
	@echo ""
	@echo "Docker版本:"
	@docker --version || echo "警告: 未安装Docker"
	@echo ""
	@echo "Docker Compose版本:"
	@docker-compose --version || echo "警告: 未安装Docker Compose"
	@echo ""
	@if command -v migrate > /dev/null; then \
		echo "Migrate已安装"; \
	else \
		echo "警告: 未安装 migrate，运行迁移前请先安装"; \
	fi
	@if command -v air > /dev/null; then \
		echo "Air已安装"; \
	else \
		echo "提示: 未安装 air，可以使用 'make dev' 需要先安装"; \
	fi

