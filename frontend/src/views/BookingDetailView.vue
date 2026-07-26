<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  ArrowLeft,
  Calendar,
  Clock,
  MapPin,
  Users,
  DollarSign,
  UserCheck,
  AlertTriangle,
} from '@lucide/vue'
import { useApi } from '../composables/useApi'
import { useToast } from '../composables/useToast'
import { useConfirm } from '../composables/useConfirm'
import { useCalendar } from '../composables/useCalendar'
import LocationMap from '../components/LocationMap.vue'

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
    : `Se aplicará la política de cancelación tardía (<48h). Recibirás un reembolso del 50% ($${(booking.value.total_price * 0.5).toFixed(2)} PAB).`

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
  return new Date(dateStr).toLocaleTimeString('es-MX', {
    hour: 'numeric',
    minute: '2-digit',
    hour12: true,
  })
}
</script>

<template>
  <div class="booking-detail-page page">
    <header class="header">
      <button class="back-btn" @click="router.back()">
        <ArrowLeft :size="20" />
      </button>
      <h1 class="header-title">Detalle de Reserva</h1>
      <div style="width: 40px"></div>
    </header>

    <div v-if="loading" class="container py-8 text-center text-secondary">
      <div class="skeleton" style="height: 120px; margin-bottom: var(--spacing-4)"></div>
      <div class="skeleton" style="height: 200px"></div>
    </div>

    <div v-else-if="booking" class="container">
      <!-- Tour Hero Image Card -->
      <div class="card hero-card">
        <div class="hero-image-wrap">
          <img
            :src="
              booking.tour_image ||
              'https://images.unsplash.com/photo-1501785888041-af3ef285b470?w=800&h=400&fit=crop'
            "
            :alt="booking.tour_title"
            class="hero-image"
          />
          <div class="hero-overlay">
            <span
              class="badge"
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

        <div class="hero-body">
          <h2 class="hero-title">{{ booking.tour_title }}</h2>
          <div class="hero-location">
            <MapPin :size="16" class="text-primary" />
            <span>{{ booking.tour_location }}</span>
          </div>
        </div>
      </div>

      <!-- Booking Info Details Grid (2x2 Grid) -->
      <div class="card detail-card">
        <h3 class="card-section-title">Información del Tour</h3>

        <div class="info-grid">
          <div class="info-cell">
            <div class="info-icon-wrap">
              <Calendar :size="18" class="icon-accent" />
            </div>
            <div class="info-text">
              <span class="info-label">Fecha</span>
              <span class="info-value">{{ formatDate(booking.schedule_start) }}</span>
            </div>
          </div>

          <div class="info-cell">
            <div class="info-icon-wrap">
              <Clock :size="18" class="icon-accent" />
            </div>
            <div class="info-text">
              <span class="info-label">Hora</span>
              <span class="info-value">{{ formatTime(booking.schedule_start) }}</span>
            </div>
          </div>

          <div class="info-cell">
            <div class="info-icon-wrap">
              <Users :size="18" class="icon-accent" />
            </div>
            <div class="info-text">
              <span class="info-label">Asistentes</span>
              <span class="info-value"
                >{{ booking.guest_count }} persona{{ booking.guest_count > 1 ? 's' : '' }}</span
              >
            </div>
          </div>

          <div class="info-cell">
            <div class="info-icon-wrap">
              <DollarSign :size="18" class="icon-accent" />
            </div>
            <div class="info-text">
              <span class="info-label">Total Pagado</span>
              <span class="info-value font-bold text-primary"
                >${{ booking.total_price.toFixed(2) }} PAB</span
              >
            </div>
          </div>
        </div>
      </div>

      <!-- Guide Info Card -->
      <div class="card detail-card">
        <h3 class="card-section-title">Tu Guía Local</h3>
        <div class="guide-profile">
          <img
            :src="
              booking.guide_avatar ||
              'https://images.unsplash.com/photo-1500648767791-00dcc994a43e?auto=format&fit=crop&w=150&q=80'
            "
            alt="Guide"
            class="guide-avatar"
          />
          <div class="guide-details">
            <span class="guide-name">{{ booking.guide_name || 'Guía Wander' }}</span>
            <span class="guide-badge">
              <UserCheck :size="13" class="text-success" /> Guía Verificado
            </span>
          </div>
        </div>
      </div>

      <!-- Meeting Point Map Card -->
      <div v-if="booking.meeting_point || booking.tour_latitude" class="card detail-card">
        <h3 class="card-section-title">Punto de Encuentro</h3>
        <p
          class="text-muted"
          style="font-size: var(--font-size-sm); margin-bottom: var(--spacing-3)"
        >
          {{ booking.meeting_point || booking.tour_location }}
        </p>
        <LocationMap
          :label="booking.meeting_point || booking.tour_location"
          :latitude="booking.tour_latitude"
          :longitude="booking.tour_longitude"
          height="250px"
          :show-directions="true"
        />
      </div>

      <!-- Add to Calendar Action -->
      <div>
        <button class="btn btn-info btn-block btn-lg" @click="handleAddToCalendar">
          <Calendar :size="18" />
          Añadir a mi Calendario (.ics)
        </button>
      </div>

      <!-- Cancellation Section -->
      <div
        v-if="booking.status !== 'cancelled' && booking.status !== 'completed'"
        class="card cancel-card"
      >
        <h3 class="cancel-title"><AlertTriangle :size="18" /> Política de Cancelación</h3>
        <p class="cancel-desc">
          Cancela con más de 48h de anticipación para un <strong>reembolso del 100%</strong>.
          Cancelaciones entre 48h y 0h reciben un <strong>reembolso del 50%</strong>.
        </p>

        <div
          class="policy-status"
          :class="isEligibleForFullRefund ? 'policy-status--success' : 'policy-status--warning'"
        >
          <span class="policy-status-text">
            {{
              isEligibleForFullRefund
                ? '✓ Cancelación gratuita disponible (Reembolso 100%)'
                : '⚠️ Cancelación tardía (<48h) — Reembolso del 50%'
            }}
          </span>
        </div>

        <button
          class="btn btn-danger-light btn-block btn-lg"
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
.booking-detail-page {
  background-color: var(--color-background);
  min-height: 100vh;
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--spacing-4);
  background: var(--color-surface);
  border-bottom: 1px solid var(--color-border-light);
  position: sticky;
  top: 0;
  z-index: var(--z-sticky);
}

