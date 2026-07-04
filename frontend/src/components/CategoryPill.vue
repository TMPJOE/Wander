<script setup lang="ts">
defineProps<{
  name: string
  slug: string
  icon: string
  active?: boolean
}>()

defineEmits<{
  select: [slug: string]
}>()

import foodImg from '../assets/categories/food.png'
import culturalImg from '../assets/categories/cultural.png'
import historicalImg from '../assets/categories/historical.png'
import adventureImg from '../assets/categories/adventure.png'
import nightlifeImg from '../assets/categories/nightlife.png'
import natureImg from '../assets/categories/nature.png'
import hikingImg from '../assets/categories/hiking.png'
import waterImg from '../assets/categories/water.png'
import photogImg from '../assets/categories/wellness.png'

const SLUG_TO_IMAGE: Record<string, string> = {
  gastronomia: foodImg,
  'cultura-historia': culturalImg,
  cultura: culturalImg,
  historia: historicalImg,
  aventura: adventureImg,
  'vida-nocturna': nightlifeImg,
  naturaleza: natureImg,
  fotografia: photogImg,
  senderismo: hikingImg,
  agua: waterImg,
}

const FALLBACK_IMAGE = waterImg

function getCategoryImage(slug: string): string {
  return SLUG_TO_IMAGE[slug] || FALLBACK_IMAGE
}
</script>

<template>
  <button
    class="category-card"
    :class="{ 'category-card--active': active }"
    @click="$emit('select', slug)"
  >
    <div
      class="category-card__background"
      :style="{ backgroundImage: `url(${getCategoryImage(slug)})` }"
    ></div>
    <div class="category-card__vignette"></div>
    <div class="category-card__content">
      <span class="category-card__label">{{ name }}</span>
    </div>
  </button>
</template>

<style scoped>
.category-card {
  position: relative;
  width: 20%;
  min-width: 80px;
  height: 120px;
  border-radius: var(--radius-xl);
  border: 2px solid transparent;
  overflow: hidden;
  cursor: pointer;
  transition: all var(--transition-normal);
  flex-shrink: 0;
}

.category-card__background {
  position: absolute;
  inset: 0;
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  transition: transform var(--transition-slow);
}

.category-card__vignette {
  position: absolute;
  inset: 0;
  background: radial-gradient(
    ellipse at center,
    transparent 0%,
    transparent 40%,
    rgba(0, 0, 0, 0.4) 100%
  );
  pointer-events: none;
}

.category-card__content {
  position: relative;
  z-index: 1;
  display: flex;
  align-items: flex-end;
  justify-content: center;
  width: 100%;
  height: 100%;
  padding: var(--spacing-6);
}

.category-card__label {
  font-size: var(--font-size-xl);
  font-weight: var(--font-weight-bold);
  color: var(--color-text-inverse);
  text-transform: capitalize;
  text-shadow: 0 2px 8px rgba(0, 0, 0, 0.6);
  text-align: center;
  width: 100%;
}

.category-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 24px rgba(0, 0, 0, 0.3);
}

.category-card:hover .category-card__background {
  transform: scale(1.05);
}

.category-card--active {
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px rgba(var(--color-primary-rgb), 0.3);
}

.category-card--active .category-card__vignette {
  background: radial-gradient(
    ellipse at center,
    transparent 0%,
    transparent 30%,
    rgba(0, 0, 0, 0.6) 100%
  );
}

@media (max-width: 768px) {
  .category-card {
    min-width: 240px;
    height: 120px;
  }

  .category-card__label {
    font-size: var(--font-size-xl);
  }
}
</style>
