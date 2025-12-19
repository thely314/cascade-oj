<script setup lang="ts">
type SubmissionStatus = 'Accepted' | 'Pending' | 'Wrong Answer'

type Submission = {
  uuid: string
  problemId: number
  userId: number
  status: SubmissionStatus
  time: string
  score: number
}

const submissions: Submission[] = [
  { uuid: 'c1a2', problemId: 401, userId: 18, status: 'Accepted', time: '2025-12-15 21:10', score: 100 },
  { uuid: 'd3b4', problemId: 402, userId: 25, status: 'Pending', time: '2025-12-15 21:08', score: 0 },
  { uuid: 'e5f6', problemId: 403, userId: 42, status: 'Wrong Answer', time: '2025-12-15 21:05', score: 30 },
]

const statusTone: Record<SubmissionStatus, string> = {
  Accepted: 'badge-live',
  Pending: 'badge-muted',
  'Wrong Answer': 'badge-dim',
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
        <button class="primary">Refresh</button>
      </div>
    </header>

    <section class="panel">
      <header class="panel-header">
        <h2>Latest Submissions</h2>
        <a href="#">View all</a>
      </header>
      <ul class="list">
        <li v-for="sub in submissions" :key="sub.uuid" class="list-item">
          <div class="list-main">
            <p class="list-title">Submission {{ sub.uuid }}</p>
            <p class="list-meta">Problem {{ sub.problemId }} · User {{ sub.userId }} · {{ sub.time }}</p>
          </div>
          <div class="list-right">
            <span class="badge" :class="statusTone[sub.status]">{{ sub.status }}</span>
            <span class="score">{{ sub.score }}</span>
            <button class="ghost">Open</button>
          </div>
        </li>
      </ul>
    </section>
  </div>
</template>

<style scoped src="./Submissions.css"></style>
