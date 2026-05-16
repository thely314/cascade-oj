<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getProblems, publishProblem, disableProblem, deleteProblem } from '../../api/admin'
import { ProblemStatus, type ProblemMetadata } from '../../api/types'
import { useRouter } from 'vue-router'
const router = useRouter()

const handleCreate = () => {
  router.push('/problems/new') 
}

const problems = ref<ProblemMetadata[]>([])
const loading = ref(false)
const error = ref('')

const fetchProblems = async () => {
  loading.value = true
  error.value = ''
  try {
    const response = await getProblems()
    console.log('后端原始返回：', response)
    problems.value = response.problems || []
  } catch (err) {
    error.value = 'Failed to load problems'
    console.error(err)
  } finally {
    loading.value = false
  }
}

// 点击发布题目
const handlePublish = async (id: number) => {
  try {
    await publishProblem(id)
    await fetchProblems()
  } catch (err) {
    console.error('Failed to publish', err)
  }
}

// 点击禁用题目
const handleDisable = async (id: number) => {
  try {
    await disableProblem(id)
    await fetchProblems() 
  } catch (err) {
    console.error('Failed to disable', err)
  }
}

const handleDelete = async (id: number) => {
  if (!confirm('确定要删除这道题吗？')) return
  try {
    await deleteProblem(id)
    await fetchProblems() // 刷新列表
  } catch (err) {
    console.error('Delete failed', err)
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
        <button class="primary" @click="handleCreate">New Problem</button>
      </div>
    </header>

    <section class="panel">
      <header class="panel-header">
        <h2>Problem Library</h2>
        <a href="#">View all</a>
      </header>
      <div v-if="error" class="error-message">{{ error }}</div>
      <ul v-else class="list">
        <li v-for="p in problems" :key="p.id" class="list-item">
          
          <div class="list-main">
            <div style="display: flex; align-items: center; gap: 8px;">
              <p class="list-title">{{p.title }}</p>
              
              <span v-if="(p.Status) === ProblemStatus.PROBLEM_STATUS_ENABLED" class="badge badge-live">已启用</span>
              <span v-else-if="(p.Status) === ProblemStatus.PROBLEM_STATUS_DISABLED" class="badge badge-muted">已禁用</span>
            </div>
            
            <p class="list-meta">ID {{p.id }} · {{p.timeLimitMs }}ms · {{p.memoryLimitMb }}MB</p>
          </div>

          <div class="list-right">
            <template v-if="(p.Status) === ProblemStatus.PROBLEM_STATUS_ENABLED">
              <button class="ghost" @click="handleDisable(p.id)">禁用</button>
            </template>
            
            <template v-else>
              <button class="primary" @click="handlePublish(p.id)">发布</button>
              <button class="ghost" @click="router.push(`/problems/edit/${p.id}`)">编辑</button>
            </template>
            
            <button class="ghost" style="color: #ff4d4f" @click="handleDelete(p.id)">删除</button>
          </div>

        </li>
      </ul>
    </section>
  </div>
</template>

<style scoped src="./Problems.css"></style>
