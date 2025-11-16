<template>
  <div class="profile-container">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1 class="page-title">个人中心</h1>
      <p class="page-description">管理您的个人信息和账户设置</p>
    </div>

    <!-- 内容区域 -->
    <div class="content-wrapper">
      <div class="cards-grid">
        <!-- 个人信息卡片 -->
        <div class="info-card">
          <t-card>
            <div class="card-header">
              <h3 class="card-title">
                <t-icon name="user" />
                个人信息
              </h3>
            </div>
            <t-loading :loading="adminStore.loading">
              <div class="info-content">
                <div
                  class="info-item"
                  v-for="item in userInfo"
                  :key="item.label"
                >
                  <div class="info-label">
                    <span>{{ item.label }}</span>
                  </div>
                  <div class="info-value">{{ item.content }}</div>
                </div>
              </div>
            </t-loading>
          </t-card>
        </div>

        <!-- 修改密码卡片 -->
        <div class="password-card">
          <t-card>
            <div class="card-header">
              <h3 class="card-title">
                <t-icon name="lock-on" />
                修改密码
              </h3>
            </div>
            <t-form
              :data="passwordForm"
              ref="passwordFormRef"
              @submit="handleUpdatePassword"
              class="password-form"
            >
              <t-form-item label="当前密码" name="old_password">
                <t-input
                  v-model="passwordForm.old_password"
                  type="password"
                  placeholder="请输入当前密码"
                  :clearable="true"
                  size="large"
                />
              </t-form-item>
              <t-form-item label="新密码" name="new_password">
                <t-input
                  v-model="passwordForm.new_password"
                  type="password"
                  placeholder="请输入新密码（至少6位）"
                  :clearable="true"
                  size="large"
                />
              </t-form-item>
              <t-form-item label="确认新密码" name="confirm_password">
                <t-input
                  v-model="passwordForm.confirm_password"
                  type="password"
                  placeholder="请再次输入新密码"
                  :clearable="true"
                  size="large"
                />
              </t-form-item>
              <t-form-item>
                <t-button
                  theme="warning"
                  type="submit"
                  :loading="adminStore.loading"
                  size="large"
                  class="submit-btn"
                >
                  <t-icon name="check" />
                  更新密码
                </t-button>
              </t-form-item>
            </t-form>
          </t-card>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from "vue";
import { useRouter } from "vue-router";
import { useAdminStore } from "@/stores/admin";
import { MessagePlugin } from "tdesign-vue-next";

const router = useRouter();
const adminStore = useAdminStore();
const passwordFormRef = ref(null);

const passwordForm = ref({
  old_password: "",
  new_password: "",
  confirm_password: "",
});

const userInfo = computed(() => {
  if (!adminStore.user) {
    return [];
  }
  return [
    { label: "用户ID", content: adminStore.user.id },
    { label: "用户名", content: adminStore.user.username },
    {
      label: "角色",
      content: adminStore.user.role === "admin" ? "管理员" : "用户",
    },
    {
      label: "创建时间",
      content: adminStore.user.created_at
        ? new Date(adminStore.user.created_at).toLocaleString("zh-CN")
        : "-",
    },
    {
      label: "更新时间",
      content: adminStore.user.updated_at
        ? new Date(adminStore.user.updated_at).toLocaleString("zh-CN")
        : "-",
    },
  ];
});

const handleUpdatePassword = async () => {
  if (!passwordForm.value.old_password) {
    MessagePlugin.warning("请输入当前密码");
    return;
  }

  if (!passwordForm.value.new_password) {
    MessagePlugin.warning("请输入新密码");
    return;
  }

  if (passwordForm.value.new_password.length < 6) {
    MessagePlugin.warning("新密码长度至少为6位");
    return;
  }

  if (passwordForm.value.new_password !== passwordForm.value.confirm_password) {
    MessagePlugin.warning("两次输入的新密码不一致");
    return;
  }

  const success = await adminStore.changePassword(
    passwordForm.value.old_password,
    passwordForm.value.new_password,
    passwordForm.value.confirm_password,
  );

  if (success) {
    // 清空表单
    passwordForm.value = {
      old_password: "",
      new_password: "",
      confirm_password: "",
    };
  }
};

const handleLogout = () => {
  adminStore.logout();
  router.push("/admin/login");
};

onMounted(async () => {
  // 如果用户信息不存在，则获取
  if (!adminStore.user) {
    await adminStore.fetchProfile();
  }
});
</script>

<style scoped>
.profile-container {
  max-width: 1400px;
  margin: 0 auto;
}

/* 页面标题 */
.page-header {
  margin-bottom: 32px;
}

.page-title {
  font-size: 28px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 8px;
}

.page-description {
  font-size: 14px;
  color: #6b7280;
}

/* 内容区域 */
.content-wrapper {
  width: 100%;
}

.cards-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(500px, 1fr));
  gap: 24px;
}

/* 卡片通用样式 */
.info-card :deep(.t-card),
.password-card :deep(.t-card) {
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  border: none;
  overflow: hidden;
  transition: all 0.15s ease;
  background-color: #ffffff;
}

.info-card :deep(.t-card:hover),
.password-card :deep(.t-card:hover) {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  transform: translateY(-2px);
}

.info-card :deep(.t-card__body),
.password-card :deep(.t-card__body) {
  padding: 24px;
}

/* 卡片头部 */
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.card-title {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
  display: flex;
  align-items: center;
  gap: 8px;
}

.card-title :deep(.t-icon) {
  font-size: 20px;
  color: #f59e0b;
}

/* 个人信息内容 */
.info-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding-bottom: 20px;
  border-bottom: 1px solid #f3f4f6;
}

.info-item:last-child {
  border-bottom: none;
  padding-bottom: 0;
}

.info-label {
  font-size: 14px;
  font-weight: 500;
  color: #6b7280;
}

.info-value {
  font-size: 15px;
  color: #1f2937;
  font-weight: 500;
}

/* 密码表单 */
.password-form {
  margin-top: 8px;
}

.password-form :deep(.t-form-item) {
  margin-bottom: 20px;
}

.password-form :deep(.t-form-item__label) {
  font-size: 14px;
  font-weight: 500;
  color: #1f2937;
  margin-bottom: 8px;
}

.password-form :deep(.t-input) {
  border-radius: 8px;
}

.submit-btn {
  width: 100%;
  margin-top: 8px;
  border-radius: 8px;
  font-weight: 500;
}

.submit-btn :deep(.t-icon) {
  margin-right: 6px;
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .cards-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .profile-container {
    padding: 0;
  }

  .page-header {
    margin-bottom: 24px;
  }

  .page-title {
    font-size: 24px;
  }

  .cards-grid {
    gap: 16px;
  }

  .info-card :deep(.t-card__body),
  .password-card :deep(.t-card__body) {
    padding: 20px;
  }

  .card-title {
    font-size: 16px;
  }
}
</style>
