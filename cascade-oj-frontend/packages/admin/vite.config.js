import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    vueDevTools(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
  },
  // --- 新增：Server 代理配置 ---
  server: {
    host: 'localhost',
    port: 15173,
    proxy: {
      '/api': {
        // TODO nginx docker 代理
        target: 'http://localhost:8080', // 这里填你队友后端的真实地址 (IP+端口)
        changeOrigin: true,
        // rewrite: (path) => path.replace(/^\/api/, '') // 如果后端接口不带 /api 前缀，就把这就行取消注释
      }
    }
  }
})
