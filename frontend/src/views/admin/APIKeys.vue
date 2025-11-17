<template>
  <div class="api-keys-container">
    <t-card class="api-keys-card">
      <div class="card-header">
        <h2>API密钥列表</h2>
        <t-button theme="primary" @click="showCreateDialog = true">
          创建密钥
        </t-button>
      </div>

      <t-loading :loading="adminStore.loading">
        <t-table
          :data="adminStore.apiKeys"
          :columns="columns"
          hover
          row-key="id"
        >
          <template #project_id="{ row }">
            <span v-if="row.project_id">{{ row.project_id }}</span>
            <span v-else class="text-muted">全部项目</span>
          </template>
          <template #expires_at="{ row }">
            <span v-if="row.expires_at">{{
              new Date(row.expires_at).toLocaleString()
            }}</span>
            <span v-else class="text-muted">永不过期</span>
          </template>
          <template #operation="{ row }">
            <t-button theme="danger" size="small" @click="handleDelete(row.id)">
              删除
            </t-button>
          </template>
        </t-table>
      </t-loading>
    </t-card>

    <!-- 创建密钥对话框 -->
    <t-dialog
      v-model:visible="showCreateDialog"
      title="创建API密钥"
      @confirm="handleCreate"
    >
      <t-form :data="createForm" ref="createFormRef">
        <t-form-item label="名称" name="name">
          <t-input v-model="createForm.name" placeholder="请输入密钥名称" />
        </t-form-item>
        <t-form-item label="所属项目" name="project_id">
          <t-space style="width: 100%">
            <t-input
              v-model="selectedProjectName"
              placeholder="留空表示所有项目"
              readonly
              style="flex: 1"
            />
            <t-button
              theme="default"
              variant="outline"
              @click="showProjectDialog = true"
            >
              选择项目
            </t-button>
            <t-button
              v-if="createForm.project_id"
              theme="default"
              variant="text"
              @click="clearProject"
            >
              清除
            </t-button>
          </t-space>
        </t-form-item>
        <t-form-item label="过期时间" name="expires_at">
          <t-date-picker
            v-model="createForm.expires_at"
            placeholder="留空表示永不过期"
          />
        </t-form-item>
      </t-form>
    </t-dialog>

    <!-- 显示密钥对话框 -->
    <t-dialog
      v-model:visible="showKeyDialog"
      title="API密钥创建成功"
      :footer="false"
    >
      <t-alert theme="warning" message="请妥善保管此密钥，它只会显示一次！" />
      <div class="key-display">
        <t-input v-model="newKey" readonly />
        <t-button @click="copyKey">复制</t-button>
      </div>
    </t-dialog>

    <!-- 项目选择对话框 -->
    <t-dialog
      v-model:visible="showProjectDialog"
      title="选择项目"
      width="600px"
      @confirm="handleSelectProject"
      @cancel="handleCancelProjectDialog"
    >
      <t-loading :loading="projectsLoading">
        <t-table
          :data="projects"
          :columns="projectColumns"
          hover
          row-key="id"
          :pagination="false"
          :selected-row-keys="
            selectedProjectInDialog ? [selectedProjectInDialog.id] : []
          "
          @row-click="handleProjectRowClick"
        >
          <template #description="{ row }">
            <span v-if="row.description" class="description-text">{{
              row.description
            }}</span>
            <span v-else class="text-muted">暂无描述</span>
          </template>
        </t-table>
        <div
          v-if="projects.length === 0 && !projectsLoading"
          class="empty-state"
        >
          <p>暂无项目，请先创建项目</p>
        </div>
      </t-loading>
      <template #footer>
        <t-space>
          <t-button theme="default" @click="showProjectDialog = false">
            取消
          </t-button>
          <t-button
            theme="primary"
            :disabled="!selectedProjectInDialog"
            @click="handleSelectProject"
          >
            确认
          </t-button>
        </t-space>
      </template>
    </t-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from "vue";
import { useAdminStore } from "@/stores/admin";
import { MessagePlugin, DialogPlugin } from "tdesign-vue-next";
import { getProjects, type Project } from "@/api/admin";

const adminStore = useAdminStore();

