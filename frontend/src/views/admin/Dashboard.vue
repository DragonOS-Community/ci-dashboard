<template>
  <div class="dashboard">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1 class="page-title">仪表盘</h1>
      <p class="page-description">查看系统运行状态和测试数据概览</p>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-grid">
      <div class="stat-card" v-for="stat in stats" :key="stat.key">
        <div class="stat-icon" :style="{ backgroundColor: stat.bgColor }">
          <t-icon :name="stat.icon" :style="{ color: stat.color }" />
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ stat.value }}</div>
          <div class="stat-label">{{ stat.label }}</div>
          <div class="stat-trend" :class="stat.trend.type">
            <t-icon :name="stat.trend.icon" />
            {{ stat.trend.text }}
          </div>
        </div>
      </div>
    </div>

    <!-- 图表区域 -->
    <div class="charts-grid">
      <!-- 测试趋势图 -->
      <div class="chart-card">
        <div class="card-header">
          <h3>测试执行趋势</h3>
          <t-radio-group v-model="trendPeriod" variant="default-filled">
            <t-radio-button value="7d">7天</t-radio-button>
            <t-radio-button value="30d">30天</t-radio-button>
            <t-radio-button value="90d">90天</t-radio-button>
          </t-radio-group>
        </div>
        <div class="chart-container">
          <!-- 这里放置图表组件，可以使用 ECharts 或其他图表库 -->
          <div class="chart-placeholder">
            <div class="chart-line"></div>
            <div class="chart-area"></div>
          </div>
        </div>
      </div>

      <!-- 成功率统计 -->
      <div class="chart-card">
        <div class="card-header">
          <h3>成功率统计</h3>
          <t-button variant="text" size="small">
            <t-icon name="download" /> 导出
          </t-button>
        </div>
        <div class="chart-container">
          <div class="success-rate">
            <div class="rate-circle">
              <div class="rate-value">{{ successRate }}%</div>
            </div>
            <div class="rate-details">
              <div class="rate-item">
                <span class="dot success"></span>
                <span>成功: {{ successCount }}</span>
              </div>
              <div class="rate-item">
                <span class="dot failed"></span>
                <span>失败: {{ failedCount }}</span>
              </div>
              <div class="rate-item">
                <span class="dot skipped"></span>
                <span>跳过: {{ skippedCount }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 最近测试任务 -->
    <div class="recent-tests-card">
      <div class="card-header">
        <h3>最近测试任务</h3>
        <router-link to="/admin/test/runs">
          <t-button variant="text" theme="primary">查看全部</t-button>
        </router-link>
      </div>
      <div class="table-container">
        <t-table
          :data="recentTests"
          :columns="tableColumns"
          hover
          :pagination="false"
          size="medium"
        >
          <template #status="{ row }">
            <t-tag :theme="getStatusTheme(row.status)" variant="light">
              {{ row.status }}
            </t-tag>
          </template>
          <template #operation="{ row }">
            <t-button variant="text" theme="primary" @click="viewDetail(row.id)">
              查看详情
            </t-button>
          </template>
        </t-table>
      </div>
    </div>

    <!-- 快捷操作 -->
    <div class="quick-actions">
      <div class="action-card" v-for="action in quickActions" :key="action.key">
        <div class="action-icon" :style="{ backgroundColor: action.bgColor }">
          <t-icon :name="action.icon" :style="{ color: action.color }" />
        </div>
        <div class="action-content">
          <h4>{{ action.title }}</h4>
          <p>{{ action.description }}</p>
          <t-button :theme="action.theme" @click="action.handler">
            {{ action.buttonText }}
          </t-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useTestRunStore } from '@/stores/testRun'

const router = useRouter()
const testRunStore = useTestRunStore()

const trendPeriod = ref('7d')

// 统计数据
const stats = ref([
  {
    key: 'total',
    label: '总测试数',
    value: '1,234',
    icon: 'chart',
    color: '#f59e0b',
    bgColor: '#fef3c7',
    trend: {
      type: 'up',
      icon: 'arrow-up',
      text: '12.5%'
    }
  },
  {
    key: 'today',
    label: '今日运行',
    value: '56',
    icon: 'play-circle',
    color: '#10b981',
    bgColor: '#d1fae5',
    trend: {
      type: 'up',
      icon: 'arrow-up',
      text: '8.2%'
    }
  },
  {
    key: 'success',
    label: '成功率',
    value: '98.5%',
    icon: 'check-circle',
    color: '#10b981',
    bgColor: '#d1fae5',
    trend: {
      type: 'up',
      icon: 'arrow-up',
      text: '2.1%'
    }
  },
  {
    key: 'avgTime',
    label: '平均耗时',
    value: '15.3m',
    icon: 'time',
    color: '#3b82f6',
    bgColor: '#dbeafe',
    trend: {
      type: 'down',
      icon: 'arrow-down',
      text: '5.3%'
    }
  }
])

// 成功率相关
const successRate = ref(98.5)
const successCount = ref(1215)
const failedCount = ref(19)
const skippedCount = ref(5)

// 最近测试数据
const recentTests = ref([
  {
    id: 1,
    branch: 'main',
    commitId: 'abc123',
    status: 'success',
    duration: '12m 30s',
    startTime: '2024-01-01 10:00:00'
  },
  {
    id: 2,
    branch: 'dev',
    commitId: 'def456',
    status: 'failed',
    duration: '8m 15s',
    startTime: '2024-01-01 09:30:00'
  },
  {
    id: 3,
    branch: 'feature/test',
    commitId: 'ghi789',
    status: 'running',
    duration: '-',
    startTime: '2024-01-01 11:00:00'
  }
])

