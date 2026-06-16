<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import {
	deleteAnnouncement,
	getAnnouncements,
	getContests,
	getProblems,
	getRanks,
	getSubmissions,
	postAnnouncement,
	putAnnouncement,
} from '../../api/admin'
import type {
	Announcement,
	ContestMetadata,
	ProblemMetadata,
	RankItem,
	SubmissionMetadata,
} from '../../api/types'
import Modal from '../../components/Modal.vue'

const stats = reactive([
	{ label: '正在进行', value: '-' },
	{ label: '已结束', value: '-' },
	{ label: '题库题目数', value: '-' },
	{ label: '总提交数', value: '-' },
])

const announcements = ref<Announcement[]>([])
const loading = ref(false)

type ExportMode = 'submissions' | 'ranks'

const createModal = reactive({
	show: false,
	title: '',
	content: '',
})

const editModal = reactive({
	show: false,
	announcementId: 0,
	title: '',
	content: '',
})

const exportModal = reactive({
	show: false,
	loadingMeta: false,
	loadingData: false,
	mode: 'submissions' as ExportMode,
	error: '',
	contests: [] as ContestMetadata[],
	problems: [] as ProblemMetadata[],
	query: {
		contestId: '',
		problemId: '',
		startTime: '',
		endTime: '',
		traversalPageSize: 50,
		previewPageSize: 8,
	},
	page: 1,
	submissions: [] as SubmissionMetadata[],
	ranks: [] as RankItem[],
})

const selectedSubmissionIds = ref<string[]>([])
const selectedRankIds = ref<number[]>([])

const parseOptionalNumber = (value: string) => {
	const normalizedValue = value.trim()
	if (!normalizedValue) return undefined

	const parsedValue = Number(normalizedValue)
	return Number.isFinite(parsedValue) ? parsedValue : undefined
}

const formatDateTime = (dateStr: string) => {
	if (!dateStr) return ''
	return new Date(dateStr).toLocaleString('zh-CN')
}

const formatToIsoOrUndefined = (value: string) => {
	if (!value) return undefined
	const date = new Date(value)
	return Number.isNaN(date.getTime()) ? undefined : date.toISOString()
}

const isSubmissionItem = (item: SubmissionMetadata | RankItem): item is SubmissionMetadata => {
	return 'submissionUuid' in item
}

const getExportItemKey = (item: SubmissionMetadata | RankItem) => {
	return isSubmissionItem(item) ? item.submissionUuid : String(item.userId)
}

const getExportItemTitle = (item: SubmissionMetadata | RankItem) => {
	return isSubmissionItem(item)
		? `Submission ${item.submissionUuid}`
		: `${item.rank}. ${item.username}`
}

const getExportItemMeta = (item: SubmissionMetadata | RankItem) => {
	return isSubmissionItem(item)
		? `Problem ${item.problemId} · User ${item.userId} · ${formatDateTime(item.submitTime)}`
		: `User ID ${item.userId} · Rank ${item.rank}`
}

const getExportItemScore = (item: SubmissionMetadata | RankItem) => {
	return `Score ${item.score}`
}

const currentSubmissionExportItems = computed(() => {
	const start = (exportModal.page - 1) * exportModal.query.previewPageSize
	const end = start + exportModal.query.previewPageSize
	return exportModal.submissions.slice(start, end)
})

const currentRankExportItems = computed(() => {
	const start = (exportModal.page - 1) * exportModal.query.previewPageSize
	const end = start + exportModal.query.previewPageSize
	return exportModal.ranks.slice(start, end)
})

const currentExportItems = computed(() => {
	return exportModal.mode === 'submissions'
		? currentSubmissionExportItems.value
		: currentRankExportItems.value
})

const exportTotalItems = computed(() => {
	return exportModal.mode === 'submissions' ? exportModal.submissions.length : exportModal.ranks.length
})

const exportPageCount = computed(() => {
	return Math.max(1, Math.ceil(exportTotalItems.value / exportModal.query.previewPageSize))
})

const hasExportSelection = computed(() => {
	return exportModal.mode === 'submissions'
		? selectedSubmissionIds.value.length > 0
		: selectedRankIds.value.length > 0
})

