<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Check, MapPin, Calendar, Clock } from '@lucide/vue'
import { useApi } from '../composables/useApi'
import { normalizeTourImages } from '../utils/tourImages'

const route = useRoute()
const router = useRouter()
const api = useApi()

const bookingId = route.params.id as string
const booking = ref<any>(null)

onMounted(async () => {
  try {
    const res = await api.get(`/bookings/${bookingId}`)
    booking.value = res.data
  } catch (e) {
    console.error(e)
  }
})

const tourImageUrl = computed(() => {
  if (!booking.value?.tour_image) return 'https://images.unsplash.com/photo-1501785888041-af3ef285b470?w=200&h=200&fit=crop'
  const imgs = normalizeTourImages(booking.value.tour_image)
  return imgs[0] || 'https://images.unsplash.com/photo-1501785888041-af3ef285b470?w=200&h=200&fit=crop'
})

const defaultAvatar = 'https://images.unsplash.com/photo-1500648767791-00dcc994a43e?auto=format&fit=crop&w=150&q=80'

function formatDate(dateStr: string) {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleDateString('es-MX', { day: 'numeric', month: 'short', year: 'numeric' })
}

function formatTime(dateStr: string) {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleTimeString('es-MX', { hour: 'numeric', minute: '2-digit', hour12: true })
}
</script>

<template>
  <div class="page bg-background min-h-screen flex items-center justify-center p-4">
    <div class="success-content" v-if="booking">
      <div class="icon-wrap">
        <Check :size="48" :stroke-width="3" class="text-success-dark" />
      </div>

      <h1 class="title">¡Reserva Confirmada!</h1>
      <p class="subtitle">
        Tu guía ha sido notificado. Prepárate para una aventura local inolvidable.
      </p>

      <div class="card booking-summary-card">
        <!-- Tour Info -->
        <div class="booking-tour-info">
          <img :src="tourImageUrl" alt="Tour" class="booking-tour-img" />
          <div class="booking-tour-text">
            <h3 class="booking-tour-title">{{ booking.tour_title }}</h3>
            <div class="booking-tour-location">
              <MapPin :size="14" class="text-light" />
              <span>{{ booking.tour_location }}</span>
            </div>
          </div>
        </div>

        <hr class="divider" />

        <!-- Date & Time -->
        <div class="booking-datetime">
          <div class="booking-date">
            <span class="label">FECHA</span>
            <div class="value">
              <Calendar :size="18" class="icon-accent" />
              <span>{{ formatDate(booking.schedule_start) }}</span>
            </div>
          </div>
          <div class="booking-time">
            <span class="label">HORA</span>
            <div class="value">
              <Clock :size="18" class="icon-accent" />
              <span>{{ formatTime(booking.schedule_start) }}</span>
            </div>
          </div>
        </div>

        <hr class="divider" />

        <!-- Guide -->
        <div class="booking-guide">
          <span class="label">TU GUÍA LOCAL</span>
          <div class="guide-info">
            <img :src="booking.guide_avatar || defaultAvatar" alt="Guide" class="guide-avatar" />
            <span class="guide-name">{{ booking.guide_name }}</span>
          </div>
        </div>
      </div>

      <div class="actions">
        <button class="btn btn-success-solid btn-block mb-4" @click="router.push('/bookings')">
          Ver Detalles de Reserva
        </button>
        <button class="btn btn-ghost btn-block btn-back" @click="router.push('/')">
          Volver a Explorar
        </button>
      </div>
    </div>
    <div v-else class="text-center p-8 text-light">Cargando...</div>
  </div>
</template>

<style scoped>
.bg-background {
  background: var(--color-background);
}
.min-h-screen {
  min-height: 100vh;
}
.p-4 {
  padding: var(--spacing-4);
}
.text-light {
  color: var(--color-text-light);
}
.mb-4 {
  margin-bottom: var(--spacing-4);
}

.success-content {
  width: 100%;
  max-width: 440px;
  background: white;
  border-radius: var(--radius-xl);
  padding: var(--spacing-8) var(--spacing-5);
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.05);
  border: 4px solid #f2f5ff; /* Slight glow border to match screenshot */
}

.icon-wrap {
  width: 80px;
  height: 80px;
  border-radius: var(--radius-full);
  background: #87f0b5; /* Light mint green */
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: var(--spacing-5);
}

.text-success-dark {
  color: #0b593f;
}

.title {
  font-size: var(--font-size-2xl);
  font-weight: 800;
  color: #0c1a2e; /* Very dark blue */
  margin-bottom: var(--spacing-2);
}

.subtitle {
  font-size: var(--font-size-base);
  color: var(--color-text-secondary);
  margin-bottom: var(--spacing-6);
  line-height: var(--line-height-relaxed);
  max-width: 90%;
}

.booking-summary-card {
  width: 100%;
  padding: var(--spacing-5);
  text-align: left;
  background: white;
  border: 1px solid #f0ded6; /* Soft coral/orange border */
  border-radius: var(--radius-lg);
  box-shadow: none;
  margin-bottom: var(--spacing-8);
}

.booking-tour-info {
  display: flex;
  align-items: flex-start;
  gap: var(--spacing-4);
}

.booking-tour-img {
  width: 64px;
  height: 64px;
  border-radius: var(--radius-md);
  object-fit: cover;
}

.booking-tour-text {
  flex: 1;
}

.booking-tour-title {
  font-size: var(--font-size-lg);
  font-weight: 700;
  color: #0c1a2e;
  line-height: var(--line-height-tight);
  margin-bottom: var(--spacing-2);
}

.booking-tour-location {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
}

.divider {
  border: none;
  border-top: 1px solid #f0ded6;
  margin: var(--spacing-4) 0;
}

.booking-datetime {
  display: flex;
  gap: var(--spacing-6);
}

.booking-date, .booking-time {
  flex: 1;
}

.label {
  display: block;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.5px;
  color: var(--color-text-light);
  margin-bottom: var(--spacing-2);
  text-transform: uppercase;
}

.value {
  display: flex;
  align-items: center;
  gap: var(--spacing-2);
  font-size: var(--font-size-base);
  color: #0c1a2e;
}

.icon-accent {
  color: var(--color-primary); /* Brown/Orange color from screenshot */
}

.booking-guide {
  display: flex;
  flex-direction: column;
}

.guide-info {
  display: flex;
  align-items: center;
  gap: var(--spacing-3);
  margin-top: 4px;
}

.guide-avatar {
  width: 32px;
  height: 32px;
  border-radius: var(--radius-full);
  object-fit: cover;
}

.guide-name {
  font-size: var(--font-size-base);
  font-weight: 700;
  color: #0c1a2e;
}

.actions {
  width: 100%;
}

.btn-success-solid {
  background-color: #0b6b47;
  color: white;
  font-weight: 700;
  padding: 14px;
  border-radius: var(--radius-md);
  border: none;
  transition: all 0.2s ease;
}

.btn-success-solid:hover {
  background-color: #085035;
  transform: translateY(-1px);
}

.btn-back {
  color: var(--color-primary-dark);
  font-weight: 700;
  background: transparent;
}
.btn-back:hover {
  background: var(--color-primary-50);
}
</style>
