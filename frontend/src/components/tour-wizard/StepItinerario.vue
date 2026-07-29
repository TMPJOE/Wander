<script setup lang="ts">
import { MapPin, Plus, Trash2, ArrowUp, ArrowDown } from '@lucide/vue'

const props = defineProps<{
  form: any
}>()

const emit = defineEmits<{
  'update:form': [value: any]
  next: []
  prev: []
}>()

function updateField(field: string, value: any) {
  emit('update:form', { ...props.form, [field]: value })
}

function addItineraryStep() {
  if (props.form.itinerary.length >= 20) return
  const newStep = {
    _localId: `local-${Date.now()}-${Math.random().toString(36).slice(2, 9)}`,
    sort_order: props.form.itinerary.length,
    title: '',
    description: '',
    duration_minutes: null,
    location_label: '',
  }
  emit('update:form', {
    ...props.form,
    itinerary: [...props.form.itinerary, newStep],
  })
}

function removeStep(index: number) {
  const updated = [...props.form.itinerary]
  updated.splice(index, 1)
  updated.forEach((step, idx) => {
    step.sort_order = idx
  })
  emit('update:form', { ...props.form, itinerary: updated })
}

function moveStep(index: number, direction: -1 | 1) {
  const newIndex = index + direction
  if (newIndex < 0 || newIndex >= props.form.itinerary.length) return
  const updated = [...props.form.itinerary]
  const temp = updated[index]
  updated[index] = updated[newIndex]
  updated[newIndex] = temp
  updated.forEach((step, idx) => {
    step.sort_order = idx
  })
  emit('update:form', { ...props.form, itinerary: updated })
}

function handleNext() {
  emit('next')
}

function handlePrev() {
  emit('prev')
}
</script>

<template>
  <div class="step-itinerario card">
    <h2 class="step-title">Itinerario del Tour</h2>

    <p class="text-sm mb-3" style="color: var(--color-text-muted)">
      Describe la secuencia ordenada de paradas durante el tour. Esto es diferente de "Qué incluye":
      aquí detallas la ruta y actividades en orden.
    </p>

    <button
      v-if="form.itinerary.length < 20"
      type="button"
      class="btn btn-primary btn-md mb-3"
      @click="addItineraryStep"
    >
      <Plus :size="16" />
      Agregar paso
    </button>

    <div v-if="form.itinerary.length === 0" class="itinerary-empty-state">
      <p>No hay pasos en el itinerario aún.</p>
      <p class="text-sm" style="color: var(--color-text-muted)">
        Haz clic en "Agregar paso" para comenzar a crear el itinerario.
      </p>
    </div>

    <div v-for="(step, index) in form.itinerary" :key="step._localId" class="itinerary-step">
      <div class="itinerary-step-header">
        <span class="itinerary-step-number">Paso {{ Number(index) + 1 }}</span>
        <div class="itinerary-step-actions">
          <button
            type="button"
            class="btn-icon"
            @click="moveStep(Number(index), -1)"
            :disabled="index === 0"
            title="Mover arriba"
          >
            <ArrowUp :size="16" />
          </button>
          <button
            type="button"
            class="btn-icon"
            @click="moveStep(Number(index), 1)"
            :disabled="index === form.itinerary.length - 1"
            title="Mover abajo"
          >
            <ArrowDown :size="16" />
          </button>
          <button
            type="button"
            class="btn-icon btn-icon-danger"
            @click="removeStep(Number(index))"
            title="Eliminar paso"
          >
            <Trash2 :size="16" />
          </button>
        </div>
      </div>

      <div class="grid-2 gap-4">
        <div class="form-group">
          <label class="form-label">Título del paso</label>
          <input
            v-model="step.title"
            type="text"
            class="form-input"
            placeholder="Ej: Visita al centro histórico"
            required
          />
        </div>

        <div class="form-group">
          <label class="form-label">Duración (minutos, opcional)</label>
          <input
            v-model.number="step.duration_minutes"
            type="number"
            class="form-input"
            placeholder="Ej: 45"
            min="0"
          />
        </div>
      </div>

      <div class="form-group">
        <label class="form-label">Descripción (opcional)</label>
        <textarea
          v-model="step.description"
          class="form-input form-textarea"
          placeholder="Describe qué sucede en esta parada..."
          rows="2"
        ></textarea>
      </div>

      <div class="form-group">
        <label class="form-label">Etiqueta de ubicación (opcional)</label>
        <input
          v-model="step.location_label"
          type="text"
          class="form-input"
          placeholder="Ej: Plaza Mayor"
        />
      </div>
    </div>

    <div class="flex justify-between mt-6">
      <button type="button" class="btn btn-ghost" @click="handlePrev">Atrás</button>
      <button type="button" class="btn btn-primary" @click="handleNext">Siguiente</button>
    </div>
  </div>
</template>

<style scoped>
.step-itinerario {
  padding: var(--spacing-6);
}

.step-title {
  font-size: var(--font-size-xl);
  font-weight: var(--font-weight-bold);
  margin-bottom: var(--spacing-4);
}

.mb-3 {
  margin-bottom: var(--spacing-3);
}

.mt-6 {
  margin-top: var(--spacing-6);
}

.text-sm {
  font-size: var(--font-size-sm);
}

.flex {
  display: flex;
}

.justify-between {
  justify-content: space-between;
}

.grid-2 {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
}

.gap-4 {
  gap: var(--spacing-4);
}

@media (max-width: 600px) {
  .grid-2 {
    grid-template-columns: 1fr;
  }
}

.itinerary-empty-state {
  padding: var(--spacing-6);
  text-align: center;
  color: var(--color-text-muted);
  background: var(--color-surface);
  border-radius: var(--radius-md);
  border: 1px dashed var(--color-border);
  margin-bottom: var(--spacing-4);
}

.itinerary-step {
  margin-bottom: var(--spacing-4);
  padding: var(--spacing-4);
  background: var(--color-surface);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-md);
}

.itinerary-step:last-child {
  margin-bottom: 0;
}

.itinerary-step-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-3);
  padding-bottom: var(--spacing-2);
  border-bottom: 1px solid var(--color-border-light);
}

.itinerary-step-number {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  color: var(--color-secondary);
}

.itinerary-step-actions {
  display: flex;
  gap: var(--spacing-1);
}

.btn-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  background: var(--color-surface);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.btn-icon:hover:not(:disabled) {
  background: var(--color-background);
  border-color: var(--color-secondary);
}

.btn-icon:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.btn-icon-danger:hover:not(:disabled) {
  background: var(--color-error);
  border-color: var(--color-error);
  color: white;
}
</style>
