<template>
  <div class="detail-container">
    <!-- 顶部导航 -->
    <header class="header">
      <div class="header-content">
        <div class="logo" @click="goBack" style="cursor: pointer;">
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
        <!-- 页面标题和操作栏 -->
        <div class="page-header">
          <div class="page-header-left">
            <t-button variant="text" theme="default" @click="goBack" class="back-btn">
              <t-icon name="chevron-left" />
              返回列表
            </t-button>
            <div class="page-title">
              <h1>测试运行详情</h1>
              <span v-if="testRunStore.currentTestRun" class="test-run-id">#{{ testRunStore.currentTestRun.id }}</span>
            </div>
          </div>
          <div class="page-header-actions">
            <t-button variant="outline" theme="default" @click="refreshData" class="refresh-btn">
              <t-icon name="refresh" />
              刷新
            </t-button>
          </div>
        </div>

        <t-loading :loading="testRunStore.loading">
          <!-- 基本信息卡片 -->
          <div v-if="testRunStore.currentTestRun" class="info-card">
            <t-card>
              <div class="card-header">
                <h2 class="card-title">基本信息</h2>
                <t-tag
                  :theme="getStatusTheme(testRunStore.currentTestRun.status)"
                  variant="light"
                  shape="round"
                  class="status-tag"
                >
                  <t-icon :name="getStatusIcon(testRunStore.currentTestRun.status)" />
                  {{ getStatusText(testRunStore.currentTestRun.status) }}
                </t-tag>
              </div>
              <div class="info-grid">
                <div class="info-item">
                  <div class="info-label">
                    <t-icon name="code-branch" />
                    分支名称
                  </div>
                  <div class="info-value">{{ testRunStore.currentTestRun.branch_name }}</div>
                </div>
                <div class="info-item">
                  <div class="info-label">
                    <t-icon name="commit" />
                    提交哈希
                  </div>
                  <div class="info-value code-value">
                    <code>{{ testRunStore.currentTestRun.commit_id }}</code>
                    <t-button
                      variant="text"
                      theme="primary"
                      size="small"
                      @click="copyCommitId"
                      class="copy-btn"
                    >
                      <t-icon name="file-copy" />
                    </t-button>
                  </div>
                </div>
                <div class="info-item">
                  <div class="info-label">
                    <t-icon name="setting" />
                    测试类型
                  </div>
                  <div class="info-value">
                    <t-tag theme="primary" variant="light" shape="round">
                      {{ testRunStore.currentTestRun.test_type }}
                    </t-tag>
                  </div>
                </div>
                <div class="info-item">
                  <div class="info-label">
                    <t-icon name="time" />
                    创建时间
                  </div>
                  <div class="info-value">{{ formatTime(testRunStore.currentTestRun.created_at) }}</div>
                </div>
                <div class="info-item" v-if="testRunStore.currentTestRun.started_at">
                  <div class="info-label">
                    <t-icon name="play-circle" />
                    开始时间
                  </div>
                  <div class="info-value">{{ formatTime(testRunStore.currentTestRun.started_at) }}</div>
                </div>
                <div class="info-item" v-if="testRunStore.currentTestRun.completed_at">
                  <div class="info-label">
                    <t-icon name="check-circle" />
                    完成时间
                  </div>
                  <div class="info-value">{{ formatTime(testRunStore.currentTestRun.completed_at) }}</div>
                </div>
              </div>
            </t-card>
          </div>

          <!-- 测例统计卡片 -->
          <div class="stats-card" v-if="testCases.length > 0">
            <t-card>
              <div class="stats-grid">
                <div class="stat-item">
                  <div class="stat-value">{{ allTestCases.length }}</div>
                  <div class="stat-label">总测例数</div>
                </div>
                <div class="stat-item success">
                  <div class="stat-value">{{ passedTestCases.length }}</div>
                  <div class="stat-label">通过</div>
                </div>
                <div class="stat-item danger">
                  <div class="stat-value">{{ failedTestCases.length }}</div>
                  <div class="stat-label">失败</div>
                </div>
                <div class="stat-item">
                  <div class="stat-value">{{ passRate }}%</div>
                  <div class="stat-label">通过率</div>
                </div>
              </div>
            </t-card>
          </div>

          <!-- 测例列表卡片 -->
          <div class="test-cases-card">
            <t-card>
              <div class="card-header">
                <h2 class="card-title">测例列表</h2>
                <t-tabs v-model="activeTab" class="filter-tabs">
                  <t-tab-panel value="all" :label="`全部 (${allTestCases.length})`">
                  </t-tab-panel>
                  <t-tab-panel value="passed" :label="`通过 (${passedTestCases.length})`">
                  </t-tab-panel>
                  <t-tab-panel value="failed" :label="`失败 (${failedTestCases.length})`">
                  </t-tab-panel>
                </t-tabs>
              </div>
              <div class="test-cases-toolbar">
                <div class="toolbar-left">
                  <t-input
                    v-model="searchKeyword"
                    placeholder="搜索测例名称..."
                    clearable
                    class="search-input"
                  >
                    <template #prefix-icon>
                      <t-icon name="search" />
                    </template>
                  </t-input>
                  <span v-if="searchKeyword.trim()" class="search-result-info">
                    找到 {{ filteredTestCases.length }} 个结果
                  </span>
                </div>
              </div>
              <div class="test-cases-content">
                <TestCaseList
                  v-if="!searchKeyword.trim() || filteredTestCases.length > 0"
                  :test-cases="paginatedTestCases"
                  :sort-field="sortField"
                  :sort-order="sortOrder"
                  @sort-change="handleSortChange"
                />
                <t-empty
                  v-if="searchKeyword.trim() && filteredTestCases.length === 0"
                  description="未找到匹配的测例"
                  :icon="'search'"
                  class="empty-state"
                />
              </div>
              <div class="test-cases-pagination" v-if="filteredTestCases.length > pageSize">
                <t-pagination
                  v-model="currentPage"
                  :total="filteredTestCases.length"
                  :page-size="pageSize"
                  :page-size-options="[10, 20, 50, 100]"
                  show-sizer
                  show-jumper
                  @page-size-change="handlePageSizeChange"
                />
              </div>
            </t-card>
          </div>

          <!-- 输出文件卡片 -->
          <div class="files-card">
            <t-card>
              <div class="card-header">
                <h2 class="card-title">输出文件</h2>
                <span class="file-count" v-if="files.length > 0">共 {{ files.length }} 个文件</span>
              </div>
              <div class="files-content">
                <t-list v-if="files.length > 0" class="file-list">
                  <t-list-item v-for="file in files" :key="file.id" class="file-item">
                    <div class="file-info">
                      <div class="file-name">
                        <t-icon name="file" />
                        <span>{{ file.filename }}</span>
                      </div>
                      <div class="file-meta">
                        <t-tag variant="light" size="small">{{ formatFileSize(file.size) }}</t-tag>
                      </div>
                    </div>
                    <div class="file-actions">
                      <t-button
                        theme="warning"
                        variant="outline"
                        size="small"
                        @click="downloadFile(file)"
                        class="download-btn"
                      >
                        <t-icon name="download" />
                        下载
                      </t-button>
                    </div>
                  </t-list-item>
                </t-list>
                <t-empty v-else description="暂无输出文件" :icon="'inbox'" />
              </div>
            </t-card>
          </div>
        </t-loading>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useTestRunStore } from '@/stores/testRun'
