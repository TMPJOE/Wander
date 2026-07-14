<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import { useApi } from '../composables/useApi'
import ConversationItem from '../components/ConversationItem.vue'
import EmptyState from '../components/EmptyState.vue'
import { MessageCircle } from '@lucide/vue'
import { useRouter } from 'vue-router'

const api = useApi()
const router = useRouter()

const conversations = ref<any[]>([])
const loading = ref(false)
let pollInterval: any;

onMounted(async () => {
  fetchConversations()
  
  pollInterval = setInterval(() => {
    fetchConversations(true)
  }, 5000)
})

async function fetchConversations(silent = false) {
  if (!silent) {
    loading.value = true
  }
  try {
    const res = await api.get('/messages/conversations')
    conversations.value = res.data || []
  } catch (e) {
    console.error('Failed to load conversations', e)
  } finally {
    if (!silent) {
      loading.value = false
    }
  }
}

onUnmounted(() => {
  if (pollInterval) clearInterval(pollInterval)
})
</script>

<template>
  <div class="messages-page bg-surface">
    <div class="header" style="padding-left: var(--content-padding); padding-right: var(--content-padding);">
      <h1 class="title">Mensajes</h1>
    </div>

    <div>
      <div v-if="loading" class="flex flex-col">
        <div v-for="i in 5" :key="i" class="p-4 border-b flex items-center gap-3" style="padding-left: var(--content-padding); padding-right: var(--content-padding);">
          <div class="skeleton w-12 h-12 rounded-full"></div>
          <div class="flex-1">
            <div class="skeleton h-4 w-1/3 mb-2"></div>
            <div class="skeleton h-3 w-2/3"></div>
          </div>
        </div>
      </div>

      <div v-else-if="conversations.length" class="conversations-list">
        <ConversationItem
          v-for="conv in conversations"
          :key="conv.user_id"
          :conversation="conv"
          style="padding-left: var(--content-padding); padding-right: var(--content-padding);"
        />
      </div>

      <div class="py-8" v-else>
        <EmptyState
          :icon="MessageCircle"
          title="Sin mensajes"
          message="No tienes conversaciones activas. Cuando reserves un tour o contactes a un guía, tus mensajes aparecerán aquí."
        >
          <button class="btn btn-primary" @click="router.push('/')">Explorar tours</button>
        </EmptyState>
      </div>
    </div>
  </div>
</template>

<style scoped>
.messages-page {
  flex: 1;
  width: 100%;
  min-height: 100vh;
  min-height: 100dvh;
}

.header {
  padding-top: var(--spacing-6);
  padding-bottom: var(--spacing-4);
  background: var(--color-surface);
  border-bottom: 1px solid var(--color-border-light);
  position: sticky;
  top: 0;
  z-index: 10;
}

.title {
  font-size: var(--font-size-2xl);
  font-weight: var(--font-weight-bold);
  letter-spacing: var(--letter-spacing-tight);
}

.container-fluid {
  width: 100%;
}
.px-0 {
  padding-left: 0;
  padding-right: 0;
}
.py-8 {
  padding-top: var(--spacing-8);
  padding-bottom: var(--spacing-8);
}
.p-4 {
  padding: var(--spacing-4);
}
.border-b {
  border-bottom: 1px solid var(--color-border-light);
}
.gap-3 {
  gap: var(--spacing-3);
}
.w-12 {
  width: 3rem;
}
.h-12 {
  height: 3rem;
}
.rounded-full {
  border-radius: var(--radius-full);
}
.h-4 {
  height: 1rem;
}
.h-3 {
  height: 0.75rem;
}
.w-1\/3 {
  width: 33.333333%;
}
.w-2\/3 {
  width: 66.666667%;
}
.mb-2 {
  margin-bottom: var(--spacing-2);
}
</style>
