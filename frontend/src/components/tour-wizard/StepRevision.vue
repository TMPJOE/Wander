<script setup lang="ts">
  import { computed } from 'vue'
  import PreviewButton from './PreviewButton.vue'

  const props = defineProps<{
    form: {
      title: string
      description: string
      category_id: number
      location: string
      latitude: number
      longitude: number
      duration_minutes: number
      price_per_person: number
      max_guests: number
      difficulty: string
      languages: string[]
      what_included: string[]
      meeting_point: string
      images: string[]
      is_published: boolean
      itinerary: Array<{
        title: string
        description: string
        duration_minutes: number | null
        location_label: string
      }>
    }
  }>()

  const emit = defineEmits<{
    prev: []
    preview: []
    submit: [data: any]
    'update:form': [data: any]
  }>()

  const categoryNames: Record<number, string> = {
    1: 'Aventura',
    2: 'Cultural',
    3: 'Gastronómico',
    4: 'Senderismo',
    5: 'Histórico',
    6: 'Naturaleza',
    7: 'Vida Nocturna',
    8: 'Fotografía',
    9: 'Agua',
  }

  const difficultyLabels: Record<string, string> = {
    easy: 'Fácil',
    moderate: 'Moderado',
    challenging: 'Desafiante',
    extreme: 'Extremo',
  }

  const durationDisplay = computed(() => {
    const mins = props.form.duration_minutes
    const hours = Math.floor(mins / 60)
    const remaining = mins % 60
    if (hours > 0 && remaining > 0) {
      return `${hours}h ${remaining}m`
    } else if (hours > 0) {
      return `${hours}h`
    } else {
      return `${mins}m`
    }
  })

  const hasRequiredFields = computed(() => {
    return (
      props.form.title.trim() !== '' &&
      props.form.description.trim() !== '' &&
      props.form.category_id > 0 &&
      props.form.price_per_person >= 0
    )
  })

  function handlePrev() {
    emit('prev')
  }

  function handleSubmit() {
    emit('submit', props.form)
  }

  function togglePublished() {
    emit('update:form', {
      ...props.form,
      is_published: !props.form.is_published,
    })
  }
</script>

