<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getAnnouncements } from '../../api/admin'
import type { Announcement } from '../../api/types'

const stats = [
	{ label: '正在进行', value: '-' },
	{ label: '已结束', value: '-' },
	{ label: '题库题目数', value: '-' },
	{ label: '总提交数', value: '-' },
]

const announcements = ref<Announcement[]>([])
const loading = ref(false)

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
</script>

<template>
	<div class="page">
		<header class="page-header">
			<div>
				<p class="eyebrow">Admin</p>
				<h1 class="title">Overview</h1>
			</div>
			<button class="primary">New Announcement</button>
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
					<button class="ghost">Open</button>
				</li>
			</ul>
		</section>
	</div>
</template>

<style scoped src="./Overview.css"></style>
