<template>
  <div class="home-container">
    <t-layout>
      <t-header>
        <div class="header-content">
          <h1>DragonOS CI Dashboard</h1>
        </div>
      </t-header>
      <t-content>
        <div class="content-wrapper">
          <!-- 筛选栏 -->
          <t-card class="filter-card">
            <t-form :data="testRunStore.filters" @submit="handleSearch">
              <t-row :gutter="16">
                <t-col :span="6">
                  <t-form-item label="分支名">
                    <t-input v-model="testRunStore.filters.branch" placeholder="输入分支名" clearable />
                  </t-form-item>
                </t-col>
                <t-col :span="6">
                  <t-form-item label="Commit ID">
                    <t-input v-model="testRunStore.filters.commitId" placeholder="输入Commit ID" clearable />
                  </t-form-item>
                </t-col>
                <t-col :span="6">
                  <t-form-item label="状态">
                    <t-select v-model="testRunStore.filters.status" clearable>
                      <t-option value="all" label="全部" />
                      <t-option value="passed" label="通过" />
                      <t-option value="failed" label="失败" />
                      <t-option value="running" label="运行中" />
                    </t-select>
                  </t-form-item>
                </t-col>
                <t-col :span="6">
                  <t-form-item label="操作">
                    <t-space>
                      <t-button theme="primary" type="submit">搜索</t-button>
                      <t-button theme="default" @click="handleReset">重置</t-button>
                    </t-space>
                  </t-form-item>
                </t-col>
              </t-row>
            </t-form>
          </t-card>

          <!-- 测试运行列表 -->
          <t-card class="list-card">
            <t-loading :loading="testRunStore.loading">
              <t-table
                :data="testRunStore.testRuns"
                :columns="columns"
                :pagination="paginationConfig"
                @page-change="handlePageChange"
                @page-size-change="handlePageSizeChange"
                hover
              >
                <template #status="{ row }">
                  <t-tag :theme="getStatusTheme(row.status)">
                    {{ getStatusText(row.status) }}
                  </t-tag>
                </template>
                <template #operation="{ row }">
                  <t-link theme="primary" @click="viewDetail(row.id)">查看详情</t-link>
                </template>
              </t-table>
            </t-loading>
          </t-card>
        </div>
      </t-content>
    </t-layout>
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

onMounted(() => {
  testRunStore.fetchTestRuns()
})
</script>

<style scoped>
.home-container {
  min-height: 100vh;
  background: #f5f5f5;
}

.header-content {
  padding: 0 24px;
  height: 100%;
  display: flex;
  align-items: center;
}

.header-content h1 {
  margin: 0;
  color: #fff;
  font-size: 20px;
}

.content-wrapper {
  padding: 24px;
  max-width: 1400px;
  margin: 0 auto;
}

.filter-card {
  margin-bottom: 16px;
}

.list-card {
  margin-top: 16px;
}
</style>

