<template>
  <div class="home-container">
    <!-- 顶部导航 -->
    <PublicHeader />

    <!-- 主内容区 -->
    <main class="main-content">
      <div class="content-wrapper">
        <!-- 页面标题 -->
        <div class="page-header">
          <h1>测试概览</h1>
          <p class="page-description">实时查看DragonOS项目的CI测试结果和状态</p>
        </div>

        <!-- Master分支统计卡片 -->
        <div class="stats-card-wrapper">
          <MasterStatsCard ref="masterStatsCard" />
        </div>

        <!-- 筛选卡片 -->
        <div class="filter-card">
          <t-card>
            <t-form
              :data="testRunStore.filters"
              @submit="handleSearch"
              class="filter-form"
            >
              <div class="filter-content">
                <div class="filter-fields">
                  <t-form-item label="分支名称" class="filter-item">
                    <t-input
                      v-model="testRunStore.filters.branch"
                      placeholder="例如：main、develop"
                      clearable
                      class="filter-input"
                    />
                  </t-form-item>
                  <t-form-item label="提交哈希" class="filter-item">
                    <t-input
                      v-model="testRunStore.filters.commitId"
                      placeholder="输入Commit哈希值"
                      clearable
                      class="filter-input"
                    />
                  </t-form-item>
                  <t-form-item label="测试状态" class="filter-item">
                    <t-select
                      v-model="testRunStore.filters.status"
                      clearable
                      placeholder="选择状态"
                      class="filter-select"
                    >
                      <t-option value="all" label="全部状态" />
                      <t-option value="passed" label="通过" />
                      <t-option value="failed" label="失败" />
                      <t-option value="running" label="运行中" />
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
                <span class="list-count"
                  >共 {{ testRunStore.total }} 条记录</span
                >
              </div>
              <t-space>
                <t-button
                  variant="outline"
                  theme="default"
                  @click="refreshData"
                >
                  <t-icon name="refresh" />
                  刷新
                </t-button>
              </t-space>
            </div>

            <t-loading :loading="testRunStore.loading">
              <t-table
                :data="testRunStore.testRuns"
                :columns="columns"
                :pagination="paginationConfig"
                @page-change="handlePageChange"
                @page-size-change="handlePageSizeChange"
                hover
                row-key="id"
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
                <template #created_at="{ row }">
                  <div class="time-cell">
                    {{ formatTime(row.created_at) }}
                  </div>
                </template>
                <template #operation="{ row }">
                  <t-button
                    variant="text"
                    theme="warning"
                    size="small"
                    @click="viewDetail(row.id)"
                  >
                    查看详情
                    <t-icon name="chevron-right" />
                  </t-button>
                </template>
              </t-table>
            </t-loading>
          </t-card>
        </div>
      </div>
    </main>

    <!-- 页脚 -->
    <Footer />
  </div>
</template>

<script setup>
import { onMounted, computed, ref } from "vue";
import { useRouter } from "vue-router";
import { useTestRunStore } from "@/stores/testRun";
import { MessagePlugin } from "tdesign-vue-next";
import MasterStatsCard from "@/components/MasterStatsCard.vue";
import PublicHeader from "@/components/PublicHeader.vue";
import Footer from "@/components/Footer.vue";

const router = useRouter();
const testRunStore = useTestRunStore();
const masterStatsCard = ref(null);

const columns = [
  { colKey: "id", title: "ID", width: 80 },
  { colKey: "branch_name", title: "分支名", width: 150 },
  { colKey: "commit_short_id", title: "Commit ID", width: 140 },
  { colKey: "test_type", title: "测试类型", width: 120 },
  { colKey: "status", title: "状态", width: 120 },
  { colKey: "created_at", title: "创建时间", width: 200 },
  { colKey: "operation", title: "操作", width: 120, fixed: "right" },
];

const emptyConfig = {
  description: "暂无测试运行记录",
  icon: "inbox",
};

const paginationConfig = computed(() => ({
  current: testRunStore.pagination.page,
  pageSize: testRunStore.pagination.pageSize,
  total: testRunStore.total,
  showJumper: true,
  showSizer: true,
}));

const getStatusTheme = (status) => {
  const themes = {
    passed: "success",
    failed: "danger",
    running: "warning",
    cancelled: "default",
  };
  return themes[status] || "default";
};

const getStatusText = (status) => {
  const texts = {
    passed: "通过",
    failed: "失败",
    running: "运行中",
    cancelled: "已取消",
  };
  return texts[status] || status;
};

