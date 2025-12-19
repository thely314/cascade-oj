<script setup lang="ts">
type ProblemStatus = 'Published' | 'Draft'

type Problem = {
  id: number
  title: string
  limits: string
  status: ProblemStatus
}

const problems: Problem[] = [
  { id: 401, title: 'Two Sum', limits: '1000ms · 128MB', status: 'Published' },
  { id: 402, title: 'Binary Search Tree', limits: '1500ms · 256MB', status: 'Draft' },
  { id: 403, title: 'Network Delay', limits: '2000ms · 256MB', status: 'Published' },
]

const statusTone: Record<ProblemStatus, string> = {
  Published: 'badge-live',
  Draft: 'badge-muted',
}
</script>

<template>
  <div class="page">
    <header class="page-header">
      <div>
        <p class="eyebrow">Admin</p>
        <h1 class="title">Problems</h1>
      </div>
      <div class="actions">
        <button class="ghost">Import</button>
        <button class="primary">New Problem</button>
      </div>
    </header>

    <section class="panel">
      <header class="panel-header">
        <h2>Problem Library</h2>
        <a href="#">View all</a>
      </header>
      <ul class="list">
        <li v-for="problem in problems" :key="problem.id" class="list-item">
          <div class="list-main">
            <p class="list-title">{{ problem.title }}</p>
            <p class="list-meta">ID {{ problem.id }} · {{ problem.limits }}</p>
          </div>
          <div class="list-right">
            <span class="badge" :class="statusTone[problem.status]">{{ problem.status }}</span>
            <button class="ghost">Edit</button>
          </div>
        </li>
      </ul>
    </section>
  </div>
</template>

<style scoped src="./Problems.css"></style>
