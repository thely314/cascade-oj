<!-- src/App.vue -->
<template>
  <router-view />
</template>

<style>
/* 全局重置样式，建议放在这里或 assets/base.css */
html, body {
  margin: 0;
  padding: 0;
  width: 100%;
  height: 100%;
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
  background-color: #0c1110;
  color: #cbd5c0;
}

*, *::before, *::after {
  box-sizing: border-box;
}
</style>