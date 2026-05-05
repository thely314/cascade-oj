<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import md5 from 'js-md5'
import { deleteUser, getUsers, updateUserInfo, updateUserPassword } from '../../api/admin'
import type { UserInfo } from '../../api/types'
import Modal from '../../components/Modal.vue'

const users = ref<UserInfo[]>([])
const loading = ref(false)
const error = ref('')

// Modal state
const editModal = reactive({
  show: false,
  userId: 0,
  username: '',
  email: ''
})

const passwordModal = reactive({
  show: false,
  userId: 0,
  username: '',
  password: ''
})

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

const onEditUser = (user: UserInfo) => {
  editModal.userId = user.userId
  editModal.username = user.username
  editModal.email = user.email
  editModal.show = true
}

const submitEdit = async () => {
  if (!editModal.username.trim() || !editModal.email.trim()) {
    window.alert('请填写完整信息')
    return
  }
  try {
    await updateUserInfo(editModal.userId, {
      username: editModal.username.trim(),
      email: editModal.email.trim(),
    })
    editModal.show = false
    await fetchUsers()
  } catch (err) {
    console.error(err)
    window.alert('修改用户信息失败')
  }
}

const onUpdatePassword = (user: UserInfo) => {
  passwordModal.userId = user.userId
  passwordModal.username = user.username
  passwordModal.password = ''
  passwordModal.show = true
}

const submitPassword = async () => {
  if (!passwordModal.password.trim()) {
    window.alert('请输入新密码')
    return
  }
  try {
    const encryptedPwd = md5.md5(passwordModal.password.trim())
    await updateUserPassword(passwordModal.userId, { password: encryptedPwd })
    passwordModal.show = false
    window.alert('密码已更新')
  } catch (err) {
    console.error(err)
    window.alert('修改密码失败')
  }
}

const onDeleteUser = async (user: UserInfo) => {
  if (!window.confirm(`确认删除用户 ${user.username}（ID ${user.userId}）？`)) return
  try {
    await deleteUser(user.userId)
    await fetchUsers()
  } catch (err) {
    console.error(err)
    window.alert('删除用户失败')
  }
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
        <button class="ghost" @click="fetchUsers" :disabled="loading">
          {{ loading ? 'Loading...' : 'Refresh' }}
        </button>
      </div>
    </header>

    <section class="panel">
      <header class="panel-header">
        <h2>User Directory</h2>
        <a href="#">View all</a>
      </header>
      <div v-if="error" class="error-message">{{ error }}</div>
      <ul v-else class="list">
        <li v-for="user in users" :key="user.userId" class="list-item">
          <div class="list-main">
            <p class="list-title">{{ user.username }}</p>
            <p class="list-meta">ID {{ user.userId }} · {{ user.email }}</p>
          </div>
          <div class="list-right">
            <button class="ghost" @click="onEditUser(user)">Edit</button>
            <button class="ghost" @click="onUpdatePassword(user)">Reset Password</button>
            <button class="ghost" @click="onDeleteUser(user)">Delete</button>
          </div>
        </li>
      </ul>
    </section>

    <!-- Edit User Modal -->
    <Modal
      title="Edit User"
      :show="editModal.show"
      @close="editModal.show = false"
      @submit="submitEdit"
    >
      <div class="form-group">
        <label>用户名</label>
        <input v-model="editModal.username" placeholder="请输入用户名" />
      </div>
      <div class="form-group">
        <label>邮箱</label>
        <input v-model="editModal.email" placeholder="请输入邮箱" />
      </div>
    </Modal>

    <!-- Reset Password Modal -->
    <Modal
      title="Reset Password"
      :show="passwordModal.show"
      @close="passwordModal.show = false"
      @submit="submitPassword"
    >
      <p style="margin-bottom: 16px; color: #666;">重置用户 <b>{{ passwordModal.username }}</b> 的密码</p>
      <div class="form-group">
        <label>新密码</label>
        <input v-model="passwordModal.password" type="password" placeholder="请输入新密码" @keyup.enter="submitPassword" />
      </div>
    </Modal>
  </div>
</template>

<style scoped src="./Users.css"></style>
