<script setup lang="ts">
import { useToast } from '../composables/useToast'

const { toasts, removeToast } = useToast()
</script>

<template>
  <div class="toast-container" aria-live="polite" aria-atomic="true">
    <TransitionGroup name="toast">
      <div
        v-for="t in toasts"
        :key="t.id"
        class="toast-item"
        :class="`toast-${t.type}`"
        @click="removeToast(t.id)"
      >
        <span class="toast-icon">
          <svg v-if="t.type === 'success'" xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"/></svg>
          <svg v-else-if="t.type === 'error'" xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><line x1="15" y1="9" x2="9" y2="15"/><line x1="9" y1="9" x2="15" y2="15"/></svg>
          <svg v-else xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><line x1="12" y1="16" x2="12" y2="12"/><line x1="12" y1="8" x2="12.01" y2="8"/></svg>
        </span>
        <span class="toast-message">{{ t.message }}</span>
        <button class="toast-close" title="Close toast" @click.stop="removeToast(t.id)">&times;</button>
      </div>
    </TransitionGroup>
  </div>
</template>

<style scoped>
.toast-container {
  position: fixed;
  top: 1.25rem;
  right: 1.25rem;
  z-index: 9999;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  max-width: 400px;
  width: calc(100vw - 2.5rem);
  pointer-events: none;
}

@media (max-width: 640px) {
  .toast-container {
    top: auto;
    bottom: calc(var(--bottom-nav-height, 65px) + 1rem);
    left: 50%;
    right: auto;
    transform: translateX(-50%);
  }
}

.toast-item {
  pointer-events: auto;
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.875rem 1rem;
  border-radius: 12px;
  background: var(--color-surface, #ffffff);
  color: var(--color-text, #1e293b);
  box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.12), 0 8px 10px -6px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(0, 0, 0, 0.06);
  cursor: pointer;
  backdrop-filter: blur(8px);
}

.toast-success {
  border-left: 4px solid #10b981;
}

.toast-success .toast-icon {
  color: #10b981;
}

.toast-error {
  border-left: 4px solid #ef4444;
}

.toast-error .toast-icon {
  color: #ef4444;
}

.toast-info {
  border-left: 4px solid #3b82f6;
}

.toast-info .toast-icon {
  color: #3b82f6;
}

.toast-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.toast-message {
  flex: 1;
  font-size: 0.875rem;
  font-weight: 500;
  line-height: 1.4;
}

.toast-close {
  background: none;
  border: none;
  font-size: 1.25rem;
  line-height: 1;
  color: #94a3b8;
  cursor: pointer;
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
}

.toast-close:hover {
  color: #475569;
}

/* Animations */
.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}

.toast-enter-from {
  opacity: 0;
  transform: translateY(-12px) scale(0.95);
}

.toast-leave-to {
  opacity: 0;
  transform: translateX(30px) scale(0.95);
}
</style>
