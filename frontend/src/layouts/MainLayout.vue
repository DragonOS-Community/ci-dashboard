<template>
  <div class="main-layout">
    <!-- 侧边栏 -->
    <aside class="sidebar" :class="{ collapsed: sidebarCollapsed }">
      <div class="sidebar-header">
        <router-link to="/" class="logo">
          <img :src="logoImage" class="logo-icon" alt="DragonOS Logo" />
          <span v-show="!sidebarCollapsed" class="logo-text">DragonOS CI</span>
        </router-link>
      </div>

      <nav class="sidebar-nav">
        <div v-for="menu in menuItems" :key="menu.key" class="menu-group">
          <div
            v-if="menu.title"
            v-show="!sidebarCollapsed"
            class="menu-group-title"
          >
            {{ menu.title }}
          </div>
          <router-link
            v-for="item in menu.items"
            :key="item.path"
            :to="item.path"
            class="menu-item"
            :class="{ active: $route.path.startsWith(item.path) }"
          >
            <t-icon :name="item.icon" class="menu-icon" />
            <span v-show="!sidebarCollapsed" class="menu-text">{{
              item.label
            }}</span>
          </router-link>
        </div>
      </nav>

      <div class="sidebar-footer">
        <t-button
          variant="text"
          shape="square"
          class="collapse-btn"
          @click="toggleSidebar"
        >
          <t-icon :name="sidebarCollapsed ? 'chevron-right' : 'chevron-left'" />
        </t-button>
      </div>
    </aside>

    <!-- 主内容区 -->
    <div class="main-content">
      <!-- 顶部导航 -->
      <header class="header">
        <div class="header-left">
          <nav class="breadcrumb-nav">
            <div class="breadcrumb-list">
              <div
                v-for="(item, index) in breadcrumbs"
                :key="`${item.path || 'current'}-${index}`"
                class="breadcrumb-item"
              >
                <router-link
                  v-if="item.path"
                  :to="item.path"
                  class="breadcrumb-link"
                >
                  <t-icon
                    v-if="index === 0"
                    name="home"
                    class="breadcrumb-icon"
                  />
                  {{ item.title }}
                </router-link>
                <span v-else class="breadcrumb-current">
                  <t-icon
                    v-if="index === 0"
                    name="home"
                    class="breadcrumb-icon"
                  />
                  {{ item.title }}
                </span>
                <t-icon
                  v-if="index < breadcrumbs.length - 1"
                  name="chevron-right"
                  class="breadcrumb-separator"
                />
              </div>
            </div>
          </nav>
        </div>

        <div class="header-right">
          <t-dropdown>
            <t-button variant="text" class="user-button">
              <t-avatar size="32">A</t-avatar>
              <span class="username">{{
                adminStore.user?.username || "Admin"
              }}</span>
              <t-icon name="chevron-down" />
            </t-button>
            <t-dropdown-menu>
              <t-dropdown-item @click="goToProfile">
                <t-icon name="user" /> 个人中心
              </t-dropdown-item>
              <t-dropdown-item @click="handleLogout">
                <t-icon name="logout" /> 退出登录
              </t-dropdown-item>
            </t-dropdown-menu>
          </t-dropdown>
        </div>
      </header>

      <!-- 页面内容 -->
      <main class="content">
        <router-view v-slot="{ Component, route }">
          <transition name="fade" mode="out-in">
            <component :is="Component" :key="route.path" />
          </transition>
        </router-view>
      </main>

      <!-- 页脚 -->
      <Footer />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from "vue";
import { useRouter, useRoute } from "vue-router";
import { useAdminStore } from "@/stores/admin";
import { MessagePlugin } from "tdesign-vue-next";
import logoImage from "@/assets/dragonos.jpeg";
import Footer from "@/components/Footer.vue";

const router = useRouter();
const route = useRoute();
const adminStore = useAdminStore();

const sidebarCollapsed = ref(false);

const toggleSidebar = () => {
  sidebarCollapsed.value = !sidebarCollapsed.value;
};

