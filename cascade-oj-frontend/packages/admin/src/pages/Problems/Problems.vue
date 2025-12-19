<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getProblems } from '../../api/admin'
import type { ProblemMetadata } from '../../api/types'

const problems = ref<ProblemMetadata[]>([])
const loading = ref(false)
const error = ref('')

const fetchProblems = async () => {
  loading.value = true
  error.value = ''
  try {
    const response = await getProblems()
    problems.value = response.problems || []
  } catch (err) {
    error.value = 'Failed to load problems'
    console.error(err)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchProblems()
})
</script>

<template>
  <div class="page">
    <header class="page-header">
      <div>
        <p class="eyebrow">Admin</p>
        <h1 class="title">Problems</h1>
      </div>
      <div class="actions">
        <button class="ghost" @click="fetchProblems" :disabled="loading">
          {{ loading ? 'Loading...' : 'Refresh' }}
        </button>
        <button class="primary">New Problem</button>
      </div>
    </header>

    <section class="panel">
      <header class="panel-header">
        <h2>Problem Library</h2>
        <a href="#">View all</a>
      </header>
      <div v-if="error" class="error-message">{{ error }}</div>
      <ul v-else class="list">
        <li v-for="problem in problems" :key="problem.id" class="list-item">
          <div class="list-main">
            <p class="list-title">{{ problem.title }}</p>
            <p class="list-meta">ID {{ problem.id }} · {{ problem.time_limit_ms }}ms · {{ problem.memory_limit_mb }}MB</p>
          </div>
          <div class="list-right">
            <button class="ghost">Edit</button>
          </div>
        </li>
      </ul>
    </section>
  </div>
</template>

<style scoped src="./Problems.css"></style>
