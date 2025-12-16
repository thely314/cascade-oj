<script setup>
const stats = [
	{ label: 'Active Users', value: '1,248' },
	{ label: 'Problems', value: '362' },
	{ label: 'Submissions Today', value: '5,421' },
	{ label: 'Pending Reviews', value: '18' },
]

const recent = [
	{ title: 'System maintenance window', time: 'Today 14:00-15:00' },
	{ title: 'New contest onboarding', time: 'Tomorrow 09:00' },
	{ title: 'Weekly log review', time: 'Fri 10:00' },
]
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

		<section class="stat-grid">
			<article v-for="item in stats" :key="item.label" class="stat-card">
				<p class="stat-label">{{ item.label }}</p>
				<p class="stat-value">{{ item.value }}</p>
			</article>
		</section>

		<section class="panel">
			<header class="panel-header">
				<h2>Recent Items</h2>
				<a href="#">View all</a>
			</header>
			<ul class="list">
				<li v-for="item in recent" :key="item.title" class="list-item">
					<div>
						<p class="list-title">{{ item.title }}</p>
						<p class="list-meta">{{ item.time }}</p>
					</div>
					<button class="ghost">Open</button>
				</li>
			</ul>
		</section>
	</div>
</template>

<style scoped src="./Overview.css"></style>
