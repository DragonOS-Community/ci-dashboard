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
          <div ref="trendChartRef" class="trend-chart"></div>
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
            <div class="rate-circle" :style="getRateCircleStyle()">
              <div class="rate-value">{{ successRate.toFixed(1) }}%</div>
            </div>
            <div class="rate-details">
              <div class="rate-item">
                <span class="dot success"></span>
                <span>成功: {{ formatNumber(successCount) }}</span>
              </div>
              <div class="rate-item">
                <span class="dot failed"></span>
                <span>失败: {{ formatNumber(failedCount) }}</span>
              </div>
              <div class="rate-item">
                <span class="dot skipped"></span>
                <span>跳过: {{ formatNumber(skippedCount) }}</span>
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
          row-key="id"
          size="medium"
        >
          <template #status="{ row }">
            <t-tag :theme="getStatusTheme(row.status)" variant="light">
              {{ formatStatus(row.status) }}
            </t-tag>
          </template>
          <template #operation="{ row }">
            <t-button
              variant="text"
              theme="primary"
              @click="viewDetail(row.id)"
            >
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
import {
  ref,
  onMounted,
  computed,
  watch,
  nextTick,
  onBeforeUnmount,
} from "vue";
import { useRouter } from "vue-router";
import { useTestRunStore } from "@/stores/testRun";
import { getDashboardStats, getDashboardTrend } from "@/api/admin";
import { MessagePlugin } from "tdesign-vue-next";
import * as echarts from "echarts";

const router = useRouter();
const testRunStore = useTestRunStore();

const trendPeriod = ref("7d");
const loading = ref(false);
const trendData = ref([]);
const trendChartRef = ref(null);
let trendChart = null;

// 保存resize事件处理函数，以便后续移除
const handleResize = () => {
  if (trendChart) {
    trendChart.resize();
  }
};

// 统计数据
const stats = ref([
  {
    key: "total",
    label: "总测试数",
    value: "0",
    icon: "chart",
    color: "#f59e0b",
    bgColor: "#fef3c7",
    trend: {
      type: "up",
      icon: "arrow-up",
      text: "0%",
    },
  },
  {
    key: "today",
    label: "今日运行",
    value: "0",
    icon: "play-circle",
    color: "#10b981",
    bgColor: "#d1fae5",
    trend: {
      type: "up",
      icon: "arrow-up",
      text: "0%",
    },
  },
  {
    key: "success",
    label: "成功率",
    value: "0%",
    icon: "check-circle",
    color: "#10b981",
    bgColor: "#d1fae5",
    trend: {
      type: "up",
      icon: "arrow-up",
      text: "0%",
    },
  },
  {
    key: "avgTime",
    label: "平均耗时",
    value: "0s",
    icon: "time",
    color: "#3b82f6",
    bgColor: "#dbeafe",
    trend: {
      type: "down",
      icon: "arrow-down",
      text: "0%",
    },
  },
]);

// 成功率相关
const successRate = ref(0);
const successCount = ref(0);
const failedCount = ref(0);
const skippedCount = ref(0);

// 格式化数字
const formatNumber = (num) => {
  return num.toLocaleString("zh-CN");
};

// 格式化时间（秒转分钟）
const formatDuration = (seconds) => {
  if (seconds < 60) {
    return `${Math.round(seconds)}s`;
  } else if (seconds < 3600) {
    const minutes = Math.floor(seconds / 60);
    const secs = Math.round(seconds % 60);
    return secs > 0 ? `${minutes}m ${secs}s` : `${minutes}m`;
  } else {
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    return minutes > 0 ? `${hours}h ${minutes}m` : `${hours}h`;
  }
};

// 计算趋势百分比
const calculateTrend = (current, previous) => {
  if (previous === 0) {
    return { type: "up", text: current > 0 ? "100%" : "0%" };
  }
  const change = ((current - previous) / previous) * 100;
  if (Math.abs(change) < 0.01) {
    return { type: "up", text: "0%" };
  }
  const type = change > 0 ? "up" : "down";
  const text = `${change > 0 ? "+" : ""}${change.toFixed(1)}%`;
  return { type, text };
};

