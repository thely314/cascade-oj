<script setup lang="ts">
import { useEditProblem } from './EditProblemLogic'
import CodeEditor from '@/components/CodeEditor/CodeEditor.vue' 
import { marked } from 'marked' 
import DOMPurify from 'dompurify'

const renderMarkdown = (text: string) => {
  const rawHtml = marked.parse(text) as string 
  return DOMPurify.sanitize(rawHtml)
}

const { 
  form, 
  loading, 
  router,
  activeLang, 
  existingLangs,
  activeTemplate,
  handleFileUpload,
  handleUpdate,
  activeFileName,
  currentLangFiles,
  switchLang,
  addFile,
  renameFile,
  deleteFile
} = useEditProblem()
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

      <!-- 测试用例上传区 -->
      <div class="upload-section panel-box" style="margin-bottom: 24px;">
        <div class="panel-header">TEST CASES (ZIP)</div>
        <div class="upload-area" style="padding: 20px; display: flex; align-items: center; gap: 15px;">
          <input type="file" @change="handleFileUpload" accept=".zip" id="zip-file" hidden />
          <label for="zip-file" class="button secondary" style="cursor: pointer; padding: 8px 16px; border: 1px solid #444; border-radius: 4px;">
            {{ loading ? 'Uploading...' : 'Choose .zip Package' }}
          </label>
          <p style="font-size: 12px; color: #888;">Only .zip allowed. Must contain .in and .ans files.</p>
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

         <!-- 右侧：多语言 & 多文件代码模板 (巅峰重构) -->
        <div class="panel-box" style="display: flex; flex-direction: column;">
          <!-- 顶层：语言切换栏 -->
          <div style="display: flex; justify-content: space-between; align-items: center; padding: 10px 15px; border-bottom: 1px solid #333; background: #1a1a1a;">
            <span style="font-size: 13px; font-weight: bold; color: #42b983;">INITIAL CODE TEMPLATES</span>
            <div class="lang-tabs" style="display: flex; gap: 6px;">
              <button 
                v-for="lang in existingLangs" 
                :key="lang"
                @click="switchLang(lang)"
                :style="{
                  padding: '4px 10px', fontSize: '11px', fontWeight: 'bold', cursor: 'pointer',
                  background: activeLang === lang ? '#42b983' : 'transparent',
                  color: activeLang === lang ? '#000' : '#888',
                  border: '1px solid #42b983', borderRadius: '4px'
                }"
              >
                {{ lang.toUpperCase() }}
              </button>
            </div>
          </div>
          
          <!-- 次层：文件 Tab 栏 -->
          <div style="display: flex; background: #1e1e1e; padding-top: 8px; padding-left: 10px; border-bottom: 1px solid #2d2d2d; gap: 4px;">
            <!-- 渲染当前语言下的所有文件 -->
            <div 
              v-for="file in currentLangFiles" 
              :key="file.name"
              @click="activeFileName = file.name"
              @dblclick="renameFile(file.name)"
              title="Double click to rename"
              :style="{
                padding: '8px 16px', fontSize: '12px', cursor: 'pointer',
                background: activeFileName === file.name ? '#2d2d2d' : 'transparent',
                borderTop: activeFileName === file.name ? '2px solid #42b983' : '2px solid transparent',
                color: activeFileName === file.name ? '#fff' : '#888',
                borderTopLeftRadius: '4px', borderTopRightRadius: '4px',
                display: 'flex', alignItems: 'center', gap: '8px'
              }"
            >
              {{ file.name }}
              <!-- 删除按钮 -->
              <span 
                @click.stop="deleteFile(file.name)" 
                style="color: #ff5555; cursor: pointer; font-size: 14px; font-weight: bold; line-height: 1;"
                title="Delete File"
              >×</span>
            </div>
            
            <!-- 添加文件按钮 -->
            <button 
              @click="addFile" 
              title="Add New File"
              style="background: transparent; border: none; color: #42b983; cursor: pointer; font-size: 18px; padding: 0 10px; font-weight: bold;"
            >
              +
            </button>
          </div>
          
          <!-- 编辑器 -->
          <CodeEditor 
            v-if="activeTemplate"
            v-model="activeTemplate.code" 
            :language="activeLang" 
            style="flex: 1; min-height: 500px;" 
          />
          <!-- 如果当前语言没有文件时的空状态占位 -->
          <div v-else style="flex: 1; display: flex; align-items: center; justify-content: center; color: #555; min-height: 500px;">
            Click the '+' button to add a file for {{ activeLang.toUpperCase() }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped src="./NewProblem.css"></style>