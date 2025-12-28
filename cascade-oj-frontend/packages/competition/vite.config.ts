import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'
import monacoEditorPlugin from 'vite-plugin-monaco-editor' // <--- 引入插件

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    vueDevTools(),
    // 注册插件，按需加载语言可以减小体积
    // 使用 (xxx as any).default 的方式来强制调用
    (monacoEditorPlugin as any).default({
      languageWorkers: ['editorWorkerService', 'typescript', 'json', 'css', 'html'], 
      // 对于 C++，通常不需要特定 Worker，Monaco 基础包里包含了高亮规则
    })
  ],
  build: {
    // 使用相对路径，避免部分插件将 root 与绝对 outDir 进行字符串拼接导致路径异常
    outDir: '../../dist/competition',
    emptyOutDir: true,
  },
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
  },

  // --- 新增：Server 代理配置 ---
  server: {
    proxy: {
      '/api': {
        // nginx docker 代理
        target: 'http://localhost:80',
        changeOrigin: true,
        // rewrite: (path) => path.replace(/^\/api/, '') // 如果后端接口不带 /api 前缀，且 nginx 未处理，就把这就行取消注释
      }
    }
  }
})