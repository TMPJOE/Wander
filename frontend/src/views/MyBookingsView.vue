<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import { useApi } from '../composables/useApi'
import { useToast } from '../composables/useToast'
import { useConfirm } from '../composables/useConfirm'
import BookingCard from '../components/BookingCard.vue'
import EmptyState from '../components/EmptyState.vue'
import { CalendarDays, Star, Send, X } from '@lucide/vue'
import { useRouter } from 'vue-router'

const api = useApi()
const toast = useToast()
const { confirm } = useConfirm()
const router = useRouter()

const bookings = ref<any[]>([])
const loading = ref(false)
const currentFilter = ref<'upcoming' | 'past' | 'cancelled'>('upcoming')

const dismissedReviewBookings = ref<Set<number>>(new Set())
const reviewRatings = ref<Record<number, number>>({})
const reviewComments = ref<Record<number, string>>({})
const submittingReview = ref<Record<number, boolean>>({})

onMounted(async () => {
  fetchMyBookings()
})

async function fetchMyBookings() {
  loading.value = true
  try {
    const res = await api.get('/bookings')
    bookings.value = res.data || []
  } catch (e) {
    console.error('Failed to load bookings', e)
    toast.error('Error al cargar reservas')
  } finally {
    loading.value = false
  }
}

const filteredBookings = computed(() => {
  const now = new Date()
  return bookings.value.filter(b => {
    const isCancelled = b.status === 'cancelled'
    const isCompleted = b.status === 'completed'
    const start = new Date(b.schedule_start)
    
    if (currentFilter.value === 'cancelled') return isCancelled
    if (currentFilter.value === 'past') return !isCancelled && (start < now || isCompleted)
    return !isCancelled && !isCompleted && start >= now
  })
})

const unreviewedPastBookings = computed(() => {
  if (currentFilter.value !== 'past') return []
  return filteredBookings.value.filter(b => !dismissedReviewBookings.value.has(b.id))
})

function setRating(bookingId: number, rating: number) {
  reviewRatings.value[bookingId] = rating
}

function dismissReview(bookingId: number) {
  dismissedReviewBookings.value.add(bookingId)
}

async function submitReview(booking: any) {
  const rating = reviewRatings.value[booking.id] || 5
  const comment = reviewComments.value[booking.id] || ''

  submittingReview.value[booking.id] = true
  try {
    await api.post(`/tours/${booking.tour_id}/reviews`, {
      rating,
      comment,
      booking_id: booking.id
    })
    toast.success('¡Gracias por tu reseña!')
    dismissReview(booking.id)
  } catch (e: any) {
    console.error(e)
    toast.error(e.response?.data?.message || 'Error al enviar reseña')
  } finally {
    submittingReview.value[booking.id] = false
  }
}

async function cancelBooking(id: number) {
  const approved = await confirm({
    title: '¿Cancelar reserva?',
    body: '¿Estás seguro de que deseas cancelar esta reserva?',
    confirmLabel: 'Sí, cancelar',
    cancelLabel: 'No',
    confirmVariant: 'danger',
  })
  if (!approved) return

  try {
    const res = await api.patch(`/bookings/${id}/cancel`)
    const msg = res.data?.message || 'Reserva cancelada exitosamente'
    toast.success(msg)
    await fetchMyBookings()
  } catch (e) {
    console.error(e)
    toast.error('Error al cancelar la reserva')
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
      <div v-if="loading" class="flex flex-col gap-4">
        <div v-for="i in 3" :key="i" class="skeleton h-32 rounded-lg"></div>
      </div>

      <!-- Post-trip Review Prompt Banner -->
      <div v-if="currentFilter === 'past' && unreviewedPastBookings.length" class="review-prompt-section mb-4">
        <div v-for="booking in unreviewedPastBookings" :key="`review-${booking.id}`" class="review-prompt-card">
          <button class="dismiss-btn" title="Descartar" @click="dismissReview(booking.id)">
            <X :size="16" />
          </button>
          <div class="review-prompt-header">
            <Star :size="18" class="text-warning fill-warning" />
            <h3 class="text-sm font-bold text-dark">Califica tu experiencia</h3>
          </div>
          <p class="text-xs text-secondary mt-1">¿Qué tal estuvo tu tour "{{ booking.tour_title }}"?</p>

          <div class="star-rating my-2">
            <button
              v-for="star in 5"
              :key="star"
              type="button"
              class="star-btn"
              @click="setRating(booking.id, star)"
            >
              <Star
                :size="20"
                :class="(reviewRatings[booking.id] || 5) >= star ? 'text-warning fill-warning' : 'text-muted'"
              />
            </button>
          </div>

          <div class="flex gap-2">
            <input
              v-model="reviewComments[booking.id]"
              type="text"
              class="form-input text-xs flex-1"
              placeholder="Escribe un comentario breve (opcional)..."
            />
            <button
              class="btn btn-primary btn-sm flex items-center gap-1"
              :disabled="submittingReview[booking.id]"
              @click="submitReview(booking)"
            >
              <Send :size="14" /> Enviar
            </button>
          </div>
        </div>
      </div>

      <div v-if="filteredBookings.length" class="bookings-list">
        <BookingCard
          v-for="booking in filteredBookings"
          :key="booking.id"
          :booking="booking"
          @cancel="cancelBooking"
        />
      </div>

      <EmptyState
        v-else-if="bookings.length === 0"
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

.review-prompt-card {
  position: relative;
  background: var(--color-surface);
  border: 1.5px solid #fef08a;
  background-color: #fefce8;
  border-radius: var(--radius-lg);
  padding: var(--spacing-4);
  box-shadow: 0 4px 12px rgba(234, 179, 8, 0.08);
}

.review-prompt-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.dismiss-btn {
  position: absolute;
  top: 10px;
  right: 10px;
  background: none;
  border: none;
  color: #94a3b8;
  cursor: pointer;
  padding: 4px;
}

.dismiss-btn:hover {
  color: #475569;
}

.star-rating {
  display: flex;
  gap: 4px;
}

.star-btn {
  background: none;
  border: none;
  cursor: pointer;
  padding: 2px;
}

.fill-warning {
  fill: #eab308;
}

.text-warning {
  color: #eab308;
}
</style>