const isCurrentPageSelected = computed(() => {
	if (currentExportItems.value.length === 0) return false
	if (exportModal.mode === 'submissions') {
		return currentExportItems.value.every((item) => isSubmissionItem(item) && selectedSubmissionIds.value.includes(item.submissionUuid))
	}
	return currentExportItems.value.every((item) => selectedRankIds.value.includes(item.userId))
})

const exportButtonLabel = computed(() => {
	return exportModal.mode === 'submissions' ? 'Export Selected Submissions' : 'Export Selected Ranks'
})

const exportPreviewTitle = computed(() => {
	return exportModal.mode === 'submissions' ? 'Submission Export Preview' : 'Contest Rank Export Preview'
})

const exportFileName = computed(() => {
	const stamp = new Date().toISOString().replaceAll(':', '-').replaceAll('.', '-')
	return exportModal.mode === 'submissions'
		? `submissions-export-${stamp}.json`
		: `contest-ranks-export-${stamp}.json`
})

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

const updateStats = (index: number, value: number) => {
	stats[index].value = String(value)
}

const fetchOverviewStats = async () => {
	const [contestReply, problemReply] = await Promise.all([getContests(), getProblems()])
	const contests = contestReply.contests || []
	const problems = problemReply.problems || []

	const ongoingCount = contests.filter((contest) => String(contest.status).toLowerCase() === 'ongoing').length
	const endedCount = contests.filter((contest) => String(contest.status).toLowerCase() === 'ended').length

	updateStats(0, ongoingCount)
	updateStats(1, endedCount)
	updateStats(2, problems.length)

	let totalSubmissions = 0
	const pageSize = 200

	for (let page = 1; page <= 200; page += 1) {
		const response = await getSubmissions({ page, pageSize })
		const batch = response.submissions || []
		totalSubmissions += batch.length

		if (batch.length < pageSize) {
			break
		}
	}

	updateStats(3, totalSubmissions)
}

const fetchRankExportData = async () => {
	const contestId = parseOptionalNumber(exportModal.query.contestId)
	if (!contestId) {
		exportModal.error = '请选择 contest'
		return
	}

	const resp = await getRanks(contestId)
	exportModal.ranks = resp.ranks || []
	selectedRankIds.value = exportModal.ranks.map((r) => r.userId)
}

const loadExportMetadata = async () => {
	if (exportModal.contests.length > 0 && exportModal.problems.length > 0) {
		return
	}

	exportModal.loadingMeta = true
	exportModal.error = ''
	try {
		const [contestReply, problemReply] = await Promise.all([getContests(), getProblems()])
		exportModal.contests = contestReply.contests || []
		exportModal.problems = problemReply.problems || []
	} catch (err) {
		console.error(err)
		exportModal.error = '加载筛选数据失败'
	} finally {
		exportModal.loadingMeta = false
	}
}

const resetExportState = (mode: ExportMode = 'submissions') => {
	exportModal.mode = mode
	exportModal.loadingData = false
	exportModal.error = ''
	exportModal.page = 1
	exportModal.query = {
		contestId: '',
		problemId: '',
		startTime: '',
		endTime: '',
		traversalPageSize: 50,
		previewPageSize: 8,
	}
	exportModal.submissions = []
	exportModal.ranks = []
	selectedSubmissionIds.value = []
	selectedRankIds.value = []
}

const openExportModal = async (mode: ExportMode = 'submissions') => {
	resetExportState(mode)
	exportModal.show = true
	await loadExportMetadata()
}

const closeExportModal = () => {
	exportModal.show = false
}

const selectExportTab = (mode: ExportMode) => {
	if (exportModal.mode === mode) return
	resetExportState(mode)
	exportModal.show = true
}

const fetchSubmissionExportData = async () => {
	const contestId = parseOptionalNumber(exportModal.query.contestId)
	const problemId = parseOptionalNumber(exportModal.query.problemId)
	const timeStart = formatToIsoOrUndefined(exportModal.query.startTime)
	const timeEnd = formatToIsoOrUndefined(exportModal.query.endTime)
	const pageSize = Math.max(1, exportModal.query.traversalPageSize)
	const collected: SubmissionMetadata[] = []

	for (let page = 1; page <= 200; page += 1) {
		const response = await getSubmissions({
			contestId,
			problemId,
			page,
			pageSize,
		})

		const batch = response.submissions || []
		if (batch.length === 0) {
			break
		}

		const filtered = batch.filter((item) => {
			const submittedAt = new Date(item.submitTime).getTime()
			if (timeStart && submittedAt < new Date(timeStart).getTime()) return false
			if (timeEnd && submittedAt > new Date(timeEnd).getTime()) return false
			return true
		})

		collected.push(...filtered)

		if (batch.length < pageSize) {
			break
		}
	}

	exportModal.submissions = collected
	selectedSubmissionIds.value = collected.map((item) => item.submissionUuid)
}


