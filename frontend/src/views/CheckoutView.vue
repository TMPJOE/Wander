<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft, Lock, ShieldCheck } from '@lucide/vue'
import {
  loadStripe,
  type Stripe,
  type StripeElements,
  type StripeCardElement,
} from '@stripe/stripe-js'
import { useApi } from '../composables/useApi'
import { useToast } from '../composables/useToast'

const route = useRoute()
const router = useRouter()
const api = useApi()
const toast = useToast()

const bookingId = computed(() => route.params.bookingId as string)

const booking = ref<any>(null)
const loading = ref(true)
const processing = ref(false)
const errorMsg = ref('')
const cardReady = ref(false)
const alreadyAuthorized = ref(false)

let stripe: Stripe | null = null
let elements: StripeElements | null = null
let cardElement: StripeCardElement | null = null
let clientSecret = ''

onMounted(async () => {
  try {
    const bookingRes = await api.get(`/bookings/${bookingId.value}`)
    booking.value = bookingRes.data
    loading.value = false
  } catch (e: any) {
    console.error('Failed to load booking', e)
    errorMsg.value = e?.response?.data?.message || 'No se pudo cargar la reserva.'
    toast.error(errorMsg.value)
    loading.value = false
    return
  }

  try {
    const intentRes = await api.post(`/payments/bookings/${bookingId.value}/intent`)
    clientSecret = intentRes.data.client_secret
    const publishableKey =
      intentRes.data.publishable_key || import.meta.env.VITE_STRIPE_PUBLISHABLE_KEY

    // If the card hold was already placed in a prior session, skip the card form.
    if (intentRes.data.already_authorized) {
      alreadyAuthorized.value = true
      return
    }

    stripe = await loadStripe(publishableKey)
    if (!stripe) {
      errorMsg.value = 'No se pudo inicializar el procesador de pagos.'
      toast.error(errorMsg.value)
      return
    }

    elements = stripe.elements()
    cardElement = elements.create('card', {
      style: {
        base: {
          fontSize: '15px',
          fontFamily: 'Inter, sans-serif',
          color: '#1a1a1a',
          '::placeholder': { color: '#888888' },
        },
        invalid: { color: '#dc2626' },
      },
    })

    await nextTick()
    if (!document.getElementById('card-element')) {
      errorMsg.value = 'No se pudo mostrar el formulario de pago.'
      toast.error(errorMsg.value)
      return
    }
    cardElement.mount('#card-element')
    cardElement.on('ready', () => (cardReady.value = true))
    cardElement.on('change', (event) => {
      errorMsg.value = event.error ? event.error.message : ''
    })
  } catch (e: any) {
    console.error('Failed to create payment intent', e)
    const backendError = e?.response?.data?.error
    const backendMessage = e?.response?.data?.message
    errorMsg.value = backendError ? `${backendMessage}: ${backendError}` : (backendMessage || 'No se pudo iniciar el pago. Intenta de nuevo.')
    toast.error(errorMsg.value)
  }
})

onBeforeUnmount(() => {
  cardElement?.unmount()
})

async function handlePay() {
  if (processing.value) return

  processing.value = true
  errorMsg.value = ''

  try {
    // If the card was already authorized in a prior session, skip confirmCardPayment.
    if (!alreadyAuthorized.value) {
      if (!stripe || !cardElement) return

      const { error, paymentIntent } = await stripe.confirmCardPayment(clientSecret, {
        payment_method: { card: cardElement },
      })

      if (error) {
        // Stripe errors (declined, auth failure, etc.) — show the message and stop.
        throw new Error(error.message || 'El pago no pudo completarse.')
      }

      // With manual capture, a successful authorization lands in requires_capture, not succeeded.
      if (paymentIntent?.status !== 'requires_capture' && paymentIntent?.status !== 'succeeded') {
        throw new Error('El pago no pudo completarse. Intenta de nuevo.')
      }
    }

    await api.post(`/payments/bookings/${bookingId.value}/confirm`)
    router.push(`/booking-success/${bookingId.value}`)
  } catch (e: any) {
    // Stripe errors have a plain .message; axios errors have .response.data.message
    errorMsg.value = e?.response?.data?.message || e?.message || 'Ocurrió un error al procesar el pago.'
    toast.error(errorMsg.value)
  } finally {
    processing.value = false
  }
}
</script>

