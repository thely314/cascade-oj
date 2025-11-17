<template>
  <NavBar v-if="!hideNav" />
  <router-view />
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import NavBar from './components/NavBar/NavBar.vue'

const route = useRoute()
const hideNav = computed(() => !!route.meta?.hideNav)
</script>

<style scoped>
:root {
  --app-bg: #f7f9fc;
}
body {
  background: var(--app-bg);
}
</style>
