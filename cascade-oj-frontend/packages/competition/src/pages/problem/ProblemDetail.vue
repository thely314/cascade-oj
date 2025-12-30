<!-- packages/competition/src/pages/problem/ProblemDetail.vue -->
<template>
  <div class="problem-container">
    <!-- 提交结果模态框 -->
    <div v-if="showResultModal" class="modal-overlay" @click="showResultModal = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3 :class="submissionResult?.status === 'accepted' ? 'text-green' : 'text-red'">
            {{ submissionResult?.status }}
          </h3>
          <button @click="showResultModal = false">✕</button>
        </div>
        <div class="modal-body">
          <div class="stat-item"><span>时间:</span> {{ submissionResult?.timeCost || '-' }}</div>
          <div class="stat-item"><span>内存:</span> {{ submissionResult?.memoryCost || '-' }}</div>
        </div>
        <div class="modal-footer">
          <button class="btn-primary" @click="showResultModal = false">确定</button>
        </div>
      </div>
    </div>

    <!-- 主体内容区 -->
    <div class="main-content" @mouseup="stopDrag" @mouseleave="stopDrag">
      
      <!-- A. 左半页面 -->
      <div class="pane left-pane" :style="{ width: leftWidth + '%' }">
        
        <!-- 1. 题目列表抽屉 -->
        <div class="problem-drawer" :class="{ open: showProblemDrawer }">
          <div class="drawer-header">
            <span>题目列表</span>
            <button class="close-drawer-btn" @click="showProblemDrawer = false">×</button>
          </div>
          <div class="drawer-list">
            <div 
              v-for="p in problemList" 
              :key="p.id" 
              class="drawer-item"
              :class="{ active: p.id === problemData.id }"
              @click="jumpToProblem(p.id)"
            >
              <span class="id">{{ p.id }}.</span>
              <span class="title">{{ p.title }}</span>
            </div>
          </div>
        </div>

        <!-- 2. 顶部导航 -->
        <div class="pane-header">
          <button class="back-btn" @click="$router.push('/home')">← 返回</button>
          <div class="tab-nav">
            <span v-for="(tab, index) in leftTabs" :key="index" 
                  :class="{ active: currentLeftTab === index }"
                  @click="currentLeftTab = index">
              {{ tab }}
            </span>
          </div>
        </div>

        <!-- 3. 内容区域 -->
        <div class="pane-body scrollable">

          <!-- Tab 0: 题目描述 -->
          <div v-if="currentLeftTab === 0" class="problem-content-placeholder">
            <h1>{{ problemData.title }}</h1>
            <div class="meta-info">
              <span>作者: {{ problemData.creator || 'Admin' }}</span> 
              <span>时间限制: {{ problemData.timeLimit }}</span>
              <span>内存限制: {{ problemData.memoryLimit }}</span>
            </div>
            <!-- Markdown 渲染区 -->
            <div class="markdown-body" v-html="descriptionHtml"></div>
            <br><br>
          </div>

          <!-- Tab 1: 提交记录 (新增) -->
          <!-- 注意：这里使用了 key，强制在切换 tab 或题目时重新渲染组件 -->
          <SubmissionTab 
            ref="submissionTabRef" 
            v-else-if="currentLeftTab === 1" 
            :contest-id="contestId" 
            :problem-id="problemData.id" 
            :key="`sub-${problemData.id}`"
          />

          <!-- Tab 2: 成绩 (显示所属比赛的排名信息) -->
          <div v-else-if="currentLeftTab === 2" class="rank-pane">
            <h3>比赛排名</h3>
            <p v-if="rankLoading">加载中...</p>
            <p v-else-if="!rankLoading && rankEntries.length === 0">暂无排名数据</p>
            <ol class="rank-list" v-else>
              <li v-for="e in rankEntries" :key="e.userId">{{ (e as any).username || e.userId }} — {{ e.score }} 分</li>
            </ol>
            <p v-if="rankError" class="error-text">{{ rankError }}</p>
          </div>

        </div>

        <!-- 4. 底部切题栏 -->
        <div class="pane-footer left-footer">
          <button class="icon-btn list-btn" @click="showProblemDrawer = !showProblemDrawer">
            ☰ 题目列表
          </button>
          <div class="pagination-ctrl">
            
            <!-- 上一题按钮 -->
            <button 
              class="nav-btn" 
              :disabled="isFirstProblem" 
              @click="handlePrevProblem"
              :style="{ opacity: isFirstProblem ? 0.5 : 1, cursor: isFirstProblem ? 'not-allowed' : 'pointer' }"
            > 
              &lt; 上一题
            </button>         

            <!-- 页码显示 -->
            <span class="page-num">{{ currentIndex }} / {{ totalProblems }}</span>

            <!-- 下一题按钮 -->
            <button 
              class="nav-btn" 
              :disabled="isLastProblem" 
              @click="handleNextProblem"
              :style="{ opacity: isLastProblem ? 0.5 : 1, cursor: isLastProblem ? 'not-allowed' : 'pointer' }"
            >
              下一题 &gt; 
            </button>
          </div>
        </div>
      </div>

      <!-- B. 分隔条 -->
      <div class="resizer" @mousedown="startDrag"></div>

      <!-- C. 右半页面 -->
      <div class="pane right-pane" :style="{ width: (100 - leftWidth) + '%' }">
        <!-- 1. Code Header (包含语言切换) -->
        <div class="pane-header code-header">
          <div class="file-tabs">
            <span v-for="file in currentFiles" :key="file.name"
                  class="file-tab" :class="{ active: currentFileName === file.name }"
                  @click="currentFileName = file.name">
              {{ file.name }}
              <span v-if="file.isReadOnly" style="margin-left:5px">🔒</span>
            </span>
          </div>
          
          <!-- 右侧工具栏 -->
          <div class="editor-actions">
            <!-- 语言选择器 -->
            <select v-model="currentLanguage" class="lang-select">
              <option v-for="lang in supportedLanguages" :key="lang" :value="lang">
                {{ lang }}
              </option>
            </select>
            <!-- 工具栏按钮 -->
            <span class="action-icon" @click="resetCode" title="重置代码">↺</span>
            <span class="action-icon" @click="copyCode" title="复制代码">📋</span>
            <span class="action-icon settings-icon" title="设置">⚙️</span>
            <span class="action-icon fullscreen-icon" title="全屏">⛶</span>
          </div>
        </div>

        <!-- 2. 编辑器 (绑定只读和语言) -->
        <div class="pane-body editor-wrapper">
          <CodeEditor 
            v-model="currentCode" 
            :language="currentFiles.find(f => f.name === currentFileName)?.lang || 'cpp'"
            :read-only="isCurrentReadOnly"
            theme="vs-dark" 
          />
        </div>

        <!-- 3. Playground 控制台 -->
        <div class="playground-panel" :class="{ open: isPlaygroundOpen }">
          <div class="playground-content">
            <div class="io-box">
              <div class="io-header">输入 (Stdin)</div>
              <textarea v-model="playgroundInput" placeholder="在此处填写输入数据..."></textarea>
            </div>
            
            <div class="io-box">
              <div class="io-header">
                <span>输出 (Stdout)</span>
              </div>
              
              <!-- 运行结果容器 -->
              <div class="output-container">
                
                <!-- 1. 顶部统计信息栏 -->
                <div v-if="runStats.time !== '-'" class="run-meta">
                  <div class="meta-status">===== 运行成功 =====</div>
                  <div class="meta-detail">
                    <span>CPU 时间: {{ runStats.time }}</span>
                    <span>内存占用: {{ runStats.memory }}</span>
                  </div>
                </div>

                <!-- 2. 调试/错误信息 (Stderr) -->
                <div v-if="runStats.stderr" class="stderr-view">
                  <div class="meta-divider text-red">===== 调试/错误信息 =====</div>
                  <pre>{{ runStats.stderr }}</pre>
                </div>

                <!-- 3. 实际输出内容 (Stdout) -->
                <div class="output-view-wrapper">
                  <div v-if="runStats.time !== '-'" class="meta-divider">===== 标准输出 =====</div>
                  <div class="output-view">{{ playgroundOutput }}</div>
                </div>
                
              </div>
            </div>
          </div>
        </div>

        <!-- 4. 底部操作按钮 -->
        <div class="pane-footer right-footer">
          <button class="playground-toggle" @click="isPlaygroundOpen = !isPlaygroundOpen">
            <span v-if="isPlaygroundOpen">⌄ 收起控制台</span>
            <span v-else>⌃ Playground</span>
          </button>
          
          <div class="action-btns">
             <button class="btn-test" :disabled="!isPlaygroundOpen || isRunning" @click="handleTestRun">
                <span v-if="isRunning">⏳ 运行中...</span>
                <span v-else>测试运行</span>
             </button>
             <button class="btn-submit" :disabled="isSubmitting" @click="handleSubmit"
                     :style="{ opacity: isSubmitting ? 0.6 : 1, cursor: isSubmitting ? 'not-allowed' : 'pointer' }">
                <span v-if="isSubmitting">提交中...</span>
                <span v-else>➤ 提交</span>
             </button>
          </div>
        </div>
      </div>

    </div>
  </div>