<template>
  <div class="page bg-surface min-h-screen">
    <header class="header">
      <button class="back-btn" @click="router.back()">
        <ArrowLeft :size="20" />
      </button>
      <h1 class="header-title">Pago</h1>
      <div style="width: 40px"></div>
    </header>

    <div class="container pb-20" v-if="loading">
      <div class="skeleton" style="height: 90px; margin-bottom: var(--spacing-4)"></div>
      <div class="skeleton" style="height: 140px"></div>
    </div>

    <div class="container pb-20" v-else-if="booking">
      <div class="summary">
        <h2 class="text-lg font-bold mb-1">{{ booking.tour_title }}</h2>
        <p class="text-sm text-secondary">
          {{ booking.guest_count }} persona{{ booking.guest_count > 1 ? 's' : '' }} ·
          {{
            new Date(booking.schedule_start).toLocaleDateString('es-MX', {
              day: 'numeric',
              month: 'short',
            })
          }}
        </p>
      </div>

      <div class="section">
        <h3 class="section-title"><Lock :size="16" /> Detalles de la tarjeta</h3>

        <!-- Already authorized from a previous session -->
        <div v-if="alreadyAuthorized" class="authorized-banner">
          <ShieldCheck :size="20" class="text-success" />
          <div>
            <p class="authorized-title">Cargo ya autorizado</p>
            <p class="authorized-sub">Tu tarjeta ya tiene un cargo autorizado. Pulsa "Confirmar reserva" para finalizar.</p>
          </div>
        </div>

        <!-- Normal card entry -->
        <div v-else class="card-input-wrap">
          <div id="card-element"></div>
        </div>

        <p v-if="errorMsg" class="form-error mt-2">{{ errorMsg }}</p>
      </div>

      <div class="security-note">
        <ShieldCheck :size="16" class="text-success" />
        <span>Pago procesado de forma segura por Stripe</span>
      </div>
    </div>

    <div class="bottom-bar" v-if="booking && !loading">
      <div class="price-summary">
        <span class="price-label">Total</span>
        <span class="price-value">${{ booking.total_price.toLocaleString('es-MX') }}</span>
      </div>
      <button
        class="btn btn-secondary-light btn-lg"
        @click="handlePay"
        :disabled="processing || (!alreadyAuthorized && !cardReady)"
      >
        {{ processing ? 'Procesando...' : alreadyAuthorized ? 'Confirmar reserva' : 'Pagar ahora' }}
      </button>
    </div>
  </div>
</template>

<style scoped>
.bg-surface {
  background: var(--color-surface);
}
.min-h-screen {
  min-height: 100vh;
}
.container {
  padding: 0 var(--content-padding);
}
.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--spacing-4);
  border-bottom: 1px solid var(--color-border-light);
  background: var(--color-surface);
  position: sticky;
  top: 0;
  z-index: 10;
}
.back-btn {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-full);
  background: var(--color-background);
}
.header-title {
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
}
.pb-20 {
  padding-bottom: 5rem;
}
.summary {
  padding: var(--spacing-4) 0;
  border-bottom: 1px solid var(--color-border-light);
}
.text-lg {
  font-size: var(--font-size-lg);
}
.font-bold {
  font-weight: var(--font-weight-bold);
}
.mb-1 {
  margin-bottom: var(--spacing-1);
}
.text-sm {
  font-size: var(--font-size-sm);
}
.text-secondary {
  color: var(--color-text-secondary);
}
.text-success {
  color: var(--color-success);
}
.mt-2 {
  margin-top: var(--spacing-2);
}
.section {
  padding: var(--spacing-5) 0;
}
.section-title {
  display: flex;
  align-items: center;
  gap: var(--spacing-2);
  font-size: var(--font-size-base);
  font-weight: var(--font-weight-semibold);
  margin-bottom: var(--spacing-4);
}
.card-input-wrap {
  padding: var(--spacing-3) var(--spacing-4);
  border: 1.5px solid var(--color-border);
  border-radius: var(--radius-md);
  background: var(--color-surface);
}
.security-note {
  display: flex;
  align-items: center;
  gap: var(--spacing-2);
  font-size: var(--font-size-xs);
  color: var(--color-text-light);
  padding: var(--spacing-3) 0;
}
.bottom-bar {
  position: fixed;
  bottom: var(--bottom-nav-height);
  left: 0;
  right: 0;
  max-width: var(--max-width);
  margin: 0 auto;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--spacing-4) var(--spacing-5);
  background: var(--color-surface);
  border-top: 1px solid var(--color-border-light);
  z-index: var(--z-sticky);
  box-shadow: 0 -4px 12px rgba(0, 0, 0, 0.05);
}
.price-summary {
  display: flex;
  flex-direction: column;
}
.price-label {
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
}
.price-value {
  font-size: var(--font-size-xl);
  font-weight: var(--font-weight-bold);
  color: var(--color-secondary);
}

@media (min-width: 1024px) {
  .bottom-bar {
    right: var(--nav-width);
  }
}
.authorized-banner {
  display: flex;
  align-items: flex-start;
  gap: var(--spacing-3);
  padding: var(--spacing-4);
  border-radius: var(--radius-md);
  background: color-mix(in srgb, var(--color-success) 10%, transparent);
  border: 1px solid color-mix(in srgb, var(--color-success) 30%, transparent);
}
.authorized-title {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-semibold);
  color: var(--color-success);
  margin-bottom: var(--spacing-1);
}
.authorized-sub {
  font-size: var(--font-size-xs);
  color: var(--color-text-secondary);
}
</style>
