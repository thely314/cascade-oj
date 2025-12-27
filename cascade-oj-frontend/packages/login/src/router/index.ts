import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router';

const routes: Array<RouteRecordRaw> = [
  { path: '/', component: () => import('../views/AuthPage.vue') }, // 根路径显示认证页面（登录/注册）
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;