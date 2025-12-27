<!-- src/components/auth/AuthLeft.vue -->
<template>
  <div class="left-section">
    <div class="content-container">
      <!-- Logo 区 -->
      <div class="brand-container">
        <div class="logo-wrapper">
          <img :src="logo" alt="Cascade Logo" class="logo-img" />
        </div>
        <h1 class="brand-name">Cascade OJ</h1>
      </div>

      <!-- Slogan 区 -->
      <div class="text-content">
        <h2 class="main-heading">
          编程艺术的殿堂<br>
          <span class="highlight">The Art of Coding</span>
        </h2>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import logo from '../../../../../../assets/cascade.svg';
</script>

<style scoped>
.left-section {
  position: relative;
  padding: 0 80px;
  /* 移除之前的背景图片设置，因为父组件已经有了 */
  height: 100%;
  display: flex;
  align-items: center; /* 垂直居中 */
}

.content-container {
  /* 可以微调这个位置，比如 margin-bottom 让它稍微往上提一点 */
  margin-bottom: 100px;
}

.brand-container {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 40px;
}

.logo-img {
  width: 48px;
  height: 48px;
}

.brand-name {
  font-size: 26px;
  font-weight: 700;
  color: #fff;
  letter-spacing: 0.5px;
  margin: 0;
}

.main-heading {
  font-size: 56px; /* 字体稍微加大一点，因为现在它是主角 */
  line-height: 1.2;
  color: #fff;
  margin: 0;
  font-weight: 800;
}

.main-heading .highlight {
  font-family: 'Georgia', serif;
  font-style: italic;
  color: #23aa8f;
  font-weight: 400;
  display: block; /* 换行显示 */
  margin-top: 8px;
}
</style>