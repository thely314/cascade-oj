<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { getAnnouncements, postAnnouncement, putAnnouncement, deleteAnnouncement } from '../../api/admin'
import type { Announcement } from '../../api/types'
import Modal from '../../components/Modal.vue'

const stats = [
	{ label: '正在进行', value: '-' },
	{ label: '已结束', value: '-' },
	{ label: '题库题目数', value: '-' },
	{ label: '总提交数', value: '-' },
]

const announcements = ref<Announcement[]>([])
const loading = ref(false)

// Modal states
const createModal = reactive({
  show: false,
  title: '',
  content: ''
})

const editModal = reactive({
  show: false,
  announcementId: 0,
  title: '',
  content: ''
})

const fetchAnnouncements = async () => {
  loading.value = true
  try {
    const response = await getAnnouncements()
    announcements.value = response.announcements || []
  } catch (err) {
    console.error(err)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchAnnouncements()
})

const onNewAnnouncement = () => {
  createModal.title = ''
  createModal.content = ''
  createModal.show = true
}

const submitCreate = async () => {
  if (!createModal.title.trim() || !createModal.content.trim()) {
    globalThis.alert('请填写完整信息')
    return
  }
  try {
    await postAnnouncement({
      title: createModal.title.trim(),
      content: createModal.content.trim(),
    })
    createModal.show = false
    await fetchAnnouncements()
  } catch (err) {
    console.error(err)
    globalThis.alert('创建公告失败')
  }
}

const onOpenAnnouncement = (announcement: Announcement) => {
  editModal.announcementId = announcement.id
  editModal.title = announcement.title
  editModal.content = announcement.content
  editModal.show = true
}

const submitEdit = async () => {
  if (!editModal.title.trim() || !editModal.content.trim()) {
    globalThis.alert('请填写完整信息')
    return
  }
  try {
    await putAnnouncement(editModal.announcementId, {
      title: editModal.title.trim(),
      content: editModal.content.trim(),
    })
    editModal.show = false
    await fetchAnnouncements()
  } catch (err) {
    console.error(err)
    globalThis.alert('修改公告失败')
  }
}

const onDeleteAnnouncement = async (announcement: Announcement) => {
  if (!globalThis.confirm(`确认删除公告 "${announcement.title}"？`)) return
  try {
    await deleteAnnouncement(announcement.id)
    await fetchAnnouncements()
  } catch (err) {
    console.error(err)
    globalThis.alert('删除公告失败')
  }
}
</script>

<template>
	<div class="page">
		<header class="page-header">
			<div>
				<p class="eyebrow">Admin</p>
				<h1 class="title">Overview</h1>
			</div>
			<button class="primary" @click="onNewAnnouncement">New Announcement</button>
		</header>

		<section class="stat-card">
			<p class="row-center-align">比赛</p>
			<article class="stat-grid">
				<article class="stat-card">
					<p class="stat-label">{{ stats[0].label }}</p>
					<p class="stat-value">{{ stats[0].value }}</p>
				</article>
				<article class="stat-card">
					<p class="stat-label">{{ stats[1].label }}</p>
					<p class="stat-value">{{ stats[1].value }}</p>
				</article>
			</article>
			<article class="stat-grid">
				<article class="stat-card">
					<p class="stat-label">{{ stats[2].label }}</p>
					<p class="stat-value">{{ stats[2].value }}</p>
				</article>
				<article class="stat-card">
					<p class="stat-label">{{ stats[3].label }}</p>
					<p class="stat-value">{{ stats[3].value }}</p>
				</article>
			</article>
		</section>

		<section class="panel">
			<header class="panel-header">
				<h2>Announcements</h2>
				<a href="#">View all</a>
			</header>
			<ul class="list">
				<li v-for="item in announcements" :key="item.id" class="list-item">
					<div>
						<p class="list-title">{{ item.title }}</p>
						<p class="list-meta">{{ item.publisherName }}</p>
					</div>
					<div class="list-actions">
						<button class="ghost" @click="onOpenAnnouncement(item)">Edit</button>
						<button class="ghost danger" @click="onDeleteAnnouncement(item)">Delete</button>
					</div>
				</li>
			</ul>
		</section>
	</div>

	<!-- Create Announcement Modal -->
	<Modal title="创建新公告" :show="createModal.show" @close="createModal.show = false" @submit="submitCreate">
		<div class="form-group">
			<label for="create-title">标题</label>
			<input id="create-title" v-model="createModal.title" type="text" placeholder="输入公告标题" />
		</div>
		<div class="form-group">
			<label for="create-content">内容</label>
			<textarea id="create-content" v-model="createModal.content" placeholder="输入公告内容" rows="6"></textarea>
		</div>
	</Modal>

	<!-- Edit Announcement Modal -->
	<Modal title="编辑公告" :show="editModal.show" @close="editModal.show = false" @submit="submitEdit">
		<div class="form-group">
			<label for="edit-title">标题</label>
			<input id="edit-title" v-model="editModal.title" type="text" placeholder="输入公告标题" />
		</div>
		<div class="form-group">
			<label for="edit-content">内容</label>
			<textarea id="edit-content" v-model="editModal.content" placeholder="输入公告内容" rows="6"></textarea>
		</div>
	</Modal>
</template>

<style scoped src="./Overview.css"></style>
