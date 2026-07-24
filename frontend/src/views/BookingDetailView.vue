<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft, Calendar, Clock, MapPin, Users, DollarSign, UserCheck, AlertTriangle } from '@lucide/vue'
import { useApi } from '../composables/useApi'
import { useToast } from '../composables/useToast'
import { useConfirm } from '../composables/useConfirm'
import { useCalendar } from '../composables/useCalendar'

const route = useRoute()
const router = useRouter()
const api = useApi()
const toast = useToast()
const { confirm } = useConfirm()
const { downloadIcs } = useCalendar()

const booking = ref<any>(null)
const loading = ref(true)
const cancelling = ref(false)

const bookingId = computed(() => route.params.id as string)

onMounted(async () => {
  fetchBooking()
})

async function fetchBooking() {
  loading.value = true
  try {
    const res = await api.get(`/bookings/${bookingId.value}`)
    booking.value = res.data
  } catch (e) {
    console.error(e)
    toast.error('Error al cargar la reserva')
  } finally {
    loading.value = false
  }
}

const hoursUntilTour = computed(() => {
  if (!booking.value?.schedule_start) return 0
  const start = new Date(booking.value.schedule_start).getTime()
  return (start - Date.now()) / (1000 * 3600)
})

const isEligibleForFullRefund = computed(() => {
  return hoursUntilTour.value > 48
})

function handleAddToCalendar() {
  if (!booking.value) return
  downloadIcs({
    title: booking.value.tour_title,
    description: `Reserva #${booking.value.id} con el guía ${booking.value.guide_name}`,
    location: booking.value.tour_location,
    start: booking.value.schedule_start,
    end: booking.value.schedule_end || booking.value.schedule_start,
  })
  toast.success('Evento de calendario descargado')
}

async function handleCancel() {
  if (!booking.value) return

  const isFull = isEligibleForFullRefund.value
  const refundText = isFull
    ? 'Recibirás un reembolso completo (100%).'
    : `Se aplicará la política de cancelación tardía (<48h). Recibirás un reembolso del 50% ($${(booking.value.total_price * 0.5).toFixed(2)} USD).`

  const approved = await confirm({
    title: '¿Cancelar reserva?',
    body: `${refundText} ¿Estás seguro de que deseas continuar?`,
    confirmLabel: 'Sí, cancelar reserva',
    cancelLabel: 'No, conservar reserva',
    confirmVariant: 'danger',
  })

  if (!approved) return

  cancelling.value = true
  try {
    const res = await api.patch(`/bookings/${booking.value.id}/cancel`)
    const result = res.data
    toast.success(result.message || 'Reserva cancelada exitosamente')
    await fetchBooking()
  } catch (e: any) {
    console.error(e)
    toast.error(e.response?.data?.message || 'Error al cancelar la reserva')
  } finally {
    cancelling.value = false
  }
}

function formatDate(dateStr: string) {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleDateString('es-MX', {
    weekday: 'long',
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  })
}

function formatTime(dateStr: string) {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleTimeString('es-MX', { hour: 'numeric', minute: '2-digit', hour12: true })
}
</script>

