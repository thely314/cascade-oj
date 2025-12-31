
<template>
	<div class="about-container">
		<main class="main-content">
			<h1 class="page-title">关于 Cascade</h1>
			<p class="lead">Cascade 意为级联、连续传递，表达知识的持续传递与共享之意。</p>

			<div class="grid">
				<section class="card">
					<h3>简介</h3>
					<p>Cascade 是一个面向语言学习与算法学习的在线测评平台，基于 Kratos 微服务架构。提供题库管理、比赛编排和排行榜功能。</p>
				</section>

				<section class="card">
					<h3>使命与愿景</h3>
					<ul>
						<li>用户友好的界面与体验。</li>
						<li>通过开源协作推动题库与评测工具的共享。</li>
						<li>作为可供参考学习的开源项目，欢迎有兴趣的开发者、学习者参与。</li>
					</ul>
				</section>

				<section class="card">
					<h3>联系方式 & 仓库</h3>
					<!-- <p>电子邮件: <a href="mailto:email">暂不提供</a></p> -->
					<p>项目仓库: <a target="_blank" rel="noopener noreferrer" href="https://github.com/thely314/cascade-oj">https://github.com/thely314/cascade-oj</a></p>
					<p>使用 go-judge 测评机: <a target="_blank" rel="noopener noreferrer" href="https://github.com/criyle/go-judge">https://github.com/criyle/go-judge</a></p>
					<p>架构参考 SASTOJ: <a target="_blank" rel="noopener noreferrer" href="https://github.com/NJUPT-SAST/sastoj">https://github.com/NJUPT-SAST/sastoj</a></p>
				</section>

				<section class="card">
					<h3>贡献与许可证</h3>
					<p>欢迎贡献。请阅读仓库的 <code>README.md</code> 以及贡献指南。</p>
					<br></br>
					<p>使用 MIT License。</p>
				</section>
			</div>

		</main>
	</div>
</template>

<script setup lang="ts">
// NavBar is provided globally by App.vue
</script>

<style scoped src="./about.css"></style>
