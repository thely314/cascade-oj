<script setup lang="ts">
type ContestStatus = 'Scheduled' | 'Running' | 'Finished'

type Contest = {
  id: number
  title: string
  start: string
  end: string
  status: ContestStatus
}

const contests: Contest[] = [
  { id: 101, title: 'Spring Challenge', start: '2025-03-12 10:00', end: '2025-03-12 14:00', status: 'Scheduled' },
  { id: 92, title: 'Winter Cup', start: '2025-01-20 18:00', end: '2025-01-20 22:00', status: 'Finished' },
  { id: 87, title: 'Weekly #87', start: '2025-12-20 19:00', end: '2025-12-20 21:00', status: 'Running' },
]

const statusTone: Record<ContestStatus, string> = {
  Scheduled: 'badge-muted',
  Running: 'badge-live',
  Finished: 'badge-dim',
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
        <button class="ghost">Refresh</button>
        <button class="primary">New Contest</button>
      </div>
    </header>

    <section class="panel">
      <header class="panel-header">
        <h2>Contest List</h2>
        <a href="#">View all</a>
      </header>
      <ul class="list">
        <li v-for="contest in contests" :key="contest.id" class="list-item">
          <div class="list-main">
            <p class="list-title">{{ contest.title }}</p>
            <p class="list-meta">ID {{ contest.id }} · {{ contest.start }} → {{ contest.end }}</p>
          </div>
          <div class="list-right">
            <span class="badge" :class="statusTone[contest.status]">{{ contest.status }}</span>
            <button class="ghost">Manage</button>
          </div>
        </li>
      </ul>
    </section>
  </div>
</template>

<style scoped src="./Contest.css"></style>
