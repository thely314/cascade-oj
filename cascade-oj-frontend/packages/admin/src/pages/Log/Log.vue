<script setup lang="ts">
import { ref, onMounted, reactive, computed } from 'vue'
import { VueDatePicker } from '@vuepic/vue-datepicker'
import '@vuepic/vue-datepicker/dist/main.css'
import { PaginationBar } from 'v-page'
import { listLogFiles, queryLogContent, downloadLogs } from '../../api/admin'
import type { ListLogFilesRequest, QueryLogContentRequest } from '../../api/types'

const logFiles = ref<string[]>([])
const selectedFiles = ref<string[]>([])
const logContent = ref<string[]>([])
const loadingFiles = ref(false)
const loadingContent = ref(false)
const error = ref('')

// 分页状态
const filePagination = reactive({
  current: 1,
  pageSize: 8,
  total: 0,
})
const filePageSizeOptions = [8, 16, 32]

const contentPagination = reactive({
  current: 1,
  pageSize: 30,
  total: 0,
})
const contentPageSizeOptions = [30, 45, 60]

const filters = reactive({
  level: '',
  timeRange: ref<[Date, Date]>(),
})

const isAllSelected = computed({
  get: () => logFiles.value.length > 0 && selectedFiles.value.length === logFiles.value.length,
  set: (value: boolean) => {
    if (value) {
      selectedFiles.value = [...logFiles.value]
    } else {
      selectedFiles.value = []
    }
  },
})

const selectAllFiles = (event: Event) => {
  const target = event.target as HTMLInputElement
  isAllSelected.value = target.checked
}

const fetchLogFiles = async () => {
  loadingFiles.value = true
  error.value = ''
  try {
    const params: ListLogFilesRequest = {
      pageSize: filePagination.pageSize,
      offset: (filePagination.current - 1) * filePagination.pageSize,
    }
    const response = await listLogFiles(params)
    logFiles.value = response.filenames || []
    filePagination.total = response.total || 0
  } catch (err) {
    error.value = 'Failed to load log files'
    console.error(err)
  } finally {
    loadingFiles.value = false
  }
}

const handleFilePageChange = (page: number) => {
  if (page === filePagination.current) {
    return
  }
  filePagination.current = page
  fetchLogFiles()
}

const handleContentPageChange = (page: number) => {
  if (page === contentPagination.current) {
    return
  }
  contentPagination.current = page
  handleQuery()
}

const handleQuery = async () => {
  if (selectedFiles.value.length !== 1) {
    alert('Please select exactly one file to query.')
    return
  }
  loadingContent.value = true
  try {
    const { timeRange } = filters
    let timeStart: string | undefined
    let timeEnd: string | undefined

    if (timeRange && timeRange.length === 2 && timeRange[0] && timeRange[1]) {
      timeStart = timeRange[0].toISOString()
      timeEnd = timeRange[1].toISOString()
    }

    const params: QueryLogContentRequest = {
      filename: selectedFiles.value[0],
      level: filters.level,
      timeStart,
      timeEnd,
      pageSize: contentPagination.pageSize,
      offset: (contentPagination.current - 1) * contentPagination.pageSize,
    }
    const response = await queryLogContent(params)
    logContent.value = response.lines || []
    contentPagination.total = response.total || 0
  } catch (err) {
    error.value = 'Failed to query log content'
    console.error(err)
  } finally {
    loadingContent.value = false
  }
}

const handleExport = async () => {
  if (selectedFiles.value.length === 0) {
    alert('Please select files to export.')
    return
  }
  try {
    const blob = await downloadLogs({ filenames: selectedFiles.value })
    if (blob?.contentType === 'application/zip' && blob.data) {
      const binaryString = atob(blob.data);
      const bytes = new Uint8Array(binaryString.length);
      for (let i = 0; i < binaryString.length; i++) {
        bytes[i] = binaryString.charCodeAt(i);
      }
      const link = document.createElement('a')
      link.href = URL.createObjectURL(new Blob([bytes], { type: 'application/zip' }))
      link.download = `logs_${new Date().toISOString()}.zip`
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
      URL.revokeObjectURL(link.href) // Clean up
    } else {
      throw new Error('Invalid response format for log export')
    }
  } catch (err) {
    console.error(err)
  }
}

onMounted(() => {
  fetchLogFiles()
})

</script>

<template>
  <div class="page">
    <header class="page-header">
      <div>
        <p class="eyebrow">Admin</p>
        <h1 class="title">System Logs</h1>
      </div>
    </header>

    <div v-if="error" class="error-message">{{ error }}</div>

    <section class="panel">
      <header class="panel-header">
        <h2>Log Files</h2>
      </header>
      <!-- Replace with your component library's table -->
      <div class="table-container">
        <table>
          <thead>
            <tr>
              <td><input type="checkbox" :checked="isAllSelected" @change="selectAllFiles" /></td>
              <td>Select ALL</td>
            </tr>
          </thead>
          <tbody>
            <tr v-for="file in logFiles" :key="file">
              <td><input type="checkbox" :value="file" v-model="selectedFiles" /></td>
              <td>{{ file }}</td>
            </tr>
          </tbody>
        </table>
      </div>
      <!-- log file pagination -->
      <span class="pagination" v-if="filePagination.total > 0">
        <PaginationBar
          language="cn"
          align="left"
          v-model:model-value="filePagination.current"
          v-model:page-size="filePagination.pageSize"
          v-model:total-row="filePagination.total"
          :page-size-menu="filePageSizeOptions"
          @update:model-value="handleFilePageChange"
        />
      </span>
    </section>

    <section class="panel">
      <header class="panel-header">
        <h2>Query Logs</h2>
      </header>
      <div class="filters">
        <header class="panel-header">
        <h3>日志级别</h3>
        </header>
        <!-- level picker -->
        <select v-model="filters.level" class="level-filter">
          <option value="">All Levels</option>
          <option value="INFO">Info</option>
          <option value="WARN">Warn</option>
          <option value="ERROR">Error</option>
          <option value="FATAL">Fatal</option>
        </select>
        <header class="panel-header">
        <h3>时间范围</h3>
        </header>
        <!-- date range picker here -->
        <div style="max-width: 50%;"><VueDatePicker v-model="filters.timeRange" range dark /></div>
        <button class="primary" @click="handleQuery" :disabled="loadingContent || selectedFiles.length !== 1">
          {{ loadingContent ? 'Querying...' : 'Query Selected File' }}
        </button>
        <button class="ghost" @click="handleExport" :disabled="selectedFiles.length === 0">
          Export Selected ({{ selectedFiles.length }})
        </button>
      </div>
      <div class="log-content-container">
        <pre v-if="!loadingContent">{{ logContent.join('\n') }}</pre>
        <div v-else>Loading content...</div>
      </div>
      <!-- Pagination for content would go here -->
      <span class="pagination" v-if="contentPagination.total > 0">
        <PaginationBar
          language="cn"
          align="left"
          v-model:model-value="contentPagination.current"
          v-model:page-size="contentPagination.pageSize"
          v-model:total-row="contentPagination.total"
          :page-size-menu="contentPageSizeOptions"
          @update:model-value="handleContentPageChange"
        />
      </span>
    </section>
  </div>
</template>

<style scoped src="./Log.css"></style>
