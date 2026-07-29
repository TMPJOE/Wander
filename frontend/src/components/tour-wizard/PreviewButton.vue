<script setup lang="ts">
import { defineProps, defineEmits } from 'vue'

const props = defineProps<{
  label: string
}>()

const emit = defineEmits<{
  click: []
}>()
</script>

<template>
  <button type="button" class="preview-button" @click="emit('click')">
    {{ props.label }}
  </button>
</template>

<style scoped>
.preview-button {
  display: flex;
  align-items: center;
  gap: var(--spacing-1);
  padding: var(--spacing-2) var(--spacing-3);
  border-radius: var(--radius-full);
  background: var(--color-secondary-50);
  color: var(--color-secondary);
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  border: 1px solid var(--color-secondary-100);
  transition: all var(--transition-fast);
  cursor: pointer;
}

.preview-button:hover {
  background: var(--color-secondary-100);
  transform: translateY(-1px);
  box-shadow: var(--shadow-sm);
}
</style>
