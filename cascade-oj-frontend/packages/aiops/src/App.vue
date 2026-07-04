<script setup lang="ts">
import { RouterView } from 'vue-router';
</script>

<template>
  <div class="content">
    <RouterView />
  </div>
</template>

<style>
:global(body) {
  margin: 0;
  background: #202020;
}

.content {
  flex: 1;
  padding: 24px 32px;
  overflow: auto;
  background: #202020;
}

</style>
