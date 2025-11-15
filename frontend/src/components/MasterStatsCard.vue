<template>
  <t-card class="master-stats-card" :class="{ 'loading': loading }">
    <div class="stats-header">
      <div class="header-left">
        <div class="branch-badge">
          <t-icon name="code-branch" />
          <span>{{ stats?.branch_name || 'master' }}</span>
        </div>
        <div class="test-type-badge">
          <t-icon name="file-code" />
          <span>{{ stats?.test_type || 'gvisor' }} 系统调用测试</span>
        </div>
      </div>
      <div class="header-right">
        <t-tag
          v-if="stats"
          :theme="getStatusTheme(stats.status)"
          variant="light"
          shape="round"
        >
          <t-icon :name="getStatusIcon(stats.status)" />
          {{ getStatusText(stats.status) }}
        </t-tag>
        <t-tag
          v-else-if="!loading && hasNoData"
          theme="default"
          variant="light"
          shape="round"
        >
          <t-icon name="info-circle" />
          等待数据
        </t-tag>
      </div>
    </div>

    <t-loading :loading="loading" size="large">
      <div v-if="stats" class="stats-content">
        <!-- 主要指标 -->
        <div class="main-metrics">
          <div class="metric-item pass-rate">
            <div class="metric-label">通过率</div>
            <div class="metric-value">
              <span class="value-number">{{ formatPassRate(stats.pass_rate) }}</span>
              <span class="value-unit">%</span>
            </div>
            <div class="metric-progress">
              <t-progress
                :percentage="stats.pass_rate"
                :theme="getPassRateTheme(stats.pass_rate)"
                :label="false"
                :size="'medium'"
              />
            </div>
          </div>

          <div class="metric-item total-cases">
            <div class="metric-label">通过测例数</div>
            <div class="metric-value">
              <span class="value-number">{{ formatNumber(stats.passed_cases) }}</span>
            </div>
            <div class="metric-detail">
              <span class="detail-item">
                共 {{ formatNumber(stats.total_cases) }} 个测例
              </span>
              <span class="detail-item failed" v-if="stats.failed_cases > 0">
                <t-icon name="close-circle" />
                {{ formatNumber(stats.failed_cases) }} 失败
              </span>
              <span class="detail-item skipped" v-if="stats.skipped_cases > 0">
                <t-icon name="jump" />
                {{ formatNumber(stats.skipped_cases) }} 跳过
              </span>
            </div>
          </div>

          <div class="metric-item duration">
            <div class="metric-label">总耗时</div>
            <div class="metric-value">
              <span class="value-number">{{ formatDuration(stats.duration) }}</span>
            </div>
          </div>
        </div>

        <!-- 详细信息 -->
        <div class="stats-footer">
          <div class="footer-item">
            <t-icon name="commit" />
            <span class="footer-label">Commit:</span>
            <span class="footer-value">{{ stats.commit_short_id }}</span>
          </div>
          <div class="footer-item">
            <t-icon name="time" />
            <span class="footer-label">测试时间:</span>
            <span class="footer-value">{{ formatTime(stats.created_at) }}</span>
          </div>
          <div class="footer-item">
            <t-button
              variant="text"
              theme="warning"
              size="small"
              @click="viewDetail"
            >
              查看详情
              <t-icon name="chevron-right" />
            </t-button>
          </div>
        </div>
      </div>

      <div v-else-if="!loading && hasNoData" class="no-data-content">
        <div class="no-data-icon">
          <t-icon name="chart" size="48px" />
        </div>
        <div class="no-data-title">暂无测试数据</div>
        <div class="no-data-description">
          <p>master 分支还没有测试运行记录</p>
          <p class="no-data-hint">测试数据将在首次 CI 运行后显示</p>
        </div>
      </div>
    </t-loading>
  </t-card>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getMasterBranchStats } from '@/api/testRun'
import type { MasterBranchStats } from '@/api/testRun'

const router = useRouter()
const stats = ref<MasterBranchStats | null>(null)
const loading = ref(false)
const hasNoData = ref(false)

type TestRunStatus = 'passed' | 'failed' | 'running' | 'cancelled'

const getStatusTheme = (status: string): string => {
  const themes: Record<TestRunStatus, string> = {
    passed: 'success',
    failed: 'danger',
    running: 'warning',
    cancelled: 'default',
  }
  return themes[status as TestRunStatus] || 'default'
}

const getStatusText = (status: string): string => {
  const texts: Record<TestRunStatus, string> = {
    passed: '通过',
    failed: '失败',
    running: '运行中',
    cancelled: '已取消',
  }
  return texts[status as TestRunStatus] || status
}

const getStatusIcon = (status: string): string => {
  const icons: Record<TestRunStatus, string> = {
    passed: 'check-circle',
    failed: 'close-circle',
    running: 'time',
    cancelled: 'stop-circle',
  }
  return icons[status as TestRunStatus] || 'question-circle'
}

const getPassRateTheme = (rate: number): string => {
  if (rate >= 95) return 'success'
  if (rate >= 80) return 'warning'
  return 'danger'
}

