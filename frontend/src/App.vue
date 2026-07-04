<script setup lang="ts">
import { onMounted } from 'vue'
import { RouterView, useRoute } from 'vue-router'
import { useAuthStore } from './stores/auth'
import { useApi } from './composables/useApi'
import BottomNav from './components/BottomNav.vue'

const authStore = useAuthStore()
const api = useApi()
const route = useRoute()

const hideBottomNav = ['login', 'register']

onMounted(async () => {
  if (authStore.token && !authStore.user) {
    try {
      const response = await api.get('/users/me')
      authStore.user = response.data
    } catch {
      authStore.logout()
    }
  }
})
</script>

<template>
  <div class="app-shell">
    <main class="app-main">
      <RouterView v-slot="{ Component, route: viewRoute }">
        <Transition name="fade" mode="out-in">
          <component :is="Component" :key="viewRoute.path" />
        </Transition>
      </RouterView>
    </main>

    <BottomNav v-if="!hideBottomNav.includes(route.name as string)" />
  </div>
</template>

<style scoped>
.app-shell {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  min-height: 100dvh;
  width: 100%;
  max-width: var(--max-width);
  margin: 0 auto;
  position: relative;
  background: var(--color-background);
}

.app-main {
  flex: 1;
  display: flex;
  flex-direction: column;
}
</style>
