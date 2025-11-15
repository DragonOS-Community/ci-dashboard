<template>
  <div class="projects-container">
    <t-card class="projects-card">
      <div class="card-header">
        <h2>项目管理</h2>
        <t-button theme="primary" @click="handleCreate"> 创建项目 </t-button>
      </div>

      <t-loading :loading="loading">
        <t-table
          :data="projects"
          :columns="columns"
          hover
          row-key="id"
          :pagination="pagination"
          @page-change="handlePageChange"
          @page-size-change="handlePageSizeChange"
        >
          <template #description="{ row }">
            <span v-if="row.description" class="description-text">{{
              row.description
            }}</span>
            <span v-else class="text-muted">暂无描述</span>
          </template>
          <template #created_at="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
          <template #updated_at="{ row }">
            {{ formatDate(row.updated_at) }}
          </template>
          <template #operation="{ row }">
            <t-space>
              <t-button
                theme="primary"
                variant="text"
                size="small"
                @click="handleEdit(row)"
              >
                编辑
              </t-button>
              <t-button
                theme="danger"
                variant="text"
                size="small"
                @click="handleDelete(row.id)"
              >
                删除
              </t-button>
            </t-space>
          </template>
        </t-table>
      </t-loading>
    </t-card>

    <!-- 创建/编辑项目对话框 -->
    <t-dialog
      v-model:visible="dialogVisible"
      :title="isEdit ? '编辑项目' : '创建项目'"
      width="600px"
      @confirm="handleSubmit"
      @cancel="handleCancel"
    >
      <t-form ref="formRef" :data="formData" :rules="rules" label-width="100px">
        <t-form-item label="项目名称" name="name">
          <t-input
            v-model="formData.name"
            placeholder="请输入项目名称"
            :maxlength="255"
            show-word-limit
          />
        </t-form-item>
        <t-form-item label="项目描述" name="description">
          <t-textarea
            v-model="formData.description"
            placeholder="请输入项目描述（可选）"
            :autosize="{ minRows: 3, maxRows: 6 }"
            :maxlength="1000"
            show-word-limit
          />
        </t-form-item>
      </t-form>
    </t-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { MessagePlugin, DialogPlugin } from "tdesign-vue-next";
import type { FormInstanceFunctions } from "tdesign-vue-next";
import {
  getProjects,
  createProject,
  updateProject,
  deleteProject,
  type Project,
  type ProjectData,
} from "@/api/admin";

const loading = ref(false);
const projects = ref<Project[]>([]);
const dialogVisible = ref(false);
const isEdit = ref(false);
const editingId = ref<number | null>(null);
const formRef = ref<FormInstanceFunctions | null>(null);

const formData = ref<ProjectData>({
  name: "",
  description: "",
});

const rules = {
  name: [
    { required: true, message: "请输入项目名称", type: "error" },
    { max: 255, message: "项目名称不能超过255个字符", type: "error" },
  ],
  description: [
    { max: 1000, message: "项目描述不能超过1000个字符", type: "error" },
  ],
};

const pagination = ref({
  current: 1,
  pageSize: 20,
  total: 0,
  showJumper: true,
  showSizePicker: true,
  pageSizeOptions: [10, 20, 50, 100],
});

const columns = [
  { colKey: "id", title: "ID", width: 80 },
  { colKey: "name", title: "项目名称", width: 200 },
  { colKey: "description", title: "项目描述", ellipsis: true },
  { colKey: "created_at", title: "创建时间", width: 180 },
  { colKey: "updated_at", title: "更新时间", width: 180 },
  { colKey: "operation", title: "操作", width: 150, fixed: "right" },
];

// 格式化日期
const formatDate = (dateString: string) => {
  if (!dateString) return "-";
  const date = new Date(dateString);
  return date.toLocaleString("zh-CN", {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
  });
};

// 获取项目列表
const fetchProjects = async () => {
  loading.value = true;
  try {
    const res = await getProjects();
    projects.value = res.data || [];
    pagination.value.total = projects.value.length;
  } catch (error: any) {
    console.error("Failed to fetch projects:", error);
    MessagePlugin.error("获取项目列表失败");
  } finally {
    loading.value = false;
  }
};