const tableColumns = [
  { colKey: 'id', title: 'ID', width: 80 },
  { colKey: 'branch', title: '分支', width: 120 },
  { colKey: 'commitId', title: 'Commit ID', width: 100 },
  { colKey: 'status', title: '状态', width: 100 },
  { colKey: 'duration', title: '执行时间', width: 120 },
  { colKey: 'startTime', title: '开始时间', width: 180 },
  { colKey: 'operation', title: '操作', width: 100, fixed: 'right' }
]

// 快捷操作
const quickActions = ref([
  {
    key: 'newTest',
    title: '创建测试任务',
    description: '快速创建新的测试运行任务',
    icon: 'add',
    color: '#ffffff',
    bgColor: '#f59e0b',
    theme: 'warning',
    buttonText: '立即创建',
    handler: () => {
      router.push('/admin/test/runs?action=new')
    }
  },
  {
    key: 'apiKey',
    title: '管理 API 密钥',
    description: '创建或管理 API 访问密钥',
    icon: 'key',
    color: '#ffffff',
    bgColor: '#3b82f6',
    theme: 'primary',
    buttonText: '管理密钥',
    handler: () => {
      router.push('/admin/system/api-keys')
    }
  },
  {
    key: 'report',
    title: '生成测试报告',
    description: '导出最近一段时间的测试报告',
    icon: 'file-excel',
    color: '#ffffff',
    bgColor: '#10b981',
    theme: 'success',
    buttonText: '生成报告',
    handler: () => {
      router.push('/admin/test/reports')
    }
  }
])

const getStatusTheme = (status) => {
  const themes = {
    success: 'success',
    failed: 'danger',
    running: 'warning',
    cancelled: 'default'
  }
  return themes[status] || 'default'
}

const viewDetail = (id) => {
  router.push(`/test-runs/${id}`)
}

onMounted(() => {
  // 加载最近测试数据
  testRunStore.fetchTestRuns({ limit: 10 }).then(data => {
    if (data?.items) {
      recentTests.value = data.items.map(item => ({
        id: item.id,
        branch: item.branch_name,
        commitId: item.commit_short_id,
        status: item.status,
        duration: item.duration || '-',
        startTime: item.created_at
      }))
    }
  })
})
</script>

<style scoped>
.dashboard {
  max-width: 1400px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 32px;
}

.page-title {
  font-size: 28px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 8px;
}

.page-description {
  font-size: 14px;
  color: #6b7280;
}

/* 统计卡片网格 */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
  gap: 24px;
  margin-bottom: 32px;
}

.stat-card {
  background-color: #ffffff;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  transition: all 0.15s ease;
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  transform: translateY(-2px);
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.stat-icon .t-icon {
  font-size: 28px;
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: #1f2937;
  line-height: 1.2;
}

.stat-label {
  font-size: 14px;
  color: #6b7280;
  margin-top: 4px;
}

.stat-trend {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  margin-top: 8px;
}

.stat-trend.up {
  color: #10b981;
}

.stat-trend.down {
  color: #ef4444;
}

/* 图表网格 */
.charts-grid {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 24px;
  margin-bottom: 32px;
}

.chart-card {
  background-color: #ffffff;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.card-header h3 {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.chart-container {
  height: 280px;
  position: relative;
}

.chart-placeholder {
  width: 100%;
  height: 100%;
  background: #f9fafb;
  border-radius: 8px;
  position: relative;
  overflow: hidden;
}

.chart-line {
  position: absolute;
  top: 50%;
  left: 10%;
  right: 10%;
  height: 2px;
  background: #f59e0b;
  transform: translateY(-50%);
}

.chart-area {
  position: absolute;
  bottom: 0;
  left: 10%;
  right: 10%;
  height: 60%;
  background: linear-gradient(to top, #fef3c7 0%, transparent 100%);
}

/* 成功率样式 */
.success-rate {
  display: flex;
  align-items: center;
  gap: 32px;
}

.rate-circle {
  width: 120px;
  height: 120px;
  border-radius: 50%;
  background: conic-gradient(#10b981 0% 98.5%, #ef4444 98.5% 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
}

.rate-circle::before {
  content: '';
  position: absolute;
  width: 100px;
  height: 100px;
  background-color: #ffffff;
  border-radius: 50%;
}

.rate-value {
  font-size: 24px;
  font-weight: 700;
  color: #1f2937;
  position: relative;
  z-index: 1;
}

.rate-details {
  flex: 1;
}

.rate-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #6b7280;
  margin-bottom: 12px;
}

.dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.dot.success {
  background-color: #10b981;
}

.dot.failed {
  background-color: #ef4444;
}

.dot.skipped {
  background-color: #9ca3af;
}

/* 最近测试卡片 */
.recent-tests-card {
  background-color: #ffffff;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  margin-bottom: 32px;
}

.table-container {
  margin-top: 24px;
}

/* 快捷操作 */
.quick-actions {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
  gap: 24px;
}

.action-card {
  background-color: #ffffff;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  display: flex;
  gap: 20px;
  transition: all 0.15s ease;
}

.action-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  transform: translateY(-2px);
}

.action-icon {
  width: 48px;
  height: 48px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.action-icon .t-icon {
  font-size: 24px;
}

.action-content h4 {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 8px;
}

.action-content p {
  font-size: 14px;
  color: #6b7280;
  margin: 0 0 16px;
  line-height: 1.5;
}

/* 响应式 */
@media (max-width: 1024px) {
  .charts-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }

  .quick-actions {
    grid-template-columns: 1fr;
  }

  .success-rate {
    flex-direction: column;
    text-align: center;
  }
}
</style>