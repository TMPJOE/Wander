<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useAuthState } from '../composables/useAuthState'
import { useApi } from '../composables/useApi'
import { useConfirm } from '../composables/useConfirm'
import LocationPickerMap from './LocationPickerMap.vue'
import { MapPin, Plus, Trash2, ArrowUp, ArrowDown } from '@lucide/vue'

const props = defineProps<{
  initialData?: any
  loading?: boolean
}>()

const emit = defineEmits<{
  submit: [data: any]
  cancel: []
  change: [data: any]
}>()

const authState = useAuthState()
const api = useApi()
const { confirm } = useConfirm()

const categories = ref<any[]>([])
const form = ref({
  title: '',
  description: '',
  category_id: 1,
  location: '',
  latitude: 0,
  longitude: 0,
  duration_minutes: 120,
  price_per_person: 0,
  max_guests: 10,
  difficulty: 'moderate',
  languages: ['Español'],
  what_included: [] as string[],
  meeting_point: '',
  images: [] as string[],
  is_published: true,
  itinerary: [] as Array<{
    id?: number
    _localId: string
    sort_order: number
    title: string
    description: string
    duration_minutes: number | null
    location_label: string
    latitude: number | null
    longitude: number | null
  }>,
})

const newIncluded = ref('')
const newLanguage = ref('')
const isUploadingImages = ref(false)
const uploadError = ref('')

onMounted(async () => {
  try {
    const catsRes = await api.get('/categories')
    categories.value = catsRes.data || []
  } catch (e) {
    console.error('Failed to load categories', e)
  }

  if (props.initialData) {
    const d = props.initialData
    form.value = {
      title: d.title || '',
      description: d.description || '',
      category_id: d.category_id || 1,
      location: d.location || '',
      latitude: d.latitude || 0,
      longitude: d.longitude || 0,
      duration_minutes: d.duration_minutes || 120,
      price_per_person: d.price_per_person || 0,
      max_guests: d.max_guests || 10,
      difficulty: d.difficulty || 'moderate',
      languages: d.languages || ['Español'],
      what_included: d.what_included || [],
      meeting_point: d.meeting_point || '',
      images: d.images || [],
      is_published: d.is_published !== undefined ? d.is_published : true,
      itinerary: (d.itinerary || []).map((item: any, idx: number) => ({
        id: item.id,
        _localId: item.id?.toString() || `local-${Date.now()}-${idx}`,
        sort_order: item.sort_order ?? idx,
        title: item.title || '',
        description: item.description || '',
        duration_minutes: item.duration_minutes ?? null,
        location_label: item.location_label || '',
        latitude: item.latitude ?? null,
        longitude: item.longitude ?? null,
      })),
    }
  }
})

// Live-emit form state so parents (e.g. TourFormView preview) can capture it.
watch(
  form,
  (val) =>
    emit('change', {
      ...val,
      what_included: [...val.what_included],
      languages: [...val.languages],
      images: [...val.images],
    }),
  { deep: true, immediate: true },
)

function addIncluded() {
  if (newIncluded.value.trim()) {
    form.value.what_included.push(newIncluded.value.trim())
    newIncluded.value = ''
  }
}

function removeIncluded(index: number) {
  form.value.what_included.splice(index, 1)
}

async function handleImageSelect(event: Event) {
  const input = event.target as HTMLInputElement
  const files = Array.from(input.files || [])

  if (!files.length) {
    return
  }

  isUploadingImages.value = true
  uploadError.value = ''

  try {
    const uploadedUrls = await Promise.all(
      files.map(async (file) => {
        const formData = new FormData()
        formData.append('image', file)

        const response = await fetch('/api/v1/uploads', {
          method: 'POST',
          headers: {
            Authorization: `Bearer ${authState.token.value || ''}`,
          },
          body: formData,
        })

        const payload = await response.json().catch(() => ({}))
        if (!response.ok) {
          throw new Error(payload?.message || 'No se pudo subir la imagen')
        }

        return payload?.data?.url || ''
      }),
    )

    form.value.images.push(...uploadedUrls.filter(Boolean))
  } catch (error) {
    console.error(error)
    uploadError.value = 'No se pudieron subir algunas imágenes.'
  } finally {
    isUploadingImages.value = false
    input.value = ''
  }
}

