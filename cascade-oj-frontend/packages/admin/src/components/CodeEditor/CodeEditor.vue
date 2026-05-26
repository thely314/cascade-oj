<template>
  <div ref="editorContainer" class="monaco-editor-container"></div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, watch, shallowRef } from 'vue'
import * as monaco from 'monaco-editor'

const props = defineProps({
  modelValue: { type: String, required: true },
  language: { type: String, default: 'cpp' },
  theme: { type: String, default: 'vs-dark' },
  readOnly: { type: Boolean, default: false } 
})

// 语言与 Monaco 编辑器高亮映射
const languageHighlightMap: Record<string, string> = {
  'c': 'c',
  'c++11': 'cpp',
  'c++11(O2)': 'cpp'
};

const emit = defineEmits(['update:modelValue', 'change'])
const editorContainer = ref<HTMLElement | null>(null)
const editorInstance = shallowRef<monaco.editor.IStandaloneCodeEditor | null>(null) // 使用 shallowRef 存储 editor 实例，避免 Vue 深度代理导致性能问题

// 初始化编辑器
const initMonaco = () => {
  if (!editorContainer.value) return

  editorInstance.value = monaco.editor.create(editorContainer.value, {
    value: props.modelValue,
    language: props.language,
    theme: props.theme,
    readOnly: props.readOnly,
    automaticLayout: false, // 手动控制 layout 以获得更好性能
    minimap: { enabled: false }, // 关闭代码缩略图
    fontSize: 14,
    fontFamily: 'Consolas, "Courier New", monospace',
    scrollBeyondLastLine: false,
    lineNumbersMinChars: 3, // 行号宽度
    padding: { top: 10, bottom: 10 }
  })

  // 监听内容变化 -> 通知父组件
  editorInstance.value.onDidChangeModelContent(() => {
    const value = editorInstance.value?.getValue() || ''
    if (value !== props.modelValue) {
      emit('update:modelValue', value)
      emit('change', value)
    }
  })
}

// 监听 readOnly 变化，动态更新编辑器配置
watch(() => props.readOnly, (newVal) => {
  editorInstance.value?.updateOptions({ readOnly: newVal })
})

// 监听父组件传入的 value 变化（例如切换文件时）
watch(
  () => props.modelValue,
  (newValue) => {
    if (editorInstance.value) {
      const currentValue = editorInstance.value.getValue()
      // 只有当内容确实不同时才设置，避免光标跳动
      if (newValue !== currentValue) {
        editorInstance.value.setValue(newValue)
      }
    }
  }
)

// 监听语言变化
watch(() => props.language, (newLang) => {
  if (editorInstance.value) {
    const model = editorInstance.value.getModel()
    if (model) {
      // console.log(monaco.languages.getLanguages());
      const langHighlights = languageHighlightMap[newLang] || 'plaintext';
      monaco.editor.setModelLanguage(model, langHighlights)
    }
  }
})

// --- 自动调整大小逻辑 ---
let resizeObserver: ResizeObserver | null = null

onMounted(() => {
  initMonaco()

  // 使用 ResizeObserver 监听容器大小变化，实现拖拽时的平滑自适应
  if (editorContainer.value && editorInstance.value) {
    resizeObserver = new ResizeObserver(() => {
      editorInstance.value?.layout()
    })
    resizeObserver.observe(editorContainer.value)
  }
})

onBeforeUnmount(() => {
  // 销毁实例，防止内存泄漏
  if (editorInstance.value) {
    editorInstance.value.dispose()
  }
  if (resizeObserver) {
    resizeObserver.disconnect()
  }
})
</script>

<style scoped>
.monaco-editor-container {
  width: 100%;
  height: 100%;
  overflow: hidden;
}
</style>