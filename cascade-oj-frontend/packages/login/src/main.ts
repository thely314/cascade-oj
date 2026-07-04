import { createApp } from 'vue'
import Toast, { POSITION, useToast } from 'vue-toastification'
import 'vue-toastification/dist/index.css'
import './style.css'
import App from './App.vue'
import router from './router'
import { initToast } from './utils/toast'

const app = createApp(App)
app.use(router)
app.use(Toast, {
  position: POSITION.TOP_RIGHT,
  timeout: 4000,
  closeOnClick: true,
  pauseOnHover: true,
  draggable: true,
  theme: 'dark',
})

app.mount('#app')
initToast(useToast())