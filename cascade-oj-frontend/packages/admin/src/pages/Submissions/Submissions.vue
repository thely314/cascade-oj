<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import Modal from '../../components/Modal.vue'
import { getSingleSubmission, getSubmissions, rejudgeSubmission } from '../../api/admin'
import type { GetSingleSubmissionReply, SubmissionMetadata } from '../../api/types'

const submissions = ref<SubmissionMetadata[]>([])
const loading = ref(false)
const error = ref('')
const detailModal = reactive({
  show: false,
  loading: false,
  rejudging: false,
  error: '',
  detail: null as GetSingleSubmissionReply | null,
})

const rejudgeState = reactive({
  submissionUuid: '',
  loading: false,
})

const filters = reactive({
  problemId: '',
  contestId: '',
  userId: '',
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
})

const pageSizeOptions = [10, 20, 50]

const statusTone: Record<string, string> = {
  Accepted: 'badge-live',
  Pending: 'badge-muted',
  'Wrong Answer': 'badge-dim',
}

const parseOptionalNumber = (value: string) => {
  const normalizedValue = value.trim()
  if (!normalizedValue) return undefined

  const parsedValue = Number(normalizedValue)
  return Number.isFinite(parsedValue) ? parsedValue : undefined
}

const hasNextPage = computed(() => submissions.value.length === pagination.pageSize)

const requestParams = computed(() => ({
  problemId: parseOptionalNumber(filters.problemId),
  contestId: parseOptionalNumber(filters.contestId),
  userId: parseOptionalNumber(filters.userId),
  page: pagination.page,
  pageSize: pagination.pageSize,
}))

const fetchSubmissions = async () => {
  loading.value = true
  error.value = ''
  try {
    const response = await getSubmissions(requestParams.value)
    submissions.value = response.submissions || []
  } catch (err) {
    error.value = 'Failed to load submissions'
    console.error(err)
  } finally {
    loading.value = false
  }
}

const openSubmissionDetail = async (submission: SubmissionMetadata) => {
  detailModal.show = true
  detailModal.loading = true
  detailModal.error = ''
  detailModal.detail = null

  try {
    detailModal.detail = await getSingleSubmission(submission.submissionUuid)
  } catch (err) {
    detailModal.error = 'Failed to load submission detail'
    console.error(err)
  } finally {
    detailModal.loading = false
  }
}

const submitRowRejudge = async (submission: SubmissionMetadata) => {
  if (rejudgeState.loading) return
  if (!globalThis.confirm(`确认重新测评提交 ${submission.submissionUuid}？原提交将保留，新结果会生成一条新的提交记录。`)) return

  rejudgeState.submissionUuid = submission.submissionUuid
  rejudgeState.loading = true

  try {
    const resp = await rejudgeSubmission(submission.submissionUuid)
    await fetchSubmissions()
    const newUuid = resp?.newSubmissionUuid || ''
    globalThis.alert(newUuid ? `已提交重新测评，新提交 UUID: ${newUuid}` : '已提交重新测评')
  } catch (err) {
    globalThis.alert('重新测评失败')
    console.error(err)
  } finally {
    rejudgeState.submissionUuid = ''
    rejudgeState.loading = false
  }
}

const closeDetailModal = () => {
  detailModal.show = false
  detailModal.loading = false
  detailModal.rejudging = false
  detailModal.error = ''
  detailModal.detail = null
}

const submitRejudge = async () => {
  if (!detailModal.detail || detailModal.loading || detailModal.rejudging) return

  detailModal.rejudging = true
  detailModal.error = ''

  try {
    const resp = await rejudgeSubmission(detailModal.detail.metadata.submissionUuid)
    closeDetailModal()
    await fetchSubmissions()
    const newUuid = (resp && (resp as any).newSubmissionUuid) || ''
    globalThis.alert(newUuid ? `Rejudge requested successfully: ${newUuid}` : 'Rejudge requested successfully')
  } catch (err) {
    detailModal.error = 'Failed to rejudge submission'
    console.error(err)
  } finally {
    detailModal.rejudging = false
  }
}

const applyFilters = async () => {
  pagination.page = 1
  await fetchSubmissions()
}

const resetFilters = async () => {
  filters.problemId = ''
  filters.contestId = ''
  filters.userId = ''
  pagination.page = 1
  await fetchSubmissions()
}

const goPrevPage = async () => {
  if (pagination.page <= 1 || loading.value) return
  pagination.page -= 1
  await fetchSubmissions()
}

const goNextPage = async () => {
  if (!hasNextPage.value || loading.value) return
  pagination.page += 1
  await fetchSubmissions()
}

const handlePageSizeChange = async (event: Event) => {
  const selectElement = event.target as HTMLSelectElement
  const nextPageSize = Number(selectElement.value)
  if (!Number.isFinite(nextPageSize)) return

  pagination.pageSize = nextPageSize
  pagination.page = 1
  await fetchSubmissions()
}

onMounted(() => {
  fetchSubmissions()
})

const formatDate = (dateStr: string) => {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleString()
}
</script>

