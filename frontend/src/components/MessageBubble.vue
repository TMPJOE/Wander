<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{
  message: {
    id: number;
    content: string;
    created_at: string;
    sender_id: number;
  };
  isMine: boolean;
}>();

const timeString = computed(() => {
  return new Date(props.message.created_at).toLocaleTimeString('es-MX', {
    hour: '2-digit', minute: '2-digit'
  });
});
</script>

<template>
  <div class="message-wrapper" :class="{ 'message-wrapper--mine': isMine }">
    <div class="message-bubble" :class="{ 'message-bubble--mine': isMine }">
      <p class="content">{{ message.content }}</p>
      <span class="time">{{ timeString }}</span>
    </div>
  </div>
</template>

<style scoped>
.message-wrapper {
  display: flex;
  width: 100%;
  margin-bottom: var(--spacing-2);
}

.message-wrapper--mine {
  justify-content: flex-end;
}

.message-bubble {
  max-width: 75%;
  padding: var(--spacing-2) var(--spacing-3);
  border-radius: var(--radius-lg);
  border-bottom-left-radius: 2px;
  background: var(--color-surface);
  border: 1px solid var(--color-border-light);
  box-shadow: var(--shadow-xs);
  display: flex;
  flex-direction: column;
}

.message-bubble--mine {
  background: var(--color-primary);
  color: white;
  border-color: var(--color-primary);
  border-bottom-left-radius: var(--radius-lg);
  border-bottom-right-radius: 2px;
}

.content {
  font-size: var(--font-size-sm);
  line-height: var(--line-height-relaxed);
  word-break: break-word;
}

.time {
  font-size: 0.65rem;
  color: var(--color-text-light);
  align-self: flex-end;
  margin-top: 2px;
}

.message-bubble--mine .time {
  color: rgba(255, 255, 255, 0.7);
}
</style>
