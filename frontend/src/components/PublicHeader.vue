<template>
  <header class="header">
    <div class="header-content">
      <div
        class="logo"
        :class="{ 'logo-clickable': showBack }"
        @click="handleLogoClick"
      >
        <img :src="logoImage" class="logo-icon" alt="DragonOS Logo" />
        <span class="logo-text">DragonOS CI Dashboard</span>
      </div>
      <t-button theme="warning" variant="outline" @click="goToLogin">
        <t-icon name="user" />
        管理员登录
      </t-button>
    </div>
  </header>
</template>

<script setup>
import { useRouter } from "vue-router";
import logoImage from "@/assets/dragonos.jpeg";

const props = defineProps({
  // 是否显示返回功能（点击logo可返回）
  showBack: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(["back"]);

const router = useRouter();

const handleLogoClick = () => {
  if (props.showBack) {
    emit("back");
  }
};

const goToLogin = () => {
  router.push("/admin/login");
};
</script>

<style scoped>
/* 顶部导航 */
.header {
  background-color: #ffffff;
  border-bottom: 1px solid #e5e7eb;
  position: sticky;
  top: 0;
  z-index: 100;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.header-content {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 32px;
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
  transition: opacity 0.2s ease;
}

.logo-clickable {
  cursor: pointer;
}

.logo-clickable:hover {
  opacity: 0.8;
}

.logo-icon {
  width: 40px;
  height: 40px;
  object-fit: contain;
  border-radius: 10px;
}

.logo-text {
  font-size: 20px;
  font-weight: 600;
  color: #1f2937;
}

/* 响应式 */
@media (max-width: 768px) {
  .header-content {
    padding: 0 16px;
  }
}
</style>