const showCreateDialog = ref(false);
const showKeyDialog = ref(false);
const showProjectDialog = ref(false);
const newKey = ref("");
const createFormRef = ref(null);
const projects = ref<Project[]>([]);
const projectsLoading = ref(false);
const selectedProjectInDialog = ref<Project | null>(null);

const createForm = ref<{
  name: string;
  project_id: number | null;
  expires_at: string | null;
}>({
  name: "",
  project_id: null,
  expires_at: null,
});

const selectedProjectName = computed(() => {
  if (!createForm.value.project_id) {
    return "";
  }
  const project = projects.value.find(
    (p) => p.id === createForm.value.project_id,
  );
  return project ? `${project.name} (ID: ${project.id})` : "";
});

const projectColumns = [
  { colKey: "id", title: "ID", width: 80 },
  { colKey: "name", title: "项目名称", width: 200 },
  { colKey: "description", title: "项目描述", ellipsis: true },
];

const columns = [
  { colKey: "id", title: "ID", width: 80 },
  { colKey: "name", title: "名称", width: 200 },
  { colKey: "project_id", title: "项目ID", width: 120 },
  { colKey: "created_at", title: "创建时间", width: 180 },
  { colKey: "last_used_at", title: "最后使用", width: 180 },
  { colKey: "expires_at", title: "过期时间", width: 180 },
  { colKey: "operation", title: "操作", width: 100, fixed: "right" },
];

const fetchProjects = async () => {
  projectsLoading.value = true;
  try {
    const res = await getProjects();
    projects.value = res.data || [];
  } catch (error) {
    console.error("Failed to fetch projects:", error);
    MessagePlugin.error("获取项目列表失败");
  } finally {
    projectsLoading.value = false;
  }
};

const handleProjectRowClick = ({ row }: { row: Project }) => {
  selectedProjectInDialog.value = row;
};

const handleSelectProject = () => {
  if (selectedProjectInDialog.value) {
    createForm.value.project_id = selectedProjectInDialog.value.id;
    selectedProjectInDialog.value = null;
    showProjectDialog.value = false;
  }
};

const handleCancelProjectDialog = () => {
  selectedProjectInDialog.value = null;
  showProjectDialog.value = false;
};

const clearProject = () => {
  createForm.value.project_id = null;
  selectedProjectInDialog.value = null;
};

const handleCreate = async () => {
  if (!createForm.value.name) {
    MessagePlugin.warning("请输入密钥名称");
    return;
  }

  const projectId = createForm.value.project_id ?? undefined;
  const expiresAt = createForm.value.expires_at
    ? new Date(createForm.value.expires_at).toISOString()
    : undefined;

  const result = await adminStore.createKey(
    createForm.value.name,
    projectId,
    expiresAt,
  );
  if (result) {
    showCreateDialog.value = false;
    newKey.value = result.api_key;
    showKeyDialog.value = true;
    createForm.value = {
      name: "",
      project_id: null,
      expires_at: null,
    };
  }
};

const handleDelete = async (id: string) => {
  DialogPlugin.confirm({
    header: "确认删除",
    body: "确定要删除此API密钥吗？",
    confirmBtn: { theme: "danger" },
    onConfirm: async () => {
      await adminStore.removeKey(id);
    },
  });
};

const copyKey = () => {
  navigator.clipboard.writeText(newKey.value);
  MessagePlugin.success("已复制到剪贴板");
};

onMounted(async () => {
  adminStore.fetchAPIKeys();
  // 获取项目列表
  await fetchProjects();
});
</script>

<style scoped>
.api-keys-container {
  padding: 0;
  background: #f9fafb;
  min-height: calc(100vh - 64px);
}

.api-keys-card {
  background: #ffffff;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  transition: box-shadow 0.15s ease;
}

.api-keys-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 1px solid #f3f4f6;
}

.card-header h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
}

.text-muted {
  color: #999;
}

.key-display {
  margin-top: 16px;
  display: flex;
  gap: 8px;
}

.description-text {
  color: #6b7280;
  font-size: 14px;
}

.empty-state {
  padding: 40px;
  text-align: center;
  color: #9ca3af;
}

:deep(.t-table__row) {
  cursor: pointer;
}

:deep(.t-table__row:hover) {
  background: #f9fafb;
}

:deep(.t-table__row.t-is-selected) {
  background: rgba(245, 158, 11, 0.1);
}
</style>
