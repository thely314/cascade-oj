<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue';
import { useToast } from 'vue-toastification';
import { PaginationBar } from 'v-page';
import Modal from '@/components/Modal.vue';
import {
  listReports,
  getReport,
  triggerAnalysis,
  getMetricsSummary,
  ackReport,
} from '@/api/aiops';
import type {
  ReportSummary,
  ListReportsReply,
  GetReportReply,
  GetMetricsSummaryReply,
} from '@/api/types';

// ==================== Report Section ====================

const reports = ref<ReportSummary[]>([]);
const reportTotal = ref(0);
const reportLoading = ref(false);
const statusFilter = ref('');
const statusOptions = [
  { label: '全部', value: '' },
  { label: 'draft', value: 'draft' },
  { label: 'published', value: 'published' },
  { label: 'acknowledged', value: 'acknowledged' },
];

const reportPagination = reactive({
  current: 1,
  pageSize: 10,
});
const reportPageSizeOptions = [5, 10, 20, 50];

const fetchReports = async () => {
  reportLoading.value = true;
  try {
    const result: ListReportsReply = await listReports({
      page: reportPagination.current,
      pageSize: reportPagination.pageSize,
      statusFilter: statusFilter.value,
    });
    reports.value = result.reports || [];
    reportTotal.value = result.total || 0;
  } catch (err) {
    console.error('Failed to load reports:', err);
  } finally {
    reportLoading.value = false;
  }
};

const handleReportPageChange = (page: number) => {
  if (page === reportPagination.current) return;
  reportPagination.current = page;
  fetchReports();
};

const onFilterChange = () => {
  reportPagination.current = 1;
  fetchReports();
};

// Report detail modal
const detailModal = reactive({
  show: false,
  loading: false,
  ackLoading: false,
  data: null as GetReportReply | null,
});

const viewReportDetail = async (report: ReportSummary) => {
  detailModal.show = true;
  detailModal.loading = true;
  detailModal.data = null;
  try {
    detailModal.data = await getReport(report.id);
  } catch (err) {
    console.error('Failed to get report detail:', err);
  } finally {
    detailModal.loading = false;
  }
};

const handleAckReport = async () => {
  if (!detailModal.data) return;
  detailModal.ackLoading = true;
  try {
    await ackReport(detailModal.data.id);
    // Refresh detail to show updated status
    detailModal.data = await getReport(detailModal.data.id);
    // Refresh list
    fetchReports();
  } catch (err) {
    console.error('Failed to acknowledge report:', err);
  } finally {
    detailModal.ackLoading = false;
  }
};

const formatTimestamp = (ts: string): string => {
  if (!ts) return '-';
  const d = new Date(ts);
  return d.toLocaleString();
};

const toast = useToast();

// Trigger analysis form
const analysisForm = reactive({
  time_range_minutes: 15,
  services: '',
});

const analysisLoading = ref(false);

const submitAnalysis = () => {
  const servicesList = analysisForm.services
    .split(',')
    .map((s) => s.trim())
    .filter((s) => s.length > 0);

  analysisLoading.value = true;

  // 即时通知：任务已触发，LLM 处理中，去 Reports 列表等结果
  toast.info('分析任务已触发，请稍后刷新 Reports 列表查看结果', {
    timeout: 5000,
  });

  // 短暂 loading 后恢复按钮
  setTimeout(() => {
    analysisLoading.value = false;
  }, 1500);

  // fire-and-forget：超时是预期行为
  triggerAnalysis({
    timeRangeMinutes: analysisForm.time_range_minutes,
    services: servicesList,
  })
    .then((result) => {
      fetchReports();
      toast.success(`分析完成，Report #${result.reportId} 已生成`);
    })
    .catch(() => {
      // 超时已在初始 toast 中提示，非超时异常静默处理
    });
};

// ==================== Metrics Section ====================

const allServices = ['public', 'user', 'admin', 'judge', 'aiops'];
const selectedServices = ref<string[]>([...allServices]);
const allServicesSelected = ref(true);

const toggleAllServices = () => {
  if (allServicesSelected.value) {
    selectedServices.value = [];
  } else {
    selectedServices.value = [...allServices];
  }
  allServicesSelected.value = !allServicesSelected.value;
};

const onServiceCheck = () => {
  allServicesSelected.value = selectedServices.value.length === allServices.length;
};

const metricsLoading = ref(false);
const metricsResult = ref<GetMetricsSummaryReply | null>(null);

const fetchMetrics = async () => {
  metricsLoading.value = true;
  metricsResult.value = null;
  try {
    metricsResult.value = await getMetricsSummary({
      services: selectedServices.value,
    });
  } catch (err) {
    console.error('Failed to get metrics:', err);
  } finally {
    metricsLoading.value = false;
  }
};

const formatJson = (raw: string): string => {
  if (!raw) return '-';
  try {
    return JSON.stringify(JSON.parse(raw), null, 2);
  } catch {
    return raw;
  }
};

onMounted(() => {
  fetchReports();
});
</script>

