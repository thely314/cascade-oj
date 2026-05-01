import { getToken } from '../utils/auth';
import type { Router, RouteLocationNormalized, NavigationGuardNext } from 'vue-router';

export function setupRouterGuard(router: Router) {
  router.beforeEach(
    async (to: RouteLocationNormalized, from: RouteLocationNormalized, next: NavigationGuardNext) => {
      const token = getToken();
      if (to.name !== 'Login' && !token) {
        next({ name: 'Login', params: { sourceApp: 'admin' }, query: { redirect: to.fullPath } });
      } else {
        next();
      }
    }
  );
}
