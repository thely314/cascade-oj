<template>
	<div class="card">
		<h3>排名</h3>

		<p v-if="loading">加载中...</p>
		<p v-else-if="!loading && entries.length === 0">暂无排名数据</p>

		<ol class="rank-list" v-else>
			<li v-for="e in entries" :key="e.userId">{{ e.userId }} — {{ e.score }} 分</li>
		</ol>
	</div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'

const props = defineProps<{ contestId?: string }>()

type Entry = { userId: string; score: number }

const entries = ref<Entry[]>([])
const loading = ref(false)

// TODO: replace this stub with a real API call.
async function fetchRanking(contestId?: string) {
	loading.value = true
	// simulate network latency and return mock data for now
	await new Promise((res) => setTimeout(res, 150))

	// If contestId is provided, you could fetch specific data.
	// Here we return deterministic mock data so the UI can be built.
	entries.value = contestId
		? [
				{ userId: `${contestId}-user1`, score: 500 },
				{ userId: `${contestId}-user2`, score: 420 },
				{ userId: `${contestId}-user3`, score: 380 },
			]
		: [
				{ userId: 'user1', score: 500 },
				{ userId: 'user2', score: 420 },
				{ userId: 'user3', score: 380 },
			]

	loading.value = false
}

onMounted(() => {
	fetchRanking(props.contestId)
})

watch(
	() => props.contestId,
	(v) => {
		// refetch when contest changes
		fetchRanking(v)
	}
)
</script>

<style scoped>
@import "../home/home.css";

.rank-list {
	padding-left: 1.2rem;
	margin: 0.5rem 0 0 0;
}

.rank-list li {
	margin: 6px 0;
	color: #e6f3ef;
}
</style>
