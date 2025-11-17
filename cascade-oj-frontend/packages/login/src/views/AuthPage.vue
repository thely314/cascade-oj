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
}

.right-section {
  width: 400px;
  background-color: #fff;
  padding: 24px;
  box-shadow: 0 0 20px rgba(0, 0, 0, 0.05);
  display: flex;  
  align-items:center ;
  justify-content: space-between;
  min-height: 100vh;
}

.auth-card {
  width: 100%;
  margin: 0;
}

.auth-card h2 {
  font-size: 24px;
  margin-top: 0;
  margin-bottom: 8px;
}

.sub-title {
  font-size: 14px;
  color: #999;
  margin-bottom: 24px;
}

.tab-buttons {
  display: flex;
  margin-bottom: 24px;
  border-bottom: 1px solid #eee;
}

.tab-buttons button {
  padding: 12px 24px;
  border: none;
  background: none;
  font-size: 16px;
  cursor: pointer;
  transition: color 0.2s;
}

.tab-buttons button.active {
  border-bottom: 2px solid #0066FF;
  color: #0066FF;
  font-weight: bold;
}
</style>