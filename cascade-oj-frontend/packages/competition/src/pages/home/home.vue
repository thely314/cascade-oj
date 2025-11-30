<template>
  <div class="home-container">
    <main class="main-content">
      <h1>欢迎来到CasCade平台</h1>
      <p>这是一个在线判题系统</p>
      <div class="feature-cards">
        <div class="card card-blue">
          <h3>题库</h3>
          <p>海量编程题目供您练习</p>
        </div>
        <router-link class="card card-white card-link" to="/competition">
          <h3>比赛</h3>
          <p>参加编程竞赛提升技能</p>
        </router-link>
        <div class="card card-purple">
          <h3>排名</h3>
          <p>查看您在社区中的排名</p>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
// NavBar is provided globally by App.vue
</script>
<style scoped src="./home.css"></style>