<template>
  <div class="booking-detail-page bg-surface min-h-screen">
    <header class="header">
      <div class="container flex items-center justify-between">
        <button class="back-btn" @click="router.back()">
          <ArrowLeft :size="20" />
        </button>
        <h1 class="title">Detalle de Reserva</h1>
        <div style="width: 36px"></div>
      </div>
    </header>

    <div v-if="loading" class="container py-8 text-center text-secondary">
      Cargando detalle...
    </div>

    <div v-else-if="booking" class="container py-6 pb-20">
      <!-- Tour Hero Image Card -->
      <div class="card hero-card overflow-hidden mb-6">
        <div class="hero-image-wrap">
          <img
            :src="booking.tour_image || 'https://images.unsplash.com/photo-1501785888041-af3ef285b470?w=800&h=400&fit=crop'"
            :alt="booking.tour_title"
            class="hero-image"
          />
          <div class="hero-overlay">
            <span
              class="badge status-badge"
              :class="{
                'badge-success': booking.status === 'confirmed',
                'badge-warning': booking.status === 'pending',
                'badge-secondary': booking.status === 'completed',
                'badge-error': booking.status === 'cancelled',
              }"
            >
              {{ booking.status.toUpperCase() }}
            </span>
            <span class="ref-code">WND-{{ booking.id }}</span>
          </div>
        </div>

        <div class="p-5">
          <h2 class="text-xl font-bold text-dark mb-2">{{ booking.tour_title }}</h2>
          <div class="flex items-center gap-2 text-sm text-secondary">
            <MapPin :size="16" class="text-primary" />
            <span>{{ booking.tour_location }}</span>
          </div>
        </div>
      </div>

      <!-- Booking Info Details Grid -->
      <div class="card p-5 mb-6 flex flex-col gap-4">
        <h3 class="text-md font-semibold border-b pb-2">Información del Tour</h3>

        <div class="info-row">
          <div class="info-item">
            <Calendar :size="18" class="icon-accent" />
            <div>
              <span class="info-label">Fecha</span>
              <span class="info-value">{{ formatDate(booking.schedule_start) }}</span>
            </div>
          </div>
        </div>

        <div class="info-row">
          <div class="info-item">
            <Clock :size="18" class="icon-accent" />
            <div>
              <span class="info-label">Hora</span>
              <span class="info-value">{{ formatTime(booking.schedule_start) }}</span>
            </div>
          </div>
        </div>

        <div class="info-row">
          <div class="info-item">
            <Users :size="18" class="icon-accent" />
            <div>
              <span class="info-label">Asistentes</span>
              <span class="info-value">{{ booking.guest_count }} persona{{ booking.guest_count > 1 ? 's' : '' }}</span>
            </div>
          </div>
          <div class="info-item">
            <DollarSign :size="18" class="icon-accent" />
            <div>
              <span class="info-label">Total Paid</span>
              <span class="info-value font-bold text-primary">${{ booking.total_price.toFixed(2) }} USD</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Guide Info Card -->
      <div class="card p-5 mb-6">
        <h3 class="text-md font-semibold border-b pb-2 mb-3">Tu Guía Local</h3>
        <div class="flex items-center gap-3">
          <img
            :src="booking.guide_avatar || 'https://images.unsplash.com/photo-1500648767791-00dcc994a43e?auto=format&fit=crop&w=150&q=80'"
            alt="Guide"
            class="guide-avatar"
          />
          <div>
            <span class="font-bold text-dark block">{{ booking.guide_name || 'Guía Wander' }}</span>
            <span class="text-xs text-secondary flex items-center gap-1">
              <UserCheck :size="12" /> Guía Verificado
            </span>
          </div>
        </div>
      </div>

      <!-- Add to Calendar Action -->
      <div class="mb-6">
        <button class="btn btn-outline btn-block flex items-center justify-center gap-2" @click="handleAddToCalendar">
          <Calendar :size="18" />
          Añadir a mi Calendario (.ics)
        </button>
      </div>

      <!-- Cancellation Section -->
      <div v-if="booking.status !== 'cancelled' && booking.status !== 'completed'" class="card cancel-card p-5">
        <h3 class="text-md font-semibold text-error mb-2 flex items-center gap-2">
          <AlertTriangle :size="18" /> Política de Cancelación
        </h3>
        <p class="text-xs text-secondary mb-4 leading-relaxed">
          Cancela con más de 48h de anticipación para un <strong>reembolso del 100%</strong>. Cancelaciones entre 48h y 0h reciben un <strong>reembolso del 50%</strong>.
        </p>

        <div class="policy-status mb-4 p-3 rounded-lg" :class="isEligibleForFullRefund ? 'bg-success-light' : 'bg-warning-light'">
          <span class="text-xs font-semibold" :class="isEligibleForFullRefund ? 'text-success' : 'text-warning'">
            {{ isEligibleForFullRefund ? '✓ Cancelación gratuita disponible (Reembolso 100%)' : '⚠️ Cancelación tardía (<48h) — Reembolso del 50%' }}
          </span>
        </div>

        <button
          class="btn btn-error-outline btn-block"
          :disabled="cancelling"
          @click="handleCancel"
        >
          {{ cancelling ? 'Cancelando...' : 'Cancelar esta reserva' }}
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.pb-20 {
  padding-bottom: 5rem;
}

.container {
  padding: 0 var(--content-padding, 1rem);
  max-width: 600px;
  margin: 0 auto;
}

.header {
  padding: 1rem;
  background: var(--color-surface);
  border-bottom: 1px solid var(--color-border-light);
  position: sticky;
  top: 0;
  z-index: 10;
}

.back-btn {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: var(--color-background);
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
}

.title {
  font-size: 1.125rem;
  font-weight: 700;
}

.card {
  background: var(--color-surface);
  border-radius: 16px;
  border: 1px solid var(--color-border-light);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.03);
}

.hero-card {
  position: relative;
}

.hero-image-wrap {
  position: relative;
  height: 180px;
}

.hero-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.hero-overlay {
  position: absolute;
  top: 12px;
  left: 12px;
  right: 12px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.status-badge {
  padding: 4px 10px;
  font-size: 0.75rem;
  font-weight: 700;
  border-radius: 8px;
}

.ref-code {
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
  color: white;
  padding: 4px 8px;
  border-radius: 6px;
  font-size: 0.75rem;
  font-weight: 600;
}

.info-row {
  display: flex;
  gap: 1.5rem;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  flex: 1;
}

.info-label {
  display: block;
  font-size: 0.75rem;
  color: var(--color-text-secondary);
  text-transform: uppercase;
  font-weight: 600;
}

.info-value {
  font-size: 0.95rem;
  font-weight: 600;
  color: var(--color-text);
}

.icon-accent {
  color: var(--color-primary);
}

.guide-avatar {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  object-fit: cover;
}

.btn-block {
  width: 100%;
  padding: 0.75rem;
  border-radius: 12px;
  font-weight: 600;
  cursor: pointer;
}

.btn-outline {
  background: transparent;
  border: 1.5px solid var(--color-primary);
  color: var(--color-primary);
}

.btn-outline:hover {
  background: var(--color-primary-50);
}

.btn-error-outline {
  background: transparent;
  border: 1.5px solid #ef4444;
  color: #ef4444;
}

.btn-error-outline:hover {
  background: #fef2f2;
}

.bg-success-light {
  background: #f0fdf4;
}

.text-success {
  color: #16a34a;
}

.bg-warning-light {
  background: #fefce8;
}

.text-warning {
  color: #ca8a04;
}

.text-error {
  color: #ef4444;
}
</style>
