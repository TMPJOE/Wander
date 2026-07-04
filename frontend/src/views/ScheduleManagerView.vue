<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useApi } from '../composables/useApi'
import { ArrowLeft } from '@lucide/vue'
import ScheduleCalendar from '../components/ScheduleCalendar.vue'

const route = useRoute()
const router = useRouter()
const api = useApi()

const tourId = computed(() => route.params.id as string)
const tour = ref<any>(null)
const schedules = ref<any[]>([])
const loading = ref(false)

onMounted(async () => {
  fetchData()
})

async function fetchData() {
  try {
    const tourRes = await api.get(`/tours/${tourId.value}`)
    tour.value = tourRes.data

    const schedRes = await api.get(`/tours/${tourId.value}/schedules`)
    // Show only future schedules
    schedules.value = (schedRes.data || []).filter((s: any) => new Date(s.start_time) > new Date())
  } catch (e) {
    console.error(e)
  }
}

async function addSchedule(start: string, end: string, spots: number) {
  loading.value = true
  try {
    await api.post(`/schedules`, {
      tour_id: Number(tourId.value),
      start_time: start,
      end_time: end,
      available_spots: spots,
    })
    await fetchData()
  } catch (e) {
    console.error(e)
    alert('Error al añadir horario')
  } finally {
    loading.value = false
  }
}

async function deleteSchedule(id: number) {
  if (!confirm('¿Seguro que deseas eliminar este horario?')) return
  loading.value = true
  try {
    await api.delete(`/schedules/${id}`)
    await fetchData()
  } catch (e) {
    console.error(e)
    alert('Error al eliminar horario. ¿Quizás tiene reservas?')
  } finally {
    loading.value = false
  }
}

async function toggleActive(id: number, isActive: boolean) {
  loading.value = true
  try {
    await api.put(`/schedules/${id}`, { is_active: isActive })
    await fetchData()
  } catch (e) {
    console.error(e)
    alert('Error al cambiar el estado del horario')
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
      <h1 class="title">Horarios del Tour</h1>
      <div style="width: 36px"></div>
    </header>

    <div class="container py-4">
      <div v-if="tour" class="mb-4">
        <h2 class="font-semibold p-4">{{ tour.title }}</h2>
      </div>

      <ScheduleCalendar
        :schedules="schedules"
        :loading="loading"
        @add="addSchedule"
        @delete="deleteSchedule"
        @toggle-active="toggleActive"
      />
    </div>
  </div>
</template>

<style scoped>
.container {
  padding: 0 var(--content-padding);
  flex: 1;
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

.bg-surface {
  background: var(--color-surface);
}
.min-h-screen {
  min-height: 100vh;
}
.py-4 {
  padding-top: var(--spacing-4);
  padding-bottom: var(--spacing-4);
}
.mb-4 {
  margin-bottom: var(--spacing-4);
}
.font-semibold {
  font-weight: var(--font-weight-semibold);
}
.p-4 {
  padding-left: var(--spacing-4);
}
</style>
