<template>
  <div class="profile-container">
    <t-layout>
      <t-header>
        <div class="header-content">
          <h1>个人面板</h1>
          <t-space>
            <t-button theme="default" @click="router.push('/admin/api-keys')"
              >API密钥管理</t-button
            >
            <span>欢迎，{{ adminStore.user?.username }}</span>
            <t-button theme="default" @click="handleLogout">退出</t-button>
          </t-space>
        </div>
      </t-header>
      <t-content>
        <div class="content-wrapper">
          <t-row :gutter="16">
            <t-col :span="12">
              <t-card title="个人信息">
                <t-loading :loading="adminStore.loading">
                  <t-descriptions :data="userInfo" :column="1" />
                </t-loading>
              </t-card>
            </t-col>
            <t-col :span="12">
              <t-card title="修改密码">
                <t-form
                  :data="passwordForm"
                  ref="passwordFormRef"
                  @submit="handleUpdatePassword"
                >
                  <t-form-item label="当前密码" name="old_password">
                    <t-input
                      v-model="passwordForm.old_password"
                      type="password"
                      placeholder="请输入当前密码"
                      :clearable="true"
                    />
                  </t-form-item>
                  <t-form-item label="新密码" name="new_password">
                    <t-input
                      v-model="passwordForm.new_password"
                      type="password"
                      placeholder="请输入新密码（至少6位）"
                      :clearable="true"
                    />
                  </t-form-item>
                  <t-form-item label="确认新密码" name="confirm_password">
                    <t-input
                      v-model="passwordForm.confirm_password"
                      type="password"
                      placeholder="请再次输入新密码"
                      :clearable="true"
                    />
                  </t-form-item>
                  <t-form-item>
                    <t-button
                      theme="primary"
                      type="submit"
                      :loading="adminStore.loading"
                    >
                      更新密码
                    </t-button>
                  </t-form-item>
                </t-form>
              </t-card>
            </t-col>
          </t-row>
        </div>
      </t-content>
    </t-layout>
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
</style>