const refreshExportData = async () => {
	exportModal.loadingData = true
	exportModal.error = ''
	exportModal.page = 1
	selectedSubmissionIds.value = []
	selectedRankIds.value = []

	try {
		if (exportModal.mode === 'submissions') {
			await fetchSubmissionExportData()
		} else {
			await fetchRankExportData()
		}
	} catch (err) {
		console.error(err)
		exportModal.error = err instanceof Error ? err.message : '加载导出数据失败'
	} finally {
		exportModal.loadingData = false
	}
}

const setCurrentPageSelection = (checked: boolean) => {
	if (exportModal.mode === 'submissions') {
		const next = new Set(selectedSubmissionIds.value)
		currentExportItems.value.forEach((item) => {
			if (!isSubmissionItem(item)) {
				return
			}
			if (checked) {
				next.add(item.submissionUuid)
			} else {
				next.delete(item.submissionUuid)
			}
		})
		selectedSubmissionIds.value = Array.from(next)
		return
	}

	const next = new Set(selectedRankIds.value)
	currentExportItems.value.forEach((item) => {
		if (checked) {
			next.add(item.userId)
		} else {
			next.delete(item.userId)
		}
	})
	selectedRankIds.value = Array.from(next)
}

const isItemSelected = (item: SubmissionMetadata | RankItem) => {
	return exportModal.mode === 'submissions'
		? isSubmissionItem(item) && selectedSubmissionIds.value.includes(item.submissionUuid)
		: selectedRankIds.value.includes(item.userId)
}

const toggleExportItem = (item: SubmissionMetadata | RankItem, checked: boolean) => {
	if (exportModal.mode === 'submissions') {
		if (!isSubmissionItem(item)) {
			return
		}
		const key = item.submissionUuid
		selectedSubmissionIds.value = checked
			? Array.from(new Set([...selectedSubmissionIds.value, key]))
			: selectedSubmissionIds.value.filter((value) => value !== key)
		return
	}

	const key = item.userId
	selectedRankIds.value = checked
		? Array.from(new Set([...selectedRankIds.value, key]))
		: selectedRankIds.value.filter((value) => value !== key)
}

const downloadExportFile = (filename: string, payload: unknown) => {
	const blob = new Blob([JSON.stringify(payload, null, 2)], { type: 'application/json;charset=utf-8' })
	const url = URL.createObjectURL(blob)
	const link = document.createElement('a')
	link.href = url
	link.download = filename
	document.body.appendChild(link)
	link.click()
	link.remove()
	URL.revokeObjectURL(url)
}

const groupItemsByPage = <T>(items: T[], pageSize: number) => {
	const normalizedPageSize = Math.max(1, pageSize)
	const pages: Array<{ page: number; items: T[] }> = []

	for (let index = 0; index < items.length; index += normalizedPageSize) {
		pages.push({
			page: pages.length + 1,
			items: items.slice(index, index + normalizedPageSize),
		})
	}

	return pages
}

