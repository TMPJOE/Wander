<script setup lang="ts">
import { computed } from 'vue';
import { useRouter } from 'vue-router';

const props = defineProps<{
  conversation: {
    user_id: number;
    user_name: string;
    user_avatar: string;
    last_message: string;
    last_at: string;
    unread_count: number;
  };
}>();

const router = useRouter();

const timeString = computed(() => {
  const d = new Date(props.conversation.last_at);
  const now = new Date();
  if (d.toDateString() === now.toDateString()) {
    return d.toLocaleTimeString('es-MX', { hour: '2-digit', minute: '2-digit' });
  }
  return d.toLocaleDateString('es-MX', { day: 'numeric', month: 'short' });
});
</script>

<template>
  <div class="conversation-item" @click="router.push(`/messages/${conversation.user_id}`)">
    <img 
      v-if="conversation.user_avatar" 
      :src="conversation.user_avatar" 
      :alt="conversation.user_name" 
      class="avatar" 
    />
    <div v-else class="avatar avatar--placeholder">
      {{ conversation.user_name[0] }}
    </div>
    
    <div class="info">
      <div class="header">
        <h3 class="name">{{ conversation.user_name }}</h3>
        <span class="time">{{ timeString }}</span>
      </div>
      <p class="last-msg" :class="{ 'last-msg--unread': conversation.unread_count > 0 }">
        {{ conversation.last_message }}
      </p>
    </div>
    
    <div v-if="conversation.unread_count > 0" class="badge-unread">
      {{ conversation.unread_count > 9 ? '9+' : conversation.unread_count }}
    </div>
  </div>
</template>

<style scoped>
.conversation-item {
  display: flex;
  align-items: center;
  padding: var(--spacing-4);
  gap: var(--spacing-3);
  background: var(--color-surface);
  border-bottom: 1px solid var(--color-border-light);
  cursor: pointer;
  transition: background var(--transition-fast);
}

.conversation-item:hover {
  background: var(--color-surface-hover);
}

.conversation-item:last-child {
  border-bottom: none;
}

.avatar {
  width: 48px;
  height: 48px;
  border-radius: var(--radius-full);
  object-fit: cover;
  flex-shrink: 0;
}

.avatar--placeholder {
  background: var(--color-primary-50);
  color: var(--color-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: var(--font-weight-semibold);
  font-size: var(--font-size-lg);
}

.info {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2px;
}

.name {
  font-size: var(--font-size-base);
  font-weight: var(--font-weight-semibold);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.time {
  font-size: var(--font-size-xs);
  color: var(--color-text-light);
  flex-shrink: 0;
}

.last-msg {
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.last-msg--unread {
  font-weight: var(--font-weight-semibold);
  color: var(--color-text);
}

.badge-unread {
  background: var(--color-primary);
  color: white;
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-bold);
  min-width: 20px;
  height: 20px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 6px;
  flex-shrink: 0;
}
</style>
