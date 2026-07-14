<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft, Calendar, Users, Info } from '@lucide/vue'
import { useApi } from '../composables/useApi'

const route = useRoute()
const router = useRouter()
const api = useApi()

const tourId = computed(() => route.params.id as string)
const tour = ref<any>(null)
const schedules = ref<any[]>([])
const selectedScheduleId = ref<number | null>(null)
const guestCount = ref(1)
const notes = ref('')
const loading = ref(false)

const selectedSchedule = computed(() =>
  schedules.value.find((s) => s.id === selectedScheduleId.value),
)

const totalPrice = computed(() => {
  if (!tour.value) return 0
  return tour.value.price_per_person * guestCount.value
})

onMounted(async () => {
  try {
    const tourRes = await api.get(`/tours/${tourId.value}`)
    tour.value = tourRes.data
  } catch (e) {
    console.error('Failed to load tour', e)
  }
  await fetchSchedules()
})

async function fetchSchedules() {
  try {
    const res = await api.get(`/tours/${tourId.value}/schedules`, { params: { active: 'true' } })
    schedules.value = (res.data || []).filter((s: any) => new Date(s.start_time) > new Date())
  } catch (e) {
    console.error(e)
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

async function handleBook() {
  if (!selectedScheduleId.value) return
  loading.value = true
  try {
    const response = await api.post('/bookings', {
      schedule_id: selectedScheduleId.value,
      guest_count: guestCount.value,
      notes: notes.value,
    })
    router.push(`/checkout/${response.data.id}`)
  } catch (e) {
    console.error('Failed to book', e)
    alert('Ocurrió un error al procesar tu reserva.')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="page bg-surface min-h-screen">
    <header class="header">
      <button class="back-btn" @click="router.back()">
        <ArrowLeft :size="20" />
      </button>
      <h1 class="header-title">Reservar tour</h1>
      <div style="width: 40px"></div>
    </header>

    <div class="container pb-20" v-if="tour">
      <div class="tour-summary">
        <h2 class="text-lg font-bold mb-1">{{ tour.title }}</h2>
        <p class="text-sm text-secondary">
          ${{ tour.price_per_person.toLocaleString('es-MX') }} por persona
        </p>
      </div>

      <div class="section">
        <h3 class="section-title"><Calendar :size="18" /> Selecciona una fecha</h3>

        <div v-if="schedules.length === 0" class="empty-state">
          No hay fechas disponibles para este tour en este momento.
        </div>

        <div v-else class="schedule-grid">
          <button
            v-for="s in schedules"
            :key="s.id"
            class="schedule-card"
            :class="{ 'schedule-card--active': selectedScheduleId === s.id }"
            @click="selectedScheduleId = s.id"
          >
            <span class="date">{{ formatDate(s.start_time) }}</span>
            <span class="time">{{ formatTime(s.start_time) }}</span>
            <span class="spots" :class="{ 'text-error': s.available_spots <= 2 }">
              {{ s.available_spots }} lugares
            </span>
          </button>
        </div>
      </div>

      <div class="section" v-if="selectedScheduleId">
        <h3 class="section-title"><Users :size="18" /> ¿Cuántas personas?</h3>

        <div class="counter">
          <button
            class="counter-btn"
            @click="guestCount > 1 && guestCount--"
            :disabled="guestCount <= 1"
          >
            -
          </button>
          <span class="counter-val">{{ guestCount }}</span>
          <button
            class="counter-btn"
            @click="guestCount < (selectedSchedule?.available_spots || 1) && guestCount++"
            :disabled="guestCount >= (selectedSchedule?.available_spots || 1)"
          >
            +
          </button>
        </div>
      </div>

      <div class="section" v-if="selectedScheduleId">
        <h3 class="section-title"><Info :size="18" /> Notas para el guía (Opcional)</h3>
        <textarea
          v-model="notes"
          class="form-input form-textarea"
          placeholder="Alergias, peticiones especiales..."
        ></textarea>
      </div>
    </div>

    <!-- Bottom Bar -->
    <div class="bottom-bar" v-if="tour && selectedScheduleId">
      <div class="price-summary">
        <span class="price-label">Total</span>
        <span class="price-value">${{ totalPrice.toLocaleString('es-MX') }}</span>
      </div>
      <button
        class="btn btn-secondary-light btn-lg"
        @click="handleBook"
        :disabled="loading"
      >
        {{ loading ? 'Procesando...' : 'Continuar al pago' }}
      </button>
    </div>
  </div>
</template>

<style scoped>
.bg-surface {
  background: var(--color-surface);
}

.container {
  padding: 0 var(--content-padding);
}

.min-h-screen {
  min-height: 100vh;
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

.tour-summary {
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
.text-error {
  color: var(--color-error);
}

.section {
  padding: var(--spacing-5) 0;
  border-bottom: 1px solid var(--color-border-light);
}

.section:last-child {
  border-bottom: none;
}

.section-title {
  display: flex;
  align-items: center;
  gap: var(--spacing-2);
  font-size: var(--font-size-base);
  font-weight: var(--font-weight-semibold);
  margin-bottom: var(--spacing-4);
}

.schedule-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
  gap: var(--spacing-3);
}

.schedule-card {
  display: flex;
  flex-direction: column;
  padding: var(--spacing-3);
  border: 1.5px solid var(--color-border);
  border-radius: var(--radius-lg);
  background: var(--color-surface);
  text-align: left;
  transition: all var(--transition-fast);
}

.schedule-card:hover {
  border-color: var(--color-primary-light);
}

.schedule-card--active {
  border-color: var(--color-primary);
  background: var(--color-primary-50);
}

.schedule-card .date {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text);
  margin-bottom: 2px;
}

.schedule-card .time {
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
  margin-bottom: var(--spacing-2);
}

.schedule-card .spots {
  font-size: var(--font-size-xs);
  color: var(--color-success);
  font-weight: var(--font-weight-medium);
}

.counter {
  display: inline-flex;
  align-items: center;
  border: 1.5px solid var(--color-border);
  border-radius: var(--radius-lg);
  padding: 4px;
}

.counter-btn {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-md);
  background: var(--color-background);
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text);
}

.counter-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.counter-val {
  width: 40px;
  text-align: center;
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-bold);
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
</style>
