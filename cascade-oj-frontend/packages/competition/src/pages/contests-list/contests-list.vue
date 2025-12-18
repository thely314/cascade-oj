
<template>
	<div class="home-container">
		<main class="main-content">
			<h1>正在进行的比赛</h1>
			<p>示例比赛卡片，点击进入对应比赛详情。</p>

			<div class="feature-cards">
				<!-- 保留一个静态示例比赛卡片 -->
				<router-link class="card card-link" to="/competition/1"><h3>示例比赛（静态）</h3><p>这是一个静态示例，后续可删除。</p></router-link>

				<!-- 动态比赛列表 -->
				<router-link v-for="c in contests" :key="c.id" class="card card-link" :to="`/competition/${c.id}`">
					<h3>{{ c.title }}</h3>
					<p>{{ c.status }} · {{ c.startTime }} — {{ c.endTime }}</p>
				</router-link>
			</div>

		</main>

	</div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getContests, type ContestMetadata } from '../../api/contests'

const contests = ref<ContestMetadata[]>([])
const loading = ref(false)
const error = ref<string | null>(null)

onMounted(async () => {
	loading.value = true
	try {
		const res = await getContests()
		contests.value = res.contests
	} catch (e: any) {
		error.value = e?.message ?? '加载比赛列表失败'
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
</style>

