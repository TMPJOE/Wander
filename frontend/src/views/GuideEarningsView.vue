<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useApi } from '../composables/useApi'
import { useAuthState } from '../composables/useAuthState'
import { ArrowLeft, TrendingUp, Clock, CheckCircle2 } from '@lucide/vue'

const router = useRouter()
const api = useApi()
const authState = useAuthState()

const period = ref<'week' | 'month' | 'year' | 'all'>('all')
const periods = [
  { key: 'week', label: 'Semana' },
  { key: 'month', label: 'Mes' },
  { key: 'year', label: 'Año' },
  { key: 'all', label: 'Todo' },
] as const

const earnings = ref<{
  total_authorized: number
  total_paid: number
  by_tour: {
    tour_id: number
    tour_title: string
    bookings: number
    revenue: number
    status: string
  }[]
} | null>(null)

const loading = ref(false)

async function fetchEarnings() {
  loading.value = true
  try {
    const q = period.value !== 'all' ? `?period=${period.value}` : ''
    const res = await api.get(`/guide/earnings${q}`)
    earnings.value = res.data
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  if (authState.user.value?.role !== 'guide') {
    router.push('/')
    return
  }
  fetchEarnings()
})

watch(period, fetchEarnings)

const fmtMoney = (n: number) =>
  `$${(n || 0).toLocaleString('es-MX', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`
const statusLabel = (s: string) => (s === 'paid' ? 'Pagado' : 'Autorizado')
</script>

<template>
  <div class="page">
    <header class="header">
      <div class="flex items-center gap-3">
        <button class="back-btn" @click="router.push('/guide/dashboard')">
          <ArrowLeft :size="20" />
        </button>
        <h1 class="title">Ganancias</h1>
      </div>
    </header>

    <div class="container py-4">
      <!-- Period Filter -->
      <div class="filters mb-4">
        <button
          v-for="p in periods"
          :key="p.key"
          class="filter-pill"
          :class="{ active: period === p.key }"
          @click="period = p.key"
        >
          {{ p.label }}
        </button>
      </div>

      <div v-if="loading" class="flex flex-col gap-3">
        <div class="skeleton h-20 rounded-lg"></div>
        <div class="skeleton h-20 rounded-lg"></div>
        <div class="skeleton h-40 rounded-lg"></div>
      </div>

      <template v-else-if="earnings">
        <!-- Summary cards: authorized vs paid -->
        <div class="summary-grid mb-5">
          <div class="summary-card summary-card--authorized">
            <div class="summary-card__icon"><Clock :size="20" /></div>
            <div class="summary-card__body">
              <span class="summary-card__value">{{ fmtMoney(earnings.total_authorized) }}</span>
              <span class="summary-card__label">
                Autorizado <TrendingUp :size="11" class="inline-icon" />
              </span>
              <span class="summary-card__hint">Esperando tu confirmación</span>
            </div>
          </div>
          <div class="summary-card summary-card--paid">
            <div class="summary-card__icon"><CheckCircle2 :size="20" /></div>
            <div class="summary-card__body">
              <span class="summary-card__value">{{ fmtMoney(earnings.total_paid) }}</span>
              <span class="summary-card__label">Cobrado</span>
              <span class="summary-card__hint">Ya capturado a tus clientes</span>
            </div>
          </div>
        </div>

        <!-- Per-tour table -->
        <div v-if="earnings.by_tour.length" class="card p-0 overflow-hidden">
          <h2 class="table-title">Desglose por tour</h2>
          <div
            v-for="(row, i) in earnings.by_tour"
            :key="`${row.tour_id}-${row.status}`"
            class="table-row"
            :class="{ 'table-row--last': i === earnings.by_tour.length - 1 }"
          >
            <div class="table-row__main">
              <span class="table-row__title">{{ row.tour_title }}</span>
              <span class="table-row__meta">{{ row.bookings }} reservas</span>
            </div>
            <div class="table-row__right">
              <span
                class="badge"
                :class="row.status === 'paid' ? 'badge-success' : 'badge-warning'"
              >
                {{ statusLabel(row.status) }}
              </span>
              <span class="table-row__revenue">{{ fmtMoney(row.revenue) }}</span>
            </div>
          </div>
        </div>

        <div v-else class="text-center py-12">
          <p class="text-secondary">Aún no tienes ingresos registrados en este período.</p>
        </div>
      </template>
    </div>
  </div>