function removeImage(index: number) {
  form.value.images.splice(index, 1)
}

function addLanguage() {
  if (newLanguage.value.trim() && !form.value.languages.includes(newLanguage.value.trim())) {
    form.value.languages.push(newLanguage.value.trim())
    newLanguage.value = ''
  }
}

function removeLanguage(index: number) {
  form.value.languages.splice(index, 1)
}

function handleSubmit() {
  // Prepare itinerary for submission (remove local IDs)
  const itineraryPayload = form.value.itinerary.map((item) => ({
    id: item.id,
    sort_order: item.sort_order,
    title: item.title,
    description: item.description,
    duration_minutes: item.duration_minutes,
    location_label: item.location_label,
    latitude: item.latitude,
    longitude: item.longitude,
  }))

  emit('submit', {
    ...form.value,
    what_included: [...form.value.what_included],
    languages: [...form.value.languages],
    images: [...form.value.images],
    itinerary: itineraryPayload,
  })
}

function addItineraryStep() {
  if (form.value.itinerary.length >= 20) return

  const newStep = {
    _localId: `local-${Date.now()}-${Math.random().toString(36).slice(2, 9)}`,
    sort_order: form.value.itinerary.length,
    title: '',
    description: '',
    duration_minutes: null,
    location_label: '',
    latitude: null,
    longitude: null,
  }

  form.value.itinerary.push(newStep)
}

async function removeStep(index: number) {
  const confirmed = await confirm({
    title: 'Eliminar paso del itinerario',
    body: '¿Estás seguro de que deseas eliminar este paso? Esta acción no se puede deshacer.',
    confirmLabel: 'Eliminar',
    cancelLabel: 'Cancelar',
    confirmVariant: 'danger',
  })

  if (!confirmed) return

  form.value.itinerary.splice(index, 1)

  // Re-index sort_order
  form.value.itinerary.forEach((step, idx) => {
    step.sort_order = idx
  })
}

function moveStep(index: number, direction: -1 | 1) {
  const newIndex = index + direction

  if (newIndex < 0 || newIndex >= form.value.itinerary.length) return

  // Swap items
  const temp = form.value.itinerary[index]
  form.value.itinerary[index] = form.value.itinerary[newIndex]!
  form.value.itinerary[newIndex] = temp!

  // Update sort_order
  form.value.itinerary.forEach((step, idx) => {
    step.sort_order = idx
  })
}
</script>

