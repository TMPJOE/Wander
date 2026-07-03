import { defineStore } from 'pinia';
import { ref } from 'vue';
import { useApi } from '../composables/useApi';

export const useMessagesStore = defineStore('messages', () => {
  const conversations = ref<any[]>([]);
  const currentMessages = ref<any[]>([]);
  const loading = ref(false);
  const api = useApi();

  async function fetchConversations() {
    loading.value = true;
    try {
      const response = await api.get('/messages/conversations');
      conversations.value = response.data;
    } catch (e) {
      console.error('Failed to fetch conversations', e);
    } finally {
      loading.value = false;
    }
  }

  async function fetchMessages(userId: string) {
    loading.value = true;
    try {
      const response = await api.get(`/messages/${userId}`);
      currentMessages.value = response.data;
    } catch (e) {
      console.error('Failed to fetch messages', e);
    } finally {
      loading.value = false;
    }
  }

  async function sendMessage(userId: string, content: string) {
    try {
      await api.post(`/messages/${userId}`, { content });
      await fetchMessages(userId);
    } catch (e) {
      console.error('Failed to send message', e);
    }
  }

  return { conversations, currentMessages, loading, fetchConversations, fetchMessages, sendMessage };
});
