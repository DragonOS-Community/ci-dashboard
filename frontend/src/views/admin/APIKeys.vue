<template>
  <div class="api-keys-container">
    <t-layout>
      <t-header>
        <div class="header-content">
          <h1>API密钥管理</h1>
          <t-space>
            <t-button theme="default" @click="router.push('/admin/profile')">个人面板</t-button>
            <span>欢迎，{{ adminStore.user?.username }}</span>
            <t-button theme="default" @click="handleLogout">退出</t-button>
          </t-space>
        </div>
      </t-header>
      <t-content>
        <div class="content-wrapper">
          <t-card>
            <div class="card-header">
              <h2>API密钥列表</h2>
              <t-button theme="primary" @click="showCreateDialog = true">创建密钥</t-button>
            </div>

            <t-loading :loading="adminStore.loading">
              <t-table :data="adminStore.apiKeys" :columns="columns" hover>
                <template #project_id="{ row }">
                  <span v-if="row.project_id">{{ row.project_id }}</span>
                  <span v-else class="text-muted">全部项目</span>
                </template>
                <template #expires_at="{ row }">
                  <span v-if="row.expires_at">{{ new Date(row.expires_at).toLocaleString() }}</span>
                  <span v-else class="text-muted">永不过期</span>
                </template>
                <template #operation="{ row }">
                  <t-button theme="danger" size="small" @click="handleDelete(row.id)">删除</t-button>
                </template>
              </t-table>
            </t-loading>
          </t-card>
        </div>
      </t-content>
    </t-layout>

    <!-- 创建密钥对话框 -->
    <t-dialog v-model:visible="showCreateDialog" title="创建API密钥" @confirm="handleCreate">
      <t-form :data="createForm" ref="createFormRef">
        <t-form-item label="名称" name="name">
          <t-input v-model="createForm.name" placeholder="请输入密钥名称" />
        </t-form-item>
        <t-form-item label="项目ID" name="project_id">
          <t-input v-model.number="createForm.project_id" placeholder="留空表示所有项目" type="number" />
        </t-form-item>
        <t-form-item label="过期时间" name="expires_at">
          <t-date-picker v-model="createForm.expires_at" placeholder="留空表示永不过期" />
        </t-form-item>
      </t-form>
    </t-dialog>

    <!-- 显示密钥对话框 -->
    <t-dialog v-model:visible="showKeyDialog" title="API密钥创建成功" :footer="false">
      <t-alert theme="warning" message="请妥善保管此密钥，它只会显示一次！" />
      <div class="key-display">
        <t-input v-model="newKey" readonly />
        <t-button @click="copyKey">复制</t-button>
      </div>
    </t-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAdminStore } from '@/stores/admin'
import { MessagePlugin } from 'tdesign-vue-next'

const router = useRouter()
const adminStore = useAdminStore()

const showCreateDialog = ref(false)
const showKeyDialog = ref(false)
const newKey = ref('')
const createFormRef = ref(null)

const createForm = ref({
  name: '',
  project_id: null,
  expires_at: null,
})

const columns = [
  { colKey: 'id', title: 'ID', width: 80 },
  { colKey: 'name', title: '名称', width: 200 },
  { colKey: 'project_id', title: '项目ID', width: 120 },
  { colKey: 'created_at', title: '创建时间', width: 180 },
  { colKey: 'last_used_at', title: '最后使用', width: 180 },
  { colKey: 'expires_at', title: '过期时间', width: 180 },
  { colKey: 'operation', title: '操作', width: 100, fixed: 'right' },
]

const handleCreate = async () => {
  if (!createForm.value.name) {
    MessagePlugin.warning('请输入密钥名称')
    return
  }

  const projectId = createForm.value.project_id || null
  const expiresAt = createForm.value.expires_at ? new Date(createForm.value.expires_at).toISOString() : null

  const result = await adminStore.createKey(createForm.value.name, projectId, expiresAt)
  if (result) {
    showCreateDialog.value = false
    newKey.value = result.api_key
    showKeyDialog.value = true
    createForm.value = {
      name: '',
      project_id: null,
      expires_at: null,
    }
  }
}

const handleDelete = async (id) => {
  const result = await MessagePlugin.confirm({
    header: '确认删除',
    body: '确定要删除此API密钥吗？',
  })
  if (result) {
    await adminStore.removeKey(id)
  }
}

const handleLogout = () => {
  adminStore.logout()
  router.push('/admin/login')
}

const copyKey = () => {
  navigator.clipboard.writeText(newKey.value)
  MessagePlugin.success('已复制到剪贴板')
}

onMounted(async () => {
  // 如果用户信息不存在，则获取
  if (!adminStore.user) {
    await adminStore.fetchProfile()
  }
  adminStore.fetchAPIKeys()
})
</script>

<style scoped>
.api-keys-container {
  min-height: 100vh;
  background: #f5f5f5;
}

.header-content {
  padding: 0 24px;
  height: 100%;
  display: flex;
  justify-content: space-between;
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

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.card-header h2 {
  margin: 0;
}

.text-muted {
  color: #999;
}

.key-display {
  margin-top: 16px;
  display: flex;
  gap: 8px;
}
</style>

