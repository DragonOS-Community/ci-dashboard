import request from '@/utils/request'
import type { AxiosPromise } from 'axios'

// 定义接口类型
export interface LoginData {
  username: string
  password: string
}

export interface APIKeyData {
  name: string
  description?: string
  expires_at?: string
}

export interface APIKey {
  id: string
  name: string
  description?: string
  key: string
  expires_at?: string
  created_at: string
  updated_at: string
}

export interface Profile {
  id: string
  username: string
  created_at: string
  updated_at: string
}

export interface UpdatePasswordData {
  old_password: string
  new_password: string
  confirm_password: string
}

// 管理员登录
export function adminLogin(data: LoginData): AxiosPromise<{
  token: string
  user: Profile
}> {
  return request({
    url: '/admin/login',
    method: 'post',
    data,
  })
}

// 获取API密钥列表
export function getAPIKeys(): AxiosPromise<APIKey[]> {
  return request({
    url: '/admin/api-keys',
    method: 'get',
  })
}

// 创建API密钥
export function createAPIKey(data: APIKeyData): AxiosPromise<APIKey> {
  return request({
    url: '/admin/api-keys',
    method: 'post',
    data,
  })
}

// 删除API密钥
export function deleteAPIKey(id: string): AxiosPromise {
  return request({
    url: `/admin/api-keys/${id}`,
    method: 'delete',
  })
}

// 获取当前用户信息
export function getProfile(): AxiosPromise<Profile> {
  return request({
    url: '/admin/profile',
    method: 'get',
  })
}

// 更新密码
export function updatePassword(data: UpdatePasswordData): AxiosPromise {
  return request({
    url: '/admin/profile/password',
    method: 'put',
    data,
  })
}