import { getTestCasesByTestRunId, getFilesByTestRunId, downloadFile as downloadFileAPI } from '@/api/testRun'
import { MessagePlugin } from 'tdesign-vue-next'
import TestCaseList from '@/components/TestCaseList.vue'

const route = useRoute()
const router = useRouter()
const testRunStore = useTestRunStore()

const activeTab = ref('all')
const testCases = ref([])
const files = ref([])
const searchKeyword = ref('')
const currentPage = ref(1)
const pageSize = ref(20)
const sortField = ref('') // 排序字段：'name' 或 'duration_ms'
const sortOrder = ref('') // 排序方向：'asc' 或 'desc'

const allTestCases = computed(() => testCases.value)
const passedTestCases = computed(() => testCases.value.filter(tc => tc.status === 'passed'))
const failedTestCases = computed(() => testCases.value.filter(tc => tc.status === 'failed'))

// 根据标签页过滤的测例
const tabFilteredTestCases = computed(() => {
  switch (activeTab.value) {
    case 'passed':
      return passedTestCases.value
    case 'failed':
      return failedTestCases.value
    default:
      return allTestCases.value
  }
})

// 根据搜索关键词过滤的测例
const filteredTestCases = computed(() => {
  let result = tabFilteredTestCases.value
  
  // 搜索过滤
  if (searchKeyword.value.trim()) {
    const keyword = searchKeyword.value.trim().toLowerCase()
    result = result.filter(tc => {
      return tc.name && tc.name.toLowerCase().includes(keyword)
    })
  }
  
  // 排序
  if (sortField.value && sortOrder.value) {
    result = [...result].sort((a, b) => {
      let aValue, bValue
      
      if (sortField.value === 'name') {
        // 按名称排序（字符串）
        aValue = (a.name || '').toLowerCase()
        bValue = (b.name || '').toLowerCase()
        if (sortOrder.value === 'asc') {
          return aValue.localeCompare(bValue, 'zh-CN')
        } else {
          return bValue.localeCompare(aValue, 'zh-CN')
        }
      } else if (sortField.value === 'duration_ms') {
        // 按耗时排序（数字）
        aValue = a.duration_ms || 0
        bValue = b.duration_ms || 0
        if (sortOrder.value === 'asc') {
          return aValue - bValue
        } else {
          return bValue - aValue
        }
      }
      return 0
    })
  }
  
  return result
})

