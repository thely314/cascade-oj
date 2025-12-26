import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'

import Overview from '../pages/Overview/Overview.vue'
import Contest from '../pages/Contest/Contest.vue'
import Problems from '../pages/Problems/Problems.vue'
import Submissions from '../pages/Submissions/Submissions.vue'
import Users from '../pages/Users/Users.vue'
import Log from '../pages/Log/Log.vue'

const routes: RouteRecordRaw[] = [
    { path: '/', name: 'overview', component: Overview },
    { path: '/contests', name: 'contests', component: Contest },
    { path: '/problems', name: 'problems', component: Problems },
    { path: '/submissions', name: 'submissions', component: Submissions },
    { path: '/users', name: 'users', component: Users },
    { path: '/logs', name: 'logs', component: Log },
]

const router = createRouter({
    history: createWebHistory(),
    routes,
})

export default router
