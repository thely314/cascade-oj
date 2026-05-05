<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { addContestUser, getContestUsers, getContests, removeContestUser } from '../../api/admin'
import type { ContestMetadata, UserInfo } from '../../api/types'
import Modal from '../../components/Modal.vue'

const contests = ref<ContestMetadata[]>([])
const selectedContestUsers = ref<UserInfo[]>([])
const selectedContestId = ref<number | null>(null)
const loading = ref(false)
const error = ref('')

// Modal state
const userModal = reactive({
  show: false,
  contestId: 0,
  userId: '',
  mode: 'add' as 'add' | 'remove'
})

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

const onViewContestUsers = async (contestId: number) => {
  try {
    const response = await getContestUsers(contestId)
    selectedContestId.value = contestId
    selectedContestUsers.value = response.users || []
  } catch (err) {
    console.error(err)
    window.alert('加载参赛用户失败')
  }
}

const onAddContestUser = (contestId: number) => {
  userModal.contestId = contestId
  userModal.userId = ''
  userModal.mode = 'add'
  userModal.show = true
}

const onRemoveContestUser = (contestId: number) => {
  userModal.contestId = contestId
  userModal.userId = ''
  userModal.mode = 'remove'
  userModal.show = true
}

const submitUserOp = async () => {
  const userId = Number(userModal.userId)
  if (!Number.isInteger(userId) || userId <= 0) {
    window.alert('请输入有效的用户 ID')
    return
  }

  try {
    if (userModal.mode === 'add') {
      await addContestUser(userModal.contestId, userId)
    } else {
      await removeContestUser(userModal.contestId, userId)
    }
    userModal.show = false
    await onViewContestUsers(userModal.contestId)
  } catch (err) {
    console.error(err)
    window.alert(userModal.mode === 'add' ? '移入赛事失败' : '移出赛事失败')
  }
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
            <p class="list-meta">ID {{ contest.id }} · {{ formatDate(contest.startTime) }} → {{ formatDate(contest.endTime) }}</p>
          </div>
          <div class="list-right">
            <span class="badge" :class="statusTone[contest.status] || 'badge-muted'">{{ contest.status }}</span>
            <button class="ghost" @click="onViewContestUsers(contest.id)">Users</button>
            <button class="ghost" @click="onAddContestUser(contest.id)">Add User</button>
            <button class="ghost" @click="onRemoveContestUser(contest.id)">Remove User</button>
          </div>
        </li>
      </ul>
    </section>

    <section class="panel" v-if="selectedContestId !== null">
      <header class="panel-header">
        <h2>Contest {{ selectedContestId }} Users</h2>
      </header>
      <ul class="list" v-if="selectedContestUsers.length > 0">
        <li v-for="user in selectedContestUsers" :key="user.userId" class="list-item">
          <div class="list-main">
            <p class="list-title">{{ user.username }}</p>
            <p class="list-meta">ID {{ user.userId }} · {{ user.email }}</p>
          </div>
        </li>
      </ul>
      <p v-else class="list-meta">暂无参赛用户</p>
    </section>

    <!-- Contest User Modal (Add/Remove) -->
    <Modal
      :title="userModal.mode === 'add' ? '移入赛事' : '移出赛事'"
      :show="userModal.show"
      @close="userModal.show = false"
      @submit="submitUserOp"
    >
      <div class="form-group">
        <label>用户 ID</label>
        <input v-model="userModal.userId" type="number" placeholder="请输入用户 ID" @keyup.enter="submitUserOp" />
      </div>
    </Modal>
  </div>
</template>

<style scoped src="./Contest.css"></style>
