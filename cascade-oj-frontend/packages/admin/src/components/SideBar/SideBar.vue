<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter, RouterLink } from 'vue-router'
import { removeToken } from '../../utils/auth'

type MenuItem = { label: string; to: string }

const menuItems: MenuItem[] = [
  { label: 'Overview', to: '/' },
  { label: 'Contests', to: '/contests' },
  { label: 'Problems', to: '/problems' },
  { label: 'Submissions', to: '/submissions' },
  { label: 'Users', to: '/users' },
  { label: 'Logs', to: '/logs' },
]

// 用于读取路由信息，只读的意图
const route = useRoute()
// 用于执行路由操作，写的意图
const router = useRouter()
const activePath = computed(() => route.path)

const handleLogout = () => {
  removeToken()
  router.push({ name: 'Login', params: { sourceApp: 'admin' } })
}
</script>

<template>
  <aside class="sideBar">
    <div class="brand"></div>
    <nav class="nav">
      <RouterLink
        v-for="item in menuItems"
        :key="item.label"
        class="nav-item"
        :to="item.to"
        :class="{ active: activePath === item.to }"
      >
        {{ item.label }}
      </RouterLink>
    </nav>
    <div class="logout-section">
      <button @click="handleLogout" class="logout-button">退出登录</button>
    </div>
  </aside>
</template>

<style scoped src="./SideBar.css"></style>
