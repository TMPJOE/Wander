<script setup lang="ts">
import { onMounted } from 'vue'
import { useBookingsStore } from '../stores/bookings'
import { useApi } from '../composables/useApi'
import BookingCard from '../components/BookingCard.vue'
import EmptyState from '../components/EmptyState.vue'
import { CalendarDays } from '@lucide/vue'
import { useRouter } from 'vue-router'

const bookingsStore = useBookingsStore()
const api = useApi()
const router = useRouter()

onMounted(async () => {
  await bookingsStore.fetchMyBookings()
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
  <div class="page">
    <div class="header">
      <h1 class="title">Mis Reservas</h1>
    </div>

    <div class="container py-4">
      <div v-if="bookingsStore.loading" class="flex flex-col gap-4">
        <div v-for="i in 3" :key="i" class="skeleton h-32 rounded-lg"></div>
      </div>

      <div v-else-if="bookingsStore.bookings.length" class="bookings-list">
        <BookingCard
          v-for="booking in bookingsStore.bookings"
          :key="booking.id"
          :booking="booking"
          @cancel="cancelBooking"
        />
      </div>

      <EmptyState
        v-else
        :icon="CalendarDays"
        title="Sin reservas"
        message="Aún no has reservado ningún tour. ¡Explora las opciones y planifica tu próxima aventura!"
      >
        <button class="btn btn-primary" @click="router.push('/')">Explorar tours</button>
      </EmptyState>
    </div>
  </div>
</template>

<style scoped>
.header {
  padding: var(--spacing-6) var(--spacing-4) var(--spacing-4);
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
  padding: var(--spacing-4);
}
</style>
