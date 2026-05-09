<script setup lang="ts">
defineProps<{
  title: string
  show: boolean
}>()

defineEmits<{
  (e: 'close'): void
  (e: 'submit'): void
}>()
</script>

<template>
  <Transition name="modal">
    <div v-if="show" class="modal-mask" @click="$emit('close')">
      <div class="modal-container" @click.stop>
        <div class="modal-header">
          <h3>{{ title }}</h3>
          <button class="close-btn" @click="$emit('close')">&times;</button>
        </div>

        <div class="modal-body">
          <slot></slot>
        </div>

        <div class="modal-footer">
          <button class="ghost" @click="$emit('close')">取消</button>
          <button class="primary" @click="$emit('submit')">确定</button>
        </div>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
.modal-mask {
  position: fixed;
  z-index: 9998;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  transition: opacity 0.3s ease;
  align-items: center;
  justify-content: center;
}

.modal-container {
  width: 500px;
  max-width: 90vw;
  background-color: #141816;
  border: 1px solid rgba(22, 163, 118, 0.12);
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.32);
  transition: all 0.3s ease;
  display: flex;
  flex-direction: column;
}

.modal-header {
  padding: 16px 20px;
  border-bottom: 1px solid rgba(22, 163, 118, 0.08);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-header h3 {
  margin: 0;
  font-size: 1.1rem;
  color: #e6f7ec;
}

.close-btn {
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  color: #8aa39b;
  padding: 0;
  line-height: 1;
}

.close-btn:hover {
  color: #cbd5c0;
}

.modal-body {
  padding: 20px;
  max-height: 70vh;
  overflow-y: auto;
  box-sizing: border-box;
}

.modal-footer {
  padding: 16px 20px;
  border-top: 1px solid rgba(22, 163, 118, 0.08);
  background-color: #0f1412;
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

/* Button styles for modal */
:deep(.ghost) {
  padding: 8px 12px;
  border-radius: 8px;
  border: 1px solid rgba(22, 163, 118, 0.2);
  background: #0f1412;
  color: #cbd5c0;
  cursor: pointer;
  transition: background 0.12s ease, border-color 0.12s ease;
  font-family: inherit;
}

:deep(.ghost:hover) {
  background: rgba(22, 163, 118, 0.08);
  border-color: rgba(22, 163, 118, 0.35);
}

:deep(.primary) {
  padding: 10px 14px;
  border: none;
  border-radius: 10px;
  background: linear-gradient(135deg, #16a34a, #1dad80);
  color: #04130a;
  font-weight: 600;
  cursor: pointer;
  box-shadow: 0 8px 16px rgba(22, 163, 118, 0.25);
  transition: transform 0.12s ease, box-shadow 0.12s ease;
  font-family: inherit;
}

:deep(.primary:hover) {
  transform: translateY(-1px);
  box-shadow: 0 10px 18px rgba(22, 163, 118, 0.32);
}

.modal-enter-from {
  opacity: 0;
}

.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .modal-container,
.modal-leave-to .modal-container {
  transform: scale(1.1);
}

/* Common form styles for modal-body */
:deep(.form-group) {
  margin-bottom: 16px;
}

:deep(.form-group label) {
  display: block;
  margin-bottom: 6px;
  font-size: 0.9rem;
  color: #e6f7ec;
  font-weight: 600;
}

:deep(.form-group input),
:deep(.form-group textarea) {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid rgba(22, 163, 118, 0.2);
  border-radius: 6px;
  background: #0f1412;
  color: #cbd5c0;
  font-family: inherit;
  box-sizing: border-box;
}

:deep(.form-group input:focus),
:deep(.form-group textarea:focus) {
  outline: none;
  border-color: #16a379;
  box-shadow: 0 0 0 2px rgba(22, 163, 118, 0.1);
}
</style>