<template>
  <form @submit.prevent="handleSubmit" class="tour-form card">
    <div class="form-group">
      <label class="form-label">Título del Tour</label>
      <input v-model="form.title" type="text" class="form-input" required />
    </div>

    <div class="form-group">
      <label class="form-label">Descripción</label>
      <textarea v-model="form.description" class="form-input form-textarea" required></textarea>
    </div>

    <div class="grid-2 gap-4">
      <div class="form-group">
        <label class="form-label">Categoría</label>
        <select v-model="form.category_id" class="form-input" required>
          <option v-for="cat in categories" :key="cat.id" :value="cat.id">
            {{ cat.name }}
          </option>
        </select>
      </div>

      <div class="form-group">
        <label class="form-label">Dificultad</label>
        <select v-model="form.difficulty" class="form-input" required>
          <option value="easy">Fácil</option>
          <option value="moderate">Moderado</option>
          <option value="challenging">Desafiante</option>
          <option value="extreme">Extremo</option>
        </select>
      </div>
    </div>

    <div class="grid-2 gap-4">
      <div class="form-group">
        <label class="form-label">Ubicación (Ciudad/Zona)</label>
        <input v-model="form.location" type="text" class="form-input" required />
      </div>

      <div class="form-group">
        <label class="form-label">Punto de encuentro</label>
        <input v-model="form.meeting_point" type="text" class="form-input" required />
      </div>
    </div>

    <div class="grid-3 gap-4">
      <div class="form-group">
        <label class="form-label">Precio (PAB)</label>
        <input v-model="form.price_per_person" type="number" class="form-input" required min="0" />
      </div>

      <div class="form-group">
        <label class="form-label">Duración (minutos)</label>
        <input v-model="form.duration_minutes" type="number" class="form-input" required min="15" />
      </div>

      <div class="form-group">
        <label class="form-label">Máx. Personas</label>
        <input v-model="form.max_guests" type="number" class="form-input" required min="1" />
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
          @keydown.enter.prevent="addLanguage"
        />
        <button type="button" class="btn btn-outline" @click="addLanguage">Agregar</button>
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
          @keydown.enter.prevent="addIncluded"
        />
        <button type="button" class="btn btn-outline" @click="addIncluded">Agregar</button>
      </div>
      <ul class="list-disc pl-5">
        <li
          v-for="(item, idx) in form.what_included"
          :key="idx"
          class="text-sm flex justify-between"
        >
          {{ item }}
          <button type="button" class="text-error" @click="removeIncluded(idx)">&times;</button>
        </li>
      </ul>
    </div>

    <div class="form-group">
      <label class="form-label">Imágenes del tour</label>
      <p class="text-sm mb-2" style="color: var(--color-text-muted)">
        Sube imágenes desde tu dispositivo. Se añadirán automáticamente al tour.
      </p>
      <input type="file" accept="image/*" multiple class="form-input" @change="handleImageSelect" />
      <p v-if="isUploadingImages" class="text-sm mt-2">Subiendo imágenes...</p>
      <p v-if="uploadError" class="text-sm mt-2 text-error">{{ uploadError }}</p>
      <div v-if="form.images.length" class="grid grid-cols-4 gap-2 mt-2">
        <div v-for="(img, idx) in form.images" :key="idx" class="relative group">
          <img :src="img" class="w-full aspect-square object-cover rounded-md" />
          <button
            type="button"
            class="absolute top-1 right-1 bg-error text-white rounded-full w-5 h-5 flex items-center justify-center text-xs"
            @click="removeImage(idx)"
          >
            &times;
          </button>
        </div>
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

    <!-- Itinerary Section -->
    <div class="form-group itinerary-section">
      <div class="itinerary-header">
        <h3 class="itinerary-title">
          <MapPin :size="20" />
          Itinerario del Tour
        </h3>
        <button
          type="button"
          class="btn btn-primary btn-md"
          @click="addItineraryStep"
          :disabled="form.itinerary.length >= 20"
        >
          <Plus :size="16" />
          Agregar paso
        </button>
      </div>

      <p class="text-sm mb-3" style="color: var(--color-text-muted)">
        Describe la secuencia ordenada de paradas durante el tour. Esto es diferente de "Qué
        incluye": aquí detallas la ruta y actividades en orden.
      </p>

      <div v-if="form.itinerary.length === 0" class="itinerary-empty-state">
        <p>No hay pasos en el itinerario aún.</p>
        <p class="text-sm" style="color: var(--color-text-muted)">
          Haz clic en "Agregar paso" para comenzar a crear el itinerario.
        </p>
      </div>

      <div v-for="(step, index) in form.itinerary" :key="step._localId" class="itinerary-step card">
        <div class="itinerary-step-header">
          <span class="itinerary-step-number">Paso {{ index + 1 }}</span>
          <div class="itinerary-step-actions">
            <button
              type="button"
              class="btn-icon"
              @click="moveStep(index, -1)"
              :disabled="index === 0"
              title="Mover arriba"
            >
              <ArrowUp :size="16" />
            </button>
            <button
              type="button"
              class="btn-icon"
              @click="moveStep(index, 1)"
              :disabled="index === form.itinerary.length - 1"
              title="Mover abajo"
            >
              <ArrowDown :size="16" />
            </button>
            <button
              type="button"
              class="btn-icon btn-icon-danger"
              @click="removeStep(index)"
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

        <div class="grid-2 gap-4">
          <div class="form-group">
            <label class="form-label">Etiqueta de ubicación (opcional)</label>
            <input
              v-model="step.location_label"
              type="text"
              class="form-input"
              placeholder="Ej: Plaza Mayor"
            />
          </div>

          <div class="form-group">
            <label class="form-label">Coordenadas (opcional)</label>
            <div class="coords-inputs">
              <input
                v-model.number="step.latitude"
                type="number"
                class="form-input"
                placeholder="Latitud"
                step="any"
              />
              <input
                v-model.number="step.longitude"
                type="number"
                class="form-input"
                placeholder="Longitud"
                step="any"
              />
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="form-group flex items-center gap-2 mt-4">
      <input type="checkbox" id="is_published" v-model="form.is_published" />
      <label for="is_published" class="form-label mb-0" style="margin-bottom: 0"
        >Publicar Tour inmediatamente</label
      >
    </div>

    <div class="flex justify-end gap-3 mt-6">
      <button type="button" class="btn btn-ghost" @click="$emit('cancel')">Cancelar</button>
      <button type="submit" class="btn btn-primary" :disabled="loading">
        {{ loading ? 'Guardando...' : 'Guardar Tour' }}
      </button>
    </div>
  </form>
