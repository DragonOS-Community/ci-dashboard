<template>
  <div class="login-container">
    <t-card class="login-card">
      <h2>管理员登录</h2>
      <t-form :data="formData" @submit="handleLogin">
        <t-form-item label="用户名" name="username">
          <t-input v-model="formData.username" placeholder="请输入用户名" />
        </t-form-item>
        <t-form-item label="密码" name="password">
          <t-input v-model="formData.password" type="password" placeholder="请输入密码" />
        </t-form-item>
        <t-form-item>
          <t-button theme="primary" type="submit" block :loading="adminStore.loading">
            登录
          </t-button>
        </t-form-item>
      </t-form>
    </t-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAdminStore } from '@/stores/admin'

const router = useRouter()
const route = useRoute()
const adminStore = useAdminStore()

const formData = ref({
  username: '',
  password: '',
})

const handleLogin = async () => {
  const success = await adminStore.login(formData.value.username, formData.value.password)
  if (success) {
    const redirect = route.query.redirect || '/admin/api-keys'
    router.push(redirect)
  }
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f5f5f5;
}

.login-card {
  width: 400px;
}

.login-card h2 {
  margin-bottom: 24px;
  text-align: center;
}
</style>

