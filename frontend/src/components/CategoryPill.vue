<script setup lang="ts">
defineProps<{
  name: string;
  slug: string;
  icon: string;
  active?: boolean;
}>();

defineEmits<{
  select: [slug: string];
}>();

const iconMap: Record<string, string> = {
  hiking: '🥾',
  cultural: '🏛️',
  food: '🍽️',
  adventure: '🧗',
  nature: '🌿',
  photography: '📷',
  historical: '🏰',
  nightlife: '🌙',
  wellness: '🧘',
  water: '🚣',
  default: '📍',
};
</script>

<template>
  <button
    class="category-pill"
    :class="{ 'category-pill--active': active }"
    @click="$emit('select', slug)"
  >
    <span class="category-pill__icon">{{ iconMap[icon] || iconMap.default }}</span>
    <span class="category-pill__label">{{ name }}</span>
  </button>
</template>

<style scoped>
.category-pill {
  display: inline-flex;
  align-items: center;
  gap: var(--spacing-2);
  padding: var(--spacing-2) var(--spacing-4);
  border-radius: var(--radius-full);
  border: 1.5px solid var(--color-border);
  background: var(--color-surface);
  white-space: nowrap;
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  color: var(--color-text-secondary);
  transition: all var(--transition-fast);
  cursor: pointer;
  flex-shrink: 0;
}

.category-pill:hover {
  border-color: var(--color-primary-light);
  color: var(--color-primary);
}

.category-pill--active {
  background: var(--color-primary);
  color: var(--color-text-inverse);
  border-color: var(--color-primary);
}

.category-pill--active:hover {
  background: var(--color-primary-dark);
  color: var(--color-text-inverse);
}

.category-pill__icon {
  font-size: 1rem;
  line-height: 1;
}

.category-pill__label {
  line-height: 1;
}
</style>
