<script setup>
const users = [
  { id: 18, username: 'alice', email: 'alice@example.com', status: 'Active' },
  { id: 25, username: 'bob', email: 'bob@example.com', status: 'Active' },
  { id: 42, username: 'charlie', email: 'charlie@example.com', status: 'Suspended' },
]

const statusTone = {
  Active: 'badge-live',
  Suspended: 'badge-dim',
}
</script>

<template>
  <div class="page">
    <header class="page-header">
      <div>
        <p class="eyebrow">Admin</p>
        <h1 class="title">Users</h1>
      </div>
      <div class="actions">
        <button class="ghost">Invite</button>
        <button class="primary">Add User</button>
      </div>
    </header>

    <section class="panel">
      <header class="panel-header">
        <h2>User Directory</h2>
        <a href="#">View all</a>
      </header>
      <ul class="list">
        <li v-for="user in users" :key="user.id" class="list-item">
          <div class="list-main">
            <p class="list-title">{{ user.username }}</p>
            <p class="list-meta">ID {{ user.id }} · {{ user.email }}</p>
          </div>
          <div class="list-right">
            <span class="badge" :class="statusTone[user.status]">{{ user.status }}</span>
            <button class="ghost">Edit</button>
          </div>
        </li>
      </ul>
    </section>
  </div>
</template>

<style scoped src="./Users.css"></style>
