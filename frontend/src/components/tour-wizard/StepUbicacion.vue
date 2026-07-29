<script setup lang="ts">
import LocationPickerMap from '../LocationPickerMap.vue'

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

function handleNext() {
  emit('next')
}

function handlePrev() {
  emit('prev')
}
</script>

<template>
  <div class="step-ubicacion card">
    <h2 class="step-title">Ubicación</h2>

    <div class="grid-2 gap-4">
      <div class="form-group">
        <label class="form-label">Ubicación (Ciudad/Zona)</label>
        <input
          v-model="form.location"
          type="text"
          class="form-input"
          required
          @update:model-value="(val) => updateField('location', val)"
        />
      </div>

      <div class="form-group">
        <label class="form-label">Punto de encuentro</label>
        <input
          v-model="form.meeting_point"
          type="text"
          class="form-input"
          required
          @update:model-value="(val) => updateField('meeting_point', val)"
        />
      </div>
    </div>

    <!-- Location Picker Map -->
    <div class="form-group">
      <label class="form-label">Ubicación exacta en el mapa</label>
      <p class="text-sm mb-2" style="color: var(--color-text-muted)">
        Seleccione el punto de encuentro exacto en el mapa. Esto ayudará a los viajeros a encontrar
        el lugar fácilmente.
      </p>
      <LocationPickerMap
        v-model:label="form.meeting_point"
        v-model:latitude="form.latitude"
        v-model:longitude="form.longitude"
      />
    </div>

    <div class="flex justify-between mt-6">
      <button type="button" class="btn btn-ghost" @click="handlePrev">Atrás</button>
      <button type="button" class="btn btn-primary" @click="handleNext">Siguiente</button>
    </div>
  </div>
</template>

<style scoped>
.step-ubicacion {
  padding: var(--spacing-6);
}

.step-title {
  font-size: var(--font-size-xl);
  font-weight: var(--font-weight-bold);
  margin-bottom: var(--spacing-4);
}

.grid-2 {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
}

@media (max-width: 600px) {
  .grid-2 {
    grid-template-columns: 1fr;
  }
}

.gap-4 {
  gap: var(--spacing-4);
}

.mb-2 {
  margin-bottom: var(--spacing-2);
}

.mt-6 {
  margin-top: var(--spacing-6);
}

.flex {
  display: flex;
}

.justify-between {
  justify-content: space-between;
}

.text-sm {
  font-size: var(--font-size-sm);
}
</style>
