import { createRouter, createWebHistory } from 'vue-router'
import Home from '../pages/home/home.vue'
import About from '../pages/about/about.vue'
import ErrorPage from '../pages/error/error.vue'
import ContestsList from '../pages/contests-list/contests-list.vue'
import Contest from '../pages/contest/contest.vue'

const routes = [
    { path: '/', redirect: '/home' },
    { path: '/home', name: 'Home', component: Home, meta: { title: '主页 - Cascade' } },
    { path: '/competition', name: 'ContestsList', component: ContestsList, meta: { title: '比赛列表 - Cascade' } },
    { path: '/competition/:id', name: 'Contest', component: Contest, meta: { title: '比赛详情 - Cascade' } },
    { path: '/about', name: 'About', component: About, meta: { title: '关于我们 - Cascade' } },
    { path: '/login', name: 'Login', component: () => import('../../../login/src/pages/AuthPage.vue'), meta: { hideNav: true, title: '登录 - Cascade' } },
    {
        path: '/contest/:contestId/problem/:id', 
        name: 'ProblemDetail',
        component: () => import('../pages/problem/ProblemDetail.vue'),
        meta: { title: '题目详情 - Cascade' }
    },
    { path: '/:pathMatch(.*)*', name: 'NotFound', component: ErrorPage, meta: { title: '未找到 - Cascade' } }
]

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes,
})

// 同步路由标题到标签页标题
router.afterEach((to) => {
    const title = (to.meta && (to.meta as Record<string, any>).title) || 'Cascade'
    if (typeof title === 'string') {
        document.title = title
    }
})

export default router
