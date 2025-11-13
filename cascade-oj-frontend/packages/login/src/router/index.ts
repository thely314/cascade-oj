import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router';
import AuthPage from '../views/AuthPage.vue';

const routes: Array<RouteRecordRaw> = [
  { path: '/', component: AuthPage }, // 根路径显示认证页面（登录/注册）
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;