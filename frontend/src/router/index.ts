import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { useAdminStore } from '@/stores/admin'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/Home.vue'),
  },
  {
    path: '/test-runs/:id',
    name: 'TestRunDetail',
    component: () => import('@/views/TestRunDetail.vue'),
    props: true,
  },
  {
    path: '/admin',
    redirect: '/admin/login',
  },
  {
    path: '/admin/login',
    name: 'AdminLogin',
    component: () => import('@/views/admin/Login.vue'),
  },
  {
    path: '/admin/api-keys',
    name: 'APIKeys',
    component: () => import('@/views/admin/APIKeys.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/admin/profile',
    name: 'Profile',
    component: () => import('@/views/admin/Profile.vue'),
    meta: { requiresAuth: true },
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

// 路由守卫
router.beforeEach((to, _from, next) => {
  const adminStore = useAdminStore()

  if (to.meta.requiresAuth && !adminStore.isAuthenticated) {
    next({ name: 'AdminLogin', query: { redirect: to.fullPath } })
  } else {
    next()
  }
})

export default router