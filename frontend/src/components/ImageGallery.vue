<script setup lang="ts">
import { ref, computed } from 'vue';
import { ChevronLeft, ChevronRight } from '@lucide/vue';

const props = defineProps<{
  images: string[];
}>();

const current = ref(0);

const hasMultiple = computed(() => props.images.length > 1);

function next() {
  if (current.value < props.images.length - 1) current.value++;
  else current.value = 0;
}

function prev() {
  if (current.value > 0) current.value--;
  else current.value = props.images.length - 1;
}
</script>

<template>
  <div class="gallery">
    <div class="gallery__viewport">
      <img
        :src="images[current] || 'https://images.unsplash.com/photo-1501785888041-af3ef285b470?w=800&h=500&fit=crop'"
        :alt="`Imagen ${current + 1}`"
        class="gallery__image"
      />
      <template v-if="hasMultiple">
        <button class="gallery__nav gallery__nav--prev" @click.stop="prev">
          <ChevronLeft :size="20" />
        </button>
        <button class="gallery__nav gallery__nav--next" @click.stop="next">
          <ChevronRight :size="20" />
        </button>
        <div class="gallery__dots">
          <span
            v-for="(_, i) in images"
            :key="i"
            class="gallery__dot"
            :class="{ 'gallery__dot--active': i === current }"
            @click.stop="current = i"
          ></span>
        </div>
      </template>
      <div class="gallery__counter" v-if="hasMultiple">
        {{ current + 1 }} / {{ images.length }}
      </div>
    </div>
  </div>
</template>

<style scoped>
.gallery {
  width: 100%;
}

.gallery__viewport {
  position: relative;
  aspect-ratio: 16 / 10;
  overflow: hidden;
  border-radius: 0;
  background: var(--color-border-light);
}

.gallery__image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: opacity var(--transition-base);
}

.gallery__nav {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  width: 36px;
  height: 36px;
  border-radius: var(--radius-full);
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-text);
  box-shadow: var(--shadow-md);
  transition: all var(--transition-fast);
  opacity: 0;
  border: none;
  cursor: pointer;
}

.gallery__viewport:hover .gallery__nav {
  opacity: 1;
}

.gallery__nav:hover {
  background: white;
  transform: translateY(-50%) scale(1.05);
}

.gallery__nav--prev {
  left: var(--spacing-3);
}

.gallery__nav--next {
  right: var(--spacing-3);
}

.gallery__dots {
  position: absolute;
  bottom: var(--spacing-3);
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  gap: var(--spacing-1);
}

.gallery__dot {
  width: 6px;
  height: 6px;
  border-radius: var(--radius-full);
  background: rgba(255, 255, 255, 0.5);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.gallery__dot--active {
  background: white;
  width: 18px;
}

.gallery__counter {
  position: absolute;
  top: var(--spacing-3);
  right: var(--spacing-3);
  background: rgba(0, 0, 0, 0.5);
  color: white;
  padding: 2px var(--spacing-2);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-xs);
  backdrop-filter: blur(4px);
}
</style>