const handleSearch = () => {
  testRunStore.pagination.page = 1;
  testRunStore.fetchTestRuns();
};

const handleReset = () => {
  testRunStore.resetFilters();
  testRunStore.fetchTestRuns();
};

const handlePageChange = (page) => {
  testRunStore.setPagination(page, testRunStore.pagination.pageSize);
  testRunStore.fetchTestRuns();
};

const handlePageSizeChange = (pageSize) => {
  testRunStore.setPagination(1, pageSize);
  testRunStore.fetchTestRuns();
};

const viewDetail = (id) => {
  router.push(`/test-runs/${id}`);
};

const getStatusIcon = (status) => {
  const icons = {
    passed: "check-circle",
    failed: "close-circle",
    running: "time",
    cancelled: "stop-circle",
  };
  return icons[status] || "question-circle";
};

const refreshData = () => {
  testRunStore.fetchTestRuns();
  if (masterStatsCard.value) {
    masterStatsCard.value.refresh();
  }
  MessagePlugin.success("数据已刷新");
};

const formatTime = (timeStr) => {
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

onMounted(() => {
  testRunStore.fetchTestRuns();
});
</script>

<style scoped>
/* 整体布局 */
.home-container {
  min-height: 100vh;
  background-color: #f9fafb;
  display: flex;
  flex-direction: column;
}

/* 主内容区 */
.main-content {
  padding: 32px;
  flex: 1;
}

.content-wrapper {
  max-width: 1400px;
  margin: 0 auto;
}

/* 页面标题 */
.page-header {
  margin-bottom: 32px;
}

.page-header h1 {
  font-size: 28px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 8px;
}

.page-description {
  font-size: 14px;
  color: #6b7280;
}

/* 统计卡片 */
.stats-card-wrapper {
  margin-bottom: 32px;
}

/* 卡片样式 */
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

/* 筛选区域 */
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

.search-btn,
.reset-btn {
  min-width: 100px;
  height: 40px;
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.2s ease;
  display: flex !important;
  align-items: center !important;
  justify-content: center !important;
  gap: 6px;
}

.search-btn :deep(.t-button__text),
.reset-btn :deep(.t-button__text) {
  display: flex !important;
  align-items: center !important;
  justify-content: center !important;
  gap: 6px;
  line-height: 1;
}

.search-btn :deep(.t-icon),
.reset-btn :deep(.t-icon) {
  font-size: 16px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.search-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(245, 158, 11, 0.3);
}

.reset-btn {
  border-color: #e5e7eb;
}

.reset-btn:hover {
  border-color: #d1d5db;
  background-color: #f9fafb;
}

/* 列表头部 */
.list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

/* 刷新按钮样式 */
.list-header :deep(.t-button) {
  display: flex !important;
  align-items: center !important;
  justify-content: center !important;
  gap: 6px;
  height: 40px;
  border-radius: 8px;
  font-weight: 500;
  border-color: #e5e7eb;
  transition: all 0.2s ease;
}

.list-header :deep(.t-button__text) {
  display: flex !important;
  align-items: center !important;
  justify-content: center !important;
  gap: 6px;
  line-height: 1;
}

.list-header :deep(.t-button .t-icon) {
  font-size: 16px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.list-header :deep(.t-button:hover) {
  border-color: #d1d5db;
  background-color: #f9fafb;
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

/* 表格样式 */
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

/* 操作按钮样式 */
:deep(.t-table .t-button--variant-text) {
  display: flex !important;
  align-items: center !important;
  justify-content: center !important;
  gap: 4px;
  line-height: 1;
}

:deep(.t-table .t-button--variant-text .t-button__text) {
  display: flex !important;
  align-items: center !important;
  justify-content: center !important;
  gap: 4px;
  line-height: 1;
}

:deep(.t-table .t-button--variant-text .t-icon) {
  font-size: 14px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* 状态标签 */
:deep(.t-tag--light-warning) {
  background-color: #fef3c7;
  color: #d97706;
  border-color: #f59e0b;
}

:deep(.t-tag--light-success) {
  background-color: #d1fae5;
  color: #065f46;
  border-color: #10b981;
}

:deep(.t-tag--light-danger) {
  background-color: #fee2e2;
  color: #991b1b;
  border-color: #ef4444;
}

/* 响应式 */
@media (max-width: 768px) {
  .main-content {
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
