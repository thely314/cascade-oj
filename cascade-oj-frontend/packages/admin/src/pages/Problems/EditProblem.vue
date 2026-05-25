<script setup lang="ts">
import { useEditProblem } from './EditProblemLogic'
import CodeEditor from '@/components/CodeEditor/CodeEditor.vue' 
import { marked } from 'marked' 
import DOMPurify from 'dompurify'

const renderMarkdown = (text: string) => {
  const rawHtml = marked.parse(text) as string 
  return DOMPurify.sanitize(rawHtml)
}

const { form, loading, handleUpdate, router } = useEditProblem()
</script>

<template>
  <div class="page">
    <header class="page-header">
      <div>
        <p class="eyebrow" @click="router.push('/problems')" style="cursor: pointer;">← Back to List</p>
        <h1 class="title">Create New Problem</h1>
      </div>
      <div class="actions">
        <button class="primary" @click="handleUpdate" :disabled="loading">
          {{ loading ? 'Saving...' : 'Save Changes' }}
        </button>
      </div>
    </header>

    <div class="new-problem-container">
      <!-- 基础信息 -->
      <div class="form-grid">
        <div class="input-group">
          <label>TITLE</label>
          <input v-model="form.title" placeholder="Problem Title..." />
        </div>
        <div class="input-group">
          <label>TIME LIMIT (MS)</label>
          <input type="number" v-model.number="form.timeLimitMs" />
        </div>
        <div class="input-group">
          <label>MEMORY LIMIT (MB)</label>
          <input type="number" v-model.number="form.memoryLimitMb" />
        </div>
      </div>

      <!-- 编辑区域 -->
      <div class="editor-section">
        <!-- 左侧：题面描述 (Markdown) -->
        <div class="panel-box">
          <div class="panel-header">DESCRIPTION (MARKDOWN)</div>
          <textarea v-model="form.description" style="height: 50%; border: none;"></textarea>
          <div class="panel-header" style="border-top: 1px solid #333">PREVIEW</div>
          <div class="preview-area" v-html="renderMarkdown(form.description)"></div>
        </div>

        <!-- 右侧：代码模板 (Monaco) -->
        <div class="panel-box">
          <div class="panel-header">INITIAL CODE TEMPLATE</div>
          <CodeEditor 
            v-model="form.codeTemplate" 
            language="cpp" 
            style="flex: 1" 
          />
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped src="./NewProblem.css"></style>