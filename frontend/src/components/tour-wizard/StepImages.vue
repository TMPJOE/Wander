<script setup lang="ts">
import { ref } from 'vue'
import { useConfirm } from '@/composables/useConfirm'
import { useApi } from '@/composables/useApi'
import { Upload, X, ImageIcon } from '@lucide/vue'
import WizardActions from './WizardActions.vue'

const props = defineProps<{
  form: {
    images: string[]
  }
}>()

const emit = defineEmits<{
  next: []
  prev: []
  submit: [data: any]
  'update:form': [data: any]
}>()

const { confirm } = useConfirm()
const api = useApi()
const isUploading = ref(false)
const uploadError = ref('')

async function handleImageUpload(event: Event) {
  const target = event.target as HTMLInputElement
  const files = target.files
  if (!files || files.length === 0) return

  isUploading.value = true
  uploadError.value = ''

  try {
    const uploadedUrls: string[] = []
    for (let i = 0; i < files.length; i++) {
      const file = files[i]
      if (!file) continue
      const formData = new FormData()
      formData.append('image', file)
      const res = await api.post('/uploads', formData)
      const url: string | undefined = res.data?.url
      if (url) uploadedUrls.push(url)
    }

    const updatedImages = [...props.form.images, ...uploadedUrls]
    emit('update:form', { ...props.form, images: updatedImages })
  } catch (err) {
    console.error('Image upload error', err)
    uploadError.value = 'Error al subir imágenes. Intente nuevamente.'
  } finally {
    isUploading.value = false
    // Allow selecting the same file again later
    target.value = ''
  }
}

async function removeImage(url: string) {
  const confirmed = await confirm({ title: '¿Eliminar imagen?', body: 'Esta acción no se puede deshacer.', confirmVariant: 'danger' })
  if (!confirmed) return

  const updatedImages = props.form.images.filter((img) => img !== url)
  emit('update:form', { ...props.form, images: updatedImages })
}
</script>

<template>
  <div class="step-imagenes">
    <h2 class="step-title">Imágenes del Tour</h2>
    <p class="step-description">
      Sube fotos atractivas de tu experiencia. Las imágenes ayudan a los viajeros a visualizar la
      aventura.
    </p>

    <!-- Upload Section -->
    <div class="upload-section">
      <label class="upload-label" :class="{ 'upload-label--disabled': isUploading }">
        <input
          type="file"
          accept="image/*"
          multiple
          :disabled="isUploading"
          @change="handleImageUpload"
        />
        <Upload class="upload-icon" :size="32" />
        <span class="upload-text">
          {{ isUploading ? 'Subiendo...' : 'Click para subir imágenes o arrastra aquí' }}
        </span>
      </label>
      <p v-if="uploadError" class="upload-error">{{ uploadError }}</p>
    </div>

    <!-- Image Grid -->
    <div v-if="form.images.length > 0" class="image-grid">
      <div v-for="(img, idx) in form.images" :key="idx" class="image-card">
        <img :src="img" :alt="'Imagen ' + (idx + 1)" class="image-card__img" />
        <button
          type="button"
          class="image-card__remove"
          @click="removeImage(img)"
          title="Eliminar imagen"
        >
          <X :size="16" />
        </button>
        <span class="image-card__index">{{ idx + 1 }}</span>
      </div>
    </div>

    <div v-else class="no-images">
      <ImageIcon :size="28" />
      <p>No hay imágenes subidas aún.</p>
    </div>

    <!-- Navigation -->
    <WizardActions step="middle" @prev="emit('prev')" @next="emit('next')" />
  </div>
</template>

<style scoped>
.step-imagenes {
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

.upload-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--spacing-2);
}

.upload-label {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: 100%;
  max-width: 400px;
  padding: var(--spacing-6);
  border: 2px dashed var(--color-border);
  border-radius: var(--radius-lg);
  cursor: pointer;
  transition: all var(--transition-fast);
  background: var(--color-background);
}

.upload-label:hover:not(.upload-label--disabled) {
  border-color: var(--color-primary);
  background: var(--color-primary-50);
}

.upload-label--disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.upload-label input {
  display: none;
}

.upload-icon {
  color: var(--color-primary);
  margin-bottom: var(--spacing-2);
}

.upload-text {
  font-size: var(--font-size-base);
  color: var(--color-text-secondary);
}

.upload-error {
  color: var(--color-danger);
  font-size: var(--font-size-sm);
}

.image-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: var(--spacing-3);
  width: 100%;
}

.image-card {
  position: relative;
  aspect-ratio: 1;
  border-radius: var(--radius-lg);
  overflow: hidden;
  box-shadow: var(--shadow-sm);
}

.image-card__img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.image-card__remove {
  position: absolute;
  top: var(--spacing-1);
  right: var(--spacing-1);
  width: 24px;
  height: 24px;
  border-radius: var(--radius-full);
  background: rgba(0, 0, 0, 0.7);
  color: white;
  border: none;
  line-height: 1;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background var(--transition-fast);
}

.image-card__remove:hover {
  background: var(--color-danger);
}

.image-card__index {
  position: absolute;
  bottom: var(--spacing-1);
  left: var(--spacing-1);
  background: rgba(0, 0, 0, 0.7);
  color: white;
  font-size: var(--font-size-xs);
  padding: 2px 6px;
  border-radius: var(--radius-md);
}

.no-images {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--spacing-2);
  text-align: center;
  padding: var(--spacing-8);
  color: var(--color-text-secondary);
}
</style>
