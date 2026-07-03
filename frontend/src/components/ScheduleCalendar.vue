<script setup lang="ts">
import { ref } from 'vue';

const props = defineProps<{
  schedules: any[];
  loading?: boolean;
}>();

const emit = defineEmits<{
  add: [start: string, end: string, spots: number];
  delete: [id: number];
}>();

const date = ref('');
const time = ref('');
const spots = ref(10);
const durationHours = ref(2);

function handleAdd() {
  if (!date.value || !time.value) return;
  
  const dateParts = date.value.split('-');
  const timeParts = time.value.split(':');
  
  const year = Number(dateParts[0]);
  const month = Number(dateParts[1]);
  const day = Number(dateParts[2]);
  
  const hours = Number(timeParts[0]);
  const minutes = Number(timeParts[1]);
  
  const startObj = new Date(year, month - 1, day, hours, minutes, 0);
  const start = startObj.toISOString();
  
  // Simple end time calculation (duration hours later)
  const end = new Date(startObj.getTime() + durationHours.value * 3600000).toISOString();
  
  emit('add', start, end, spots.value);
  
  // reset
  time.value = '';
}
</script>

<template>
  <div class="schedule-calendar card p-6">
    <h3 class="text-lg font-semibold mb-4">Agregar Horario</h3>
    
    <form @submit.prevent="handleAdd" class="add-form mb-6">
      <div class="schedule-form-row">
        <div class="form-group date-col">
          <label class="form-label">Fecha</label>
          <input v-model="date" type="date" class="form-input" required />
        </div>
        <div class="form-group time-col">
          <label class="form-label">Hora</label>
          <input v-model="time" type="time" class="form-input" required />
        </div>
        <div class="form-group spots-col">
          <label class="form-label">Lugares</label>
          <input v-model.number="spots" type="number" class="form-input" min="1" required />
        </div>
        <div class="form-group btn-col">
          <button type="submit" class="btn btn-primary w-full" :disabled="loading">
            Añadir
          </button>
        </div>
      </div>
    </form>

    <h3 class="text-lg font-semibold mb-3">Horarios Programados</h3>
    
    <div v-if="schedules.length === 0" class="text-center py-6 text-secondary text-sm">
      No hay horarios programados.
    </div>
    
    <div v-else class="schedules-list">
      <div v-for="s in schedules" :key="s.id" class="schedule-item">
        <div class="schedule-info">
          <span class="font-semibold">{{ new Date(s.start_time).toLocaleDateString('es-MX', { weekday: 'short', day: 'numeric', month: 'short' }) }}</span>
          <span class="text-secondary ml-2">{{ new Date(s.start_time).toLocaleTimeString('es-MX', { hour: '2-digit', minute: '2-digit' }) }}</span>
        </div>
        <div class="schedule-spots">
          <span class="badge badge-success">{{ s.available_spots }} disponibles</span>
        </div>
        <button class="btn btn-ghost text-error ml-auto" @click="emit('delete', s.id)">
          &times;
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.p-6 { padding: var(--spacing-6); }
.py-6 { padding-top: var(--spacing-6); padding-bottom: var(--spacing-6); }
.mb-4 { margin-bottom: var(--spacing-4); }
.mb-6 { margin-bottom: var(--spacing-6); }
.mb-3 { margin-bottom: var(--spacing-3); }
.ml-2 { margin-left: var(--spacing-2); }
.ml-auto { margin-left: auto; }
.text-lg { font-size: var(--font-size-lg); }
.text-sm { font-size: var(--font-size-sm); }
.font-semibold { font-weight: var(--font-weight-semibold); }
.text-secondary { color: var(--color-text-secondary); }
.text-center { text-align: center; }
.text-error { color: var(--color-error); }
.w-full { width: 100%; }

.schedule-form-row {
  display: flex;
  flex-wrap: wrap;
  align-items: flex-end;
  gap: var(--spacing-3);
}
.date-col, .time-col {
  flex: 1;
  min-width: 130px;
}
.spots-col {
  width: 90px;
}
.btn-col {
  flex: 1;
  min-width: 100px;
}

.schedules-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-2);
}

.schedule-item {
  display: flex;
  align-items: center;
  padding: var(--spacing-3);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-md);
  background: var(--color-background);
}
</style>
