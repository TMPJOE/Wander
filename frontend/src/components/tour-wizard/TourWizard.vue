<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import StepDetalles from './StepDetalles.vue'
import StepUbicacion from './StepUbicacion.vue'
import StepItinerario from './StepItinerario.vue'
import StepImagenes from './StepImagenes.vue'
import StepRevision from './StepRevision.vue'

const props = defineProps<{
  initialData?: any
}>()

const emit = defineEmits<{
  submit: [data: any]
  change: [data: any]
}>()

const route = useRoute()
const currentStep = ref(0)

// Shared form state - same shape as TourForm.vue
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
  languages: ['Español'] as string[],
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
  }>,
})

const steps = [
  { name: 'Detalles', component: StepDetalles },
  { name: 'Ubicación', component: StepUbicacion },
  { name: 'Itinerario', component: StepItinerario },
  { name: 'Imágenes', component: StepImagenes },
  { name: 'Revisión', component: StepRevision },
]

const currentComponent = computed(() => steps[currentStep.value].component)

const progress = computed(() => ((currentStep.value + 1) / steps.length) * 100)

// Initialize form from initialData or sessionStorage
onMounted(() => {
  if (props.initialData) {
    // Edit path: use API data
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
      })),
    }
  } else {
    // New-tour path: check sessionStorage for preview restoration
    const previewReturn = sessionStorage.getItem('wander_preview_return')
    if (previewReturn === route.fullPath) {
      const previewData = sessionStorage.getItem('wander_tour_preview')
      if (previewData) {
        try {
          const parsed = JSON.parse(previewData)
          form.value = {
            ...form.value,
            ...parsed,
            what_included: [...(parsed.what_included || [])],
            languages: [...(parsed.languages || ['Español'])],
            images: [...(parsed.images || [])],
            itinerary: (parsed.itinerary || []).map((item: any, idx: number) => ({
              ...item,
              _localId: item._localId || `local-${Date.now()}-${idx}`,
            })),
          }
        } catch (e) {
          console.error('Failed to parse preview data', e)
        }
      }
      // Consume the return path but keep preview data for subsequent previews
      sessionStorage.removeItem('wander_preview_return')
    }
  }
})

// Watch form changes and emit to parent
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

function nextStep() {
  if (currentStep.value < steps.length - 1) {
    currentStep.value++
  }
}

function prevStep() {
  if (currentStep.value > 0) {
    currentStep.value--
  }
}

function goToStep(index: number) {
  currentStep.value = index
}

function handleSubmit() {
  const itineraryPayload = form.value.itinerary.map((item) => ({
    id: item.id,
    sort_order: item.sort_order,
    title: item.title,
    description: item.description,
    duration_minutes: item.duration_minutes,
    location_label: item.location_label,
  }))

  emit('submit', {
    ...form.value,
    what_included: [...form.value.what_included],
    languages: [...form.value.languages],
    images: [...form.value.images],
    itinerary: itineraryPayload,
  })
}

// Expose form and methods to child steps via provide/inject pattern or slots
defineExpose({
  form,
  currentStep,
  steps,
  nextStep,
  prevStep,
  goToStep,
  handleSubmit,
})
</script>

<template>
  <div class="tour-wizard">
    <!-- Progress Header -->
    <header class="wizard-header">
      <div class="wizard-progress">
        <div class="wizard-progress__bar" :style="{ width: `${progress}%` }"></div>
      </div>
      <nav class="wizard-steps">
        <button
          v-for="(step, idx) in steps"
          :key="step.name"
          type="button"
          class="wizard-step"
          :class="{
            'wizard-step--active': idx === currentStep,
            'wizard-step--completed': idx < currentStep,
          }"
          @click="goToStep(idx)"
        >
          <span class="wizard-step__number">{{ idx + 1 }}</span>
          <span class="wizard-step__name">{{ step.name }}</span>
        </button>
      </nav>
    </header>

    <!-- Step Content -->
    <div class="wizard-content">
      <component
        :is="currentComponent"
        :form="form"
        @update:form="(val: any) => (form.value = val)"
        @next="nextStep"
        @prev="prevStep"
        @submit="handleSubmit"
      />
    </div>
  </div>
</template>

<style scoped>
.tour-wizard {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-4);
}

.wizard-header {
  padding: var(--spacing-4);
  background: var(--color-surface);
  border-bottom: 1px solid var(--color-border-light);
}

.wizard-progress {
  height: 4px;
  background: var(--color-background);
  border-radius: var(--radius-full);
  overflow: hidden;
  margin-bottom: var(--spacing-3);
}

.wizard-progress__bar {
  height: 100%;
  background: var(--color-primary);
  transition: width var(--transition-normal);
}

.wizard-steps {
  display: flex;
  gap: var(--spacing-2);
  overflow-x: auto;
  padding-bottom: var(--spacing-1);
}

.wizard-step {
  display: flex;
  align-items: center;
  gap: var(--spacing-2);
  padding: var(--spacing-2) var(--spacing-3);
  background: var(--color-background);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  cursor: pointer;
  transition: all var(--transition-fast);
  white-space: nowrap;
}

.wizard-step--active {
  border-color: var(--color-primary);
  background: var(--color-primary-50);
}

.wizard-step--completed {
  border-color: var(--color-primary);
  background: var(--color-primary);
  color: white;
}

.wizard-step__number {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border-radius: var(--radius-full);
  background: var(--color-background);
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-semibold);
}

.wizard-step--active .wizard-step__number,
.wizard-step--completed .wizard-step__number {
  background: var(--color-primary);
  color: white;
}

.wizard-step--completed .wizard-step__number {
  background: white;
  color: var(--color-primary);
}

.wizard-step__name {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
}

.wizard-content {
  padding: var(--spacing-4) 0;
}
</style>
