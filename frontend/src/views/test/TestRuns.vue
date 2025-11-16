<template>
  <div class="page-container">
    <div class="page-header">
      <h1 class="page-title">测试运行管理</h1>
      <p class="page-description">
        管理和查看所有测试运行任务，支持删除和设置可见性
      </p>
    </div>

    <!-- 筛选卡片 -->
    <div class="filter-card">
      <t-card>
        <t-form :data="filters" @submit="handleSearch" class="filter-form">
          <div class="filter-content">
            <div class="filter-fields">
              <t-form-item label="分支名称" class="filter-item">
                <t-input
                  v-model="filters.branch"
                  placeholder="例如：main、develop"
                  clearable
                  class="filter-input"
                />
              </t-form-item>
              <t-form-item label="提交哈希" class="filter-item">
                <t-input
                  v-model="filters.commit_id"
                  placeholder="输入Commit哈希值"
                  clearable
                  class="filter-input"
                />
              </t-form-item>
              <t-form-item label="测试类型" class="filter-item">
                <t-select
                  v-model="filters.test_type"
                  clearable
                  placeholder="选择测试类型"
                  class="filter-select"
                >
                  <t-option value="" label="全部类型" />
                  <t-option value="gvisor" label="gvisor" />
                </t-select>
              </t-form-item>
              <t-form-item label="测试状态" class="filter-item">
                <t-select
                  v-model="filters.status"
                  clearable
                  placeholder="选择状态"
                  class="filter-select"
                >
                  <t-option value="" label="全部状态" />
                  <t-option value="passed" label="通过" />
                  <t-option value="failed" label="失败" />
                  <t-option value="running" label="运行中" />
                  <t-option value="cancelled" label="已取消" />
                </t-select>
              </t-form-item>
            </div>
            <div class="filter-actions">
              <t-button theme="warning" type="submit" class="search-btn">
                <t-icon name="search" />
                搜索
              </t-button>
              <t-button
                variant="outline"
                theme="default"
                @click="handleReset"
                class="reset-btn"
              >
                <t-icon name="refresh" />
                重置
              </t-button>
            </div>
          </div>
        </t-form>
      </t-card>
    </div>

    <!-- 测试运行列表 -->
    <div class="list-card">
      <t-card>
        <div class="list-header">
          <div class="list-title">
            <h3>测试运行记录</h3>
            <span class="list-count">共 {{ total }} 条记录</span>
          </div>
          <t-space>
            <t-button variant="outline" theme="default" @click="refreshData">
              <t-icon name="refresh" />
              刷新
            </t-button>
          </t-space>
        </div>

        <t-loading :loading="loading">
          <t-table
            :data="testRuns"
            :columns="columns"
            :pagination="paginationConfig"
            @page-change="handlePageChange"
            @page-size-change="handlePageSizeChange"
            hover
            row-key="id"
            size="medium"
            :empty="emptyConfig"
          >
            <template #branch_name="{ row }">
              <div class="branch-cell">
                <t-icon name="code-branch" />
                <span>{{ row.branch_name }}</span>
              </div>
            </template>
            <template #commit_short_id="{ row }">
              <div class="commit-cell">
                <t-icon name="commit" />
                <code>{{ row.commit_short_id }}</code>
              </div>
            </template>
            <template #test_type="{ row }">
              <t-tag theme="primary" variant="light" shape="round">
                {{ row.test_type }}
              </t-tag>
            </template>
            <template #status="{ row }">
              <t-tag
                :theme="getStatusTheme(row.status)"
                variant="light"
                shape="round"
              >
                <t-icon :name="getStatusIcon(row.status)" />
                {{ getStatusText(row.status) }}
              </t-tag>
            </template>
            <template #is_public="{ row }">
              <t-tag
                :theme="row.is_public ? 'success' : 'warning'"
                variant="light"
                shape="round"
              >
                <t-icon
                  :name="row.is_public ? 'browse' : 'browse-off'"
                  style="margin-right: 4px"
                />
                {{ row.is_public ? "公开" : "私有" }}
              </t-tag>
            </template>
            <template #created_at="{ row }">
              <div class="time-cell">
                {{ formatTime(row.created_at) }}
              </div>
            </template>
            <template #operation="{ row }">
              <t-space>
                <t-button
                  variant="text"
                  theme="primary"
                  size="small"
                  @click="toggleVisibility(row)"
                >
                  <t-icon
                    :name="row.is_public ? 'browse-off' : 'browse'"
                    style="margin-right: 4px"
                  />
                  {{ row.is_public ? "设为私有" : "设为公开" }}
                </t-button>
                <t-button
                  variant="text"
                  theme="danger"
                  size="small"
                  @click="handleDelete(row)"
                >
                  <t-icon name="delete" style="margin-right: 4px" />
                  删除
                </t-button>
              </t-space>
            </template>
          </t-table>
        </t-loading>
      </t-card>
    </div>

    <!-- 删除确认对话框 -->
    <t-dialog
      v-model:visible="deleteDialogVisible"
      header="确认删除"
      :confirm-btn="{ theme: 'danger' }"
      @confirm="confirmDelete"
      @cancel="cancelDelete"
    >
      <div class="delete-dialog-content">
        <t-alert theme="error" message="这是一个危险操作，删除后无法恢复！" />
        <div class="delete-info">
          <p><strong>测试运行ID：</strong>{{ deleteTarget?.id }}</p>
          <p><strong>分支名称：</strong>{{ deleteTarget?.branch_name }}</p>
          <p><strong>提交哈希：</strong>{{ deleteTarget?.commit_short_id }}</p>
          <p><strong>测试类型：</strong>{{ deleteTarget?.test_type }}</p>
        </div>
        <p class="delete-warning">
          ⚠️ 删除此测试运行将同时删除所有关联的测例和输出文件，此操作不可逆！
        </p>
      </div>
    </t-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { MessagePlugin } from "tdesign-vue-next";
