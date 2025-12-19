<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getLogs } from '../../api/admin'
import type { LogEntry } from '../../api/types'

const logs = ref<LogEntry[]>([])
const loading = ref(false)
const error = ref('')

const fetchLogs = async () => {
  loading.value = true
  error.value = ''
  try {
    const response = await getLogs({})
    logs.value = response.logs || []
  } catch (err) {
    error.value = 'Failed to load logs'
    console.error(err)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchLogs()
})

const formatDate = (dateStr: string) => {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleString()
}
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
        <button class="primary" @click="fetchLogs" :disabled="loading">
          {{ loading ? 'Loading...' : 'Refresh' }}
        </button>
      </div>
    </header>

    <section class="panel">
      <header class="panel-header">
        <h2>Recent Logs</h2>
        <a href="#">View all</a>
      </header>
      <div v-if="error" class="error-message">{{ error }}</div>
      <ul v-else class="list">
        <li v-for="log in logs" :key="log.id" class="list-item">
          <div class="list-main">
            <p class="list-title">{{ log.message }}</p>
            <p class="list-meta">ID {{ log.id }} · {{ formatDate(log.timestamp) }}</p>
          </div>
          <button class="ghost">Details</button>
        </li>
      </ul>
    </section>
  </div>
</template>

<style scoped src="./Log.css"></style>
