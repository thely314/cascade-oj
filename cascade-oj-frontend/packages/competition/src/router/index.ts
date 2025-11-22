import { createRouter, createWebHistory } from 'vue-router'
import Home from '../pages/home/home.vue'
import About from '../pages/about/about.vue'
import ErrorPage from '../pages/error/error.vue'

// 从 monorepo 中直接引入 login 包的 AuthPage
import Login from '../../../login/src/views/AuthPage.vue'

const routes = [
    { path: '/', redirect: '/home' },
    { path: '/home', name: 'Home', component: Home },
    { path: '/about', name: 'About', component: About },
    { path: '/login', name: 'Login', component: Login, meta: { hideNav: true } },
    { path: '/:pathMatch(.*)*', name: 'NotFound', component: ErrorPage }
]

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes,
})

export default router
