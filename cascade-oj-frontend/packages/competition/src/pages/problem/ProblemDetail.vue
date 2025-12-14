<!-- packages/competition/src/pages/problem/ProblemDetail.vue -->
<template>
  <div class="problem-container">
    <!-- 提交结果模态框 -->
    <div v-if="showResultModal" class="modal-overlay" @click="showResultModal = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3 :class="submissionResult?.status === 'Accepted' ? 'text-green' : 'text-red'">
            {{ submissionResult?.status }}
          </h3>
          <button @click="showResultModal = false">✕</button>
        </div>
        <div class="modal-body">
          <div class="stat-item"><span>时间:</span> {{ submissionResult?.time || '-' }}</div>
          <div class="stat-item"><span>内存:</span> {{ submissionResult?.memory || '-' }}</div>
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

        <!-- 3. 题目内容 -->
        <div class="pane-body scrollable">
          <div class="problem-content-placeholder">
            <h1>{{ problemData.title }}</h1>
            <div class="meta-info">
              <span>时间限制: {{ problemData.timeLimit }}</span>
              <span>内存限制: {{ problemData.memoryLimit }}</span>
            </div>
            <!-- Markdown 渲染区 -->
            <div class="markdown-body" v-html="descriptionHtml"></div>
            <br><br>
          </div>
        </div>

        <!-- 4. 底部切题栏 -->
        <div class="pane-footer left-footer">
          <button class="icon-btn list-btn" @click="showProblemDrawer = !showProblemDrawer">
            ☰ 题目列表
          </button>
          <div class="pagination-ctrl">
            <button class="nav-btn" @click="handlePrevProblem"> &lt; 上一题</button>
            <span class="page-num">{{ problemData.id }} / {{ problemList.length || '...' }}</span>
            <button class="nav-btn" @click="handleNextProblem">下一题 &gt; </button>
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
            <span class="settings-icon">⚙️</span>
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
                <!-- 1. 顶部统计信息栏 (仅当有数据时显示) -->
                <div v-if="runStats.time !== '-'" class="run-meta">
                  <div class="meta-status">===== 运行成功 =====</div>
                  <div class="meta-detail">
                    <span>CPU 时间: {{ runStats.time }}</span>
                    <span>内存占用: {{ runStats.memory }}</span>
                  </div>
                  <div class="meta-divider">===== 程序输出 =====</div>
                </div>

                <!-- 2. 实际输出内容 -->
                <div class="output-view">{{ playgroundOutput }}</div>
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
import { useProblemDetail } from './ProblemDetail'; // 引入抽离的逻辑

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
  showResultModal, submissionResult,
  runStats
} = useProblemDetail();
</script>

<style scoped src="./ProblemDetail.css"></style> 