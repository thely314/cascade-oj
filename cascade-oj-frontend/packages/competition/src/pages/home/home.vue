<template>
  <div class="home-container">
    <main class="main-content">
      <h1 class="hero-title">
        欢迎来到 Cascade 平台<br />
        <span class="highlight">The Art of Coding</span>
      </h1>
      <div class="hero-accent" />
      <p class="lead">这是一个在线判题系统</p>
      <div class="feature-cards">

        <div class="card card-blue info">
          <h3>比赛</h3>
          <p>加入比赛，展示你的编程实力</p>
        </div>

        <div class="card card-purple notice">
          <h3>排名</h3>
          <p>查看比赛中各位选手的排名情况</p>
        </div>

        <router-link class="card card-white card-link" to="/competition">
          <h3> 马上开始！</h3>
          <p>在此查看正在进行的比赛</p>
        </router-link>        

      </div>
      <Announcement />
    </main>
  </div>
</template>

<script setup lang="ts">
// NavBar is provided globally by App.vue
import Announcement from '../announcement/announcement.vue'
</script>

<style scoped src="./home.css"></style>