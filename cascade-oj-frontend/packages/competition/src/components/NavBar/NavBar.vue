<template>
  <nav class="navbar">
    <div class="nav-container">
      <!-- 网站Logo -->
      <div class="nav-logo">
        <span>CasCade</span>
      </div>
      
      <!-- 导航菜单 -->
      <ul class="nav-menu">
        <li class="nav-item">
          <span class="nav-link">题库</span>
        </li>
        <li class="nav-item">
          <span class="nav-link">比赛</span>
        </li>
        <li class="nav-item">
          <span class="nav-link">排名</span>
        </li>
        <li class="nav-item">
          <span class="nav-link">讨论</span>
        </li>
      </ul>
      
      <!-- 用户功能区 -->
      <div class="nav-user">
        <span class="login-btn">登录</span>
        <span class="register-btn">注册</span>
      </div>
    </div>
  </nav>
</template>

<script setup>
// 暂时不添加任何功能
</script>

<style scoped src="./NavBar.css"></style>