<template>
  <div class="config-container">
    <t-card class="config-card">
      <div class="card-header">
        <h2>系统配置</h2>
        <t-button theme="default" variant="outline" @click="handleRefresh">
          <template #icon>
            <t-icon name="refresh" />
          </template>
          刷新
        </t-button>
      </div>

      <t-loading :loading="loading">
        <div class="config-content">
          <!-- 允许上传测试输出文件配置 -->
          <t-card class="config-item-card" shadow>
            <div class="config-item">
              <div class="config-item-header">
                <div class="config-item-title">
                  <h3>允许上传测试输出文件</h3>
                  <p class="config-item-desc">
                    控制是否允许通过API上传测试输出文件。关闭后，上传文件接口将拒绝所有请求。
                  </p>
                </div>
                <t-switch
                  v-model="allowUploadOutputFiles"
                  :loading="saving"
                  @change="handleToggleUploadFiles"
                />
              </div>
            </div>
          </t-card>

          <!-- 其他配置项可以在这里添加 -->
        </div>
      </t-loading>
    </t-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { MessagePlugin } from "tdesign-vue-next";
import {
  getSystemConfig,
  updateSystemConfig,
  type SystemConfig,
} from "@/api/admin";

const loading = ref(false);
const saving = ref(false);
const allowUploadOutputFiles = ref(false);

// 加载配置
const loadConfig = async () => {
  loading.value = true;
  try {
    const res = await getSystemConfig("allow_upload_output_files");
    allowUploadOutputFiles.value = res.data.value === "true";
  } catch (error) {
    console.error("Failed to load config:", error);
    // 如果配置不存在，使用默认值false
    allowUploadOutputFiles.value = false;
  } finally {
    loading.value = false;
  }
};

// 切换上传文件开关
const handleToggleUploadFiles = async (value: boolean) => {
  saving.value = true;
  try {
    await updateSystemConfig("allow_upload_output_files", {
      value: value ? "true" : "false",
      description: "是否允许上传测试输出文件",
    });
    MessagePlugin.success(
      value ? "已允许上传测试输出文件" : "已禁止上传测试输出文件",
    );
  } catch (error) {
    console.error("Failed to update config:", error);
    MessagePlugin.error("更新配置失败");
    // 恢复原值
    allowUploadOutputFiles.value = !value;
  } finally {
    saving.value = false;
  }
};

// 刷新配置
const handleRefresh = () => {
  loadConfig();
};

// 组件挂载时加载配置
onMounted(() => {
  loadConfig();
});
</script>

<style scoped>
.config-container {
  padding: 32px;
  background-color: #f9fafb;
  min-height: calc(100vh - 64px);
}

.config-card {
  background-color: #ffffff;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
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
  font-size: 24px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.config-content {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.config-item-card {
  background-color: #ffffff;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  transition: box-shadow 0.15s ease;
}

.config-item-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.config-item {
  padding: 8px 0;
}

.config-item-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 24px;
}

.config-item-title h3 {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 8px 0;
}

.config-item-desc {
  font-size: 14px;
  color: #6b7280;
  margin: 0;
  line-height: 1.6;
}

@media (max-width: 768px) {
  .config-container {
    padding: 16px;
  }

  .config-item-header {
    flex-direction: column;
    gap: 16px;
  }

  .card-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }
}
</style>
