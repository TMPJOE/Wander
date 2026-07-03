<script setup lang="ts">
import { PackageOpen } from '@lucide/vue';

defineProps<{
  title?: string;
  message?: string;
  icon?: typeof PackageOpen;
}>();
</script>

<template>
  <div class="empty-state">
    <div class="empty-state__icon-wrap">
      <component :is="icon || PackageOpen" :size="48" :stroke-width="1.2" />
    </div>
    <h3 class="empty-state__title">{{ title || 'Nada por aquí' }}</h3>
    <p class="empty-state__message">{{ message || 'Aún no hay contenido para mostrar.' }}</p>
    <div class="empty-state__actions">
      <slot />
    </div>
  </div>
</template>

<style scoped>
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  padding: var(--spacing-12) var(--spacing-6);
  gap: var(--spacing-3);
}

.empty-state__icon-wrap {
  width: 80px;
  height: 80px;
  border-radius: var(--radius-full);
  background: var(--color-primary-50);
  color: var(--color-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: var(--spacing-2);
}

.empty-state__title {
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text);
}

.empty-state__message {
  font-size: var(--font-size-sm);
  color: var(--color-text-light);
  max-width: 280px;
  line-height: var(--line-height-relaxed);
}

.empty-state__actions {
  margin-top: var(--spacing-4);
}
</style>
