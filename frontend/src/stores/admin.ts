import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { adminLogin, getAPIKeys, createAPIKey, deleteAPIKey, getProfile, updatePassword } from '@/api/admin'
import { MessagePlugin } from 'tdesign-vue-next'
import type { APIKey, Profile, APIKeyData } from '@/api/admin'

export interface User {
  id: string
  username: string
  role?: string
}

export const useAdminStore = defineStore('admin', () => {
  const user = ref<User | null>(null)
  const token = ref<string>(localStorage.getItem('admin_token') || '')
  const apiKeys = ref<APIKey[]>([])
  const loading = ref<boolean>(false)

  // 计算属性
  const isAuthenticated = computed(() => !!token.value)

  // 登录
  async function login(username: string, password: string): Promise<boolean> {
    loading.value = true
    try {
      const res = await adminLogin({ username, password })
      token.value = res.data.token
      user.value = res.data.user
      localStorage.setItem('admin_token', token.value)
      // 登录成功后获取完整的用户信息
      await fetchProfile()
      MessagePlugin.success('登录成功')
      return true
    } catch (error) {
      MessagePlugin.error('登录失败')
      return false
    } finally {
      loading.value = false
    }
  }

  // 登出
  function logout(): void {
    token.value = ''
    user.value = null
    localStorage.removeItem('admin_token')
  }

  // 获取API密钥列表
  async function fetchAPIKeys(): Promise<void> {
    loading.value = true
    try {
      const res = await getAPIKeys()
      apiKeys.value = res.data || []
    } catch (error) {
      console.error('Failed to fetch API keys:', error)
    } finally {
      loading.value = false
    }
  }

  // 创建API密钥
  async function createKey(name: string, description?: string, expiresAt?: string): Promise<APIKey | null> {
    loading.value = true
    try {
      const data: APIKeyData = { name }
      if (description) data.description = description
      if (expiresAt) data.expires_at = expiresAt

      const res = await createAPIKey(data)
      MessagePlugin.success('API密钥创建成功')
      await fetchAPIKeys()
      return res.data
    } catch (error) {
      MessagePlugin.error('创建API密钥失败')
      return null
    } finally {
      loading.value = false
    }
  }

  // 删除API密钥
  async function removeKey(id: string): Promise<void> {
    loading.value = true
    try {
      await deleteAPIKey(id)
      MessagePlugin.success('API密钥删除成功')
      await fetchAPIKeys()
    } catch (error) {
      MessagePlugin.error('删除API密钥失败')
    } finally {
      loading.value = false
    }
  }

  // 获取用户信息
  async function fetchProfile(): Promise<Profile | null> {
    loading.value = true
    try {
      const res = await getProfile()
      user.value = res.data
      return res.data
    } catch (error) {
      console.error('Failed to fetch profile:', error)
      return null
    } finally {
      loading.value = false
    }
  }

  // 更新密码
  async function changePassword(oldPassword: string, newPassword: string, confirmPassword: string): Promise<boolean> {
    loading.value = true
    try {
      await updatePassword({
        old_password: oldPassword,
        new_password: newPassword,
        confirm_password: confirmPassword,
      })
      MessagePlugin.success('密码更新成功')
      return true
    } catch (error) {
      MessagePlugin.error('密码更新失败')
      return false
    } finally {
      loading.value = false
    }
  }

  return {
    user,
    token,
    apiKeys,
    loading,
    isAuthenticated,
    login,
    logout,
    fetchAPIKeys,
    createKey,
    removeKey,
    fetchProfile,
    changePassword,
  }
})