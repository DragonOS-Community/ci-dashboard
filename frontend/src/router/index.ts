import {
  createRouter,
  createWebHistory,
  type RouteRecordRaw,
} from "vue-router";
import { useAdminStore } from "@/stores/admin";

const routes: RouteRecordRaw[] = [
  // 公开页面 - 首页和测试详情
  {
    path: "/",
    name: "Home",
    component: () => import("@/views/Home.vue"),
  },
  {
    path: "/test-runs/:id",
    name: "TestRunDetail",
    component: () => import("@/views/TestRunDetail.vue"),
    props: true,
  },

  // 登录页面
  {
    path: "/admin/login",
    name: "AdminLogin",
    component: () => import("@/views/admin/Login.vue"),
    meta: { hideLayout: true },
  },

  // 管理后台主布局
  {
    path: "/admin",
    component: () => import("@/layouts/MainLayout.vue"),
    meta: { requiresAuth: true },
    children: [
      {
        path: "",
        redirect: "/admin/dashboard",
      },
      // 仪表盘
      {
        path: "dashboard",
        name: "Dashboard",
        component: () => import("@/views/admin/Dashboard.vue"),
        meta: { title: "仪表盘" },
      },

      // 测试管理
      {
        path: "test",
        children: [
          {
            path: "overview",
            name: "TestOverview",
            component: () => import("@/views/test/Overview.vue"),
            meta: { title: "测试概览" },
          },
          {
            path: "runs",
            name: "TestRuns",
            component: () => import("@/views/test/TestRuns.vue"),
            meta: { title: "测试运行" },
          },
          {
            path: "reports",
            name: "TestReports",
            component: () => import("@/views/test/Reports.vue"),
            meta: { title: "测试报告" },
          },
        ],
      },

      // 系统管理
      {
        path: "system",
        children: [
          {
            path: "api-keys",
            name: "APIKeys",
            component: () => import("@/views/admin/APIKeys.vue"),
            meta: { title: "API密钥管理" },
          },
          {
            path: "projects",
            name: "Projects",
            component: () => import("@/views/system/Projects.vue"),
            meta: { title: "项目管理" },
          },
          {
            path: "users",
            name: "Users",
            component: () => import("@/views/system/Users.vue"),
            meta: { title: "用户管理" },
          },
        ],
      },

      // 监控中心
      {
        path: "monitor",
        children: [
          {
            path: "system",
            name: "SystemMonitor",
            component: () => import("@/views/monitor/SystemMonitor.vue"),
            meta: { title: "系统监控" },
          },
          {
            path: "statistics",
            name: "Statistics",
            component: () => import("@/views/monitor/Statistics.vue"),
            meta: { title: "性能统计" },
          },
        ],
      },

      // 系统设置
      {
        path: "settings",
        children: [
          {
            path: "profile",
            name: "Profile",
            component: () => import("@/views/admin/Profile.vue"),
            meta: { title: "个人中心" },
          },
          {
            path: "config",
            name: "Config",
            component: () => import("@/views/settings/Config.vue"),
            meta: { title: "系统配置" },
          },
        ],
      },
    ],
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

// 路由守卫
router.beforeEach((to, from, next) => {
  const adminStore = useAdminStore();

  // 检查是否需要认证
  if (to.meta.requiresAuth && !adminStore.isAuthenticated) {
    next({ name: "AdminLogin", query: { redirect: to.fullPath } });
  }
  // 如果已登录且访问登录页，重定向到仪表盘
  else if (to.name === "AdminLogin" && adminStore.isAuthenticated) {
    next({ name: "Dashboard" });
  }
  // 默认情况
  else {
    next();
  }
});

// 路由后置钩卫 - 设置页面标题
router.afterEach((to) => {
  if (to.meta.title) {
    document.title = `${to.meta.title} - DragonOS CI Dashboard`;
  } else {
    document.title = "DragonOS CI Dashboard";
  }
});

export default router;
