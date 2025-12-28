<template>
  <section class="announcements">
    <div class="announcements-inner">
      <h2>目前可以公开的情报...</h2>
      <ul>
        <li v-for="(item, idx) in items" :key="idx">
          <time>{{ item.date }}</time>
          <div class="announce-body">
            <strong>{{ item.title }}</strong>
            <p>{{ item.content }}</p>
          </div>
        </li>
      </ul>
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getAnnouncements, type Announcement } from '../../api/announcement'

const items = ref<Array<Announcement & { date?: string }>>([])
const loading = ref(false)
const error = ref<string | null>(null)

onMounted(async () => {
  loading.value = true
  try {
    const res = await getAnnouncements()
    // 后端目前无日期字段，这里用占位或从内容中解析；若后端增加 createdAt 可直接映射
    const dynamic = res.announcements.map(a => ({ ...a, date: '' }))
    items.value = dynamic
  } catch (e: any) {
    error.value = e?.message ?? '加载公告失败'
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.announcements {
  max-width: 1200px;
  margin: 36px auto 0;
  padding: 0 16px;
}
.announcements-inner {
  background: #0f1413;
  border-radius: 8px;
  padding: 18px;
  border: 1px solid rgba(255,255,255,0.03);
  box-shadow: 0 6px 14px rgba(0,0,0,0.5);
}
.announcements h2 {
  color: #1dad80;
  margin: 0 0 12px;
  font-size: 1.25rem;
}
.announcements ul {
  list-style: none;
  margin: 0;
  padding: 0;
  display: grid;
  gap: 10px;
}
.announcements li {
  display: flex;
  gap: 12px;
  align-items: flex-start;
  padding: 12px;
  border-radius: 6px;
  background: rgba(255,255,255,0.01);
}
.announcements time {
  color: #9fb5ab;
  font-size: 0.85rem;
  min-width: 84px;
}
.announce-body { text-align: left; }
.announce-body strong { color: #cfeee0; display:block; text-align: left; }
.announcements p { margin: 6px 0 0; color: #9fb5ab; font-size: 0.95rem; }

@media (max-width: 600px) {
  .announcements li { flex-direction: column; }
  .announcements time { min-width: auto; }
}
</style>
