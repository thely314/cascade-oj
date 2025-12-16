import { createRouter, createWebHistory } from 'vue-router'
import Home from '../pages/home/home.vue'
import About from '../pages/about/about.vue'
import ErrorPage from '../pages/error/error.vue'
import ContestsList from '../pages/contests-list/contests-list.vue'
import Contest from '../pages/contest/contest.vue'

// 从 monorepo 中直接引入 login 包的 AuthPage
import Login from '../../../login/src/views/AuthPage.vue'

const routes = [
    { path: '/', redirect: '/problem/:id' },
    { path: '/home', name: 'Home', component: Home },
    { path: '/competition', name: 'ContestsList', component: ContestsList },
    { path: '/competition/:id', name: 'Contest', component: Contest },
    { path: '/about', name: 'About', component: About },
    { path: '/login', name: 'Login', component: Login, meta: { hideNav: true } },
    {
        path: '/problem/:id', // 提取 id 作为 API 参数
        name: 'ProblemDetail',
        component: () => import('../pages/problem/ProblemDetail.vue')
    },
    { path: '/:pathMatch(.*)*', name: 'NotFound', component: ErrorPage }
]

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes,
})

export default router
