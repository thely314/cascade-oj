import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import router from './router'; // 引入路由

createApp(App)
  .use(router) // 注册路由
  .mount('#app');