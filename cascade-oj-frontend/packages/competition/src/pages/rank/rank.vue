<template>
		<div class="card card-blue info">
		<h3>排名</h3>

		<p v-if="loading">加载中...</p>
		<p v-else-if="!loading && entries.length === 0">暂无排名数据</p>

		<ol class="rank-list" v-else>
			<li v-for="e in entries" :key="e.userId">{{ (e as any).username || e.userId }} — {{ e.score }} 分</li>
		</ol>
	</div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { getRanks, type RankItem } from '../../api/rank'

const props = defineProps<{ contestId?: string }>()

const entries = ref<RankItem[]>([])
const loading = ref(false)
const error = ref<string | null>(null)

async function fetchRanking(contestId?: string) {
	if (!contestId) return
	loading.value = true
	try {
		const res = await getRanks(contestId)
		entries.value = res.ranks || []
	} catch (e: any) {
		error.value = e?.message ?? '加载排名失败'
	} finally {
		loading.value = false
	}
}

onMounted(() => {
	fetchRanking(props.contestId)
})

watch(
	() => props.contestId,
	(v) => {
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
	color: inherit;
}

/* 禁用本卡片的 hover 交互效果*/
.card.card-blue.info:hover {
	transform: none !important;
	box-shadow: 0 4px 6px rgba(0, 0, 0, 0.6) !important;
	cursor: default !important;
}
</style>
