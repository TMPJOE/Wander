<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import { useBookingsStore } from '../stores/bookings'
import { useApi } from '../composables/useApi'
import BookingCard from '../components/BookingCard.vue'
import EmptyState from '../components/EmptyState.vue'
import { CalendarDays } from '@lucide/vue'
import { useRouter } from 'vue-router'

const bookingsStore = useBookingsStore()
const api = useApi()
const router = useRouter()

const currentFilter = ref<'upcoming' | 'past' | 'cancelled'>('upcoming')

onMounted(async () => {
  await bookingsStore.fetchMyBookings()
})

const filteredBookings = computed(() => {
  const now = new Date()
  return bookingsStore.bookings.filter(b => {
    const isCancelled = b.status === 'cancelled'
    const isCompleted = b.status === 'completed'
    const start = new Date(b.schedule_start)
    
    if (currentFilter.value === 'cancelled') return isCancelled
    if (currentFilter.value === 'past') return !isCancelled && (start < now || isCompleted)
    return !isCancelled && !isCompleted && start >= now
  })
})

async function cancelBooking(id: number) {
  if (!confirm('¿Estás seguro de que deseas cancelar esta reserva?')) return
  try {
    await api.patch(`/bookings/${id}/cancel`)
    await bookingsStore.fetchMyBookings()
  } catch (e) {
    console.error(e)
    alert('Error al cancelar la reserva')
  }
}
</script>

<template>
  <div class="bookings-page bg-surface">
    <div class="header px-content">
      <h1 class="title">Mis Reservas</h1>
      <div class="filters">
        <button 
          class="filter-pill" 
          :class="{ active: currentFilter === 'upcoming' }"
          @click="currentFilter = 'upcoming'"
        >Próximas</button>
        <button 
          class="filter-pill" 
          :class="{ active: currentFilter === 'past' }"
          @click="currentFilter = 'past'"
        >Pasadas</button>
        <button 
          class="filter-pill" 
          :class="{ active: currentFilter === 'cancelled' }"
          @click="currentFilter = 'cancelled'"
        >Canceladas</button>
      </div>
    </div>

    <div class="px-content py-4">
      <div v-if="bookingsStore.loading" class="flex flex-col gap-4">
        <div v-for="i in 3" :key="i" class="skeleton h-32 rounded-lg"></div>
      </div>

      <div v-else-if="filteredBookings.length" class="bookings-list">
        <BookingCard
          v-for="booking in filteredBookings"
          :key="booking.id"
          :booking="booking"
          @cancel="cancelBooking"
        />
      </div>

      <EmptyState
        v-else-if="bookingsStore.bookings.length === 0"
        :icon="CalendarDays"
        title="Sin reservas"
        message="Aún no has reservado ningún tour. ¡Explora las opciones y planifica tu próxima aventura!"
      >
        <button class="btn btn-primary" @click="router.push('/')">Explorar tours</button>
      </EmptyState>

      <EmptyState
        v-else
        :icon="CalendarDays"
        title="No hay resultados"
        message="No tienes reservas en esta categoría."
      />
    </div>
  </div>
</template>

<style scoped>
.bookings-page {
  flex: 1;
  width: 100%;
  min-height: 100vh;
  min-height: 100dvh;
  padding-bottom: calc(var(--bottom-nav-height) + var(--spacing-4));
}

.px-content {
  padding-left: var(--content-padding);
  padding-right: var(--content-padding);
}

.header {
  padding-top: var(--spacing-6);
  padding-bottom: var(--spacing-2);
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

.filters {
  display: flex;
  gap: var(--spacing-3);
  margin-top: var(--spacing-4);
  overflow-x: auto;
  padding-bottom: var(--spacing-2);
}

.filters::-webkit-scrollbar {
  display: none;
}
.filters {
  -ms-overflow-style: none;  /* IE and Edge */
  scrollbar-width: none;  /* Firefox */
}

.filter-pill {
  padding: 8px 18px;
  border-radius: 20px;
  font-size: 13px;
  font-weight: 700;
  border: none;
  background: var(--color-primary-50);
  color: var(--color-primary);
  cursor: pointer;
  transition: all 0.2s ease;
  white-space: nowrap;
}

.filter-pill:hover {
  background: var(--color-primary-100);
}

.filter-pill.active {
  background: var(--color-primary);
  color: var(--color-text-inverse);
}

.py-4 {
  padding-top: var(--spacing-4);
  padding-bottom: var(--spacing-4);
}
.gap-4 {
  gap: var(--spacing-4);
}
.h-32 {
  height: 8rem;
}
.rounded-lg {
  border-radius: var(--radius-lg);
}

.bookings-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-4);
  padding: var(--spacing-4) 0;
}
</style>
