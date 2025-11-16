import request from "@/utils/request";
import type { AxiosPromise } from "axios";

// 定义接口类型
export interface LoginData {
  username: string;
  password: string;
}

export interface APIKeyData {
  name: string;
  description?: string;
  expires_at?: string;
}

export interface APIKey {
  id: string;
  name: string;
  description?: string;
  key: string;
  expires_at?: string;
  created_at: string;
  updated_at: string;
}

export interface Profile {
  id: string;
  username: string;
  created_at: string;
  updated_at: string;
}

export interface UpdatePasswordData {
  old_password: string;
  new_password: string;
  confirm_password: string;
}

// 管理员登录
export function adminLogin(data: LoginData): AxiosPromise<{
  token: string;
  user: Profile;
}> {
  return request({
    url: "/admin/login",
    method: "post",
    data,
  });
}

// 获取API密钥列表
export function getAPIKeys(): AxiosPromise<APIKey[]> {
  return request({
    url: "/admin/api-keys",
    method: "get",
  });
}

// 创建API密钥
export function createAPIKey(data: APIKeyData): AxiosPromise<APIKey> {
  return request({
    url: "/admin/api-keys",
    method: "post",
    data,
  });
}

// 删除API密钥
export function deleteAPIKey(id: string): AxiosPromise {
  return request({
    url: `/admin/api-keys/${id}`,
    method: "delete",
  });
}

// 获取当前用户信息
export function getProfile(): AxiosPromise<Profile> {
  return request({
    url: "/admin/profile",
    method: "get",
  });
}

// 更新密码
export function updatePassword(data: UpdatePasswordData): AxiosPromise {
  return request({
    url: "/admin/profile/password",
    method: "put",
    data,
  });
}

// 仪表板统计数据接口
export interface DashboardStats {
  total_tests: number;
  today_runs: number;
  success_rate: number;
  avg_duration: number;
  total_tests_prev: number;
  today_runs_prev: number;
  success_rate_prev: number;
  avg_duration_prev: number;
  success_count: number;
  failed_count: number;
  skipped_count: number;
}

// 趋势数据接口
export interface TrendData {
  date: string;
  count: number;
}

// 获取仪表板统计数据
export function getDashboardStats(): AxiosPromise<DashboardStats> {
  return request({
    url: "/admin/dashboard/stats",
    method: "get",
  });
}

// 获取仪表板趋势数据
export function getDashboardTrend(days: number = 7): AxiosPromise<TrendData[]> {
  return request({
    url: "/admin/dashboard/trend",
    method: "get",
    params: { days },
  });
}

// 项目相关接口
export interface Project {
  id: number;
  name: string;
  description?: string;
  created_at: string;
  updated_at: string;
}

export interface ProjectData {
  name: string;
  description?: string;
}

// 获取项目列表
export function getProjects(): AxiosPromise<Project[]> {
  return request({
    url: "/admin/projects",
    method: "get",
  });
}

// 根据ID获取项目
export function getProjectById(id: number): AxiosPromise<Project> {
  return request({
    url: `/admin/projects/${id}`,
    method: "get",
  });
}

// 创建项目
export function createProject(data: ProjectData): AxiosPromise<Project> {
  return request({
    url: "/admin/projects",
    method: "post",
    data,
  });
}

// 更新项目
export function updateProject(
  id: number,
  data: ProjectData,
): AxiosPromise<Project> {
  return request({
    url: `/admin/projects/${id}`,
    method: "put",
    data,
  });
}

// 删除项目
export function deleteProject(id: number): AxiosPromise {
  return request({
    url: `/admin/projects/${id}`,
    method: "delete",
  });
}

// 测试运行管理接口
export interface TestRun {
  id: number;
  project_id: number;
  branch_name: string;
  commit_id: string;
  commit_short_id: string;
  test_type: string;
  status: string;
  is_public: boolean;
  started_at?: string;
  completed_at?: string;
  created_at: string;
  project?: {
    id: number;
    name: string;
  };
}

export interface TestRunListResponse {
  test_runs: TestRun[];
  total: number;
  page: number;
  page_size: number;
}

// 获取测试运行列表（管理员接口，包含私有记录）
export function getTestRunsAdmin(params?: {
  page?: number;
  page_size?: number;
  branch?: string;
  commit_id?: string;
  test_type?: string;
  status?: string;
}): AxiosPromise<TestRunListResponse> {
  return request({
    url: "/admin/test-runs",
    method: "get",
    params,
  });
}

// 删除测试运行
export function deleteTestRun(id: number): AxiosPromise {
  return request({
    url: `/admin/test-runs/${id}`,
    method: "delete",
  });
}

// 更新测试运行可见性
export function updateTestRunVisibility(
  id: number,
  isPublic: boolean,
): AxiosPromise<TestRun> {
  return request({
    url: `/admin/test-runs/${id}/visibility`,
    method: "put",
    data: { is_public: isPublic },
  });
}