// 创建项目
const handleCreate = () => {
  isEdit.value = false;
  editingId.value = null;
  formData.value = {
    name: "",
    description: "",
  };
  dialogVisible.value = true;
};

// 编辑项目
const handleEdit = (row: Project) => {
  isEdit.value = true;
  editingId.value = row.id;
  formData.value = {
    name: row.name,
    description: row.description || "",
  };
  dialogVisible.value = true;
};

// 提交表单
const handleSubmit = async () => {
  const valid = await formRef.value?.validate();
  if (!valid) {
    return;
  }

  try {
    if (isEdit.value && editingId.value !== null) {
      // 更新项目
      await updateProject(editingId.value, formData.value);
      MessagePlugin.success("项目更新成功");
    } else {
      // 创建项目
      await createProject(formData.value);
      MessagePlugin.success("项目创建成功");
    }
    dialogVisible.value = false;
    await fetchProjects();
  } catch (error: any) {
    const errorMsg = error.response?.data?.message || "操作失败，请稍后重试";
    MessagePlugin.error(errorMsg);
  }
};

// 删除项目
const handleDelete = async (id: number) => {
  DialogPlugin.confirm({
    header: "确认删除",
    body: "确定要删除此项目吗？删除后，该项目下的所有测试运行记录和API密钥也将被删除。",
    confirmBtn: { theme: "danger" },
    onConfirm: async () => {
      try {
        await deleteProject(id);
        MessagePlugin.success("项目删除成功");
        await fetchProjects();
      } catch (error: any) {
        const errorMsg =
          error.response?.data?.message || "删除失败，请稍后重试";
        MessagePlugin.error(errorMsg);
      }
    },
  });
};

// 取消对话框
const handleCancel = () => {
  formRef.value?.reset();
  dialogVisible.value = false;
};

// 分页变化
const handlePageChange = (page: number) => {
  pagination.value.current = page;
  // 这里可以实现服务端分页，目前是前端分页
};

const handlePageSizeChange = (size: number) => {
  pagination.value.pageSize = size;
  pagination.value.current = 1;
};

onMounted(() => {
  fetchProjects();
});
</script>

<style scoped>
.projects-container {
  padding: 24px;
  background: #f9fafb;
  min-height: calc(100vh - 64px);
}

.projects-card {
  background: #ffffff;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  transition: box-shadow 0.15s ease;
}

.projects-card:hover {
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

.description-text {
  color: #6b7280;
  font-size: 14px;
}

.text-muted {
  color: #9ca3af;
  font-size: 14px;
  font-style: italic;
}

:deep(.t-table) {
  background: #ffffff;
}

:deep(.t-table__header) {
  background: #fafafa;
}

:deep(.t-table__row:hover) {
  background: #f9fafb;
}

:deep(.t-button--theme-primary) {
  background: #f59e0b;
  border-color: #f59e0b;
}

:deep(.t-button--theme-primary:hover) {
  background: #d97706;
  border-color: #d97706;
}

:deep(.t-button--variant-text.t-button--theme-primary) {
  color: #d97706;
  background: transparent;
}

:deep(.t-button--variant-text.t-button--theme-primary:hover) {
  color: #f59e0b;
  background: rgba(245, 158, 11, 0.1);
}

:deep(.t-dialog__header) {
  border-bottom: 1px solid #f3f4f6;
}

:deep(.t-form-item__label) {
  color: #1f2937;
  font-weight: 500;
}

:deep(.t-input),
:deep(.t-textarea) {
  border-radius: 8px;
}

:deep(.t-input:focus),
:deep(.t-textarea:focus) {
  border-color: #f59e0b;
  box-shadow: 0 0 0 2px rgba(245, 158, 11, 0.1);
}

@media (max-width: 768px) {
  .projects-container {
    padding: 16px;
  }

  .card-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }

  .card-header h2 {
    font-size: 16px;
  }
}
</style>