// 加载仪表板数据
const loadDashboardData = async () => {
  loading.value = true;
  try {
    // 加载统计数据
    const statsRes = await getDashboardStats();
    const dashboardStats = statsRes.data;

    // 更新统计卡片
    const totalTrend = calculateTrend(
      dashboardStats.total_tests,
      dashboardStats.total_tests_prev,
    );
    const todayTrend = calculateTrend(
      dashboardStats.today_runs,
      dashboardStats.today_runs_prev,
    );
    const successTrend = calculateTrend(
      dashboardStats.success_rate,
      dashboardStats.success_rate_prev,
    );
    const durationTrend = calculateTrend(
      dashboardStats.avg_duration,
      dashboardStats.avg_duration_prev,
    );

    stats.value = [
      {
        key: "total",
        label: "总测试数",
        value: formatNumber(dashboardStats.total_tests),
        icon: "chart",
        color: "#f59e0b",
        bgColor: "#fef3c7",
        trend: {
          type: totalTrend.type,
          icon: totalTrend.type === "up" ? "arrow-up" : "arrow-down",
          text: totalTrend.text,
        },
      },
      {
        key: "today",
        label: "今日运行",
        value: formatNumber(dashboardStats.today_runs),
        icon: "play-circle",
        color: "#10b981",
        bgColor: "#d1fae5",
        trend: {
          type: todayTrend.type,
          icon: todayTrend.type === "up" ? "arrow-up" : "arrow-down",
          text: todayTrend.text,
        },
      },
      {
        key: "success",
        label: "成功率",
        value: `${dashboardStats.success_rate.toFixed(1)}%`,
        icon: "check-circle",
        color: "#10b981",
        bgColor: "#d1fae5",
        trend: {
          type: successTrend.type,
          icon: successTrend.type === "up" ? "arrow-up" : "arrow-down",
          text: successTrend.text,
        },
      },
      {
        key: "avgTime",
        label: "平均耗时",
        value: formatDuration(dashboardStats.avg_duration),
        icon: "time",
        color: "#3b82f6",
        bgColor: "#dbeafe",
        trend: {
          type: durationTrend.type,
          icon: durationTrend.type === "up" ? "arrow-up" : "arrow-down",
          text: durationTrend.text,
        },
      },
    ];

    // 更新成功率相关数据
    successRate.value = dashboardStats.success_rate;
    successCount.value = dashboardStats.success_count;
    failedCount.value = dashboardStats.failed_count;
    skippedCount.value = dashboardStats.skipped_count;
  } catch (error) {
    console.error("Failed to load dashboard stats:", error);
    MessagePlugin.error("加载仪表板数据失败");
  } finally {
    loading.value = false;
  }
};

// 初始化趋势图表
const initTrendChart = () => {
  if (!trendChartRef.value) return;

  if (trendChart) {
    trendChart.dispose();
  }

  trendChart = echarts.init(trendChartRef.value);

  // 设置默认配置
  const option = {
    grid: {
      left: "3%",
      right: "4%",
      bottom: "3%",
      top: "10%",
      containLabel: true,
    },
    xAxis: {
      type: "category",
      boundaryGap: false,
      data: [],
      axisLine: {
        lineStyle: {
          color: "#E5E7EB",
        },
      },
      axisLabel: {
        color: "#6B7280",
        fontSize: 12,
      },
    },
    yAxis: {
      type: "value",
      axisLine: {
        lineStyle: {
          color: "#E5E7EB",
        },
      },
      axisLabel: {
        color: "#6B7280",
        fontSize: 12,
      },
      splitLine: {
        lineStyle: {
          color: "#F3F4F6",
        },
      },
    },
    series: [
      {
        name: "测试运行数",
        type: "line",
        smooth: true,
        data: [],
        itemStyle: {
          color: "#F59E0B",
        },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            {
              offset: 0,
              color: "rgba(245, 158, 11, 0.3)",
            },
            {
              offset: 1,
              color: "rgba(245, 158, 11, 0.05)",
            },
          ]),
        },
        lineStyle: {
          width: 2,
          color: "#F59E0B",
        },
      },
    ],
    tooltip: {
      trigger: "axis",
      backgroundColor: "rgba(255, 255, 255, 0.95)",
      borderColor: "#E5E7EB",
      borderWidth: 1,
      textStyle: {
        color: "#1F2937",
      },
      axisPointer: {
        type: "line",
        lineStyle: {
          color: "#F59E0B",
        },
      },
    },
  };

  trendChart.setOption(option);

  // 响应式调整
  window.addEventListener("resize", handleResize);
};

// 更新趋势图表数据
const updateTrendChart = () => {
  if (!trendChart || !trendData.value || trendData.value.length === 0) {
    return;
  }

  // 格式化日期
  const dates = trendData.value.map((item) => {
    const date = new Date(item.date);
    const month = date.getMonth() + 1;
    const day = date.getDate();
    return `${month}/${day}`;
  });

  const counts = trendData.value.map((item) => item.count);

  trendChart.setOption({
    xAxis: {
      data: dates,
    },
    series: [
      {
        data: counts,
      },
    ],
  });
};

// 加载趋势数据
const loadTrendData = async () => {
  const days =
    trendPeriod.value === "7d" ? 7 : trendPeriod.value === "30d" ? 30 : 90;
  try {
    const trendRes = await getDashboardTrend(days);
    trendData.value = trendRes.data;
    await nextTick();
    updateTrendChart();
  } catch (error) {
    console.error("Failed to load trend data:", error);
  }
};

// 监听趋势周期变化
watch(
  trendPeriod,
  () => {
    loadTrendData();
  },
  { immediate: false },
);

