<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { ArrowLeft, Send } from '@lucide/vue';
import { useMessagesStore } from '../stores/messages';
import { useAuthStore } from '../stores/auth';
import { useApi } from '../composables/useApi';
import MessageBubble from '../components/MessageBubble.vue';

const route = useRoute();
const router = useRouter();
const messagesStore = useMessagesStore();
const authStore = useAuthStore();
const api = useApi();

const userId = computed(() => route.params.userId as string);
const newMessage = ref('');
const messagesContainer = ref<HTMLElement | null>(null);
const otherUser = ref<any>(null);

let pollInterval: any;

onMounted(async () => {
  await fetchOtherUser();
  await messagesStore.fetchMessages(userId.value);
  scrollToBottom();
  
  // Basic polling for new messages every 5 seconds
  pollInterval = setInterval(async () => {
    await messagesStore.fetchMessages(userId.value);
  }, 5000);
});

onUnmounted(() => {
  clearInterval(pollInterval);
});

async function fetchOtherUser() {
  try {
    const res = await api.get(`/users/${userId.value}`);
    otherUser.value = res.data;
  } catch (e) {
    console.error(e);
  }
}

async function sendMessage() {
  if (!newMessage.value.trim()) return;
  const content = newMessage.value;
  newMessage.value = '';
  await messagesStore.sendMessage(userId.value, content);
  scrollToBottom();
}

function scrollToBottom() {
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight;
    }
  });
}
</script>

<template>
  <div class="chat-page">
    <header class="chat-header">
      <button class="back-btn" @click="router.back()">
        <ArrowLeft :size="20" />
      </button>
      
      <div class="user-info" v-if="otherUser">
        <img 
          v-if="otherUser.avatar_url" 
          :src="otherUser.avatar_url" 
          class="avatar" 
        />
        <div v-else class="avatar avatar--placeholder">
          {{ otherUser.first_name[0] }}
        </div>
        <div>
          <h2 class="name">{{ otherUser.first_name }} {{ otherUser.last_name }}</h2>
          <span class="status">{{ otherUser.role === 'guide' ? 'Guía local' : 'Viajero' }}</span>
        </div>
      </div>
      <div v-else class="user-info">Cargando...</div>
    </header>

    <div class="chat-body" ref="messagesContainer">
      <div class="messages-list">
        <MessageBubble
          v-for="msg in messagesStore.currentMessages"
          :key="msg.id"
          :message="msg"
          :isMine="msg.sender_id === authStore.user?.id"
        />
      </div>
    </div>

    <footer class="chat-footer">
      <form class="compose" @submit.prevent="sendMessage">
        <input 
          v-model="newMessage" 
          type="text" 
          class="compose-input" 
          placeholder="Escribe un mensaje..."
        />
        <button 
          type="submit" 
          class="send-btn" 
          :disabled="!newMessage.trim()"
        >
          <Send :size="18" />
        </button>
      </form>
    </footer>
  </div>
</template>

<style scoped>
.chat-page {
  display: flex;
  flex-direction: column;
  height: 100vh;
  height: 100dvh;
  background: var(--color-background);
}

.chat-header {
  display: flex;
  align-items: center;
  gap: var(--spacing-3);
  padding: var(--spacing-3) var(--spacing-4);
  background: var(--color-surface);
  border-bottom: 1px solid var(--color-border-light);
  z-index: 10;
}

.back-btn {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-full);
  margin-left: -8px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: var(--spacing-3);
}

.avatar {
  width: 40px;
  height: 40px;
  border-radius: var(--radius-full);
  object-fit: cover;
}

.avatar--placeholder {
  background: var(--color-primary-50);
  color: var(--color-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: var(--font-weight-semibold);
}

.name {
  font-size: var(--font-size-base);
  font-weight: var(--font-weight-semibold);
  line-height: 1.2;
}

.status {
  font-size: var(--font-size-xs);
  color: var(--color-text-light);
}

.chat-body {
  flex: 1;
  overflow-y: auto;
  padding: var(--spacing-4);
}

.messages-list {
  display: flex;
  flex-direction: column;
}

.chat-footer {
  padding: var(--spacing-3) var(--spacing-4);
  background: var(--color-surface);
  border-top: 1px solid var(--color-border-light);
  padding-bottom: calc(var(--spacing-3) + env(safe-area-inset-bottom));
}

.compose {
  display: flex;
  align-items: center;
  gap: var(--spacing-2);
  background: var(--color-background);
  border-radius: var(--radius-full);
  padding: 4px;
  padding-left: var(--spacing-4);
  border: 1px solid var(--color-border);
}

.compose-input {
  flex: 1;
  border: none;
  background: transparent;
  padding: var(--spacing-2) 0;
  outline: none;
  font-size: var(--font-size-sm);
}

.send-btn {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-primary);
  color: white;
  border-radius: var(--radius-full);
  transition: opacity var(--transition-fast);
}

.send-btn:disabled {
  opacity: 0.5;
  background: var(--color-text-light);
}
</style>