<template>
  <div class="page">
    <!-- ==================== Report Section ==================== -->
    <section class="section">
      <h2 class="section-title">Reports</h2>

      <!-- Filter bar -->
      <div class="toolbar">
        <div class="filter-group">
          <label class="filter-label">状态过滤：</label>
          <select v-model="statusFilter" class="filter-select" @change="onFilterChange">
            <option v-for="opt in statusOptions" :key="opt.value" :value="opt.value">
              {{ opt.label }}
            </option>
          </select>
        </div>
        <button class="btn-refresh" :disabled="reportLoading" @click="fetchReports">
          {{ reportLoading ? '加载中...' : 'Refresh' }}
        </button>
      </div>

      <!-- Report table -->
      <div class="table-container">
        <table class="data-table">
          <thead>
            <tr>
              <th>ID</th>
              <th>标题</th>
              <th>噪声数</th>
              <th>真实风险数</th>
              <th>状态</th>
              <th>创建时间</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="reports.length === 0 && !reportLoading">
              <td colspan="6" class="empty-cell">暂无数据</td>
            </tr>
            <tr
              v-for="report in reports"
              :key="report.id"
              class="clickable-row"
              @click="viewReportDetail(report)"
            >
              <td>{{ report.id }}</td>
              <td>{{ report.title }}</td>
              <td>{{ report.noiseCount }}</td>
              <td>{{ report.realRiskCount }}</td>
              <td>
                <span class="status-tag" :class="'status-' + report.status">
                  {{ report.status }}
                </span>
              </td>
              <td>{{ formatTimestamp(report.createdAt) }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <span v-if="reportTotal > 0" class="pagination">
        <PaginationBar
          language="cn"
          align="left"
          v-model:model-value="reportPagination.current"
          v-model:page-size="reportPagination.pageSize"
          v-model:total-row="reportTotal"
          :page-size-menu="reportPageSizeOptions"
          @update:model-value="handleReportPageChange"
        />
      </span>

      <!-- Trigger analysis form -->
      <div class="analysis-form">
        <h3 class="subsection-title">触发手动分析</h3>
        <div class="form-row">
          <div class="form-group inline">
            <label>时间范围（分钟）：</label>
            <input
              v-model.number="analysisForm.time_range_minutes"
              type="number"
              min="1"
              class="form-input short"
            />
          </div>
          <div class="form-group inline">
            <label>服务（逗号分隔，留空=全部）：</label>
            <input
              v-model="analysisForm.services"
              type="text"
              placeholder="public,user,admin,judge,aiops"
              class="form-input"
            />
          </div>
          <button class="btn-primary" :disabled="analysisLoading" @click="submitAnalysis">
            {{ analysisLoading ? '分析中...' : 'Analysis' }}
          </button>
        </div>
      </div>
    </section>

    <!-- ==================== Metrics Section ==================== -->
    <section class="section">
      <h2 class="section-title">Metrics</h2>

      <div class="toolbar">
        <div class="checkbox-group">
          <label class="checkbox-item">
            <input
              type="checkbox"
              :checked="allServicesSelected"
              @change="toggleAllServices"
            />
            <span>全部</span>
          </label>
          <label
            v-for="svc in allServices"
            :key="svc"
            class="checkbox-item"
          >
            <input
              type="checkbox"
              :value="svc"
              v-model="selectedServices"
              @change="onServiceCheck"
            />
            <span>{{ svc }}</span>
          </label>
        </div>
        <button class="btn-primary" :disabled="metricsLoading" @click="fetchMetrics">
          {{ metricsLoading ? '查询中...' : '查询 Metrics' }}
        </button>
      </div>

      <!-- Metrics result -->
      <div v-if="metricsResult" class="result-card">
        <p><strong>Prometheus 状态:</strong>
          <span :class="metricsResult.prometheusStatus === 'ok' ? 'text-ok' : 'text-warn'">
            {{ metricsResult.prometheusStatus }}
          </span>
        </p>
        <p><strong>查询时间:</strong> {{ formatTimestamp(metricsResult.queriedAt) }}</p>
        <details>
          <summary>Metrics 快照 (JSON)</summary>
          <pre class="json-preview">{{ formatJson(metricsResult.metricsSnapshotJson) }}</pre>
        </details>
      </div>
    </section>

    <!-- ==================== Report Detail Modal ==================== -->
    <Modal
      title="报告详情"
      :show="detailModal.show"
      submit-text="Acknowledge"
      cancel-text="关闭"
      :submit-disabled="detailModal.ackLoading || (detailModal.data?.status === 'acknowledged')"
      @close="detailModal.show = false"
      @submit="handleAckReport"
    >
      <div v-if="detailModal.loading" class="modal-loading">加载中...</div>
      <div v-else-if="detailModal.data" class="report-detail">
        <p><strong>ID:</strong> {{ detailModal.data.id }}</p>
        <p><strong>标题:</strong> {{ detailModal.data.title }}</p>
        <p><strong>摘要:</strong> {{ detailModal.data.summary }}</p>
        <p><strong>状态:</strong>
          <span class="status-tag" :class="'status-' + detailModal.data.status">
            {{ detailModal.data.status }}
          </span>
        </p>
        <p><strong>噪声数:</strong> {{ detailModal.data.noiseCount }}</p>
        <p><strong>真实风险数:</strong> {{ detailModal.data.realRiskCount }}</p>
        <p><strong>创建时间:</strong> {{ formatTimestamp(detailModal.data.createdAt) }}</p>
        <p><strong>确认时间:</strong> {{ formatTimestamp(detailModal.data.acknowledgedAt) }}</p>
        <details v-if="detailModal.data.rootCausesJson">
          <summary>根因分析 (JSON)</summary>
          <pre class="json-preview">{{ formatJson(detailModal.data.rootCausesJson) }}</pre>
        </details>
        <details v-if="detailModal.data.suggestedActionsJson">
          <summary>建议操作 (JSON)</summary>
          <pre class="json-preview">{{ formatJson(detailModal.data.suggestedActionsJson) }}</pre>
        </details>
        <details v-if="detailModal.data.fullReportMarkdown">
          <summary>完整报告 (Markdown)</summary>
          <pre class="markdown-preview">{{ detailModal.data.fullReportMarkdown }}</pre>
        </details>
      </div>
      <div v-else class="modal-loading">加载失败</div>
    </Modal>
  </div>
</template>

<style scoped src="./AIOpsPage.css"></style>