</template>

<script setup lang="ts">
import CodeEditor from '@/components/CodeEditor/CodeEditor.vue';
import SubmissionTab from '@/components/SubmissionTab/SubmissionTab.vue';
import { ref, watch, nextTick, onMounted } from 'vue';
import { useProblemDetail } from './ProblemDetail'; // 引入抽离的逻辑
import { getRanks, type RankItem } from '../../api/rank';

// 解构逻辑层导出的状态与方法
const {
  leftTabs, currentLeftTab, isPlaygroundOpen, loading, isSubmitting, isRunning,
  problemData, descriptionHtml,
  currentFiles, currentFileName, currentCode, isCurrentReadOnly,
  supportedLanguages, currentLanguage,
  playgroundInput, playgroundOutput,
  leftWidth, startDrag, stopDrag,
  handleTestRun, handleSubmit,
  showProblemDrawer, problemList,
  jumpToProblem, handlePrevProblem, handleNextProblem,
  currentIndex, isFirstProblem, isLastProblem, totalProblems,
  showResultModal, submissionResult,
  runStats,
  contestId,
  copyCode,
  resetCode 
} = useProblemDetail();

// 引用提交列表组件
const submissionTabRef = ref();
// 简单的比赛排名数据（文本列表）
const rankEntries = ref<RankItem[]>([]);
const rankLoading = ref(false);
const rankError = ref<string | null>(null);

