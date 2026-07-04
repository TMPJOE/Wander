<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useApi } from '../composables/useApi'
import { ArrowLeft, Check, X, Calendar } from '@lucide/vue'
import EmptyState from '../components/EmptyState.vue'

const router = useRouter()
const api = useApi()

const bookings = ref<any[]>([])
const loading = ref(true)
const filter = ref('all') // all, pending, confirmed, completed

onMounted(async () => {
  fetchBookings()
})

async function fetchBookings() {
  loading.value = true
  try {
    const res = await api.get('/bookings')
    bookings.value = res.data || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const filteredBookings = computed(() => {
  if (filter.value === 'all') return bookings.value
  return bookings.value.filter((b) => b.status === filter.value)
})

async function confirmBooking(id: number) {
  try {
    await api.patch(`/bookings/${id}/confirm`)
    await fetchBookings()
  } catch (e) {
    console.error(e)
    alert('Error al confirmar reserva')
  }
}

async function rejectBooking(id: number) {
  try {
    await api.patch(`/bookings/${id}/reject`)
    await fetchBookings()
  } catch (e) {
    console.error(e)
    alert('Error al rechazar reserva')
  }
}

async function completeBooking(id: number) {
  try {
    await api.patch(`/bookings/${id}/complete`)
    await fetchBookings()
  } catch (e) {
    console.error(e)
    alert('Error al marcar como completada')
  }
}

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString('es-MX', {
    weekday: 'short',
    day: 'numeric',
    month: 'short',
  })
}
function formatTime(dateStr: string) {
  return new Date(dateStr).toLocaleTimeString('es-MX', {
    hour: '2-digit',
    minute: '2-digit',
  })
}
</script>

<template>
  <div class="page bg-surface min-h-screen">
    <header class="header">
      <div class="flex items-center gap-3">
        <button class="back-btn" @click="router.push('/guide/dashboard')">
          <ArrowLeft :size="20" />
        </button>
        <h1 class="title">Reservas Recibidas</h1>
      </div>
    </header>

    <div class="container">
      <!-- Filters -->
      <div class="filters mb-4">
        <select v-model="filter" class="form-input">
          <option value="all">Todas las reservas</option>
          <option value="pending">Pendientes</option>
          <option value="confirmed">Confirmadas</option>
          <option value="completed">Completadas</option>
          <option value="cancelled">Canceladas</option>
        </select>
      </div>

      <div v-if="loading" class="flex flex-col gap-4">
        <div v-for="i in 3" :key="i" class="skeleton h-32 rounded-lg"></div>
      </div>

      <div v-else-if="filteredBookings.length" class="flex flex-col gap-4">
        <div v-for="booking in filteredBookings" :key="booking.id" class="card">
          <div class="flex justify-between items-start mb-2 border-b pb-2">
            <div>
              <span
                class="badge"
                :class="
                  booking.status === 'confirmed'
                    ? 'badge-success'
                    : booking.status === 'pending'
                      ? 'badge-warning'
                      : 'badge-secondary'
                "
              >
                {{ booking.status }}
              </span>
              <p class="text-xs text-light mt-1">Reserva #{{ booking.id }}</p>
            </div>
            <div class="text-right">
              <p class="font-bold text-primary">
                ${{ booking.total_price.toLocaleString('es-MX') }}
              </p>
              <p class="text-xs text-secondary">
                {{ booking.guest_count }} persona{{ booking.guest_count > 1 ? 's' : '' }}
              </p>
            </div>
          </div>

          <div class="py-2">
            <h3 class="font-semibold text-sm">{{ booking.tour_title }}</h3>
            <div class="flex items-center gap-2 mt-1 text-sm text-secondary">
              <Calendar :size="14" />
              <span
                >{{ formatDate(booking.schedule_start) }} a las
                {{ formatTime(booking.schedule_start) }}</span
              >
            </div>
          </div>

          <div class="flex items-center gap-2 mt-3 pt-3 border-t">
            <!-- Traveler info placeholder -->
            <div class="flex-1 flex items-center gap-2">
              <div
                class="w-8 h-8 rounded-full bg-secondary-50 text-secondary flex items-center justify-center font-semibold text-xs"
              >
                V
              </div>
              <span class="text-sm font-medium">
                Viajero: {{ booking.user_name || `#${booking.user_id}` }}
              </span>
            </div>

            <div class="flex gap-2" v-if="booking.status === 'pending'">
              <button
                class="btn btn-error btn-xs flex items-center gap-1"
                @click="rejectBooking(booking.id)"
              >
                <X :size="14" /> Rechazar
              </button>
              <button
                class="btn btn-success btn-xs flex items-center gap-1"
                @click="confirmBooking(booking.id)"
              >
                <Check :size="14" /> Confirmar
              </button>
            </div>
            <div v-else-if="booking.status === 'confirmed'">
              <button
                class="btn btn-primary btn-xs flex items-center gap-1"
                @click="completeBooking(booking.id)"
              >
                <Check :size="14" /> Marcar completado
              </button>
            </div>
          </div>
        </div>
      </div>

      <EmptyState
        v-else
        :icon="Calendar"
        title="No hay reservas"
        message="No se encontraron reservas con los filtros seleccionados."
      />
    </div>
  </div>
</template>

<style scoped>
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-4);
  background: var(--color-surface);
  border-bottom: 1px solid var(--color-border-light);
  position: sticky;
  top: 0;
  z-index: 10;
}

