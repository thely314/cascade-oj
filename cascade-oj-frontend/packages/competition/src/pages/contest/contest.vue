
<template>
	<div class="home-container">
		<main class="main-content">
			<h1>比赛详情</h1>
			<p>当前比赛 ID：{{ idStr }}</p>

			<div class="contest-grid">
				<div class="left">
					<div class="problem-list">
							<router-link v-for="p in problems" :key="p.id" class="card problem-item card-link" :to="`/competition/${idStr.value}/problem/${p.id}`">
								<h3>{{ p.title }}</h3>
								<p>时间限制：{{ p.timeLimitMs }}ms · 内存限制：{{ p.memoryLimitMb }}MB</p>
						</router-link>
					</div>
				</div>

				<div class="right">
					<Rank :contestId="idStr.value" />
				</div>
			</div>

		</main>
	</div>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router'
import { computed, ref, onMounted } from 'vue'
import Rank from '../rank/rank.vue'
import { getContest, getContestProblems, type GetSingleContestReply, type ProblemMetadata } from '../../api/contest'

const route = useRoute()
const id = computed(() => String(route.params.id ?? '未知'))
// 字符串形式，用于在表达式中安全构建 URL
const idStr = computed(() => String(id.value))

const contest = ref<GetSingleContestReply | null>(null)
// 预置一个静态示例题目，后续可删除
const problems = ref<ProblemMetadata[]>([
	{ id: 'example', title: '示例题目（静态）', timeLimitMs: 1000, memoryLimitMb: 256 }
])
const loading = ref(false)
const error = ref<string | null>(null)

onMounted(async () => {
	loading.value = true
	try {
		const [contestRes, problemsRes] = await Promise.all([
			getContest(idStr.value),
			getContestProblems(idStr.value),
		])
		contest.value = contestRes
		problems.value = [problems.value[0], ...problemsRes.problems]
	} catch (e: any) {
		error.value = e?.message ?? '加载比赛信息失败'
	} finally {
		loading.value = false
	}
})
</script>

<style scoped>
.home-container {
	min-height: 100vh;
	display: flex;
	flex-direction: column;
	background: #202020;
	/* 页面背景：稍微调亮的深色 */
	/* 页面背景由纯黑(#000)调亮为 #0d1310 */
	color: #cbd5c0;
	/* 全局文字：浅灰/米色，便于黑底阅读 */
}

.main-content {
	flex: 1;
	padding: 40px 20px;
	text-align: center;
	max-width: 1200px;
	margin: 0 auto;
	width: 100%;
}

h1 {
	color: #1dad80;
	/* 主色：绿色 */
	margin-bottom: 40px;
	font-size: 2.5rem;
}

p {
	color: #ced7db;
	/* 次要文字：浅灰 */
	font-size: 1.2rem;
	margin-bottom: 60px;
}

.feature-cards {
	display: grid;
	grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
	gap: 40px;
	margin-top: 50px;
}

.card {
	background: #141816;
	/* 卡片背景：深色，和页面黑色区分 */
	padding: 30px 20px;
	border-radius: 8px;
	box-shadow: 0 4px 6px rgba(0, 0, 0, 0.6);
	border: 1px solid rgba(30, 134, 68, 0.06);
	/* 细微绿边，增强层次 */
	transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.card:hover {
	transform: translateY(-5px);
	box-shadow: 0 6px 18px rgba(22, 163, 142, 0.12);
	/* hover 带绿色光晕 */
}

.card h3 {
	color: #1dad80;
	/* 卡片标题使用绿色 */
	margin-bottom: 15px;
	font-size: 1.5rem;
}

.card p {
	color: #98cdda;
	margin: 0;
	font-size: 1rem;
}

/* 去除卡片中链接的默认下划线，并让链接继承卡片样式 */
.card-link {
	color: inherit;
	text-decoration: none;
	display: block;
	height: 100%;
}

.card-link:focus {
	outline: none;
}

.card-link h3, .card-link p {
	margin: 0;
}

.contest-grid {
	display: flex;
	gap: 24px;
	align-items: flex-start;
	margin-top: 24px;
}
.left { flex: 2; }
.right { flex: 1; min-width: 260px; }
.problem-list { display: flex; flex-direction: column; gap: 12px; }
.problem-item { text-align: left; }

@media (max-width: 768px) {
	.contest-grid { flex-direction: column; }
	.right { min-width: auto; }
}
</style>

