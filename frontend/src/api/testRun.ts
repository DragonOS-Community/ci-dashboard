import request from '@/utils/request'
import type { AxiosPromise } from 'axios'

// 定义接口类型
export interface TestRunParams {
  page?: number
  page_size?: number
  project_name?: string
  status?: string
}

export interface TestRunData {
  project_name: string
  branch?: string
  commit_hash?: string
  test_type?: string
}

export interface TestCase {
  id: string
  name: string
  status: string
  duration?: number
  error_msg?: string
  created_at: string
  updated_at: string
}

export interface TestFile {
  id: string
  filename: string
  filepath: string
  size: number
  file_type: string
  created_at: string
}

// 获取测试运行列表
export function getTestRuns(params: TestRunParams): AxiosPromise {
  return request({
    url: '/test-runs',
    method: 'get',
    params,
  })
}

// 获取测试运行详情
export function getTestRunById(id: string): AxiosPromise {
  return request({
    url: `/test-runs/${id}`,
    method: 'get',
  })
}

// 创建测试运行（需要API Key）
export function createTestRun(data: TestRunData): AxiosPromise {
  return request({
    url: '/test-runs',
    method: 'post',
    data,
  })
}

// 获取测例列表
export function getTestCasesByTestRunId(id: string): AxiosPromise<TestCase[]> {
  return request({
    url: `/test-runs/${id}/test-cases`,
    method: 'get',
  })
}

// 获取文件列表
export function getFilesByTestRunId(id: string): AxiosPromise<TestFile[]> {
  return request({
    url: `/test-runs/${id}/files`,
    method: 'get',
  })
}

// 下载文件
export function downloadFile(testRunId: string, fileId: string): AxiosPromise<Blob> {
  return request({
    url: `/test-runs/${testRunId}/output-files/${fileId}`,
    method: 'get',
    responseType: 'blob',
  })
}

// 上传文件（需要API Key）
export function uploadFile(testRunId: string, file: File): AxiosPromise {
  const formData = new FormData()
  formData.append('file', file)
  return request({
    url: `/test-runs/${testRunId}/output-files`,
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  })
}