.container {
  padding: var(--spacing-4);
}

.card {
  background: var(--color-surface);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-lg);
  padding: var(--spacing-4);
}

.back-btn {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-full);
  background: var(--color-background);
}

.title {
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-bold);
}

.bg-surface {
  background: var(--color-surface);
}
.min-h-screen {
  min-height: 100vh;
}
.py-4 {
  padding-top: var(--spacing-4);
  padding-bottom: var(--spacing-4);
}
.py-2 {
  padding-top: var(--spacing-2);
  padding-bottom: var(--spacing-2);
}
.pb-2 {
  padding-bottom: var(--spacing-2);
}
.pt-3 {
  padding-top: var(--spacing-3);
}
.mb-4 {
  margin-bottom: var(--spacing-4);
}
.mb-2 {
  margin-bottom: var(--spacing-2);
}
.mt-1 {
  margin-top: 2px;
}
.mt-3 {
  margin-top: var(--spacing-3);
}

.flex {
  display: flex;
}
.flex-col {
  flex-direction: column;
}
.flex-1 {
  flex: 1;
}
.items-center {
  align-items: center;
}
.items-start {
  align-items: flex-start;
}
.justify-between {
  justify-content: space-between;
}
.justify-center {
  justify-content: center;
}
.gap-3 {
  gap: var(--spacing-3);
}
.gap-4 {
  gap: var(--spacing-4);
}
.gap-2 {
  gap: var(--spacing-2);
}
.gap-1 {
  gap: var(--spacing-1);
}

.text-right {
  text-align: right;
}
.font-bold {
  font-weight: var(--font-weight-bold);
}
.font-semibold {
  font-weight: var(--font-weight-semibold);
}
.font-medium {
  font-weight: var(--font-weight-medium);
}
.text-xs {
  font-size: var(--font-size-xs);
}
.text-sm {
  font-size: var(--font-size-sm);
}
.text-lg {
  font-size: var(--font-size-lg);
}

.text-primary {
  color: var(--color-primary);
}
.text-secondary {
  color: var(--color-text-secondary);
}
.text-light {
  color: var(--color-text-light);
}
.bg-secondary-50 {
  background: var(--color-secondary-50);
}

.h-32 {
  height: 8rem;
}
.w-8 {
  width: 2rem;
}
.h-8 {
  height: 2rem;
}
.rounded-lg {
  border-radius: var(--radius-lg);
}
.rounded-full {
  border-radius: var(--radius-full);
}

.border-b {
  border-bottom: 1px solid var(--color-border-light);
}
.border-t {
  border-top: 1px solid var(--color-border-light);
}

.btn-error {
  background: var(--color-error);
  color: white;
  border-color: var(--color-error);
}
.btn-success {
  background: var(--color-success);
  color: white;
  border-color: var(--color-success);
}
</style>
