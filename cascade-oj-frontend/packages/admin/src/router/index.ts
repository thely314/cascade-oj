import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'

import Overview from '../pages/Overview/Overview.vue'
import Contest from '../pages/Contest/Contest.vue'
import Problems from '../pages/Problems/Problems.vue'
import Submissions from '../pages/Submissions/Submissions.vue'
import Users from '../pages/Users/Users.vue'
import Log from '../pages/Log/Log.vue'
import NewProblem from '../pages/Problems/NewProblem.vue'
import { setupRouterGuard } from './guard'

const routes: RouteRecordRaw[] = [
    { path: '/', name: 'overview', component: Overview, meta: { title: '总览 - Cascade' } },
    { path: '/contests', name: 'contests', component: Contest, meta: { title: '赛事 - Cascade' } },
    { path: '/problems', name: 'problems', component: Problems, meta: { title: '题目 - Cascade' } },
    { path: '/problems/new', name: 'new-problem', component: NewProblem, meta: { title: '新建题目 - Cascade' } },
    { path: '/submissions', name: 'submissions', component: Submissions, meta: { title: '提交 - Cascade' } },
    { path: '/users', name: 'users', component: Users, meta: { title: '用户 - Cascade' } },
    { path: '/logs', name: 'logs', component: Log, meta: { title: '日志 - Cascade' } },
    {
        path: '/:sourceApp/login',
        name: 'Login',
        component: () => import('../../../login/src/pages/AuthPage.vue'),
        props: true,
        meta: { title: '登录 - Cascade' }
    },
    { 
        path: '/problems/edit/:id', 
        name: 'edit-problem', 
        component: () => import('../pages/Problems/EditProblem.vue'),
        meta: { title: '编辑题目 - Cascade' } 
    },
]

const router = createRouter({
    history: createWebHistory('/admin'),
    routes,
})

setupRouterGuard(router);

// 同步路由标题到标签页标题
router.afterEach((to) => {
    const title = (to.meta && (to.meta as Record<string, any>).title) || 'Cascade'
    if (typeof title === 'string') {
        document.title = title
    }
})

export default router
