import { createRouter, createWebHistory } from 'vue-router';
import AIOpsPage from '@/pages/AIOpsPage.vue';

const router = createRouter({
  history: createWebHistory('/aiops'),
  routes: [
    {
      path: '/',
      name: 'aiops',
      component: AIOpsPage,
      meta: { title: 'AIOps - Cascade' },
    },
  ],
});

router.afterEach((to) => {
  const title = (to.meta?.title as string) || 'AIOps - Cascade';
  document.title = title;
});

export default router;
