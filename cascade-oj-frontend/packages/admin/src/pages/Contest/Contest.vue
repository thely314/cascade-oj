<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getContests } from '../../api/admin'
import type { ContestMetadata } from '../../api/types'

const contests = ref<ContestMetadata[]>([])
const loading = ref(false)
const error = ref('')

const fetchContests = async () => {
  loading.value = true
  error.value = ''
  try {
    const response = await getContests()
    contests.value = response.contests || []
  } catch (err) {
    error.value = 'Failed to load contests'
    console.error(err)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchContests()
})

const statusTone: Record<string, string> = {
  Scheduled: 'badge-muted',
  Running: 'badge-live',
  Finished: 'badge-dim',
  // Add default or map other statuses if needed
}

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
        <h1 class="title">Contests</h1>
      </div>
      <div class="actions">
        <button class="ghost" @click="fetchContests" :disabled="loading">
          {{ loading ? 'Loading...' : 'Refresh' }}
        </button>
        <button class="primary">New Contest</button>
      </div>
    </header>

    <section class="panel">
      <header class="panel-header">
        <h2>Contest List</h2>
        <a href="#">View all</a>
      </header>
      <div v-if="error" class="error-message">{{ error }}</div>
      <ul v-else class="list">
        <li v-for="contest in contests" :key="contest.id" class="list-item">
          <div class="list-main">
            <p class="list-title">{{ contest.title }}</p>
            <p class="list-meta">ID {{ contest.id }} · {{ formatDate(contest.start_time) }} → {{ formatDate(contest.end_time) }}</p>
          </div>
          <div class="list-right">
            <span class="badge" :class="statusTone[contest.status] || 'badge-muted'">{{ contest.status }}</span>
            <button class="ghost">Manage</button>
          </div>
        </li>
      </ul>
    </section>
  </div>
</template>

<style scoped src="./Contest.css"></style>
