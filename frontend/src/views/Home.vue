<template>
  <div class="home-container">
    <!-- 顶部导航 -->
    <header class="header">
      <div class="header-content">
        <div class="logo">
          <div class="logo-icon">🐉</div>
          <span class="logo-text">DragonOS CI Dashboard</span>
        </div>
        <t-button theme="warning" variant="outline" @click="goToLogin">
          <t-icon name="user" />
          管理员登录
        </t-button>
      </div>
    </header>

    <!-- 主内容区 -->
    <main class="main-content">
      <div class="content-wrapper">
        <!-- 页面标题 -->
        <div class="page-header">
          <h1>测试运行列表</h1>
          <p class="page-description">查看所有测试运行记录和状态</p>
        </div>

        <!-- 筛选卡片 -->
        <div class="filter-card">
          <t-card>
            <div class="filter-header">
              <h3>筛选条件</h3>
              <t-button variant="text" @click="handleReset">
                <t-icon name="refresh" />
                重置
              </t-button>
            </div>
            <t-form :data="testRunStore.filters" @submit="handleSearch" class="filter-form">
              <t-row :gutter="16">
                <t-col :span="6">
                  <t-form-item label="分支名">
                    <t-input
                      v-model="testRunStore.filters.branch"
                      placeholder="输入分支名"
                      clearable
                      prefix-icon="search"
                    />
                  </t-form-item>
                </t-col>
                <t-col :span="6">
                  <t-form-item label="Commit ID">
                    <t-input
                      v-model="testRunStore.filters.commitId"
                      placeholder="输入Commit ID"
                      clearable
                      prefix-icon="code"
                    />
                  </t-form-item>
                </t-col>
                <t-col :span="6">
                  <t-form-item label="状态">
                    <t-select v-model="testRunStore.filters.status" clearable placeholder="选择状态">
                      <t-option value="all" label="全部" />
                      <t-option value="passed" label="通过" />
                      <t-option value="failed" label="失败" />
                      <t-option value="running" label="运行中" />
                    </t-select>
                  </t-form-item>
                </t-col>
                <t-col :span="6">
                  <t-form-item>
                    <t-space>
                      <t-button theme="warning" type="submit">
                        <t-icon name="search" />
                        搜索
                      </t-button>
                    </t-space>
                  </t-form-item>
                </t-col>
              </t-row>
            </t-form>
          </t-card>
        </div>

        <!-- 测试运行列表 -->
        <div class="list-card">
          <t-card>
            <div class="list-header">
              <div class="list-title">
                <h3>测试运行记录</h3>
                <span class="list-count">共 {{ testRunStore.total }} 条记录</span>
              </div>
              <t-button variant="text" theme="warning" @click="exportData">
                <t-icon name="download" />
                导出
              </t-button>
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
              >
                <template #status="{ row }">
                  <t-tag
                    :theme="getStatusTheme(row.status)"
                    variant="light"
                    :shape="'rounded'"
                  >
                    <t-icon :name="getStatusIcon(row.status)" />
                    {{ getStatusText(row.status) }}
                  </t-tag>
                </template>
                <template #operation="{ row }">
                  <t-button
                    variant="text"
                    theme="warning"
                    @click="viewDetail(row.id)"
                  >
                    查看详情
                  </t-button>
                </template>
              </t-table>
            </t-loading>
          </t-card>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useTestRunStore } from '@/stores/testRun'
import { MessagePlugin } from 'tdesign-vue-next'

const router = useRouter()
const testRunStore = useTestRunStore()

const columns = [
  { colKey: 'id', title: 'ID', width: 80 },
  { colKey: 'branch_name', title: '分支名', width: 150 },
  { colKey: 'commit_short_id', title: 'Commit ID', width: 120 },
  { colKey: 'test_type', title: '测试类型', width: 100 },
  { colKey: 'status', title: '状态', width: 100 },
  { colKey: 'created_at', title: '创建时间', width: 180 },
  { colKey: 'operation', title: '操作', width: 100, fixed: 'right' },
]

const paginationConfig = computed(() => ({
  current: testRunStore.pagination.page,
  pageSize: testRunStore.pagination.pageSize,
  total: testRunStore.total,
  showJumper: true,
  showSizer: true,
}))

const getStatusTheme = (status) => {
  const themes = {
    passed: 'success',
    failed: 'danger',
    running: 'warning',
    cancelled: 'default',
  }
  return themes[status] || 'default'
}

const getStatusText = (status) => {
  const texts = {
    passed: '通过',
    failed: '失败',
    running: '运行中',
    cancelled: '已取消',
  }
  return texts[status] || status
}

const handleSearch = () => {
  testRunStore.pagination.page = 1
  testRunStore.fetchTestRuns()
}

const handleReset = () => {
  testRunStore.resetFilters()
  testRunStore.fetchTestRuns()
}

const handlePageChange = (page) => {
  testRunStore.setPagination(page, testRunStore.pagination.pageSize)
  testRunStore.fetchTestRuns()
}

const handlePageSizeChange = (pageSize) => {
  testRunStore.setPagination(1, pageSize)
  testRunStore.fetchTestRuns()
}

const viewDetail = (id) => {
  router.push(`/test-runs/${id}`)
}

const goToLogin = () => {
  router.push('/admin/login')
}

const getStatusIcon = (status) => {
  const icons = {
    passed: 'check-circle',
    failed: 'close-circle',
    running: 'time',
    cancelled: 'stop-circle',
  }
  return icons[status] || 'question-circle'
}

const exportData = () => {
  // 导出功能待实现
  console.log('导出数据')
}

onMounted(() => {
  testRunStore.fetchTestRuns()
})
</script>

<style scoped>
/* 整体布局 */
.home-container {
  min-height: 100vh;
  background-color: #f9fafb;
}

/* 顶部导航 */
.header {
  background-color: #ffffff;
  border-bottom: 1px solid #e5e7eb;
  position: sticky;
  top: 0;
  z-index: 100;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.header-content {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 32px;
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo-icon {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #fcd34d 0%, #f59e0b 100%);
  border-radius: 10px;
  font-size: 24px;
}

.logo-text {
  font-size: 20px;
  font-weight: 600;
  color: #1f2937;
}

/* 主内容区 */
.main-content {
  padding: 32px;
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
}

.filter-card :deep(.t-card__body),
.list-card :deep(.t-card__body) {
  padding: 24px;
}

/* 筛选区域 */
.filter-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.filter-header h3 {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.filter-form {
  margin-top: 16px;
}

/* 列表头部 */
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

/* 表格样式 */
:deep(.t-table) {
  background-color: #ffffff;
}

:deep(.t-table th) {
  background-color: #f9fafb;
  font-weight: 600;
  color: #374151;
}

:deep(.t-table td) {
  border-bottom: 1px solid #f3f4f6;
}

:deep(.t-table tr:hover td) {
  background-color: #fef3c7;
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
  .header-content {
    padding: 0 16px;
  }

  .main-content {
    padding: 16px;
  }

  .filter-header,
  .list-header {
    flex-direction: column;
    gap: 16px;
    align-items: flex-start;
  }
}
</style>