// 分页后的测例
const paginatedTestCases = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredTestCases.value.slice(start, end)
})

// 监听搜索关键词变化，重置到第一页
watch(searchKeyword, () => {
  currentPage.value = 1
})

// 监听标签页变化，重置到第一页
watch(activeTab, () => {
  currentPage.value = 1
})

// 监听排序变化，重置到第一页
watch([sortField, sortOrder], () => {
  currentPage.value = 1
})

// 处理每页显示数量变化
const handlePageSizeChange = (size) => {
  pageSize.value = size
  currentPage.value = 1
}

// 处理排序变化
const handleSortChange = (sortInfo) => {
  // TDesign 表格的 sort-change 事件可能返回数组或对象
  let sortData = null
  if (Array.isArray(sortInfo) && sortInfo.length > 0) {
    sortData = sortInfo[0]
  } else if (sortInfo && typeof sortInfo === 'object') {
    sortData = sortInfo
  }
  
  if (sortData && sortData.sortBy) {
    sortField.value = sortData.sortBy
    sortOrder.value = sortData.descending ? 'desc' : 'asc'
  } else {
    sortField.value = ''
    sortOrder.value = ''
  }
}

const passRate = computed(() => {
  if (allTestCases.value.length === 0) return 0
  return Math.round((passedTestCases.value.length / allTestCases.value.length) * 100)
})

const getStatusText = (status) => {
  const texts = {
    passed: '通过',
    failed: '失败',
    running: '运行中',
    cancelled: '已取消',
  }
  return texts[status] || status
}

const getStatusTheme = (status) => {
  const themes = {
    passed: 'success',
    failed: 'danger',
    running: 'warning',
    cancelled: 'default',
  }
  return themes[status] || 'default'
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

const formatTime = (timeStr) => {
  if (!timeStr) return '-'
  const date = new Date(timeStr)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  })
}

const formatFileSize = (bytes) => {
  if (!bytes) return '0 B'
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(2) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(2) + ' MB'
}

