<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useApi } from '../composables/useApi'
import { Map, Calendar, DollarSign, Star, ChevronRight } from '@lucide/vue'

const authStore = useAuthStore()
const router = useRouter()
const api = useApi()

const stats = ref({
  total_bookings: 0,
  total_revenue: 0,
  active_tours: 0,
  average_rating: 0,
})

const recentBookings = ref<any[]>([])

onMounted(async () => {
  if (authStore.user?.role !== 'guide') {
    router.push('/')
    return
  }

  try {
    const statsRes = await api.get('/guide/stats')
    const data = statsRes.data || {}
    stats.value.total_bookings = data.total_bookings || 0
    stats.value.total_revenue = data.total_revenue || 0
    stats.value.active_tours = data.published_tours || 0
    stats.value.average_rating = data.avg_rating || 0

    const bookingsRes = await api.get('/bookings')
    recentBookings.value = (bookingsRes.data || []).slice(0, 5)
  } catch (e) {
    console.error(e)
  }
})
</script>

<template>
  <div class="page">
    <header class="header">
      <h1 class="title">Panel de Guía</h1>
      <p class="subtitle">Hola, {{ authStore.user?.first_name }}</p>
    </header>

    <div class="container py-4">
      <!-- Stats Grid -->
      <div class="stats-grid mb-6">
        <div class="stat-card">
          <div class="stat-icon bg-primary-light text-primary"><DollarSign :size="20" /></div>
          <div class="stat-info">
            <span class="stat-value">${{ stats.total_revenue.toLocaleString('es-MX') }}</span>
            <span class="stat-label">Ingresos</span>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon bg-success-light text-success"><Calendar :size="20" /></div>
          <div class="stat-info">
            <span class="stat-value">{{ stats.total_bookings }}</span>
            <span class="stat-label">Reservas</span>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon bg-warning-light text-warning"><Star :size="20" /></div>
          <div class="stat-info">
            <span class="stat-value">{{ stats.average_rating }}</span>
            <span class="stat-label">Rating</span>
          </div>
        </div>
        <div class="stat-card clickable-card" @click="router.push('/guide/tours')">
          <div class="stat-icon bg-secondary-light text-secondary"><Map :size="20" /></div>
          <div class="stat-info">
            <span class="stat-value">{{ stats.active_tours }}</span>
            <span class="stat-label">Tours activos</span>
          </div>
        </div>
      </div>

      <!-- Quick Actions -->
      <div class="mb-6">
        <h2 class="text-lg font-semibold mb-3">Acciones Rápidas</h2>
        <div class="grid grid-cols-2 gap-3">
          <button
            class="btn btn-outline flex-col h-auto py-4 gap-2"
            @click="router.push('/guide/tours/new')"
          >
            <Map :size="24" class="text-primary" />
            <span>Crear Tour</span>
          </button>
          <button
            class="btn btn-outline flex-col h-auto py-4 gap-2"
            @click="router.push('/guide/bookings')"
          >
            <Calendar :size="24" class="text-primary" />
            <span>Ver Reservas</span>
          </button>
        </div>
      </div>

      <!-- Recent Bookings -->
      <div>
        <div class="flex justify-between items-center mb-3">
          <h2 class="text-lg font-semibold">Reservas Recientes</h2>
          <button class="btn btn-ghost btn-sm" @click="router.push('/guide/bookings')">
            Ver todas
          </button>
        </div>

        <div class="card p-0 overflow-hidden">
          <div v-if="recentBookings.length === 0" class="p-4 text-center text-secondary text-sm">
            No tienes reservas recientes.
          </div>
          <div v-else>
            <div
              v-for="booking in recentBookings"
              :key="booking.id"
              class="flex items-center justify-between p-4 border-b last:border-b-0 cursor-pointer hover-bg"
              @click="router.push('/guide/bookings')"
            >
              <div>
                <p class="font-semibold text-sm">{{ booking.tour_title }}</p>
                <p class="text-xs text-secondary mt-1">
                  {{ new Date(booking.schedule_start).toLocaleDateString('es-MX') }} •
                  {{ booking.guest_count }} pax
                </p>
              </div>
              <div class="flex items-center gap-3">
                <span
                  class="badge"
                  :class="booking.status === 'confirmed' ? 'badge-success' : 'badge-warning'"
                >
                  {{ booking.status }}
                </span>
                <ChevronRight :size="16" class="text-light" />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.header {
  padding: var(--spacing-6) var(--spacing-4) var(--spacing-4);
  background: var(--color-primary-50);
}

.container {
  padding: var(--spacing-4);
}

.title {
  font-size: var(--font-size-2xl);
  font-weight: var(--font-weight-bold);
  letter-spacing: var(--letter-spacing-tight);
}

.subtitle {
  color: var(--color-text-secondary);
  margin-top: 2px;
}

.py-4 {
  padding-top: var(--spacing-4);
  padding-bottom: var(--spacing-4);
}
.mb-3 {
  margin-bottom: var(--spacing-3);
}
.mb-6 {
  margin-bottom: var(--spacing-6);
}
.mt-1 {
  margin-top: 2px;
}
.text-lg {
  font-size: var(--font-size-lg);
}
.text-sm {
  font-size: var(--font-size-sm);
}
.text-xs {
  font-size: var(--font-size-xs);
}
.font-semibold {
  font-weight: var(--font-weight-semibold);
}
.text-secondary {
  color: var(--color-text-secondary);
}
.text-light {
  color: var(--color-text-light);
}
.text-primary {
  color: var(--color-primary);
}
.text-success {
  color: var(--color-success);
}
.text-warning {
  color: var(--color-warning);
}

.bg-primary-light {
  background: var(--color-primary-50);
}
.bg-success-light {
  background: #dcfce7;
}
.bg-warning-light {
  background: #fef9c3;
}
.bg-secondary-light {
  background: var(--color-secondary-50);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: var(--spacing-3);
}

.stat-card {
  background: var(--color-surface);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-lg);
  padding: var(--spacing-3);
  display: flex;
  align-items: center;
  gap: var(--spacing-3);
  box-shadow: var(--shadow-xs);
}

.clickable-card {
  cursor: pointer;
  background: var(--color-primary-100);
  border: 2px solid var(--color-primary);
}

.stat-icon {
  width: 40px;
  height: 40px;
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.stat-info {
  display: flex;
  flex-direction: column;
}

.stat-value {
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-bold);
  line-height: 1.2;
}

.stat-label {
  font-size: var(--font-size-xs);
  color: var(--color-text-secondary);
}

.grid {
  display: grid;
}
.grid-cols-2 {
  grid-template-columns: repeat(2, 1fr);
}
.gap-3 {
  gap: var(--spacing-3);
}
.gap-2 {
  gap: var(--spacing-2);
}
.flex-col {
  flex-direction: column;
}
.h-auto {
  height: auto;
}
.py-4 {
  padding-top: var(--spacing-4);
  padding-bottom: var(--spacing-4);
}
.p-0 {
  padding: 0;
}
.p-4 {
  padding: var(--spacing-4);
}
.flex {
  display: flex;
}
.justify-between {
  justify-content: space-between;
}
.items-center {
  align-items: center;
}
.border-b {
  border-bottom: 1px solid var(--color-border-light);
}
.last\:border-b-0:last-child {
  border-bottom: none;
}
.cursor-pointer {
  cursor: pointer;
}
.hover-bg:hover {
  background: var(--color-background);
}
.overflow-hidden {
  overflow: hidden;
}
</style>