async function fetchContestRank(id?: string) {
  if (!id) return;
  rankLoading.value = true;
  rankError.value = null;
  try {
    const res = await getRanks(id);
    rankEntries.value = res.ranks || [];
  } catch (e: any) {
    rankError.value = e?.message ?? '加载排名失败';
  } finally {
    rankLoading.value = false;
  }
}

onMounted(() => {
  const cid = (contestId as any)?.value ?? (contestId as any);
  fetchContestRank(cid);
});


// 监听提交结果，一旦有新结果(提交成功)，刷新列表
watch(submissionResult, (newVal) => {
  if (newVal) {
    // 如果当前没在“提交”Tab，自动切过去（可选，看体验偏好）
    // currentLeftTab.value = 1; 
    
    // 等待 DOM 更新后刷新列表
    nextTick(() => {
      if (currentLeftTab.value === 1 && submissionTabRef.value) {
        submissionTabRef.value.refresh();
      }
    });

    // 提交后刷新比赛排名数据
    const cid = (contestId as any)?.value ?? (contestId as any);
    fetchContestRank(cid);
  }
});

// 模块切换时刷新排名信息（无论切到哪个模块都更新一次）
watch(currentLeftTab, () => {
  const cid = (contestId as any)?.value ?? (contestId as any);
  fetchContestRank(cid);
});

</script>

<style scoped src="./ProblemDetail.css"></style>
<style scoped>
.rank-pane {
  padding: 20px;
  color: #cbd5c0;
}
.rank-list {
  padding-left: 1.2rem;
  margin: 0.5rem 0 0 0;
}
.rank-list li {
  margin: 6px 0;
}
.error-text { color: #e57373; }
</style>