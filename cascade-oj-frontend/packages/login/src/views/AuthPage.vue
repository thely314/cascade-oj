<!-- src/views/AuthPage.vue -->
<template>
  <div class="auth-container">
    <div class="auth-wrapper">
      <!-- 左侧只负责显示 Logo 和文字 -->
      <AuthLeft class="auth-left-col" />
      
      <!-- 右侧表单 -->
      <div class="auth-right-col">
        <div class="auth-box">
          <div class="auth-header">
            <h2>{{ isLogin ? '欢迎来到 Cascade' : '加入 Cascade' }}</h2>
            <p class="sub-title">
              {{ isLogin ? '登录以继续您的编程之旅' : '注册账号，开启您的编程之旅' }}
            </p>
          </div>

          <!-- Tab 切换 -->
          <div class="auth-tabs">
            <div 
              class="tab-item" 
              :class="{ active: isLogin }" 
              @click="isLogin = true"
            >
              登录
            </div>
            <div 
              class="tab-item" 
              :class="{ active: !isLogin }" 
              @click="isLogin = false"
            >
              注册
            </div>
            <!-- 滑动滑块 -->
            <div class="tab-slider" :class="{ right: !isLogin }"></div>
          </div>

          <Transition name="fade" mode="out-in">
            <component :is="isLogin ? LoginForm : RegisterForm" />
          </Transition>
        </div>
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
  min-height: 100vh;
  width: 100vw;
  /* 核心修改：将背景移到这里，并铺满全屏 */
  background-color: #0c1110;
  background-image: url('https://assets.codepen.io/1462889/pat-back.svg');
  background-position: center;
  background-repeat: no-repeat;
  background-size: cover; /* 确保波浪线覆盖整个屏幕 */
  
  display: flex;
  align-items: center;
  justify-content: center;
}

.auth-wrapper {
  display: flex;
  width: 100%;
  height: 100vh;
  overflow: hidden;
}

.auth-left-col {
  flex: 1.2;
  display: flex;
  flex-direction: column;
  justify-content: center;
  /* 背景透明，让父组件的波浪透过来 */
  background: transparent; 
  @media (max-width: 900px) {
    display: none;
  }
}

.auth-right-col {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px;
  /* 核心修改：右侧背景改为非常淡的透明黑，或者完全透明，消除割裂感 */
  background: rgba(12, 17, 16, 0.3); 
  backdrop-filter: blur(10px); /* 稍微加一点毛玻璃，防止文字和波浪线重叠看不清 */
}

.auth-box {
  width: 100%;
  max-width: 420px;
}

.auth-header {
  margin-bottom: 32px;
  text-align: left;
}

.auth-header h2 {
  font-size: 32px;
  font-weight: 700;
  color: #e3ece9;
  margin: 0 0 8px 0;
}

.sub-title {
  color: #6c7c7a;
  font-size: 15px;
}

/* Tab 样式保持不变 */
.auth-tabs {
  position: relative;
  display: flex;
  background: #151b19;
  padding: 4px;
  border-radius: 12px;
  margin-bottom: 32px;
  border: 1px solid rgba(255, 255, 255, 0.05);
}

.tab-item {
  flex: 1;
  text-align: center;
  padding: 10px 0;
  font-size: 14px;
  font-weight: 600;
  color: #6c7c7a;
  cursor: pointer;
  z-index: 2;
  transition: color 0.3s;
}

.tab-item.active {
  color: #e3ece9;
}

.tab-slider {
  position: absolute;
  top: 4px;
  left: 4px;
  width: calc(50% - 4px);
  height: calc(100% - 8px);
  background: #23aa8f;
  border-radius: 8px;
  transition: transform 0.3s cubic-bezier(0.16, 1, 0.3, 1);
  z-index: 1;
  box-shadow: 0 2px 10px rgba(35, 170, 143, 0.3);
}

.tab-slider.right {
  transform: translateX(100%);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>