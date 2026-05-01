<script setup lang="ts">
import Sidebar from './components/SideBar/SideBar.vue'
import { RouterView, useRoute } from 'vue-router'

const route = useRoute();
</script>

<template>
  <div class="layout">
    <Sidebar v-if="route.name !== 'Login'" />
    <main class="content">
      <RouterView />
    </main>
  </div>
</template>

<style scoped>
:global(body) {
  margin: 0;
  background: #202020;
  color: #cbd5c0;
  font-family: "Segoe UI", "Helvetica Neue", Arial, sans-serif;
}

.layout {
  display: flex;
  min-height: 100vh;
  background: #202020;
}

.content {
  flex: 1;
  padding: 24px 32px;
  overflow: auto;
  background: #202020;
}

</style>
