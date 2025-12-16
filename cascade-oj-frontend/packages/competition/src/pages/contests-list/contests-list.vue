
<template>
	<div class="home-container">
		<main class="main-content">
			<h1>正在进行的比赛</h1>
			<p>示例比赛卡片，点击进入对应比赛详情。</p>

			<div class="feature-cards">
				<router-link class="card card-link" to="/competition/1"><h3>示例比赛 1</h3><p>为期 2 小时，涵盖算法基础。</p></router-link>
				<router-link class="card card-link" to="/competition/2"><h3>示例比赛 2</h3><p>面向进阶题目，含数据结构。</p></router-link>
				<router-link class="card card-link" to="/competition/3"><h3>示例比赛 3</h3><p>团队赛（占位），协作解题。</p></router-link>

				<router-link class="card card-link" to="/competition/4"><h3>示例比赛 4</h3><p>周末练习赛，欢迎参与。</p></router-link>
				<router-link class="card card-link" to="/competition/5"><h3>示例比赛 5</h3><p>入门题目集合，适合新手。</p></router-link>
				<router-link class="card card-link" to="/competition/6"><h3>示例比赛 6</h3><p>挑战赛，限时高难度题。</p></router-link>

				<router-link class="card card-link" to="/competition/7"><h3>示例比赛 7</h3><p>趣味题目，轻松练习。</p></router-link>
				<router-link class="card card-link" to="/competition/8"><h3>示例比赛 8</h3><p>月度排行榜赛（占位）。</p></router-link>
				<router-link class="card card-link" to="/competition/9"><h3>示例比赛 9</h3><p>专题训练：图论与搜索。</p></router-link>
			</div>

		</main>

	</div>
</template>

<script setup lang="ts">
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

