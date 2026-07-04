<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { CheckCircle2, ArrowRight } from '@lucide/vue'
import { useApi } from '../composables/useApi'

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
</script>

<template>
  <div class="page bg-surface min-h-screen flex items-center justify-center">
    <div class="success-content" v-if="booking">
      <div class="icon-wrap">
        <CheckCircle2 :size="64" class="text-success" />
      </div>

      <h1 class="title">¡Reserva confirmada!</h1>
      <p class="subtitle">
        Tu lugar en <strong>{{ booking.tour_title }}</strong> está asegurado.
      </p>

      <div class="card details-card">
        <p class="text-sm text-secondary mb-2">Detalles de reserva #{{ booking.id }}</p>
        <div class="detail-row">
          <span class="text-light">Fecha</span>
          <span class="font-medium">{{
            new Date(booking.schedule_start).toLocaleDateString('es-MX')
          }}</span>
        </div>
        <div class="detail-row">
          <span class="text-light">Personas</span>
          <span class="font-medium">{{ booking.guest_count }}</span>
        </div>
        <div class="detail-row">
          <span class="text-light">Total pagado</span>
          <span class="font-medium text-primary"
            >${{ booking.total_price.toLocaleString('es-MX') }}</span
          >
        </div>
      </div>

      <div class="actions">
        <button class="btn btn-secondary btn-block mb-3" @click="router.push('/bookings')">
          Ver mis reservas
          <ArrowRight :size="18" />
        </button>
        <button class="btn btn-ghost btn-block" @click="router.push('/')">Volver al inicio</button>
      </div>
    </div>
    <div v-else class="text-center p-8">Cargando...</div>
  </div>
</template>

<style scoped>
.bg-surface {
  background: var(--color-surface);
}
.min-h-screen {
  min-height: 100vh;
}
.text-success {
  color: var(--color-success);
}
.text-secondary {
  color: var(--color-text-secondary);
}
.text-light {
  color: var(--color-text-light);
}
.text-primary {
  color: var(--color-secondary);
}
.font-medium {
  font-weight: var(--font-weight-medium);
}
.mb-2 {
  margin-bottom: var(--spacing-2);
}
.mb-3 {
  margin-bottom: var(--spacing-3);
}

.success-content {
  width: 100%;
  max-width: 400px;
  padding: var(--spacing-6);
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
}

.icon-wrap {
  width: 96px;
  height: 96px;
  border-radius: var(--radius-full);
  background: var(--color-success-bg);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: var(--spacing-6);
}

.title {
  font-size: var(--font-size-2xl);
  font-weight: var(--font-weight-bold);
  margin-bottom: var(--spacing-2);
}

.subtitle {
  font-size: var(--font-size-base);
  color: var(--color-text-secondary);
  margin-bottom: var(--spacing-8);
  line-height: var(--line-height-relaxed);
}

.details-card {
  width: 100%;
  padding: var(--spacing-4);
  text-align: left;
  background: var(--color-background);
  border: 1px solid var(--color-border);
  box-shadow: none;
  margin-bottom: var(--spacing-8);
}

.detail-row {
  display: flex;
  justify-content: space-between;
  padding: var(--spacing-2) 0;
  border-bottom: 1px solid var(--color-border-light);
  font-size: var(--font-size-sm);
}

.detail-row:last-child {
  border-bottom: none;
  padding-bottom: 0;
}

.actions {
  width: 100%;
}
</style>
