<template>
  <nav class="navbar">
    <div class="nav-container">
      <!-- 网站Logo -->
      <div class="nav-logo">
        <span class="logo-text">CasCade</span>
      </div>
      
      <!-- 导航菜单 -->
      <ul class="nav-menu">
        <li class="nav-item">
          <router-link class="nav-link" to="/home">主页</router-link>
        </li>
        <li class="nav-item">
          <router-link class="nav-link" to="/competition">比赛</router-link>
        </li>
        <li class="nav-item">
          <router-link class="nav-link" to="/about">关于我们</router-link>
        </li>
      </ul>
      
      <!-- 用户功能区 -->
      <div class="nav-user">
        <router-link class="login-btn" to="/login">登录/注册</router-link>
      </div>
    </div>
  </nav>
</template>

<script setup>
// 暂时不添加任何功能
</script>

<style scoped src="./NavBar.css"></style>