// 组件卸载时销毁图表
onBeforeUnmount(() => {
  if (trendChart) {
    trendChart.dispose();
    trendChart = null;
  }
  window.removeEventListener("resize", handleResize);
});

// 最近测试数据
const recentTests = ref([
  {
    id: 1,
    branch: "main",
    commitId: "abc123",
    status: "success",
    duration: "12m 30s",
    startTime: "2024-01-01 10:00:00",
  },
  {
    id: 2,
    branch: "dev",
    commitId: "def456",
    status: "failed",
    duration: "8m 15s",
    startTime: "2024-01-01 09:30:00",
  },
  {
    id: 3,
    branch: "feature/test",
    commitId: "ghi789",
    status: "running",
    duration: "-",
    startTime: "2024-01-01 11:00:00",
  },
]);

const tableColumns = [
  { colKey: "id", title: "ID", width: 80 },
  { colKey: "branch", title: "分支", width: 120 },
  { colKey: "commitId", title: "Commit ID", width: 100 },
  { colKey: "status", title: "状态", width: 100 },
  { colKey: "duration", title: "执行时间", width: 120 },
  { colKey: "startTime", title: "开始时间", width: 180 },
  { colKey: "operation", title: "操作", width: 100, fixed: "right" },
];

// 快捷操作
const quickActions = ref([
  {
    key: "newTest",
    title: "创建测试任务",
    description: "快速创建新的测试运行任务",
    icon: "add",
    color: "#ffffff",
    bgColor: "#f59e0b",
    theme: "warning",
    buttonText: "立即创建",
    handler: () => {
      router.push("/admin/test/runs?action=new");
    },
  },
  {
    key: "apiKey",
    title: "管理 API 密钥",
    description: "创建或管理 API 访问密钥",
    icon: "key",
    color: "#ffffff",
    bgColor: "#3b82f6",
    theme: "primary",
    buttonText: "管理密钥",
    handler: () => {
      router.push("/admin/system/api-keys");
    },
  },
  {
    key: "report",
    title: "生成测试报告",
    description: "导出最近一段时间的测试报告",
    icon: "file-excel",
    color: "#ffffff",
    bgColor: "#10b981",
    theme: "success",
    buttonText: "生成报告",
    handler: () => {
      router.push("/admin/test/reports");
    },
  },
]);

const getStatusTheme = (status) => {
  const themes = {
    passed: "success",
    failed: "danger",
    running: "warning",
    cancelled: "default",
    success: "success", // 兼容旧数据
  };
  return themes[status] || "default";
};

// 格式化状态显示文本
const formatStatus = (status) => {
  const statusMap = {
    passed: "成功",
    failed: "失败",
    running: "运行中",
    cancelled: "已取消",
  };
  return statusMap[status] || status;
};

const viewDetail = (id) => {
  router.push(`/test-runs/${id}`);
};

// 计算成功率圆环样式
const getRateCircleStyle = () => {
  const total = successCount.value + failedCount.value + skippedCount.value;
  if (total === 0) {
    return {
      background: "conic-gradient(#9ca3af 0% 100%)",
    };
  }
  const successPercent = (successCount.value / total) * 100;
  const failedPercent = (failedCount.value / total) * 100;
  const skippedPercent = (skippedCount.value / total) * 100;

  const successEnd = successPercent;
  const failedEnd = successPercent + failedPercent;

  return {
    background: `conic-gradient(#10b981 0% ${successEnd}%, #ef4444 ${successEnd}% ${failedEnd}%, #9ca3af ${failedEnd}% 100%)`,
  };
};

onMounted(async () => {
  // 初始化图表
  await nextTick();
  initTrendChart();

  // 加载仪表板数据
  await loadDashboardData();
  await loadTrendData();

  // 加载最近测试数据
  try {
    await testRunStore.fetchTestRuns();
    if (testRunStore.testRuns && testRunStore.testRuns.length > 0) {
      recentTests.value = testRunStore.testRuns.slice(0, 10).map((item) => {
        // 格式化持续时间
        let duration = "-";
        if (item.completed_at && item.started_at) {
          const start = new Date(item.started_at).getTime();
          const end = new Date(item.completed_at).getTime();
          const diff = Math.floor((end - start) / 1000);
          duration = formatDuration(diff);
        }

        // 格式化状态（后端返回的是passed/failed/running/cancelled，前端显示需要映射）
        let status = item.status;
        // 保持原状态，getStatusTheme会处理映射

        return {
          id: item.id,
          branch: item.branch_name || "-",
          commitId:
            item.commit_short_id || item.commit_id?.substring(0, 7) || "-",
          status: status,
          duration: duration,
          startTime: item.created_at || "-",
        };
      });
    }
  } catch (error) {
    console.error("Failed to load recent tests:", error);
  }
});
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

.trend-chart {
  width: 100%;
  height: 100%;
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
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
}

.rate-circle::before {
  content: "";
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