<template>
  <div class="step-revision">
    <h2 class="step-title">Revisión Final</h2>
    <p class="step-description">Revisa toda la información antes de guardar tu tour.</p>

    <!-- Summary Cards -->
    <div class="summary-grid">
      <!-- Detalles -->
      <div class="summary-card">
        <h3 class="summary-card__title">Detalles</h3>
        <div class="summary-row">
          <span class="label">Título:</span>
          <span class="value">{{ form.title || '—' }}</span>
        </div>
        <div class="summary-row">
          <span class="label">Categoría:</span>
          <span class="value">{{ categoryNames[form.category_id] || '—' }}</span>
        </div>
        <div class="summary-row">
          <span class="label">Dificultad:</span>
          <span class="value">{{ difficultyLabels[form.difficulty] || '—' }}</span>
        </div>
        <div class="summary-row">
          <span class="label">Duración:</span>
          <span class="value">{{ durationDisplay }}</span>
        </div>
        <div class="summary-row">
          <span class="label">Precio:</span>
          <span class="value">${{ form.price_per_person }}</span>
        </div>
        <div class="summary-row">
          <span class="label">Máx. viajeros:</span>
          <span class="value">{{ form.max_guests }}</span>
        </div>
      </div>

      <!-- Ubicación -->
      <div class="summary-card">
        <h3 class="summary-card__title">Ubicación</h3>
        <div class="summary-row">
          <span class="label">Ciudad/Zona:</span>
          <span class="value">{{ form.location || '—' }}</span>
        </div>
        <div class="summary-row">
          <span class="label">Punto de encuentro:</span>
          <span class="value">{{ form.meeting_point || '—' }}</span>
        </div>
        <div v-if="form.latitude && form.longitude" class="summary-row">
          <span class="label">Coordenadas:</span>
          <span class="value">{{ form.latitude.toFixed(4) }}, {{ form.longitude.toFixed(4) }}</span>
        </div>
      </div>

      <!-- Idiomas -->
      <div class="summary-card">
        <h3 class="summary-card__title">Idiomas</h3>
        <div class="tags-list">
          <span v-for="lang in form.languages" :key="lang" class="tag">
            {{ lang }}
          </span>
        </div>
      </div>

      <!-- Qué incluye -->
      <div class="summary-card">
        <h3 class="summary-card__title">Qué Incluye</h3>
        <ul v-if="form.what_included.length > 0" class="checklist">
          <li v-for="(item, idx) in form.what_included" :key="idx">
            {{ item }}
          </li>
        </ul>
        <p v-else class="empty-text">No especificado</p>
      </div>

      <!-- Itinerario -->
      <div class="summary-card summary-card--full">
        <h3 class="summary-card__title">Itinerario</h3>
        <div v-if="form.itinerary.length > 0" class="itinerary-list">
          <div v-for="(item, idx) in form.itinerary" :key="idx" class="itinerary-item">
            <div class="itinerary-item__header">
              <span class="itinerary-item__number">{{ idx + 1 }}</span>
              <h4 class="itinerary-item__title">{{ item.title }}</h4>
              <span v-if="item.duration_minutes" class="itinerary-item__duration">
                {{ item.duration_minutes }}m
              </span>
            </div>
            <p class="itinerary-item__description">{{ item.description }}</p>
            <p v-if="item.location_label" class="itinerary-item__location">
              📍 {{ item.location_label }}
            </p>
          </div>
        </div>
        <p v-else class="empty-text">No especificado</p>
      </div>

      <!-- Imágenes -->
      <div class="summary-card summary-card--full">
        <h3 class="summary-card__title">Imágenes</h3>
        <div v-if="form.images.length > 0" class="image-preview-grid">
          <img
            v-for="(img, idx) in form.images"
            :key="idx"
            :src="img"
            :alt="'Imagen ' + (idx + 1)"
            class="image-preview"
          />
        </div>
        <p v-else class="empty-text">No hay imágenes</p>
      </div>
    </div>

    <!-- Publish Toggle -->
    <div class="publish-section">
      <label class="toggle-label">
        <input type="checkbox" :checked="form.is_published" @change="togglePublished" />
        <span class="toggle-text">
          {{ form.is_published ? 'Publicado (visible para viajeros)' : 'Borrador (oculto)' }}
        </span>
      </label>
    </div>

    <!-- Actions -->
    <div class="step-actions">
      <button type="button" class="btn btn-outline" @click="handlePrev">← Atrás</button>
      <div class="step-actions__right">
        <button
          type="button"
          class="btn btn-secondary"
          :disabled="!hasRequiredFields"
          @click="emit('preview')"
         >
          👁 Previsualizar como viajero
        </button>
        <button
          type="button"
          class="btn btn-primary"
          :disabled="!hasRequiredFields"
          @click="handleSubmit"
        >
          Guardar Tour
        </button>
      </div>
    </div>

    <p v-if="!hasRequiredFields" class="validation-error">
      Completa los campos requeridos (título, descripción, categoría y precio) para continuar.
    </p>
  </div>
</template>