<template>
  <div class="page">
    <header class="page-header">
      <div>
        <p class="eyebrow">Admin</p>
        <h1 class="title">Submissions</h1>
      </div>
      <div class="actions">
        <button class="primary" @click="fetchSubmissions" :disabled="loading">
          {{ loading ? 'Loading...' : 'Refresh' }}
        </button>
      </div>
    </header>

    <section class="panel">
      <header class="panel-header">
        <div>
          <h2>Latest Submissions</h2>
          <p class="panel-desc">按题目、比赛和用户过滤提交记录，并按页查看最近结果。</p>
        </div>
        <a href="#" @click.prevent="resetFilters">Reset filters</a>
      </header>

      <form class="filters" @submit.prevent="applyFilters">
        <label class="filter-field">
          <span>Problem ID</span>
          <input v-model="filters.problemId" type="number" min="1" placeholder="Optional" />
        </label>
        <label class="filter-field">
          <span>Contest ID</span>
          <input v-model="filters.contestId" type="number" min="1" placeholder="Optional" />
        </label>
        <label class="filter-field">
          <span>User ID</span>
          <input v-model="filters.userId" type="number" min="1" placeholder="Optional" />
        </label>
        <label class="filter-field page-size-field">
          <span>Page Size</span>
          <select :value="pagination.pageSize" @change="handlePageSizeChange">
            <option v-for="size in pageSizeOptions" :key="size" :value="size">{{ size }}</option>
          </select>
        </label>
        <div class="filter-actions">
          <button type="button" class="ghost" @click="resetFilters" :disabled="loading">Clear</button>
          <button type="submit" class="primary" :disabled="loading">Search</button>
        </div>
      </form>

      <div v-if="error" class="error-message">{{ error }}</div>
      <div v-else-if="loading" class="empty-state">Loading submissions...</div>
      <div v-else-if="submissions.length === 0" class="empty-state">No submissions found for the current filters.</div>
      <template v-else>
      <ul class="list">
        <li v-for="sub in submissions" :key="sub.submissionUuid" class="list-item">
          <div class="list-main">
            <p class="list-title">Submission {{ sub.submissionUuid }}</p>
            <p class="list-meta">Problem {{ sub.problemId }} · User {{ sub.userId }} · {{ formatDate(sub.submitTime) }}</p>
          </div>
          <div class="list-right">
            <span class="badge" :class="statusTone[sub.status] || 'badge-muted'">{{ sub.status }}</span>
            <span class="score">{{ sub.score }}</span>
            <button class="ghost" type="button" @click="openSubmissionDetail(sub)">Open</button>
            <button
              class="ghost"
              type="button"
              @click="submitRowRejudge(sub)"
              :disabled="rejudgeState.loading && rejudgeState.submissionUuid === sub.submissionUuid"
            >
              {{ rejudgeState.loading && rejudgeState.submissionUuid === sub.submissionUuid ? 'Rejudging...' : 'Rejudge' }}
            </button>
          </div>
        </li>
      </ul>

      <footer class="pagination">
        <div class="pagination-info">
          <span>Page {{ pagination.page }}</span>
          <span>·</span>
          <span>{{ submissions.length }} items</span>
        </div>
        <div class="pagination-actions">
          <button class="ghost" type="button" @click="goPrevPage" :disabled="loading || pagination.page === 1">Previous</button>
          <button class="ghost" type="button" @click="goNextPage" :disabled="loading || !hasNextPage">Next</button>
        </div>
      </footer>
      </template>
    </section>

    <Modal
      title="Submission Detail"
      :show="detailModal.show"
      cancel-text="Close"
      :submit-text="detailModal.rejudging ? 'Rejudging...' : 'Rejudge'"
      :submit-disabled="detailModal.loading || detailModal.rejudging || !detailModal.detail"
      @close="closeDetailModal"
      @submit="submitRejudge"
    >
      <div v-if="detailModal.loading" class="detail-state">Loading submission detail...</div>
      <div v-else-if="detailModal.error && !detailModal.detail" class="detail-state error-state">{{ detailModal.error }}</div>
      <div v-else-if="detailModal.detail" class="detail-body">
        <div class="detail-summary">
          <div>
            <p class="detail-label">Submission UUID</p>
            <p class="detail-value">{{ detailModal.detail.metadata.submissionUuid }}</p>
          </div>
          <div>
            <p class="detail-label">Status</p>
            <p class="detail-value">
              <span class="badge" :class="statusTone[detailModal.detail.metadata.status] || 'badge-muted'">
                {{ detailModal.detail.metadata.status }}
              </span>
            </p>
          </div>
          <div>
            <p class="detail-label">Score</p>
            <p class="detail-value">{{ detailModal.detail.metadata.score }}</p>
          </div>
          <div>
            <p class="detail-label">Problem / User</p>
            <p class="detail-value">{{ detailModal.detail.metadata.problemId }} / {{ detailModal.detail.metadata.userId }}</p>
          </div>
          <div>
            <p class="detail-label">Submitted At</p>
            <p class="detail-value">{{ formatDate(detailModal.detail.metadata.submitTime) }}</p>
          </div>
          <div>
            <p class="detail-label">Runtime / Memory</p>
            <p class="detail-value">{{ detailModal.detail.timeCost }} ms / {{ detailModal.detail.memoryCost }} KB</p>
          </div>
        </div>

        <div class="detail-section">
          <p class="detail-label">Language</p>
          <p class="detail-value">{{ detailModal.detail.language }}</p>
        </div>

        <div class="detail-section">
          <p class="detail-label">Code</p>
          <pre class="detail-code">{{ detailModal.detail.code }}</pre>
        </div>

        <div class="detail-section">
          <p class="detail-label">Case Results</p>
          <ul class="case-list">
            <li v-for="(caseResult, index) in detailModal.detail.caseResults.cases" :key="index" class="case-item">
              <span>Case {{ index + 1 }}</span>
              <span class="badge" :class="statusTone[caseResult.status] || 'badge-muted'">{{ caseResult.status }}</span>
              <span>{{ caseResult.score }}</span>
              <span>{{ caseResult.timeCost }} ms</span>
              <span>{{ caseResult.memoryCost }} KB</span>
            </li>
          </ul>
        </div>

        <div v-if="detailModal.error" class="detail-state error-state">{{ detailModal.error }}</div>
      </div>
    </Modal>
  </div>
</template>

<style scoped src="./Submissions.css"></style>
