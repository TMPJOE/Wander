<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { Compass, CalendarDays, MessageCircle, User, LayoutDashboard } from '@lucide/vue'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

interface NavItem {
  name: string
  icon: typeof Compass
  route: string
  label: string
  requiresAuth?: boolean
  requiresGuide?: boolean
}

const tabs = computed<NavItem[]>(() => {
  const items: NavItem[] = [
    { name: 'explore', icon: Compass, route: '/', label: 'Explorar' },
    {
      name: 'my-bookings',
      icon: CalendarDays,
      route: '/bookings',
      label: 'Reservas',
      requiresAuth: true,
    },
    {
      name: 'messages',
      icon: MessageCircle,
      route: '/messages',
      label: 'Mensajes',
      requiresAuth: true,
    },
    { name: 'profile', icon: User, route: '/profile', label: 'Perfil', requiresAuth: true },
  ]

  if (authStore.isGuide) {
    items.splice(1, 0, {
      name: 'guide-dashboard',
      icon: LayoutDashboard,
      route: '/guide/dashboard',
      label: 'Panel',
      requiresGuide: true,
    })
  }

  return items
})

function isActive(tab: NavItem): boolean {
  if (tab.route === '/') return route.path === '/'
  return route.path.startsWith(tab.route)
}

function navigate(tab: NavItem) {
  if (tab.requiresAuth && !authStore.isAuthenticated) {
    router.push({ name: 'login', query: { redirect: tab.route } })
    return
  }
  router.push(tab.route)
}
</script>

<template>
  <nav class="bottom-nav">
    <button
      v-for="tab in tabs"
      :key="tab.name"
      class="bottom-nav__item"
      :class="{ 'bottom-nav__item--active': isActive(tab) }"
      @click="navigate(tab)"
    >
      <component :is="tab.icon" :size="22" :stroke-width="isActive(tab) ? 2.5 : 1.8" />
      <span class="bottom-nav__label">{{ tab.label }}</span>
    </button>
  </nav>
</template>

<style scoped>
.bottom-nav {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  height: var(--bottom-nav-height);
  background: rgba(var(--color-primary-rgb), 0.08);
  border-top: 1px solid var(--color-primary-100);
  display: flex;
  align-items: center;
  justify-content: space-around;
  padding: 0 var(--spacing-2);
  z-index: var(--z-sticky);
  backdrop-filter: blur(12px);
}

.bottom-nav__item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  padding: var(--spacing-1) var(--spacing-3);
  border-radius: var(--radius-lg);
  color: var(--color-text-light);
  transition: all var(--transition-fast);
  -webkit-tap-highlight-color: transparent;
  position: relative;
}

.bottom-nav__item::before {
  content: '';
  position: absolute;
  top: -1px;
  left: 50%;
  transform: translateX(-50%) scaleX(0);
  width: 24px;
  height: 2px;
  background: var(--color-primary);
  border-radius: var(--radius-full);
  transition: transform var(--transition-spring);
}

.bottom-nav__item--active {
  color: var(--color-primary);
}

.bottom-nav__item--active::before {
  transform: translateX(-50%) scaleX(1);
}

.bottom-nav__label {
  font-size: 0.625rem;
  font-weight: var(--font-weight-medium);
  letter-spacing: var(--letter-spacing-wide);
}

.bottom-nav__item--active .bottom-nav__label {
  font-weight: var(--font-weight-bold);
}

/* ─── Desktop: convert the fixed bottom bar into a left sidebar ─── */
@media (min-width: 1024px) {
  .bottom-nav {
    border-left: 1px solid var(--color-border-light);
    position: sticky;
    top: 0;
    align-self: flex-start;
    left: auto;
    right: auto;
    bottom: auto;
    width: var(--nav-width);
    height: 100vh;
    height: 100dvh;
    flex-direction: column;
    align-items: stretch;
    justify-content: flex-start;
    gap: var(--spacing-1);
    padding: var(--spacing-6) var(--spacing-3);
    border-top: none;
    border-right: 1px solid var(--color-border-light);
    background: rgba(255, 255, 255, 0.92);
    backdrop-filter: blur(12px);
    z-index: var(--z-sticky);
  }

  .bottom-nav__item {
    flex-direction: row;
    justify-content: flex-start;
    gap: var(--spacing-3);
    padding: var(--spacing-3);
    width: 100%;
  }

  /* Active indicator switches from a top underline to a left edge bar. */
  .bottom-nav__item::before {
    top: 50%;
    left: 0;
    width: 3px;
    height: 60%;
    transform: translateY(-50%) scaleX(0);
    border-radius: var(--radius-full);
  }

  .bottom-nav__item--active::before {
    transform: translateY(-50%) scaleX(1);
  }

  .bottom-nav__item--active {
    background: var(--color-primary-50);
  }

  .bottom-nav__label {
    font-size: var(--font-size-sm);
    font-weight: var(--font-weight-medium);
  }
}
</style>
