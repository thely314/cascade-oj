<script setup lang="ts">
withDefaults(defineProps<{
  title: string;
  show: boolean;
  cancelText?: string;
  submitText?: string;
  submitDisabled?: boolean;
}>(), {
  cancelText: '取消',
  submitText: '确定',
  submitDisabled: false,
});

defineEmits<{
  (e: 'close'): void;
  (e: 'submit'): void;
}>();
</script>

<template>
  <Transition name="modal">
    <div v-if="show" class="modal-mask" @click="$emit('close')">
      <div class="modal-container" @click.stop>
        <div class="modal-header">
          <h3>{{ title }}</h3>
          <button type="button" class="close-btn" @click="$emit('close')">&times;</button>
        </div>
        <div class="modal-body">
          <slot></slot>
        </div>
        <div class="modal-footer">
          <button type="button" class="ghost" @click="$emit('close')">{{ cancelText }}</button>
          <button type="button" class="primary" :disabled="submitDisabled" @click="$emit('submit')">{{ submitText }}</button>
        </div>
      </div>
    </div>
  </Transition>
</template>

<style scoped src="./Modal.css"></style>
