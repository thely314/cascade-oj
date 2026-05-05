import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router';

const routes: Array<RouteRecordRaw> = [
  {
    path: '/:sourceApp/login',
    name: 'Login',
    component: () => import('../pages/AuthPage.vue'),
    props: true // 将路由参数作为 props 传递给组件
  }
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

export default router;