
<template>
	<div class="home-container">
		<main class="main-content">
			<h1>比赛详情</h1>
			<p>当前比赛 ID：{{ idStr }}</p>

			<div class="contest-grid">
				<div class="left">
					<div class="problem-list">
							<template v-if="loading">
								<div class="card problem-item" v-for="i in 3" :key="i" style="opacity:.6">
									<h3>加载中...</h3>
									<p>正在获取题目列表</p>
								</div>
							</template>
							<p v-else-if="error" class="error-text">{{ error }}</p>
							<p v-else-if="!problems.length" class="empty-text">暂无题目</p>
							<router-link v-else v-for="p in problems" :key="p.id" class="card problem-item card-link" :to="`/contest/${idStr.valueOf()}/problem/${p.id}`">
								<h3>{{ p.title }}</h3>
								<p>时间限制：{{ p.timeLimitMs }}ms · 内存限制：{{ p.memoryLimitMb }}MB</p>
						</router-link>
					</div>
				</div>

				<div class="right">
					<div class="side-list">
						<div class="card join-card">
							<h3>加入比赛</h3>
							<p v-if="joined">已加入该比赛，祝你取得好成绩！</p>
							<p v-else>成功加入比赛后，即可参与排名与提交。</p>
							<div class="actions">
								<button v-if="!joined" class="action-btn" :disabled="joinLoading" @click="handleJoin">
									{{ joinLoading ? '加入中...' : '加入比赛' }}
								</button>
								<button v-else class="action-btn secondary" :disabled="quitLoading" @click="handleQuit">
									{{ quitLoading ? '退出中...' : '退出比赛' }}
								</button>
							</div>
							<p v-if="joinError" class="error-text">{{ joinError }}</p>
						</div>

						<Rank :contestId="idStr.valueOf()" />
					</div>
				</div>
			</div>

		</main>
	</div>
</template>

<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'
import { computed, ref, onMounted } from 'vue'
import Rank from '../rank/rank.vue'
import { getContest, getContestProblems, joinContest, quitContest, getJoinStatus, type GetSingleContestReply, type ProblemMetadata } from '../../api/contest'

const route = useRoute()
const id = computed(() => String(route.params.id ?? '未知'))
// 字符串形式，用于在表达式中安全构建 URL
const idStr = computed(() => String(id.value))

const contest = ref<GetSingleContestReply | null>(null)
const problems = ref<ProblemMetadata[]>([])
const loading = ref(false)
const error = ref<string | null>(null)

// 加入比赛相关状态
const joined = ref(false)
const joinLoading = ref(false)
const quitLoading = ref(false)
const joinError = ref<string | null>(null)
const router = useRouter()

async function handleJoin() {
	joinError.value = null
	joinLoading.value = true
	try {
		// 若未登录，引导到登录页
		const token = localStorage.getItem('cascade_token')
		if (!token) {
			joinLoading.value = false
			router.push('/login')
			return
		}

		const res = await joinContest(idStr.value)
		joined.value = Boolean(res?.isJoin)
		if (joined.value) {
			localStorage.setItem(`cascade_joined_${idStr.value}`, '1')
		}
		if (!joined.value) {
			joinError.value = '加入失败，请稍后重试'
		}
	} catch (e: any) {
		joinError.value = e?.message ?? '加入比赛失败'
	} finally {
		joinLoading.value = false
	}
}

async function handleQuit() {
	joinError.value = null
	quitLoading.value = true
	try {
		const res = await quitContest(idStr.value)
		joined.value = Boolean(res?.isJoin)
		if (!joined.value) {
			localStorage.removeItem(`cascade_joined_${idStr.value}`)
		}
	} catch (e: any) {
		joinError.value = e?.message ?? '退出比赛失败'
	} finally {
		quitLoading.value = false
	}
}

onMounted(async () => {
	loading.value = true
	try {
		const [contestRes, problemsRes] = await Promise.all([
			getContest(idStr.value),
			getContestProblems(idStr.value),
		])
		contest.value = contestRes
			problems.value = problemsRes.problems || []
		// 先尝试服务端查询加入状态（占位接口），失败则回退到本地存储
		try {
			const status = await getJoinStatus(idStr.value)
			joined.value = Boolean(status?.isJoin)
			if (joined.value) {
				localStorage.setItem(`cascade_joined_${idStr.value}`, '1')
			} else {
				localStorage.removeItem(`cascade_joined_${idStr.value}`)
			}
		} catch {
			joined.value = localStorage.getItem(`cascade_joined_${idStr.value}`) === '1'
		}
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
.side-list { display: flex; flex-direction: column; gap: 12px; }
.actions { margin-top: 12px; display: flex; gap: 8px; justify-content: center; }
.action-btn { padding: 8px 16px; border-radius: 6px; border: 1px solid rgba(30, 134, 68, 0.2); background: #1dad80; color: #0b0f0d; cursor: pointer; }
.action-btn.secondary { background: #141816; color: #cbd5c0; }
.action-btn:disabled { opacity: 0.7; cursor: not-allowed; }
.error-text { color: #e57373; margin-top: 8px; }
.problem-list { display: flex; flex-direction: column; gap: 12px; }
.problem-item { text-align: left; }

@media (max-width: 768px) {
	.contest-grid { flex-direction: column; }
	.right { min-width: auto; }
}
</style>