const submitExport = () => {
	if (!hasExportSelection.value) {
		globalThis.alert('请先选择要导出的数据')
		return
	}

	if (exportModal.mode === 'submissions') {
		const items = exportModal.submissions.filter((item) => selectedSubmissionIds.value.includes(item.submissionUuid))
		const pageSize = Math.max(1, exportModal.query.traversalPageSize)
		const pages = groupItemsByPage(items, pageSize)
		downloadExportFile(exportFileName.value, {
			type: 'submissions',
			exportedAt: new Date().toISOString(),
			filters: {
				contestId: parseOptionalNumber(exportModal.query.contestId),
				problemId: parseOptionalNumber(exportModal.query.problemId),
				startTime: formatToIsoOrUndefined(exportModal.query.startTime),
				endTime: formatToIsoOrUndefined(exportModal.query.endTime),
			},
			paging: {
				pageSize,
				totalItems: items.length,
				totalPages: pages.length,
			},
			pages,
		})
		exportModal.show = false
		return
	}

	const items = exportModal.ranks.filter((item) => selectedRankIds.value.includes(item.userId))
	const pageSize = Math.max(1, exportModal.query.traversalPageSize)
	const pages = groupItemsByPage(items, pageSize)
	downloadExportFile(exportFileName.value, {
		type: 'contest-ranks',
		exportedAt: new Date().toISOString(),
		filters: {
			contestId: parseOptionalNumber(exportModal.query.contestId),
		},
		paging: {
			pageSize,
			totalItems: items.length,
			totalPages: pages.length,
		},
		pages,
	})
	exportModal.show = false
}

const onNewAnnouncement = () => {
	createModal.title = ''
	createModal.content = ''
	createModal.show = true
}

const submitCreate = async () => {
	if (!createModal.title.trim() || !createModal.content.trim()) {
		globalThis.alert('请填写完整信息')
		return
	}
	try {
		await postAnnouncement({
			title: createModal.title.trim(),
			content: createModal.content.trim(),
		})
		createModal.show = false
		await fetchAnnouncements()
	} catch (err) {
		console.error(err)
		globalThis.alert('创建公告失败')
	}
}

const onOpenAnnouncement = (announcement: Announcement) => {
	editModal.announcementId = announcement.id
	editModal.title = announcement.title
	editModal.content = announcement.content
	editModal.show = true
}

const submitEdit = async () => {
	if (!editModal.title.trim() || !editModal.content.trim()) {
		globalThis.alert('请填写完整信息')
		return
	}
	try {
		await putAnnouncement(editModal.announcementId, {
			title: editModal.title.trim(),
			content: editModal.content.trim(),
		})
		editModal.show = false
		await fetchAnnouncements()
	} catch (err) {
		console.error(err)
		globalThis.alert('修改公告失败')
	}
}

const onDeleteAnnouncement = async (announcement: Announcement) => {
	if (!globalThis.confirm(`确认删除公告 "${announcement.title}"？`)) return
	try {
		await deleteAnnouncement(announcement.id)
		await fetchAnnouncements()
	} catch (err) {
		console.error(err)
		globalThis.alert('删除公告失败')
	}
}

onMounted(() => {
	void Promise.all([fetchAnnouncements(), fetchOverviewStats()])
})
</script>