const copyCommitId = async () => {
  if (!testRunStore.currentTestRun) return
  try {
    await navigator.clipboard.writeText(testRunStore.currentTestRun.commit_id)
    MessagePlugin.success('已复制到剪贴板')
  } catch (error) {
    MessagePlugin.error('复制失败')
  }
}

const downloadFile = async (file) => {
  if (!testRunStore.currentTestRun) return
  try {
    const response = await downloadFileAPI(testRunStore.currentTestRun.id, file.id)
    const url = window.URL.createObjectURL(new Blob([response.data]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', file.filename)
    document.body.appendChild(link)
    link.click()
    link.remove()
    window.URL.revokeObjectURL(url)
    MessagePlugin.success('下载成功')
  } catch (error) {
    MessagePlugin.error('下载失败')
  }
}

const goBack = () => {
  router.push('/')
}

const goToLogin = () => {
  router.push('/admin/login')
}

const refreshData = async () => {
  await fetchData()
  MessagePlugin.success('数据已刷新')
}

const fetchData = async () => {
  const id = parseInt(route.params.id)
  await testRunStore.fetchTestRunById(id)

  try {
    const testCasesRes = await getTestCasesByTestRunId(id)
    testCases.value = testCasesRes.data || []

    const filesRes = await getFilesByTestRunId(id)
    files.value = filesRes.data || []
  } catch (error) {
    console.error('Failed to fetch data:', error)
    MessagePlugin.error('加载数据失败')
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
/* 整体布局 */
.detail-container {
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
  transition: opacity 0.2s ease;
}

.logo:hover {
  opacity: 0.8;
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

/* 页面标题区域 */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 32px;
  flex-wrap: wrap;
  gap: 16px;
}

.page-header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.back-btn {
  color: #6b7280;
  transition: all 0.2s ease;
}

.back-btn:hover {
  color: #f59e0b;
  background-color: #fef9f3;
}

.page-title {
  display: flex;
  align-items: baseline;
  gap: 12px;
}

.page-title h1 {
  font-size: 28px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.test-run-id {
  font-size: 16px;
  color: #9ca3af;
  font-weight: 400;
}

.page-header-actions {
  display: flex;
  gap: 12px;
}

.refresh-btn {
  border-radius: 8px;
  transition: all 0.2s ease;
}

.refresh-btn:hover {
  border-color: #f59e0b;
  color: #f59e0b;
}

/* 主内容区 */
.main-content {
  padding: 32px;
}

.content-wrapper {
  max-width: 1400px;
  margin: 0 auto;
}

/* 卡片通用样式 */
.info-card,
.stats-card,
.test-cases-card,
.files-card {
  margin-bottom: 24px;
}

.info-card :deep(.t-card),
.stats-card :deep(.t-card),
.test-cases-card :deep(.t-card),
.files-card :deep(.t-card) {
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  border: none;
  overflow: hidden;
  transition: all 0.3s ease;
}

.info-card :deep(.t-card:hover),
.stats-card :deep(.t-card:hover),
.test-cases-card :deep(.t-card:hover),
.files-card :deep(.t-card:hover) {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.info-card :deep(.t-card__body),
.stats-card :deep(.t-card__body),
.test-cases-card :deep(.t-card__body),
.files-card :deep(.t-card__body) {
  padding: 24px;
}

/* 卡片头部 */
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  flex-wrap: wrap;
  gap: 16px;
}

.card-title {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.status-tag {
  font-weight: 500;
}

/* 基本信息网格 */
.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 24px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.info-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  font-weight: 500;
  color: #6b7280;
}

.info-label :deep(.t-icon) {
  font-size: 16px;
  color: #9ca3af;
}

.info-value {
  font-size: 15px;
  color: #1f2937;
  font-weight: 500;
}

.code-value {
  display: flex;
  align-items: center;
  gap: 8px;
}

.code-value code {
  font-family: 'Monaco', 'Menlo', 'Consolas', monospace;
  font-size: 13px;
  background: #f3f4f6;
  padding: 6px 12px;
  border-radius: 6px;
  color: #1f2937;
  flex: 1;
  word-break: break-all;
}

.copy-btn {
  flex-shrink: 0;
  color: #6b7280;
  transition: all 0.2s ease;
}

.copy-btn:hover {
  color: #f59e0b;
  background-color: #fef9f3;
}

/* 统计卡片 */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 24px;
}

.stat-item {
  text-align: center;
  padding: 20px;
  background: #fafafa;
  border-radius: 8px;
  transition: all 0.2s ease;
}

.stat-item:hover {
  background: #f3f4f6;
  transform: translateY(-2px);
}

.stat-item.success {
  background: #d1fae5;
}

.stat-item.success:hover {
  background: #a7f3d0;
}

.stat-item.danger {
  background: #fee2e2;
}

.stat-item.danger:hover {
  background: #fecaca;
}

.stat-value {
  font-size: 32px;
  font-weight: 700;
  color: #1f2937;
  margin-bottom: 8px;
}

.stat-item.success .stat-value {
  color: #065f46;
}

.stat-item.danger .stat-value {
  color: #991b1b;
}

.stat-label {
  font-size: 14px;
  color: #6b7280;
  font-weight: 500;
}

/* 测例列表 */
.filter-tabs {
  flex-shrink: 0;
}

.test-cases-toolbar {
  margin-top: 16px;
  margin-bottom: 16px;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-wrap: wrap;
}

.search-input {
  max-width: 400px;
}

.search-result-info {
  font-size: 14px;
  color: #6b7280;
  white-space: nowrap;
}

.test-cases-content {
  margin-top: 16px;
}

.empty-state {
  margin: 40px 0;
}

.test-cases-pagination {
  margin-top: 24px;
  display: flex;
  justify-content: center;
}

/* 文件列表 */
.file-count {
  font-size: 14px;
  color: #9ca3af;
}

.files-content {
  margin-top: 16px;
}

.file-list {
  background: transparent;
}

.file-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  border-radius: 8px;
  transition: all 0.2s ease;
  border: 1px solid transparent;
}

.file-item:hover {
  background-color: #fef9f3;
  border-color: #fde68a;
}

.file-info {
  display: flex;
  flex-direction: column;
  gap: 8px;
  flex: 1;
  min-width: 0;
}

.file-name {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 15px;
  font-weight: 500;
  color: #1f2937;
}

.file-name :deep(.t-icon) {
  color: #9ca3af;
  font-size: 18px;
}

.file-name span {
  word-break: break-all;
}

.file-meta {
  display: flex;
  gap: 8px;
}

.file-actions {
  flex-shrink: 0;
  margin-left: 16px;
}

.download-btn {
  border-radius: 8px;
  transition: all 0.2s ease;
}

.download-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(245, 158, 11, 0.3);
}

/* 表格样式优化 */
.test-cases-content :deep(.t-table) {
  background-color: #ffffff;
}

.test-cases-content :deep(.t-table th) {
  background-color: #fafafa;
  font-weight: 600;
  color: #374151;
  font-size: 14px;
}

.test-cases-content :deep(.t-table td) {
  border-bottom: 1px solid #f3f4f6;
  padding: 16px;
}

.test-cases-content :deep(.t-table tr:hover td) {
  background-color: #fef9f3;
}

/* 状态标签样式 */
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

:deep(.t-tag--light-warning) {
  background-color: #fef3c7;
  color: #d97706;
  border-color: #f59e0b;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .header-content {
    padding: 0 16px;
  }

  .main-content {
    padding: 16px;
  }

  .page-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .page-title h1 {
    font-size: 24px;
  }

  .info-grid {
    grid-template-columns: 1fr;
    gap: 20px;
  }

  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 16px;
  }

  .stat-value {
    font-size: 24px;
  }

  .card-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .file-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }

  .file-actions {
    margin-left: 0;
    width: 100%;
  }

  .download-btn {
    width: 100%;
  }
}
</style>

