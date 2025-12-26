<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, RouterLink } from 'vue-router'

type MenuItem = { label: string; to: string }

const menuItems: MenuItem[] = [
  { label: 'Overview', to: '/' },
  { label: 'Contests', to: '/contests' },
  { label: 'Problems', to: '/problems' },
  { label: 'Submissions', to: '/submissions' },
  { label: 'Users', to: '/users' },
  { label: 'Logs', to: '/logs' },
]

const route = useRoute()
const activePath = computed(() => route.path)
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
  </aside>
</template>

<style scoped src="./SideBar.css"></style>