</template>

<style scoped>
.page {
  flex: 1;
  width: 100%;
  min-height: 100vh;
  min-height: 100dvh;
  background: var(--color-surface);
}

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

.container {
  padding: 0 var(--content-padding);
}

.py-4 {
  padding-top: var(--spacing-4);
  padding-bottom: var(--spacing-4);
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
  -ms-overflow-style: none; /* IE and Edge */
  scrollbar-width: none; /* Firefox */
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

/* Summary cards */
.summary-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: var(--spacing-3);
}
.summary-card {
  border-radius: var(--radius-lg);
  padding: var(--spacing-4);
  display: flex;
  align-items: flex-start;
  gap: var(--spacing-3);
  border: 1px solid var(--color-border-light);
}
.summary-card--authorized {
  background: #fffbeb;
  border-color: #fde68a;
}
.summary-card--paid {
  background: #ecfdf5;
  border-color: #a7f3d0;
}
.summary-card__icon {
  width: 40px;
  height: 40px;
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  background: var(--color-surface);
}
.summary-card--authorized .summary-card__icon {
  color: #b45309;
}
.summary-card--paid .summary-card__icon {
  color: #047857;
}
.summary-card__body {
  display: flex;
  flex-direction: column;
}
.summary-card__value {
  font-size: var(--font-size-xl);
  font-weight: var(--font-weight-bold);
  line-height: 1.2;
}
.summary-card__label {
  font-size: var(--font-size-xs);
  color: var(--color-text-secondary);
  margin-top: 2px;
}
.summary-card__hint {
  font-size: 10px;
  color: var(--color-text-light);
  margin-top: 4px;
}
.inline-icon {
  display: inline-flex;
  vertical-align: middle;
  margin-left: 2px;
  opacity: 0.7;
}

/* Table */
.card {
  background: var(--color-surface);
}
.p-0 {
  padding: 0;
}
.overflow-hidden {
  overflow: hidden;
}
.table-title {
  font-size: var(--font-size-base);
  font-weight: var(--font-weight-semibold);
  padding: var(--spacing-3) var(--spacing-4);
  border-bottom: 1px solid var(--color-border-light);
}
.table-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--spacing-3) var(--spacing-4);
  border-bottom: 1px solid var(--color-border-light);
}
.table-row--last {
  border-bottom: none;
}
.table-row__main {
  display: flex;
  flex-direction: column;
}
.table-row__title {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-semibold);
}
.table-row__meta {
  font-size: var(--font-size-xs);
  color: var(--color-text-secondary);
  margin-top: 2px;
}
.table-row__right {
  display: flex;
  align-items: center;
  gap: var(--spacing-3);
}
.table-row__revenue {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-bold);
  color: var(--color-secondary);
}

/* misc used */
.flex {
  display: flex;
}
.items-center {
  align-items: center;
}
.gap-3 {
  gap: var(--spacing-3);
}
.flex-col {
  flex-direction: column;
}
.text-secondary {
  color: var(--color-text-secondary);
}
.text-center {
  text-align: center;
}
.py-12 {
  padding-top: 3rem;
  padding-bottom: 3rem;
}
.mb-4 {
  margin-bottom: var(--spacing-4);
}
.mb-5 {
  margin-bottom: var(--spacing-5);
}
.h-20 {
  height: 5rem;
}
.h-40 {
  height: 10rem;
}
.rounded-lg {
  border-radius: var(--radius-lg);
}
</style>