.back-btn {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-full);
  background: var(--color-background);
  transition: background-color var(--transition-fast);
}

.back-btn:hover {
  background: var(--color-border-light);
}

.header-title {
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
}

.container {
  padding: var(--spacing-5) var(--content-padding);
  max-width: var(--max-width);
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: var(--spacing-5);
}

.card {
  background: var(--color-surface);
  border-radius: var(--radius-xl);
  border: 1px solid var(--color-border-light);
  box-shadow: var(--shadow-sm);
  overflow: hidden;
}

.hero-card {
  position: relative;
}

.hero-image-wrap {
  position: relative;
  height: 200px;
}

.hero-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.hero-overlay {
  position: absolute;
  top: var(--spacing-3);
  left: var(--spacing-3);
  right: var(--spacing-3);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.ref-code {
  background: rgba(0, 0, 0, 0.65);
  backdrop-filter: blur(8px);
  color: white;
  padding: 4px 10px;
  border-radius: var(--radius-sm);
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-semibold);
  letter-spacing: 0.5px;
}

.hero-body {
  padding: var(--spacing-4) var(--spacing-5);
}

.hero-title {
  font-size: var(--font-size-xl);
  font-weight: var(--font-weight-bold);
  color: var(--color-text);
  margin-bottom: var(--spacing-2);
  line-height: var(--line-height-tight);
}

.hero-location {
  display: flex;
  align-items: center;
  gap: var(--spacing-2);
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
}

.detail-card {
  padding: var(--spacing-5);
}

.card-section-title {
  font-size: var(--font-size-base);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text);
  padding-bottom: var(--spacing-3);
  border-bottom: 1px solid var(--color-border-light);
  margin-bottom: var(--spacing-4);
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: var(--spacing-4);
}

@media (max-width: 480px) {
  .info-grid {
    grid-template-columns: 1fr;
    gap: var(--spacing-3);
  }
}

.info-cell {
  display: flex;
  align-items: center;
  gap: var(--spacing-3);
  padding: var(--spacing-2) 0;
}

.info-icon-wrap {
  width: 38px;
  height: 38px;
  border-radius: var(--radius-md);
  background: var(--color-primary-50);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.info-text {
  display: flex;
  flex-direction: column;
}

.info-label {
  font-size: var(--font-size-xs);
  color: var(--color-text-light);
  text-transform: uppercase;
  font-weight: var(--font-weight-semibold);
  letter-spacing: 0.5px;
}

.info-value {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text);
  line-height: var(--line-height-tight);
}

.icon-accent {
  color: var(--color-primary);
}

.font-bold {
  font-weight: var(--font-weight-bold);
}

.guide-profile {
  display: flex;
  align-items: center;
  gap: var(--spacing-4);
}

.guide-avatar {
  width: 48px;
  height: 48px;
  border-radius: var(--radius-full);
  object-fit: cover;
  border: 2px solid var(--color-border-light);
}

.guide-details {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.guide-name {
  font-size: var(--font-size-base);
  font-weight: var(--font-weight-bold);
  color: var(--color-text);
}

.guide-badge {
  font-size: var(--font-size-xs);
  color: var(--color-text-secondary);
  display: flex;
  align-items: center;
  gap: 4px;
}

.cancel-card {
  padding: var(--spacing-5);
}

.cancel-title {
  font-size: var(--font-size-base);
  font-weight: var(--font-weight-semibold);
  color: var(--color-error-dark);
  margin-bottom: var(--spacing-2);
  display: flex;
  align-items: center;
  gap: var(--spacing-2);
}

.cancel-desc {
  font-size: var(--font-size-xs);
  color: var(--color-text-secondary);
  margin-bottom: var(--spacing-4);
  line-height: var(--line-height-relaxed);
}

.policy-status {
  margin-bottom: var(--spacing-4);
  padding: var(--spacing-3) var(--spacing-4);
  border-radius: var(--radius-md);
  border: 1px solid transparent;
}

.policy-status--success {
  background: var(--color-success-bg);
  border-color: color-mix(in srgb, var(--color-success) 20%, transparent);
  color: var(--color-success);
}

.policy-status--warning {
  background: var(--color-warning-bg);
  border-color: color-mix(in srgb, var(--color-warning) 20%, transparent);
  color: #b45309;
}

.policy-status-text {
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-semibold);
}
</style>
