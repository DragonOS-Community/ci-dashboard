<template>
  <div class="main-layout">
    <!-- 侧边栏 -->
    <aside class="sidebar" :class="{ collapsed: sidebarCollapsed }">
      <div class="sidebar-header">
        <div class="logo">
          <div class="logo-icon">🐉</div>
          <span v-show="!sidebarCollapsed" class="logo-text">DragonOS CI</span>
        </div>
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
          <t-breadcrumb>
            <t-breadcrumb-item v-for="item in breadcrumbs" :key="item.path">
              <router-link v-if="item.path" :to="item.path">{{
                item.title
              }}</router-link>
              <span v-else>{{ item.title }}</span>
            </t-breadcrumb-item>
          </t-breadcrumb>
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
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from "vue";
import { useRouter, useRoute } from "vue-router";
import { useAdminStore } from "@/stores/admin";
import { MessagePlugin } from "tdesign-vue-next";

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
      { path: "/admin/system/users", label: "用户管理", icon: "user" },
    ],
  },
  {
    key: "monitor",
    title: "监控中心",
    items: [
      { path: "/admin/monitor/system", label: "系统监控", icon: "monitor" },
      {
        path: "/admin/monitor/statistics",
        label: "性能统计",
        icon: "chart-bar",
      },
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
  const result = [{ title: "首页", path: "/admin/dashboard" }];

  let currentPath = "";
  const routeMap = {
    dashboard: "仪表盘",
    test: "测试",
    system: "系统",
    monitor: "监控",
    settings: "设置",
    overview: "概览",
    runs: "运行",
    reports: "报告",
    "api-keys": "API密钥",
    projects: "项目",
    users: "用户",
    system: "监控",
    statistics: "统计",
    profile: "个人中心",
    config: "配置",
  };

  paths.forEach((path) => {
    currentPath += `/${path}`;
    const title = routeMap[path] || path;
    if (result.length > 0) {
      result[result.length - 1].path = currentPath;
    }
    if (paths.indexOf(path) < paths.length - 1) {
      result.push({ title, path: "" });
    } else {
      result.push({ title, path: "" });
    }
  });

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
}

.logo-icon {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #fcd34d 0%, #f59e0b 100%);
  border-radius: 10px;
  font-size: 24px;
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
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.user-button {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border-radius: 8px;
}

.username {
  font-size: 14px;
  color: #1f2937;
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
