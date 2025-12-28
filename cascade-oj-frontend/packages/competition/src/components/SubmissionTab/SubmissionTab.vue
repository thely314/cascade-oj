<!-- src/pages/problem/components/SubmissionTab.vue -->
<template>
  <div class="submission-list">
    <div v-if="loading" class="loading-state">加载记录中...</div>
    
    <div v-else-if="list.length === 0" class="empty-state">
      暂无提交记录
    </div>

    <table v-else class="data-table">
      <thead>
        <tr>
          <th>状态</th>
          <th>分数</th>
          <th>语言</th>
          <th>耗时</th>
          <th>内存</th>
          <th>提交时间</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="item in list" :key="item.submissionUuid">
          <!-- 状态：根据不同结果显示不同颜色 -->
          <td>
            <span :class="getStatusClass(item.status)">{{ item.status }}</span>
          </td>
          <td>{{ item.score }}</td>
          <td>{{ item.language || '-' }}</td>
          <td>{{ item.timeCost ? item.timeCost + 'ms' : '0ms' }}</td>
          <td>{{ item.memoryCost ? (item.memoryCost / 1024).toFixed(1) + 'MB' : '0MB' }}</td>
          <td class="time-col">{{ formatTime(item.submitTime) }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { fetchSubmissions, type SubmissionItem } from '@/api/problem';

// 接收父组件传来的参数
const props = defineProps<{
  problemId: string;
  contestId: string;
}>();

const list = ref<SubmissionItem[]>([]);
const loading = ref(true);

const loadData = async () => {
  loading.value = true;
  try {
    list.value = await fetchSubmissions(props.contestId, props.problemId);
  } catch (e) {
    console.error(e);
  } finally {
    loading.value = false;
  }
};

// [新增] 暴露方法
defineExpose({
  refresh: loadData
});

const getStatusClass = (status: string) => {
  if (status === 'accepted') return 'status-ac';
  if (status === 'pending') return 'status-run';
  return 'status-wa';
};

const formatTime = (isoStr: string) => {
  return new Date(isoStr).toLocaleString();
};

onMounted(() => {
  loadData();
});
</script>

<style scoped>
.submission-list {
  padding: 0;
  width: 100%;
}

.loading-state, .empty-state {
  padding: 40px;
  text-align: center;
  color: #888;
}

/* 表格样式 */
.data-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 13px;
  color: #e0e0e0;
}

.data-table th {
  text-align: left;
  padding: 12px 15px;
  background-color: #252526;
  color: #888;
  font-weight: normal;
  border-bottom: 1px solid #333;
}

.data-table td {
  padding: 12px 15px;
  border-bottom: 1px solid #2c2c2c;
}

.data-table tr:hover {
  background-color: #2a2a2a;
}

/* 状态颜色 */
.status-ac { color: #2ecc71; font-weight: bold; }
.status-wa { color: #e74c3c; font-weight: bold; }
.status-run { color: #f1c40f; font-weight: bold; }

.time-col { color: #888; font-size: 12px; }
</style>