const menuItems = [
  {
    key: "dashboard",
    items: [{ path: "/admin/dashboard", label: "仪表盘", icon: "dashboard" }],
  },
  {
    key: "test",
    title: "测试管理",
    items: [
      { path: "/admin/test/overview", label: "测试概览", icon: "chart-line" },
      { path: "/admin/test/runs", label: "测试运行", icon: "play-circle" },
      { path: "/admin/test/reports", label: "测试报告", icon: "file-text" },
    ],
  },
  {
    key: "system",
    title: "系统管理",
    items: [
      { path: "/admin/system/api-keys", label: "API密钥", icon: "key" },
      { path: "/admin/system/projects", label: "项目管理", icon: "folder" },
    ],
  },
  {
    key: "settings",
    title: "系统设置",
    items: [
      {
        path: "/admin/settings/profile",
        label: "个人中心",
        icon: "user-circle",
      },
      { path: "/admin/settings/config", label: "系统配置", icon: "setting" },
    ],
  },
];

// 面包屑导航
const breadcrumbs = computed(() => {
  const paths = route.path.split("/").filter(Boolean);
  const result = [];

  // 路由映射表
  const routeMap = {
    admin: "管理后台",
    dashboard: "仪表盘",
    test: "测试管理",
    system: "系统管理",
    settings: "系统设置",
    overview: "测试概览",
    runs: "测试运行",
    reports: "测试报告",
    "api-keys": "API密钥",
    projects: "项目管理",
    profile: "个人中心",
    config: "系统配置",
  };

  // 构建面包屑路径
  let currentPath = "";
  paths.forEach((path, index) => {
    currentPath += `/${path}`;
    const title = routeMap[path] || path;
    const isLast = index === paths.length - 1;

    result.push({
      title,
      path: isLast ? null : currentPath, // 最后一项不设置路径
    });
  });

  // 如果路径为空或只有 admin，添加首页
  if (
    result.length === 0 ||
    (result.length === 1 && result[0].title === "管理后台")
  ) {
    return [{ title: "首页", path: "/admin/dashboard" }];
  }

  return result;
});

const goToProfile = () => {
  router.push("/admin/settings/profile");
};

const handleLogout = async () => {
  try {
    await adminStore.logout();
    MessagePlugin.success("退出登录成功");
    router.push("/admin/login");
  } catch (error) {
    MessagePlugin.error("退出登录失败");
  }
};
</script>

<style scoped>
.main-layout {
  display: flex;
  height: 100vh;
  background-color: #f9fafb;
}

/* 侧边栏样式 */
.sidebar {
  width: 240px;
  background-color: #fafafa;
  border-right: 1px solid #e5e7eb;
  display: flex;
  flex-direction: column;
  transition: width 0.3s ease;
  position: relative;
}

.sidebar.collapsed {
  width: 64px;
}

.sidebar-header {
  padding: 24px 20px;
  border-bottom: 1px solid #e5e7eb;
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
  text-decoration: none;
  cursor: pointer;
  transition: opacity 0.2s ease;
}

.logo:hover {
  opacity: 0.8;
}

.logo-icon {
  width: 40px;
  height: 40px;
  object-fit: contain;
  border-radius: 10px;
}

.logo-text {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
}

.sidebar-nav {
  flex: 1;
  padding: 16px 0;
  overflow-y: auto;
}

.menu-group {
  margin-bottom: 24px;
}

.menu-group-title {
  padding: 8px 20px;
  font-size: 12px;
  color: #9ca3af;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.menu-item {
  display: flex;
  align-items: center;
  padding: 10px 20px;
  color: #6b7280;
  text-decoration: none;
  transition: all 0.2s ease;
  position: relative;
}

.menu-item:hover {
  background-color: #f3f4f6;
  color: #1f2937;
}

.menu-item.active {
  background-color: #fef3c7;
  color: #d97706;
}

.menu-item.active::before {
  content: "";
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 3px;
  background-color: #f59e0b;
}

.menu-icon {
  font-size: 20px;
  min-width: 20px;
}

.menu-text {
  margin-left: 12px;
  font-size: 14px;
}

.sidebar-footer {
  padding: 16px;
  border-top: 1px solid #e5e7eb;
}

.collapse-btn {
  width: 100%;
  justify-content: center;
  color: #6b7280;
}

