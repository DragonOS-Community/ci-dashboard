import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getTestRuns, getTestRunById } from '@/api/testRun'
import type { TestCase, TestFile } from '@/api/testRun'

// 定义测试运行接口
export interface TestRun {
  id: string
  project_name: string
  branch?: string
  commit_hash?: string
  status: string
  start_time?: string
  end_time?: string
  duration?: number
  created_at: string
  updated_at: string
  test_cases?: TestCase[]
  files?: TestFile[]
}

export interface TestRunFilters {
  branch: string
  commitId: string
  startTime: string
  endTime: string
  status: string
  testCaseName: string
}

export interface Pagination {
  page: number
  pageSize: number
}

export const useTestRunStore = defineStore('testRun', () => {
  const testRuns = ref<TestRun[]>([])
  const currentTestRun = ref<TestRun | null>(null)
  const loading = ref<boolean>(false)
  const total = ref<number>(0)
  const pagination = ref<Pagination>({
    page: 1,
    pageSize: 20,
  })
  const filters = ref<TestRunFilters>({
    branch: '',
    commitId: '',
    startTime: '',
    endTime: '',
    status: 'all',
    testCaseName: '',
  })

  // 获取测试运行列表
  async function fetchTestRuns(): Promise<void> {
    loading.value = true
    try {
      const params: any = {
        page: pagination.value.page,
        page_size: pagination.value.pageSize,
      }

      if (filters.value.branch) {
        params.branch = filters.value.branch
      }
      if (filters.value.commitId) {
        params.commit_id = filters.value.commitId
      }
      if (filters.value.startTime) {
        params.start_time = filters.value.startTime
      }
      if (filters.value.endTime) {
        params.end_time = filters.value.endTime
      }
      if (filters.value.status && filters.value.status !== 'all') {
        params.status = filters.value.status
      }
      if (filters.value.testCaseName) {
        params.test_case_name = filters.value.testCaseName
      }

      const res = await getTestRuns(params)
      testRuns.value = res.data.test_runs || []
      total.value = res.data.total || 0
    } catch (error) {
      console.error('Failed to fetch test runs:', error)
    } finally {
      loading.value = false
    }
  }

  // 获取测试运行详情
  async function fetchTestRunById(id: string): Promise<void> {
    loading.value = true
    try {
      const res = await getTestRunById(id)
      currentTestRun.value = res.data
    } catch (error) {
      console.error('Failed to fetch test run:', error)
    } finally {
      loading.value = false
    }
  }

  // 重置筛选条件
  function resetFilters(): void {
    filters.value = {
      branch: '',
      commitId: '',
      startTime: '',
      endTime: '',
      status: 'all',
      testCaseName: '',
    }
    pagination.value.page = 1
  }

  // 设置分页
  function setPagination(page: number, pageSize: number): void {
    pagination.value.page = page
    pagination.value.pageSize = pageSize
  }

  return {
    testRuns,
    currentTestRun,
    loading,
    total,
    pagination,
    filters,
    fetchTestRuns,
    fetchTestRunById,
    resetFilters,
    setPagination,
  }
})