<template>
	<div class="page">
		<header class="page-header">
			<div>
				<p class="eyebrow">Admin</p>
				<h1 class="title">Overview</h1>
			</div>
			<div class="actions">
				<button class="ghost" @click="openExportModal('submissions')">Export Data</button>
				<button class="primary" @click="onNewAnnouncement">New Announcement</button>
			</div>
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
					<div class="list-actions">
						<button class="ghost" @click="onOpenAnnouncement(item)">Edit</button>
						<button class="ghost danger" @click="onDeleteAnnouncement(item)">Delete</button>
					</div>
				</li>
			</ul>
		</section>
	</div>

	<Modal
		title="Export Data"
		:show="exportModal.show"
		cancel-text="Close"
		:submit-text="exportButtonLabel"
		:submit-disabled="exportModal.loadingMeta || exportModal.loadingData || !hasExportSelection"
		@close="closeExportModal"
		@submit="submitExport"
	>
		<div class="export-modal">
			<div class="export-tabs">
				<button class="tab" :class="{ active: exportModal.mode === 'submissions' }" type="button" @click="selectExportTab('submissions')">
					Submissions
				</button>
				<button class="tab" :class="{ active: exportModal.mode === 'ranks' }" type="button" @click="selectExportTab('ranks')">
					Contest Ranks
				</button>
			</div>

			<div class="export-grid">
				<label class="form-group">
					<span>Contest</span>
					<select v-model="exportModal.query.contestId">
						<option value="">All contests</option>
						<option v-for="contest in exportModal.contests" :key="contest.id" :value="String(contest.id)">
							{{ contest.id }} · {{ contest.title }}
						</option>
					</select>
				</label>

				<label v-if="exportModal.mode === 'submissions'" class="form-group">
					<span>Problem</span>
					<select v-model="exportModal.query.problemId">
						<option value="">All problems</option>
						<option v-for="problem in exportModal.problems" :key="problem.id" :value="String(problem.id)">
							{{ problem.id }} · {{ problem.title }}
						</option>
					</select>
				</label>

				<label v-if="exportModal.mode === 'submissions'" class="form-group">
					<span>Start Time</span>
					<input v-model="exportModal.query.startTime" type="datetime-local" />
				</label>

				<label v-if="exportModal.mode === 'submissions'" class="form-group">
					<span>End Time</span>
					<input v-model="exportModal.query.endTime" type="datetime-local" />
				</label>

				<label class="form-group">
					<span>Traversal Page Size</span>
					<input v-model.number="exportModal.query.traversalPageSize" type="number" min="1" max="500" />
				</label>

				<label class="form-group">
					<span>Preview Page Size</span>
					<input v-model.number="exportModal.query.previewPageSize" type="number" min="1" max="50" />
				</label>
			</div>

			<div class="export-actions">
				<button class="ghost" type="button" @click="refreshExportData" :disabled="exportModal.loadingMeta || exportModal.loadingData">
					{{ exportModal.loadingData ? 'Loading...' : 'Load Matching Data' }}
				</button>
				<p class="export-hint">
					{{ exportModal.mode === 'submissions' ? '按 contest / problem / 时间范围遍历提交记录并导出选中项。' : '按 contest 加载排名并按页选择导出。' }}
				</p>
			</div>

			<div v-if="exportModal.error" class="error-message">{{ exportModal.error }}</div>
			<div v-else-if="exportModal.loadingMeta || exportModal.loadingData" class="empty-state">Loading export data...</div>
			<template v-else>
				<div class="export-summary">
					<span>{{ exportPreviewTitle }}</span>
					<span>{{ exportTotalItems }} records</span>
				</div>

				<div class="export-select-all">
					<label>
						<input type="checkbox" :checked="isCurrentPageSelected" @change="setCurrentPageSelection(($event.target as HTMLInputElement).checked)" />
						Select current page
					</label>
				</div>

				<ul class="export-list">
					<li v-for="item in currentExportItems" :key="getExportItemKey(item)" class="export-item">
						<label class="export-item-main">
							<input type="checkbox" :checked="isItemSelected(item)" @change="toggleExportItem(item, ($event.target as HTMLInputElement).checked)" />
							<div>
								<p class="export-item-title">{{ getExportItemTitle(item) }}</p>
								<p class="export-item-meta">{{ getExportItemMeta(item) }}</p>
							</div>
						</label>
						<span class="export-item-score">{{ getExportItemScore(item) }}</span>
					</li>
				</ul>

				<div class="export-pagination">
					<button class="ghost" type="button" :disabled="exportModal.page <= 1" @click="exportModal.page -= 1">Previous</button>
					<span>Page {{ exportModal.page }} / {{ exportPageCount }}</span>
					<button class="ghost" type="button" :disabled="exportModal.page >= exportPageCount" @click="exportModal.page += 1">Next</button>
				</div>
			</template>
		</div>
	</Modal>

	<!-- Create Announcement Modal -->
	<Modal title="创建新公告" :show="createModal.show" @close="createModal.show = false" @submit="submitCreate">
		<div class="form-group">
			<label for="create-title">标题</label>
			<input id="create-title" v-model="createModal.title" type="text" placeholder="输入公告标题" />
		</div>
		<div class="form-group">
			<label for="create-content">内容</label>
			<textarea id="create-content" v-model="createModal.content" placeholder="输入公告内容" rows="6"></textarea>
		</div>
	</Modal>

	<!-- Edit Announcement Modal -->
	<Modal title="编辑公告" :show="editModal.show" @close="editModal.show = false" @submit="submitEdit">
		<div class="form-group">
			<label for="edit-title">标题</label>
			<input id="edit-title" v-model="editModal.title" type="text" placeholder="输入公告标题" />
		</div>
		<div class="form-group">
			<label for="edit-content">内容</label>
			<textarea id="edit-content" v-model="editModal.content" placeholder="输入公告内容" rows="6"></textarea>
		</div>
	</Modal>
</template>

<style scoped src="./Overview.css"></style>
