<script setup lang="ts">
const logs = [
  { id: 301, message: 'Admin login success', time: '2025-12-15 20:20' },
  { id: 302, message: 'Contest #87 updated', time: '2025-12-15 20:05' },
  { id: 303, message: 'Problem #402 published', time: '2025-12-15 19:58' },
]
</script>

<template>
  <div class="page">
    <header class="page-header">
      <div>
        <p class="eyebrow">Admin</p>
        <h1 class="title">Logs</h1>
      </div>
      <div class="actions">
        <button class="ghost">Export</button>
        <button class="primary">Refresh</button>
      </div>
    </header>

    <section class="panel">
      <header class="panel-header">
        <h2>Recent Logs</h2>
        <a href="#">View all</a>
      </header>
      <ul class="list">
        <li v-for="log in logs" :key="log.id" class="list-item">
          <div class="list-main">
            <p class="list-title">{{ log.message }}</p>
            <p class="list-meta">ID {{ log.id }} · {{ log.time }}</p>
          </div>
          <button class="ghost">Details</button>
        </li>
      </ul>
    </section>
  </div>
</template>

<style scoped src="./Log.css"></style>