/* 主内容区样式 */
.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  min-height: 0;
}

.header {
  height: 64px;
  background-color: #ffffff;
  border-bottom: 1px solid #e5e7eb;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 32px;
}

.header-left {
  display: flex;
  align-items: center;
  flex: 1;
}

/* 面包屑导航样式 */
.breadcrumb-nav {
  display: flex;
  align-items: center;
}

.breadcrumb-list {
  display: flex;
  align-items: center;
  gap: 4px;
  flex-wrap: wrap;
}

.breadcrumb-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.breadcrumb-link {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  color: #6b7280;
  text-decoration: none;
  padding: 4px 8px;
  border-radius: 6px;
  transition: all 0.2s ease;
  font-weight: 400;
}

.breadcrumb-link:hover {
  color: #f59e0b;
  background-color: #fef9f3;
}

.breadcrumb-current {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  color: #1f2937;
  font-weight: 500;
  padding: 4px 8px;
}

.breadcrumb-icon {
  font-size: 16px;
  color: #9ca3af;
  flex-shrink: 0;
}

.breadcrumb-link:hover .breadcrumb-icon {
  color: #f59e0b;
}

.breadcrumb-current .breadcrumb-icon {
  color: #f59e0b;
}

.breadcrumb-separator {
  font-size: 14px;
  color: #d1d5db;
  margin: 0 4px;
  flex-shrink: 0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .breadcrumb-nav {
    overflow-x: auto;
    -webkit-overflow-scrolling: touch;
  }

  .breadcrumb-list {
    min-width: max-content;
  }

  .breadcrumb-link,
  .breadcrumb-current {
    font-size: 13px;
    padding: 3px 6px;
    white-space: nowrap;
  }

  .breadcrumb-icon {
    font-size: 14px;
  }

  .breadcrumb-separator {
    font-size: 12px;
    margin: 0 2px;
  }
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.user-button {
  display: flex !important;
  align-items: center !important;
  justify-content: center !important;
  gap: 10px;
  padding: 0 12px !important;
  border-radius: 8px;
  background-color: #ffffff;
  border: 1px solid #e5e7eb;
  transition: all 0.2s ease;
  height: 40px;
  line-height: 1;
  min-width: auto;
}

.user-button :deep(.t-button__text) {
  display: flex !important;
  align-items: center !important;
  gap: 10px;
  height: 100%;
}

.user-button:hover {
  background-color: #fef9f3;
  border-color: #f59e0b;
}

.user-button:hover :deep(.t-icon) {
  color: #f59e0b;
}

.username {
  font-size: 14px;
  color: #1f2937;
  font-weight: 500;
  display: flex;
  align-items: center;
  line-height: 1;
}

.user-button:hover .username {
  color: #f59e0b;
}

.user-button :deep(.t-icon) {
  color: #6b7280;
  transition: color 0.2s ease;
  font-size: 16px;
  flex-shrink: 0;
}

/* 头像样式优化 */
.user-button :deep(.t-avatar) {
  flex-shrink: 0;
}

.user-button :deep(.t-avatar__inner),
.user-button :deep(.t-avatar > div) {
  background: linear-gradient(135deg, #fcd34d 0%, #f59e0b 100%) !important;
  color: #ffffff !important;
  font-weight: 600;
  font-size: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* 下拉菜单样式 */
:deep(.t-dropdown__menu) {
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  border: 1px solid #e5e7eb;
  padding: 8px;
  background-color: #ffffff;
  margin-top: 8px;
}

:deep(.t-dropdown__menu:hover) {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

:deep(.t-dropdown-item) {
  border-radius: 8px;
  padding: 10px 12px;
  margin: 2px 0;
  transition: all 0.15s ease;
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #1f2937;
}

:deep(.t-dropdown-item:hover) {
  background-color: #fef9f3;
  color: #f59e0b;
}

:deep(.t-dropdown-item .t-icon) {
  font-size: 16px;
  color: #6b7280;
}

:deep(.t-dropdown-item:hover .t-icon) {
  color: #f59e0b;
}

.content {
  flex: 1;
  padding: 32px;
  overflow-y: auto;
}

/* 过渡动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
