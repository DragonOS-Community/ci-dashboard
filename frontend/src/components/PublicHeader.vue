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
        <div class="header-actions">
          <t-button
            theme="default"
            variant="outline"
            @click="goToCommunity"
            class="community-btn"
          >
            <t-icon name="logo-github" />
            社区主页
          </t-button>
        </div>
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

const goToCommunity = () => {
  window.open("https://github.com/DragonOS-Community/DragonOS", "_blank");
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

.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.community-btn {
  border-radius: 8px;
  transition: all 0.2s ease;
  height: 40px;
  display: flex !important;
  align-items: center !important;
  justify-content: center !important;
  gap: 6px;
  border-color: #e5e7eb;
  color: #6b7280;
  font-weight: 500;
}

.community-btn :deep(.t-button__text) {
  display: flex !important;
  align-items: center !important;
  justify-content: center !important;
  gap: 6px;
  line-height: 1;
}

.community-btn :deep(.t-icon) {
  font-size: 16px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.community-btn:hover {
  border-color: #d1d5db;
  background-color: #f9fafb;
  color: #1f2937;
}

/* 响应式 */
@media (max-width: 768px) {
  .header-content {
    padding: 0 16px;
  }

  .header-actions {
    gap: 8px;
  }

  .community-btn :deep(.t-button__text) {
    font-size: 13px;
  }
}
</style>
