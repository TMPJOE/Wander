<script setup lang="ts">
import { computed } from 'vue';
import { Star } from '@lucide/vue';

const props = defineProps<{
  rating: number;
  max?: number;
  size?: number;
  interactive?: boolean;
}>();

const emit = defineEmits<{
  rate: [value: number];
}>();

const max = computed(() => props.max || 5);
const size = computed(() => props.size || 16);

function starType(index: number): 'full' | 'half' | 'empty' {
  if (index <= Math.floor(props.rating)) return 'full';
  if (index === Math.ceil(props.rating) && props.rating % 1 >= 0.25) return 'half';
  return 'empty';
}
</script>

<template>
  <div
    class="star-rating"
    :class="{ 'star-rating--interactive': interactive }"
    role="img"
    :aria-label="`${rating} de ${max} estrellas`"
  >
    <button
      v-for="i in max"
      :key="i"
      class="star-rating__star"
      :class="`star-rating__star--${starType(i)}`"
      :style="{ width: `${size}px`, height: `${size}px` }"
      :disabled="!interactive"
      type="button"
      @click="interactive && emit('rate', i)"
    >
      <Star :size="size" :stroke-width="starType(i) === 'empty' ? 1.5 : 0" :fill="starType(i) !== 'empty' ? 'currentColor' : 'none'" />
    </button>
  </div>
</template>

<style scoped>
.star-rating {
  display: inline-flex;
  align-items: center;
  gap: 1px;
}

.star-rating__star {
  color: var(--color-star-empty);
  transition: color var(--transition-fast), transform var(--transition-fast);
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.star-rating__star--full,
.star-rating__star--half {
  color: var(--color-star);
}

.star-rating--interactive .star-rating__star {
  cursor: pointer;
}

.star-rating--interactive .star-rating__star:hover {
  transform: scale(1.2);
  color: var(--color-star);
}
</style>