</template>

<style scoped>
.tour-form {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-4);
  padding: var(--spacing-6);
}

.grid-2 {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
}
.grid-3 {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
}
.grid-cols-4 {
  grid-template-columns: repeat(4, minmax(0, 1fr));
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
.mt-2 {
  margin-top: var(--spacing-2);
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
.items-center {
  align-items: center;
}
.text-sm {
  font-size: var(--font-size-sm);
}
.text-error {
  color: var(--color-error);
}
.bg-error {
  background: var(--color-error);
}
.text-white {
  color: white;
}
.rounded-md {
  border-radius: var(--radius-md);
}
.rounded-full {
  border-radius: var(--radius-full);
}
.cursor-pointer {
  cursor: pointer;
}
.relative {
  position: relative;
}
.absolute {
  position: absolute;
}
.top-1 {
  top: 4px;
}
.right-1 {
  right: 4px;
}
.w-5 {
  width: 20px;
}
.h-5 {
  height: 20px;
}
.w-full {
  width: 100%;
}
.aspect-square {
  aspect-ratio: 1/1;
}
.object-cover {
  object-fit: cover;
}
.list-disc {
  list-style-type: disc;
}
.pl-5 {
  padding-left: 1.25rem;
}

/* Itinerary section styles */
.itinerary-section {
  margin-top: var(--spacing-6);
  padding: var(--spacing-4);
  background: var(--color-background);
  border-radius: var(--radius-md);
  border: 1px solid var(--color-border-light);
}

.itinerary-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-3);
}

.itinerary-title {
  display: flex;
  align-items: center;
  gap: var(--spacing-2);
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-primary);
}

.itinerary-empty-state {
  padding: var(--spacing-6);
  text-align: center;
  color: var(--color-text-muted);
  background: var(--color-surface);
  border-radius: var(--radius-md);
  border: 1px dashed var(--color-border);
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

.coords-inputs {
  display: flex;
  gap: var(--spacing-2);
}

.coords-inputs .form-input {
  flex: 1;
}

.mb-3 {
  margin-bottom: var(--spacing-3);
}

@media (max-width: 600px) {
  .itinerary-header {
    flex-direction: column;
    align-items: flex-start;
    gap: var(--spacing-2);
  }

  .itinerary-step-header {
    flex-wrap: wrap;
    gap: var(--spacing-2);
  }

  .coords-inputs {
    flex-direction: column;
  }
}
</style>
