<script setup lang="ts">
import { ref, watch } from 'vue';
import { X, SlidersHorizontal } from '@lucide/vue';

const props = defineProps<{
  open: boolean;
}>();

const emit = defineEmits<{
  close: [];
  apply: [filters: FilterValues];
}>();

export interface FilterValues {
  difficulty: string;
  min_price: string;
  max_price: string;
  location: string;
}

const difficulty = ref('');
const minPrice = ref('');
const maxPrice = ref('');
const location = ref('');

const difficulties = [
  { value: '', label: 'Todos' },
  { value: 'easy', label: 'Fácil' },
  { value: 'moderate', label: 'Moderado' },
  { value: 'challenging', label: 'Desafiante' },
  { value: 'extreme', label: 'Extremo' },
];

function apply() {
  emit('apply', {
    difficulty: difficulty.value,
    min_price: minPrice.value,
    max_price: maxPrice.value,
    location: location.value,
  });
  emit('close');
}

function reset() {
  difficulty.value = '';
  minPrice.value = '';
  maxPrice.value = '';
  location.value = '';
}

watch(() => props.open, (val) => {
  if (val) {
    document.body.style.overflow = 'hidden';
  } else {
    document.body.style.overflow = '';
  }
});
</script>

<template>
  <Teleport to="body">
    <Transition name="fade">
      <div v-if="open" class="drawer-overlay" @click="$emit('close')"></div>
    </Transition>
    <Transition name="slide-drawer">
      <div v-if="open" class="drawer">
        <div class="drawer__header">
          <h3 class="drawer__title">
            <SlidersHorizontal :size="18" />
            Filtros
          </h3>
          <button class="drawer__close" @click="$emit('close')">
            <X :size="20" />
          </button>
        </div>

        <div class="drawer__body">
          <div class="filter-section">
            <label class="filter-label">Dificultad</label>
            <div class="filter-chips">
              <button
                v-for="d in difficulties"
                :key="d.value"
                class="filter-chip"
                :class="{ 'filter-chip--active': difficulty === d.value }"
                @click="difficulty = d.value"
              >
                {{ d.label }}
              </button>
            </div>
          </div>

          <div class="filter-section">
            <label class="filter-label">Rango de precio</label>
            <div class="price-range">
              <input
                v-model="minPrice"
                type="number"
                class="form-input"
                placeholder="Mín"
                min="0"
              />
              <span class="price-separator">—</span>
              <input
                v-model="maxPrice"
                type="number"
                class="form-input"
                placeholder="Máx"
                min="0"
              />
            </div>
          </div>

          <div class="filter-section">
            <label class="filter-label">Ubicación</label>
            <input
              v-model="location"
              type="text"
              class="form-input"
              placeholder="Ej: Ciudad de México"
            />
          </div>
        </div>

        <div class="drawer__footer">
          <button class="btn btn-outline" @click="reset">Limpiar</button>
          <button class="btn btn-primary flex-1" @click="apply">Aplicar Filtros</button>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.drawer-overlay {
  position: fixed;
  inset: 0;
  background: var(--color-overlay);
  z-index: var(--z-drawer);
}

.drawer {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  max-height: 85vh;
  background: var(--color-surface);
  border-radius: var(--radius-2xl) var(--radius-2xl) 0 0;
  z-index: calc(var(--z-drawer) + 1);
  display: flex;
  flex-direction: column;
  box-shadow: var(--shadow-xl);
}

.drawer__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--spacing-5) var(--spacing-6);
  border-bottom: 1px solid var(--color-divider);
}

.drawer__title {
  display: flex;
  align-items: center;
  gap: var(--spacing-2);
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
}

.drawer__close {
  color: var(--color-text-light);
  padding: var(--spacing-1);
}

.drawer__body {
  padding: var(--spacing-6);
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: var(--spacing-6);
}

.filter-section {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-3);
}

.filter-label {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text);
}

.filter-chips {
  display: flex;
  flex-wrap: wrap;
  gap: var(--spacing-2);
}

.filter-chip {
  padding: var(--spacing-2) var(--spacing-3);
  border-radius: var(--radius-full);
  border: 1.5px solid var(--color-border);
  font-size: var(--font-size-sm);
  background: transparent;
  cursor: pointer;
  transition: all var(--transition-fast);
}

.filter-chip:hover {
  border-color: var(--color-primary-light);
}

.filter-chip--active {
  background: var(--color-primary);
  color: white;
  border-color: var(--color-primary);
}

.price-range {
  display: flex;
  align-items: center;
  gap: var(--spacing-3);
}

.price-separator {
  color: var(--color-text-light);
  flex-shrink: 0;
}

.drawer__footer {
  display: flex;
  gap: var(--spacing-3);
  padding: var(--spacing-4) var(--spacing-6);
  border-top: 1px solid var(--color-divider);
}

/* Transitions */
.slide-drawer-enter-active,
.slide-drawer-leave-active {
  transition: transform var(--transition-base);
}

.slide-drawer-enter-from,
.slide-drawer-leave-to {
  transform: translateY(100%);
}
</style>
