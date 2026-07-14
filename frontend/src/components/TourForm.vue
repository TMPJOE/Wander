<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthState } from '../composables/useAuthState'
import { useApi } from '../composables/useApi'

const props = defineProps<{
  initialData?: any
  loading?: boolean
}>()

const emit = defineEmits<{
  submit: [data: any]
  cancel: []
}>()

const authState = useAuthState()
const api = useApi()

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
    }
  }
})

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
  emit('submit', form.value)
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
</style>
