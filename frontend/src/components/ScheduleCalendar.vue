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
  const start = new Date(`${date.value}T${time.value}:00`).toISOString();
  // Simple end time calculation (duration hours later)
  const end = new Date(new Date(start).getTime() + durationHours.value * 3600000).toISOString();
  
  emit('add', start, end, spots.value);
  
  // reset
  time.value = '';
}
</script>

<template>
  <div class="schedule-calendar card p-6">
    <h3 class="text-lg font-semibold mb-4">Agregar Horario</h3>
    
    <form @submit.prevent="handleAdd" class="add-form mb-6">
      <div class="grid grid-cols-4 gap-3">
        <div class="form-group col-span-2 sm:col-span-1">
          <label class="form-label">Fecha</label>
          <input v-model="date" type="date" class="form-input" required />
        </div>
        <div class="form-group col-span-2 sm:col-span-1">
          <label class="form-label">Hora</label>
          <input v-model="time" type="time" class="form-input" required />
        </div>
        <div class="form-group col-span-2 sm:col-span-1">
          <label class="form-label">Lugares</label>
          <input v-model="spots" type="number" class="form-input" min="1" required />
        </div>
        <div class="form-group col-span-2 sm:col-span-1 flex items-end">
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

.grid { display: grid; }
.grid-cols-4 { grid-template-columns: repeat(4, 1fr); }
.gap-3 { gap: var(--spacing-3); }
.col-span-2 { grid-column: span 2; }
@media (min-width: 640px) {
  .sm\:col-span-1 { grid-column: span 1; }
}
.flex { display: flex; }
.items-end { align-items: flex-end; }

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
