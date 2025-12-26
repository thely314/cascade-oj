<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getSubmissions } from '../../api/admin'
import type { SubmissionMetadata } from '../../api/types'

const submissions = ref<SubmissionMetadata[]>([])
const loading = ref(false)
const error = ref('')

const fetchSubmissions = async () => {
  loading.value = true
  error.value = ''
  try {
    const response = await getSubmissions({})
    submissions.value = response.submissions || []
  } catch (err) {
    error.value = 'Failed to load submissions'
    console.error(err)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchSubmissions()
})

const statusTone: Record<string, string> = {
  Accepted: 'badge-live',
  Pending: 'badge-muted',
  'Wrong Answer': 'badge-dim',
  // Add other statuses as needed
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
        <h1 class="title">Submissions</h1>
      </div>
      <div class="actions">
        <button class="ghost">Rejudge Pending</button>
        <button class="primary" @click="fetchSubmissions" :disabled="loading">
          {{ loading ? 'Loading...' : 'Refresh' }}
        </button>
      </div>
    </header>

    <section class="panel">
      <header class="panel-header">
        <h2>Latest Submissions</h2>
        <a href="#">View all</a>
      </header>
      <div v-if="error" class="error-message">{{ error }}</div>
      <ul v-else class="list">
        <li v-for="sub in submissions" :key="sub.submission_uuid" class="list-item">
          <div class="list-main">
            <p class="list-title">Submission {{ sub.submission_uuid }}</p>
            <p class="list-meta">Problem {{ sub.problem_id }} · User {{ sub.user_id }} · {{ formatDate(sub.submit_time) }}</p>
          </div>
          <div class="list-right">
            <span class="badge" :class="statusTone[sub.status] || 'badge-muted'">{{ sub.status }}</span>
            <span class="score">{{ sub.score }}</span>
            <button class="ghost">Open</button>
          </div>
        </li>
      </ul>
    </section>
  </div>
</template>

<style scoped src="./Submissions.css"></style>