const formatPassRate = (rate: number): string => {
  return rate.toFixed(1)
}

const formatNumber = (num: number): string => {
  return num.toLocaleString()
}

const formatDuration = (ms: number): string => {
  if (ms < 1000) return `${ms}ms`
  if (ms < 60000) return `${(ms / 1000).toFixed(1)}s`
  const minutes = Math.floor(ms / 60000)
  const seconds = Math.floor((ms % 60000) / 1000)
  if (minutes < 60) {
    return seconds > 0 ? `${minutes}m ${seconds}s` : `${minutes}m`
  }
  const hours = Math.floor(minutes / 60)
  const mins = minutes % 60
  return mins > 0 ? `${hours}h ${mins}m` : `${hours}h`
}

const formatTime = (timeStr: string): string => {
  const date = new Date(timeStr)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}

const viewDetail = () => {
  if (stats.value) {
    router.push(`/test-runs/${stats.value.test_run_id}`)
  }
}

const fetchStats = async () => {
  loading.value = true
  hasNoData.value = false
  try {
    const res = await getMasterBranchStats()
    stats.value = res.data
    hasNoData.value = false
  } catch (error: any) {
    // 静默处理 404 错误，显示友好的提示
    if (error?.response?.status === 404 || error?.response?.status === 400) {
      stats.value = null
      hasNoData.value = true
    } else {
      // 其他错误也显示无数据提示，而不是错误信息
      console.error('Failed to fetch master branch stats:', error)
      stats.value = null
      hasNoData.value = true
    }
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchStats()
})

defineExpose({
  refresh: fetchStats,
})
</script>

<style scoped>
.master-stats-card {
  border-radius: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  border: none;
  overflow: hidden;
  background: linear-gradient(135deg, #ffffff 0%, #fef9f3 100%);
  transition: all 0.3s ease;
}

.master-stats-card:hover {
  box-shadow: 0 4px 16px rgba(245, 158, 11, 0.15);
  transform: translateY(-2px);
}

.master-stats-card :deep(.t-card__body) {
  padding: 32px;
}

.stats-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 32px;
  padding-bottom: 24px;
  border-bottom: 2px solid #f3f4f6;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.branch-badge,
.test-type-badge {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background: linear-gradient(135deg, #fef3c7 0%, #fde68a 100%);
  border-radius: 20px;
  font-size: 14px;
  font-weight: 600;
  color: #92400e;
}

.test-type-badge {
  background: linear-gradient(135deg, #dbeafe 0%, #bfdbfe 100%);
  color: #1e40af;
}

.header-right {
  display: flex;
  align-items: center;
}

.stats-content {
  display: flex;
  flex-direction: column;
  gap: 32px;
}

.main-metrics {
  display: grid;
  grid-template-columns: 2fr 2fr 1fr;
  gap: 32px;
}

.metric-item {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.metric-label {
  font-size: 14px;
  color: #6b7280;
  font-weight: 500;
}

.metric-value {
  display: flex;
  align-items: baseline;
  gap: 4px;
}

.value-number {
  font-size: 36px;
  font-weight: 700;
  color: #1f2937;
  line-height: 1;
}

.pass-rate .value-number {
  background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.total-cases .value-number {
  color: #3b82f6;
}

.duration .value-number {
  color: #10b981;
}

.value-unit {
  font-size: 20px;
  font-weight: 600;
  color: #9ca3af;
}

.metric-progress {
  margin-top: 8px;
}

.metric-detail {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  margin-top: 8px;
}

.detail-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  font-weight: 500;
}

.detail-item.success {
  color: #059669;
}

.detail-item.failed {
  color: #dc2626;
}

.detail-item.skipped {
  color: #6b7280;
}

.stats-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 24px;
  border-top: 1px solid #f3f4f6;
  flex-wrap: wrap;
  gap: 16px;
}

.footer-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #6b7280;
}

.footer-label {
  font-weight: 500;
}

.footer-value {
  color: #1f2937;
  font-weight: 600;
  font-family: 'Monaco', 'Menlo', monospace;
}

.no-data-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 32px;
  text-align: center;
}

.no-data-icon {
  width: 80px;
  height: 80px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #fef3c7 0%, #fde68a 100%);
  border-radius: 50%;
  margin-bottom: 24px;
  color: #d97706;
}

.no-data-title {
  font-size: 20px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 12px;
}

.no-data-description {
  color: #6b7280;
  font-size: 14px;
  line-height: 1.6;
}

.no-data-description p {
  margin: 4px 0;
}

.no-data-hint {
  color: #9ca3af;
  font-size: 13px;
  margin-top: 8px;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .main-metrics {
    grid-template-columns: 1fr 1fr;
  }

  .duration {
    grid-column: 1 / -1;
  }
}

@media (max-width: 768px) {
  .master-stats-card :deep(.t-card__body) {
    padding: 24px;
  }

  .stats-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }

  .main-metrics {
    grid-template-columns: 1fr;
    gap: 24px;
  }

  .stats-footer {
    flex-direction: column;
    align-items: flex-start;
  }

  .value-number {
    font-size: 28px;
  }
}
</style>

