<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getUsers } from '../../api/admin'
import type { UserInfo } from '../../api/types'

const users = ref<UserInfo[]>([])
const loading = ref(false)
const error = ref('')

const fetchUsers = async () => {
  loading.value = true
  error.value = ''
  try {
    const response = await getUsers({})
    users.value = response.users || []
  } catch (err) {
    error.value = 'Failed to load users'
    console.error(err)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchUsers()
})
</script>

<template>
  <div class="page">
    <header class="page-header">
      <div>
        <p class="eyebrow">Admin</p>
        <h1 class="title">Users</h1>
      </div>
      <div class="actions">
        <button class="ghost" @click="fetchUsers" :disabled="loading">
          {{ loading ? 'Loading...' : 'Refresh' }}
        </button>
        <button class="primary">Add User</button>
      </div>
    </header>

    <section class="panel">
      <header class="panel-header">
        <h2>User Directory</h2>
        <a href="#">View all</a>
      </header>
      <div v-if="error" class="error-message">{{ error }}</div>
      <ul v-else class="list">
        <li v-for="user in users" :key="user.user_id" class="list-item">
          <div class="list-main">
            <p class="list-title">{{ user.username }}</p>
            <p class="list-meta">ID {{ user.user_id }} · {{ user.email }}</p>
          </div>
          <div class="list-right">
            <button class="ghost">Edit</button>
          </div>
        </li>
      </ul>
    </section>
  </div>
</template>

<style scoped src="./Users.css"></style>
