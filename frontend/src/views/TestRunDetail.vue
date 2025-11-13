<template>
  <div class="detail-container">
    <t-layout>
      <t-header>
        <div class="header-content">
          <t-link theme="primary" @click="goBack">
            <t-icon name="chevron-left" />
            返回列表
          </t-link>
          <h1>测试运行详情</h1>
        </div>
      </t-header>
      <t-content>
        <div class="content-wrapper">
          <t-loading :loading="testRunStore.loading">
            <t-card v-if="testRunStore.currentTestRun" class="info-card">
              <t-descriptions :data="infoData" :column="2" />
            </t-card>

            <t-card class="test-cases-card" title="测例列表">
              <t-tabs v-model="activeTab">
                <t-tab-panel value="all" label="全部">
                  <TestCaseList :test-cases="allTestCases" />
                </t-tab-panel>
                <t-tab-panel value="passed" label="通过">
                  <TestCaseList :test-cases="passedTestCases" />
                </t-tab-panel>
                <t-tab-panel value="failed" label="失败">
                  <TestCaseList :test-cases="failedTestCases" />
                </t-tab-panel>
              </t-tabs>
            </t-card>

            <t-card class="files-card" title="输出文件">
              <t-list v-if="files.length > 0">
                <t-list-item v-for="file in files" :key="file.id">
                  <div class="file-item">
                    <span>{{ file.filename }}</span>
                    <t-space>
                      <t-tag>{{ formatFileSize(file.file_size) }}</t-tag>
                      <t-button size="small" @click="downloadFile(file)">下载</t-button>
                    </t-space>
                  </div>
                </t-list-item>
              </t-list>
              <t-empty v-else description="暂无文件" />
            </t-card>
          </t-loading>
        </div>
      </t-content>
    </t-layout>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
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

const allTestCases = computed(() => testCases.value)
const passedTestCases = computed(() => testCases.value.filter(tc => tc.status === 'passed'))
const failedTestCases = computed(() => testCases.value.filter(tc => tc.status === 'failed'))

const infoData = computed(() => {
  const run = testRunStore.currentTestRun
  if (!run) return []
  return [
    { label: 'ID', content: run.id },
    { label: '分支名', content: run.branch_name },
    { label: 'Commit ID', content: run.commit_id },
    { label: '测试类型', content: run.test_type },
    { label: '状态', content: getStatusText(run.status) },
    { label: '创建时间', content: new Date(run.created_at).toLocaleString() },
    { label: '开始时间', content: run.started_at ? new Date(run.started_at).toLocaleString() : '-' },
    { label: '完成时间', content: run.completed_at ? new Date(run.completed_at).toLocaleString() : '-' },
  ]
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

const formatFileSize = (bytes) => {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(2) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(2) + ' MB'
}

const downloadFile = async (file) => {
  try {
    const response = await downloadFileAPI(testRunStore.currentTestRun.id, file.id)
    const url = window.URL.createObjectURL(new Blob([response.data]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', file.filename)
    document.body.appendChild(link)
    link.click()
    link.remove()
    MessagePlugin.success('下载成功')
  } catch (error) {
    MessagePlugin.error('下载失败')
  }
}

const goBack = () => {
  router.push('/')
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
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.detail-container {
  min-height: 100vh;
  background: #f5f5f5;
}

.header-content {
  padding: 0 24px;
  height: 100%;
  display: flex;
  align-items: center;
  gap: 16px;
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

.info-card,
.test-cases-card,
.files-card {
  margin-bottom: 16px;
}

.file-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}
</style>

