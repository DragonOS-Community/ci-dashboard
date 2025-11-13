# 上传测试结果 API 文档

## 概述

本文档描述如何通过 API 上传测试结果到 DragonOS CI Dashboard。所有上传接口都需要使用 API Key 进行认证。

**Base URL**: `http://your-domain/api/v1`

**认证方式**: Bearer Token (API Key)

---

## 1. 创建测试运行（上传测试结果）

创建一次测试运行记录，可以同时上传测试用例结果。

### 接口信息

- **URL**: `/test-runs`
- **方法**: `POST`
- **认证**: 需要 API Key
- **Content-Type**: `application/json`

### 请求头

```
Authorization: Bearer YOUR_API_KEY
Content-Type: application/json
```

### 请求参数

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `branch_name` | string | 是 | Git分支名称，如 `main`、`dev` |
| `commit_id` | string | 是 | Commit ID（最少8位，支持完整或短ID） |
| `test_type` | string | 否 | 测试类型，默认为 `gvisor`（仅支持 `gvisor`） |
| `status` | string | 否 | 测试运行状态：`passed`、`failed`、`running`、`cancelled` |
| `test_cases` | array | 否 | 测试用例列表（见下表） |

**说明**：
- `commit_short_id` 由系统自动从 `commit_id` 截取前10位生成，无需传递
- `project_id` 使用默认项目ID（DragonOS项目ID为1），无需传递

#### test_cases 字段说明

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `name` | string | 是 | 测试用例名称 |
| `status` | string | 是 | 测试状态：`passed`、`failed`、`skipped` |
| `duration_ms` | number | 否 | 执行时长（毫秒） |
| `error_log` | string | 否 | 错误日志内容（最大长度2048字符） |
| `debug_log` | string | 否 | 调试日志内容（最大长度2048字符） |

### 请求示例

```bash
curl -X POST "http://localhost:8080/api/v1/test-runs" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "branch_name": "main",
    "commit_id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0",
    "test_type": "gvisor",
    "test_cases": [
      {
        "name": "test_basic_functionality",
        "status": "passed",
        "duration_ms": 2500,
        "debug_log": "Test completed successfully"
      },
      {
        "name": "test_edge_cases",
        "status": "failed",
        "duration_ms": 1500,
        "error_log": "Assertion failed at line 42",
        "debug_log": "Expected value: 100, Actual: 99"
      },
      {
        "name": "test_performance",
        "status": "skipped",
        "duration_ms": 0
      }
    ],
    "status": "passed"
  }'
```

### 响应格式

#### 成功响应 (200)

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 123,
    "project_id": 1,
    "branch_name": "main",
    "commit_id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0",
    "commit_short_id": "a1b2c3d4e5",
    "test_type": "gvisor",
    "status": "failed",
    "started_at": "2024-01-15T10:00:00Z",
    "completed_at": "2024-01-15T10:05:00Z",
    "created_at": "2024-01-15T10:00:00Z",
    "updated_at": "2024-01-15T10:05:00Z"
  }
}
```

#### 错误响应 (400/401/500)

```json
{
  "code": 400,
  "message": "Invalid request parameters",
  "data": null
}
```

### 状态码说明

- `200`: 成功创建测试运行
- `400`: 请求参数错误（可能原因：commit_id少于8位、test_type不是gvisor、日志长度超过2048字符）
- `401`: 未授权（API Key无效或缺失）
- `500`: 服务器内部错误

---

## 2. 上传测试输出文件（可选）

为已创建的测试运行上传输出文件（如日志文件、截图等）。

### 接口信息

- **URL**: `/test-runs/{test_run_id}/output-files`
- **方法**: `POST`
- **认证**: 需要 API Key
- **Content-Type**: `multipart/form-data`

### 请求头

```
Authorization: Bearer YOUR_API_KEY
Content-Type: multipart/form-data
```

### 路径参数

| 参数 | 类型 | 说明 |
|------|------|------|
| `test_run_id` | number | 测试运行ID（从创建测试运行接口返回） |

### 请求参数

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `file` | file | 是 | 要上传的文件 |

### 请求示例

```bash
curl -X POST "http://your-domain/api/v1/test-runs/123/output-files" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -F "file=@/path/to/test-output.log"
```

### 响应格式

#### 成功响应 (200)

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 456,
    "test_run_id": 123,
    "filename": "test-output.log",
    "file_path": "/storage/test-runs/123/test-output.log",
    "file_size": 10240,
    "mime_type": "text/plain",
    "created_at": "2024-01-15T10:05:30Z"
  }
}
```

#### 错误响应

```json
{
  "code": 400,
  "message": "File size exceeds limit",
  "data": null
}
```

---

## 完整使用流程示例

### 步骤1: 创建测试运行并上传结果

```bash
# 创建测试运行
RESPONSE=$(curl -X POST "http://your-domain/api/v1/test-runs" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "branch_name": "main",
    "commit_id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0",
    "test_type": "gvisor",
    "test_cases": [
      {
        "name": "test_basic",
        "status": "passed",
        "duration_ms": 2000
      }
    ]
  }')

# 提取测试运行ID
TEST_RUN_ID=$(echo $RESPONSE | jq -r '.data.id')
echo "Test Run ID: $TEST_RUN_ID"
```

### 步骤2: 上传测试输出文件（可选）

```bash
# 上传文件
curl -X POST "http://your-domain/api/v1/test-runs/$TEST_RUN_ID/output-files" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -F "file=@test-output.log"
```

---

## 注意事项

1. **API Key 获取**: 需要在后台管理系统创建 API Key
2. **项目ID**: 系统自动使用默认 DragonOS 项目ID（ID为1），无需传递
3. **Commit ID**: 
   - 最少8位字符
   - 系统会自动截取前10位作为 `commit_short_id`
   - 支持完整40位SHA-1或短ID
4. **测试类型**: `test_type` 仅支持 `gvisor`，其他值将被拒绝
5. **日志长度限制**: `error_log` 和 `debug_log` 最大长度为2048字符，超出将被截断或拒绝
6. **文件大小限制**: 上传文件大小受服务器配置限制（默认配置请查看配置文件）
7. **状态自动推断**: 如果不指定 `status`，系统会根据 `test_cases` 的状态自动推断：
   - 有任何一个 `failed` → `failed`
   - 全部 `passed` → `passed`
   - 其他情况 → `passed`（默认）
8. **时间戳**: `started_at` 和 `completed_at` 由系统自动设置

---

## 错误码参考

| 错误码 | 说明 | 解决方案 |
|--------|------|----------|
| 200 | 成功 | - |
| 400 | 请求参数错误 | 检查必填字段和参数格式 |
| 401 | 未授权 | 检查 API Key 是否正确 |
| 404 | 资源不存在 | 检查测试运行ID是否存在 |
| 500 | 服务器错误 | 联系管理员 |

---

## 快速开始

最简单的上传示例（仅创建测试运行）：

```bash
curl -X POST "http://your-domain/api/v1/test-runs" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "branch_name": "main",
    "commit_id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0"
  }'
```