import {
  getTestRunsAdmin,
  deleteTestRun,
  updateTestRunVisibility,
  type TestRun,
} from "@/api/admin";

const loading = ref(false);
const testRuns = ref<TestRun[]>([]);
const total = ref(0);
const pagination = ref({
  page: 1,
  pageSize: 20,
});

const filters = ref({
  branch: "",
  commit_id: "",
  test_type: "",
  status: "",
});

const deleteDialogVisible = ref(false);
const deleteTarget = ref<TestRun | null>(null);

const columns = [
  { colKey: "id", title: "ID", width: 80 },
  { colKey: "branch_name", title: "分支名", width: 150 },
  { colKey: "commit_short_id", title: "Commit ID", width: 140 },
  { colKey: "test_type", title: "测试类型", width: 120 },
  { colKey: "status", title: "状态", width: 120 },
  { colKey: "is_public", title: "可见性", width: 100 },
  { colKey: "created_at", title: "创建时间", width: 200 },
  { colKey: "operation", title: "操作", width: 200, fixed: "right" },
];

const emptyConfig = {
  description: "暂无测试运行记录",
  icon: "inbox",
};

const paginationConfig = computed(() => ({
  current: pagination.value.page,
  pageSize: pagination.value.pageSize,
  total: total.value,
  showJumper: true,
  showSizer: true,
}));

const getStatusTheme = (status: string) => {
  const themes: Record<string, string> = {
    passed: "success",
    failed: "danger",
    running: "warning",
    cancelled: "default",
  };
  return themes[status] || "default";
};

const getStatusText = (status: string) => {
  const texts: Record<string, string> = {
    passed: "通过",
    failed: "失败",
    running: "运行中",
    cancelled: "已取消",
  };
  return texts[status] || status;
};

const getStatusIcon = (status: string) => {
  const icons: Record<string, string> = {
    passed: "check-circle",
    failed: "close-circle",
    running: "time",
    cancelled: "stop-circle",
  };
  return icons[status] || "question-circle";
};

const fetchTestRuns = async () => {
  loading.value = true;
  try {
    const params: any = {
      page: pagination.value.page,
      page_size: pagination.value.pageSize,
    };
    if (filters.value.branch) {
      params.branch = filters.value.branch;
    }
    if (filters.value.commit_id) {
      params.commit_id = filters.value.commit_id;
    }
    if (filters.value.test_type) {
      params.test_type = filters.value.test_type;
    }
    if (filters.value.status) {
      params.status = filters.value.status;
    }

    const response = await getTestRunsAdmin(params);
    testRuns.value = response.data.test_runs || [];
    total.value = response.data.total || 0;
  } catch (error: any) {
    MessagePlugin.error(
      error.response?.data?.message || "获取测试运行列表失败",
    );
  } finally {
    loading.value = false;
  }
};

const handleSearch = () => {
  pagination.value.page = 1;
  fetchTestRuns();
};

const handleReset = () => {
  filters.value = {
    branch: "",
    commit_id: "",
    test_type: "",
    status: "",
  };
  pagination.value.page = 1;
  fetchTestRuns();
};

const handlePageChange = (page: number) => {
  pagination.value.page = page;
  fetchTestRuns();
};

const handlePageSizeChange = (pageSize: number) => {
  pagination.value.pageSize = pageSize;
  pagination.value.page = 1;
  fetchTestRuns();
};

const refreshData = () => {
  fetchTestRuns();
  MessagePlugin.success("数据已刷新");
};

const formatTime = (timeStr: string) => {
  const date = new Date(timeStr);
  return date.toLocaleString("zh-CN", {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
  });
};

const toggleVisibility = async (row: TestRun) => {
  try {
    await updateTestRunVisibility(row.id, !row.is_public);
    MessagePlugin.success(`已${row.is_public ? "设为私有" : "设为公开"}`);
    // 更新本地数据
    row.is_public = !row.is_public;
  } catch (error: any) {
    MessagePlugin.error(error.response?.data?.message || "更新可见性失败");
  }
};

const handleDelete = (row: TestRun) => {
  deleteTarget.value = row;
  deleteDialogVisible.value = true;
};

