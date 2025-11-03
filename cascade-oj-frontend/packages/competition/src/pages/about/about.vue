
<template>
	<div class="about-container">
		<NavBar />

		<main class="main-content">
			<h1 class="page-title">关于 CasCade</h1>
			<p class="lead">一个轻量、开源的在线题库与竞赛平台。以下内容为示例，方便你快速替换成项目的实际信息。</p>

			<div class="grid">
				<section class="card">
					<h3>简介</h3>
					<p>CasCade 是一个面向算法学习与竞赛训练的平台。它的目标是降低搭建和维护竞赛系统的门槛，提供清晰的题库管理、比赛编排和排行榜功能。</p>
				</section>

				<section class="card">
					<h3>使命与愿景</h3>
					<ul>
						<li>为教师和竞赛组织者提供简单可靠的比赛工具。</li>
						<li>为学习者提供高质量、有结构的练习路径。</li>
						<li>通过开源协作推动题库与评测工具的共享。</li>
					</ul>
				</section>

				<section class="card">
					<h3>团队（或贡献者）</h3>
					<ul class="team-list">
						<li><strong>姓名 A</strong> — 项目维护、架构设计</li>
						<li><strong>姓名 B</strong> — 前端与交互</li>
						<li><strong>姓名 C</strong> — 后端与评测器</li>
					</ul>
				</section>

				<section class="card">
					<h3>联系方式 & 仓库</h3>
					<p>电子邮件: <a href="mailto:your-email@example.com">your-email@example.com</a></p>
					<p>项目仓库: <a target="_blank" rel="noopener noreferrer" href="https://github.com/your/repo">https://github.com/your/repo</a></p>
				</section>

				<section class="card">
					<h3>贡献与许可证</h3>
					<p>欢迎贡献。请阅读仓库根目录的 <code>README.md</code> 以及贡献指南。许可示例：MIT / Apache-2.0。</p>
				</section>

				<section class="card">
					<h3>快速编辑提示</h3>
					<ol>
						<li>替换“简介”和“使命”中的占位文字为团队的真实描述。</li>
						<li>在“团队”部分加入头像与社交链接（可用 <code>&lt;img&gt;</code> 与 flex 布局）。</li>
						<li>在“联系方式”中列出邮件、群组和 issue 路径，便于用户反馈。</li>
					</ol>
				</section>
			</div>

		</main>
	</div>
</template>

<script setup>
import NavBar from '../../components/NavBar/NavBar.vue'
</script>

<style scoped src="./about.css"></style>
