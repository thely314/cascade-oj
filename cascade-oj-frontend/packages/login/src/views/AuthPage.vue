<template>
  <div class="auth-container">
    <AuthLeft />
    <div class="right-section">
      <div class="auth-card">
        <h2>欢迎来到 Cascade OJ</h2>
        <p class="sub-title">开始您的编程之旅</p>
        
        <div class="tab-buttons">
          <button 
            :class="{ active: isLogin }" 
            @click="isLogin = true"
          >
            登录
          </button>
          
          <button 
            :class="{ active: !isLogin }" 
            @click="isLogin = false"
          >
            注册
          </button>
        </div>
        
        <component 
          :is="isLogin ? LoginForm : RegisterForm" 
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import AuthLeft from '../components/auth/AuthLeft.vue';
import LoginForm from '../components/auth/LoginForm.vue';
import RegisterForm from '../components/auth/RegisterForm.vue';

const isLogin = ref(true);
</script>

<style scoped>
.auth-container {
  display: flex;
  min-height: 100vh;
  background: #161616; /* 与 home 页面一致的深色背景 */
  color: #cbd5c0;
}

.right-section {
  width: 480px; /* 稍微加大 auth 页右侧最大宽度 */
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  padding: 40px 24px;
  box-sizing: border-box; /* include padding in width to prevent overflow */
}

.auth-card {
  width: 100%;
  background: #141816; /* 卡片深色 */
  border-radius: 10px;
  padding: 28px;
  box-shadow: 0 6px 18px rgba(22, 163, 142, 0.06);
  border: 1px solid rgba(30,134,68,0.06);
  box-sizing: border-box; /* ensure padding doesn't enlarge card beyond container */
}

.auth-card h2 {
  font-size: 22px;
  margin: 0 0 6px 0;
  color: #1dad80; /* 主色绿 */
}

.sub-title {
  font-size: 13px;
  color: #9fbebb;
  margin-bottom: 20px;
}

.tab-buttons {
  display: flex;
  margin-bottom: 20px;
  gap: 8px;
}

.tab-buttons button {
  padding: 8px 16px;
  border: 1px solid transparent;
  background: transparent;
  color: #cbd5c0;
  font-size: 15px;
  cursor: pointer;
  border-radius: 6px;
}

.tab-buttons button.active {
  background: rgba(29,173,128,0.12);
  color: #1dad80;
  border-color: rgba(29,173,128,0.16);
  font-weight: 600;
}

.tab-buttons button:hover {
  background: rgba(255,255,255,0.02);
}
</style>