<style scoped>
  .step-revision {
    display: flex;
    flex-direction: column;
    gap: var(--spacing-4);
    padding: var(--spacing-4);
  }

  .step-title {
    font-size: var(--font-size-xl);
    font-weight: var(--font-weight-bold);
    color: var(--color-text);
  }

  .step-description {
    font-size: var(--font-size-base);
    color: var(--color-text-secondary);
    margin-bottom: var(--spacing-2);
  }

  .summary-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
    gap: var(--spacing-3);
  }

  .summary-card {
    background: var(--color-surface);
    border: 1px solid var(--color-border-light);
    border-radius: var(--radius-lg);
    padding: var(--spacing-4);
  }

  .summary-card--full {
    grid-column: 1 / -1;
  }

  .summary-card__title {
    font-size: var(--font-size-lg);
    font-weight: var(--font-weight-semibold);
    color: var(--color-text);
    margin-bottom: var(--spacing-3);
    padding-bottom: var(--spacing-2);
    border-bottom: 1px solid var(--color-border);
  }

  .summary-row {
    display: flex;
    justify-content: space-between;
    padding: var(--spacing-1) 0;
    font-size: var(--font-size-sm);
  }

  .summary-row .label {
    color: var(--color-text-secondary);
    font-weight: var(--font-weight-medium);
  }

  .summary-row .value {
    color: var(--color-text);
    font-weight: var(--font-weight-normal);
  }

  .tags-list {
    display: flex;
    flex-wrap: wrap;
    gap: var(--spacing-2);
  }

  .tag {
    display: inline-block;
    padding: 4px 12px;
    background: var(--color-primary-50);
    color: var(--color-primary);
    border-radius: var(--radius-full);
    font-size: var(--font-size-xs);
    font-weight: var(--font-weight-medium);
  }

  .checklist {
    list-style: none;
    padding: 0;
    margin: 0;
  }

  .checklist li {
    padding: var(--spacing-1) 0;
    padding-left: var(--spacing-4);
    position: relative;
    font-size: var(--font-size-sm);
  }

  .checklist li::before {
    content: '✓';
    position: absolute;
    left: 0;
    color: var(--color-success);
    font-weight: var(--font-weight-bold);
  }

  .empty-text {
    color: var(--color-text-secondary);
    font-size: var(--font-size-sm);
    font-style: italic;
  }

  .itinerary-list {
    display: flex;
    flex-direction: column;
    gap: var(--spacing-3);
  }

  .itinerary-item {
    background: var(--color-background);
    border-radius: var(--radius-md);
    padding: var(--spacing-3);
  }

  .itinerary-item__header {
    display: flex;
    align-items: center;
    gap: var(--spacing-2);
    margin-bottom: var(--spacing-2);
  }

  .itinerary-item__number {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 24px;
    height: 24px;
    background: var(--color-primary);
    color: white;
    border-radius: var(--radius-full);
    font-size: var(--font-size-xs);
    font-weight: var(--font-weight-bold);
  }

  .itinerary-item__title {
    font-size: var(--font-size-base);
    font-weight: var(--font-weight-semibold);
    color: var(--color-text);
    flex: 1;
  }

  .itinerary-item__duration {
    font-size: var(--font-size-xs);
    color: var(--color-text-secondary);
    background: var(--color-background);
    padding: 2px 8px;
    border-radius: var(--radius-md);
  }

  .itinerary-item__description {
    font-size: var(--font-size-sm);
    color: var(--color-text-secondary);
    margin-bottom: var(--spacing-1);
  }

  .itinerary-item__location {
    font-size: var(--font-size-xs);
    color: var(--color-text-secondary);
  }

  .image-preview-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
    gap: var(--spacing-2);
  }

  .image-preview {
    width: 100%;
    aspect-ratio: 1;
    object-fit: cover;
    border-radius: var(--radius-md);
    box-shadow: var(--shadow-sm);
  }

  .publish-section {
    background: var(--color-surface);
    border: 1px solid var(--color-border-light);
    border-radius: var(--radius-lg);
    padding: var(--spacing-4);
    margin-top: var(--spacing-2);
  }

  .toggle-label {
    display: flex;
    align-items: center;
    gap: var(--spacing-2);
    cursor: pointer;
  }

  .toggle-label input {
    width: 18px;
    height: 18px;
    cursor: pointer;
  }

  .toggle-text {
    font-size: var(--font-size-base);
    font-weight: var(--font-weight-medium);
    color: var(--color-text);
  }

  .step-actions {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: var(--spacing-4);
  }

  .step-actions__right {
    display: flex;
    gap: var(--spacing-2);
  }

  .validation-error {
    color: var(--color-danger);
    font-size: var(--font-size-sm);
    text-align: center;
    margin-top: var(--spacing-2);
  }
</style>