const confirmDelete = async () => {
  if (!deleteTarget.value) return;

  try {
    await deleteTestRun(deleteTarget.value.id);
    MessagePlugin.success("删除成功");
    deleteDialogVisible.value = false;
    deleteTarget.value = null;
    // 如果当前页没有数据了，返回上一页
    if (testRuns.value.length === 1 && pagination.value.page > 1) {
      pagination.value.page--;
    }
    fetchTestRuns();
  } catch (error: any) {
    MessagePlugin.error(error.response?.data?.message || "删除失败");
  }
};

const cancelDelete = () => {
  deleteDialogVisible.value = false;
  deleteTarget.value = null;
};

onMounted(() => {
  fetchTestRuns();
});
</script>

<style scoped>
.page-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 24px;
}

.page-header {
  margin-bottom: 32px;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 8px;
}

.page-description {
  font-size: 14px;
  color: #6b7280;
}

.filter-card,
.list-card {
  margin-bottom: 24px;
}

.filter-card :deep(.t-card),
.list-card :deep(.t-card) {
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  border: none;
  overflow: hidden;
  transition: all 0.3s ease;
}

.filter-card :deep(.t-card:hover),
.list-card :deep(.t-card:hover) {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.filter-card :deep(.t-card__body),
.list-card :deep(.t-card__body) {
  padding: 24px;
}

.filter-form {
  margin: 0;
}

.filter-content {
  display: flex;
  align-items: flex-end;
  gap: 24px;
  flex-wrap: wrap;
}

.filter-fields {
  display: flex;
  gap: 20px;
  flex: 1;
  min-width: 0;
  flex-wrap: wrap;
}

.filter-item {
  flex: 1;
  min-width: 200px;
  margin-bottom: 0;
}

.filter-item :deep(.t-form-item__label) {
  font-size: 14px;
  font-weight: 500;
  color: #374151;
  margin-bottom: 8px;
}

.filter-input,
.filter-select {
  width: 100%;
}

.filter-input :deep(.t-input),
.filter-select :deep(.t-select) {
  border-radius: 8px;
  border: 1px solid #e5e7eb;
  transition: all 0.2s ease;
}

.filter-input :deep(.t-input:hover),
.filter-select :deep(.t-select:hover) {
  border-color: #f59e0b;
}

.filter-input :deep(.t-input--focused),
.filter-select :deep(.t-select--focused) {
  border-color: #f59e0b;
  box-shadow: 0 0 0 3px rgba(245, 158, 11, 0.1);
}

.filter-actions {
  display: flex;
  gap: 12px;
  align-items: flex-end;
  flex-shrink: 0;
}

.search-btn {
  min-width: 100px;
  height: 40px;
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.2s ease;
}

.search-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(245, 158, 11, 0.3);
}

.reset-btn {
  min-width: 100px;
  height: 40px;
  border-radius: 8px;
  font-weight: 500;
  border-color: #e5e7eb;
  transition: all 0.2s ease;
}

.reset-btn:hover {
  border-color: #d1d5db;
  background-color: #f9fafb;
}

.list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.list-title {
  display: flex;
  align-items: baseline;
  gap: 12px;
}

.list-title h3 {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.list-count {
  font-size: 14px;
  color: #9ca3af;
}

:deep(.t-table) {
  background-color: #ffffff;
}

:deep(.t-table th) {
  background-color: #fafafa;
  font-weight: 600;
  color: #374151;
  font-size: 14px;
}

:deep(.t-table td) {
  border-bottom: 1px solid #f3f4f6;
  padding: 16px;
}

:deep(.t-table tr:hover td) {
  background-color: #fef9f3;
}

.branch-cell,
.commit-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.branch-cell span {
  font-weight: 500;
  color: #1f2937;
}

.commit-cell code {
  font-family: "Monaco", "Menlo", monospace;
  font-size: 13px;
  background: #f3f4f6;
  padding: 4px 8px;
  border-radius: 4px;
  color: #1f2937;
}

.time-cell {
  color: #6b7280;
  font-size: 13px;
}

.delete-dialog-content {
  padding: 16px 0;
}

.delete-info {
  margin: 16px 0;
  padding: 16px;
  background-color: #f9fafb;
  border-radius: 8px;
}

.delete-info p {
  margin: 8px 0;
  font-size: 14px;
  color: #374151;
}

.delete-warning {
  margin-top: 16px;
  padding: 12px;
  background-color: #fef2f2;
  border-left: 4px solid #ef4444;
  border-radius: 4px;
  color: #991b1b;
  font-size: 14px;
  line-height: 1.6;
}

/* 响应式 */
@media (max-width: 768px) {
  .page-container {
    padding: 16px;
  }

  .filter-content {
    flex-direction: column;
    align-items: stretch;
    gap: 16px;
  }

  .filter-fields {
    flex-direction: column;
    gap: 16px;
  }

  .filter-item {
    min-width: 100%;
  }

  .filter-actions {
    width: 100%;
    justify-content: stretch;
  }

  .search-btn,
  .reset-btn {
    flex: 1;
  }

  .list-header {
    flex-direction: column;
    gap: 16px;
    align-items: flex-start;
  }
}
</style>
