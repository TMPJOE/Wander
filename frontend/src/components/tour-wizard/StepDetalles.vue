<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  form: any
}>()

const emit = defineEmits<{
  'update:form': [value: any]
  next: []
}>()

function updateField(field: string, value: any) {
  emit('update:form', { ...props.form, [field]: value })
}

const newIncluded = computed({
  get: () => '',
  set: (val: string) => {
    if (val.trim()) {
      emit('update:form', {
        ...props.form,
        what_included: [...props.form.what_included, val.trim()],
      })
    }
  },
})

const newLanguage = computed({
  get: () => '',
  set: (val: string) => {
    if (val.trim() && !props.form.languages.includes(val.trim())) {
      emit('update:form', {
        ...props.form,
        languages: [...props.form.languages, val.trim()],
      })
    }
  },
})

function removeIncluded(index: number) {
  const updated = [...props.form.what_included]
  updated.splice(index, 1)
  emit('update:form', { ...props.form, what_included: updated })
}

function removeLanguage(index: number) {
  const updated = [...props.form.languages]
  updated.splice(index, 1)
  emit('update:form', { ...props.form, languages: updated })
}

function handleNext() {
  emit('next')
}
</script>

<template>
  <div class="step-detalles card">
    <h2 class="step-title">Detalles del Tour</h2>

    <div class="form-group">
      <label class="form-label">Título del Tour</label>
      <input
        v-model="form.title"
        type="text"
        class="form-input"
        required
        @update:model-value="(val) => updateField('title', val)"
      />
    </div>

    <div class="form-group">
      <label class="form-label">Descripción</label>
      <textarea
        v-model="form.description"
        class="form-input form-textarea"
        required
        @update:model-value="(val) => updateField('description', val)"
      ></textarea>
    </div>

    <div class="grid-2 gap-4">
      <div class="form-group">
        <label class="form-label">Categoría</label>
        <!-- Categories will be loaded dynamically - for now using placeholder -->
        <select
          v-model="form.category_id"
          class="form-input"
          required
          @change="(e: any) => updateField('category_id', Number(e.target.value))"
        >
          <option :value="1">Aventura</option>
          <option :value="2">Cultural</option>
          <option :value="3">Gastronomía</option>
          <option :value="4">Senderismo</option>
          <option :value="5">Histórico</option>
          <option :value="6">Naturaleza</option>
          <option :value="7">Vida nocturna</option>
          <option :value="8">Fotografía</option>
          <option :value="9">Agua</option>
        </select>
      </div>

      <div class="form-group">
        <label class="form-label">Dificultad</label>
        <select
          v-model="form.difficulty"
          class="form-input"
          required
          @change="(e: any) => updateField('difficulty', e.target.value)"
        >
          <option value="easy">Fácil</option>
          <option value="moderate">Moderado</option>
          <option value="challenging">Desafiante</option>
          <option value="extreme">Extremo</option>
        </select>
      </div>
    </div>

    <div class="grid-3 gap-4">
      <div class="form-group">
        <label class="form-label">Precio (PAB)</label>
        <input
          v-model.number="form.price_per_person"
          type="number"
          class="form-input"
          required
          min="0"
          @change="(e: any) => updateField('price_per_person', Number(e.target.value))"
        />
      </div>

      <div class="form-group">
        <label class="form-label">Duración (minutos)</label>
        <input
          v-model.number="form.duration_minutes"
          type="number"
          class="form-input"
          required
          min="15"
          @change="(e: any) => updateField('duration_minutes', Number(e.target.value))"
        />
      </div>

      <div class="form-group">
        <label class="form-label">Máx. Personas</label>
        <input
          v-model.number="form.max_guests"
          type="number"
          class="form-input"
          required
          min="1"
          @change="(e: any) => updateField('max_guests', Number(e.target.value))"
        />
      </div>
    </div>

    <div class="form-group">
      <label class="form-label">Idiomas</label>
      <div class="flex gap-2 mb-2">
        <input
          v-model="newLanguage"
          type="text"
          class="form-input"
          placeholder="Ej: Inglés"
          @keydown.enter.prevent="($event.target as HTMLInputElement).value = ''"
        />
        <button type="button" class="btn btn-outline" @click="newLanguage = 'dummy'">
          Agregar
        </button>
      </div>
      <div class="flex flex-wrap gap-2">
        <span
          v-for="(lang, idx) in form.languages"
          :key="idx"
          class="badge badge-secondary cursor-pointer"
          @click="removeLanguage(idx)"
        >
          {{ lang }} &times;
        </span>
      </div>
    </div>

    <div class="form-group">
      <label class="form-label">¿Qué incluye?</label>
      <div class="flex gap-2 mb-2">
        <input
          v-model="newIncluded"
          type="text"
          class="form-input"
          placeholder="Ej: Equipo de seguridad"
          @keydown.enter.prevent="($event.target as HTMLInputElement).value = ''"
        />
        <button type="button" class="btn btn-outline" @click="newIncluded = 'dummy'">
          Agregar
        </button>
      </div>
      <ul class="list-disc pl-5">
        <li v-for="(item, idx) in form.what_included" :key="idx" class="text-sm flex justify-between">
          {{ item }}
          <button type="button" class="text-error" @click="removeIncluded(idx)">&times;</button>
        </li>
      </ul>
    </div>

    <div class="flex justify-end mt-6">
      <button type="button" class="btn btn-primary" @click="handleNext">Siguiente</button>
    </div>
  </div>
</template>

<style scoped>
.step-detalles {
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

.grid-3 {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
}

@media (max-width: 600px) {
  .grid-2,
  .grid-3 {
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

.flex-wrap {
  flex-wrap: wrap;
}

.justify-between {
  justify-content: space-between;
}

.justify-end {
  justify-content: flex-end;
}

.text-sm {
  font-size: var(--font-size-sm);
}

.text-error {
  color: var(--color-error);
}

.list-disc {
  list-style-type: disc;
}

.pl-5 {
  padding-left: 1.25rem;
}

.cursor-pointer {
  cursor: pointer;
}
</style>
