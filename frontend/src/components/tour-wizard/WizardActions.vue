<script setup lang="ts">
import { computed } from 'vue'
import { ArrowLeft, ArrowRight, Eye, Save } from '@lucide/vue'

type WizardPosition = 'first' | 'middle' | 'last'

const props = withDefaults(
  defineProps<{
    step: WizardPosition
    showPreview?: boolean
    canProceed?: boolean
    submitting?: boolean
    submitLabel?: string
    nextLabel?: string
    previewLabel?: string
  }>(),
  {
    showPreview: false,
    canProceed: true,
    submitting: false,
    submitLabel: 'Guardar Tour',
    nextLabel: 'Siguiente',
    previewLabel: 'Previsualizar como viajero',
  },
)

const emit = defineEmits<{
  prev: []
  next: []
  submit: []
  preview: []
}>()

const isFirst = computed(() => props.step === 'first')
const isLast = computed(() => props.step === 'last')
const showBack = computed(() => !isFirst.value)
const showPreviewBtn = computed(() => isLast.value && props.showPreview)
const primaryLabel = computed(() => (isLast.value ? props.submitLabel : props.nextLabel))

function handlePrev() {
  emit('prev')
}
function handlePrimary() {
  if (isLast.value) {
    emit('submit')
  } else {
    emit('next')
  }
}
function handlePreview() {
  emit('preview')
}
</script>

<template>
  <div class="wizard-actions">
    <button
      v-if="showPreviewBtn"
      type="button"
      class="btn btn-secondary wizard-actions__preview"
      :disabled="!canProceed"
      @click="handlePreview"
    >
      <Eye :size="18" />
      <span>{{ previewLabel }}</span>
    </button>

    <button
      type="button"
      class="btn btn-primary wizard-actions__primary"
      :disabled="!canProceed || submitting"
      @click="handlePrimary"
    >
      <Save v-if="isLast" :size="18" :class="{ 'animate-spin': submitting }" />
      <span v-if="isLast && submitting">Guardando...</span>
      <span v-else>{{ primaryLabel }}</span>
      <ArrowRight v-if="!isLast" :size="18" />
    </button>

    <button
      v-if="showBack"
      type="button"
      class="btn btn-ghost wizard-actions__back"
      @click="handlePrev"
    >
      <ArrowLeft :size="18" />
      <span>Atrás</span>
    </button>
  </div>
</template>

<style scoped>
.wizard-actions {
  display: flex;
  flex-wrap: wrap;
  align-content: flex-end;
  align-items: center;
  gap: var(--spacing-2);
  justify-content: space-between;
  margin-top: var(--spacing-6);
}

.wizard-actions__preview {
  order: 1;
}

.wizard-actions__primary {
  order: 2;
  margin-left: auto;
}

.wizard-actions__back {
  order: 3;
}

/* Desktop: collapse preview+primary into a right cluster when present */
@media (min-width: 641px) {
  .wizard-actions__preview {
    margin-left: auto;
  }

  .wizard-actions__primary {
    margin-left: 0;
  }
}

/* Mobile: full-width stacked rows — primary first, back below */
@media (max-width: 640px) {
  .wizard-actions {
    flex-direction: column;
    align-items: stretch;
  }

  .wizard-actions .btn {
    width: 100%;
    margin-left: 0;
  }

  .wizard-actions__preview {
    order: 1;
  }

  .wizard-actions__primary {
    order: 2;
  }

  .wizard-actions__back {
    order: 3;
  }
}

.animate-spin {
  animation: wizard-spin 0.8s linear infinite;
}

@keyframes wizard-spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
