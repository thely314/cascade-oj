<template>
  <div class="left-section">
    <div class="logo-area">
      <div class="logo">
        <img :src="logo" alt="Cascade logo" class="logo-img" />
      </div>
      <h1>Cascade OJ</h1>
    </div>
    <p class="slogan">编程艺术的殿堂 · The Art of Coding</p>
    
    <!-- 垂直排列的功能条目 -->
    <div class="features">
      
      <div class="feature-item">
        <div class="feature-icon blue">
          <svg viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" class="icon">
            <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5" stroke="currentColor" stroke-width="2" fill="none"/>
          </svg>
        </div>
        <div class="feature-text">
          <h3>竞赛训练</h3>
          <p>ACM/ICPC、蓝桥杯等各类编程竞赛题目</p>
        </div>
      </div>
      
      <div class="feature-item">
        <div class="feature-icon green">
          <svg viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" class="icon">
            <path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z" fill="currentColor"/>
          </svg>
        </div>
        <div class="feature-text">
          <h3>社区交流</h3>
          <p>与全球程序员一起讨论算法与数据结构</p>
        </div>
      </div>
      
      <div class="feature-item">
        <div class="feature-icon purple">
          <svg viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" class="icon">
            <path d="M17 18a5 5 0 0 0-10 0h10zm2-8a5 5 0 0 0-10 0h10zm2-8a5 5 0 0 0-10 0h10z" fill="currentColor"/>
          </svg>
        </div>
        <div class="feature-text">
          <h3>实时评测</h3>
          <p>毫秒级代码评测，支持多种编程语言</p>
        </div>
      </div>
    </div>

    <div class="stats-card">
      <h3>平台统计</h3>
      <div class="stats">
        
        <div class="stat">
          <span class="number blue">10,000+</span>
          <span class="label">题目数量</span>
        </div>
        
        <div class="stat">
          <span class="number green">50,000+</span>
          <span class="label">注册用户</span>
        </div>
        
        <div class="stat">
          <span class="number purple">1M+</span>
          <span class="label">代码提交</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import logo from '../../../../../../assets/cascade.svg';
</script>

<style scoped>
.left-section {
  flex: 1; /* 左侧固定宽度视觉 */
  width: 550px;
  padding: 40px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  background-color: #0c1110; /* 深色背景与主页一致 */
}

.logo-area {
  display: flex;
  align-items: center;
  margin-bottom: 24px;
}

.logo {
  width: 40px;
  height: 40px;
  background-color: transparent;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 12px;
}

.logo svg {
  width: 24px;
  height: 24px;
}

.logo-img {
  width: 40px;
  height: 40px;
  display: block;
  object-fit: contain;
}

.logo-area h1 {
  font-size: 28px;
  color: #1dad80;
  margin: 0;
}

.slogan {
  font-size: 16px;
  color: #9fbebb;
  margin-bottom: 32px;
}

.features {
  display: flex;
  flex-direction: column;
  gap: 24px;
  margin-bottom: 32px;
}

.feature-item {
  display: flex;
  align-items: flex-start;
  gap: 16px;
}

.feature-icon {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.feature-icon.blue {
  background-color: #1dadad;
}

.feature-icon.green {
  background-color: #1dad8b;
}

.feature-icon.purple {
  background-color: #b06fd3;
}

.icon {
  width: 24px;
  height: 24px;
}

.feature-text h3 {
  font-size: 18px;
  color: #50e0b0e3;
  margin: 0 0 8px 0;
}

.feature-text p {
  font-size: 14px;
  color: #9fbebb;
  margin: 0;
}

.stats-card {
  background-color: rgba(255,255,255,0.02);
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 4px 12px rgba(0,0,0,0.6);
  width: 100%;
  max-width: 800px;
  align-self: flex-start;
}

.stats-card h3 {
  font-size: 18px;
  margin-bottom: 24px;
  margin-top: 0;
}

.stats {
  display: flex;
  justify-content: space-between;
}

.stat {
  text-align: center;
}

.number {
  font-size: 24px;
  font-weight: bold;
  display: block;
}

.number.blue {
  color: #1dadad;
}

.number.green {
  color: #1dad8b;
}

.number.purple {
  color: #b06fd3;
}

.label {
  font-size: 14px;
  color: #9fbebb;
}
</style>