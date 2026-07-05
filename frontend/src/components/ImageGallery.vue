<script setup lang="ts">
import { ref, computed, watch, onBeforeUnmount } from 'vue'
import { ChevronLeft, ChevronRight, Maximize2, X } from '@lucide/vue'

const props = defineProps<{
  images: string[]
}>()

const current = ref(0)
const isExpanded = ref(false)

const hasMultiple = computed(() => props.images.length > 1)
const activeImage = computed(() => {
  return props.images[current.value] || 'https://images.unsplash.com/photo-1501785888041-af3ef285b470?w=800&h=500&fit=crop'
})

function next() {
  if (current.value < props.images.length - 1) current.value++
  else current.value = 0
}

function prev() {
  if (current.value > 0) current.value--
  else current.value = props.images.length - 1
}

function openPreview() {
  if (props.images.length) isExpanded.value = true
}

function closePreview() {
  isExpanded.value = false
}

watch(isExpanded, (open) => {
  if (typeof document === 'undefined') return
  document.body.style.overflow = open ? 'hidden' : ''
  document.body.style.touchAction = open ? 'none' : ''
})

onBeforeUnmount(() => {
  if (typeof document === 'undefined') return
  document.body.style.overflow = ''
  document.body.style.touchAction = ''
})
</script>

<template>
  <div class="gallery">
    <div class="gallery__viewport" @click="openPreview">
      <img
        :src="activeImage"
        :alt="`Imagen ${current + 1}`"
        class="gallery__image"
      />
      <button class="gallery__expand" @click.stop="openPreview" aria-label="Ampliar galería">
        <Maximize2 :size="16" />
      </button>
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

    <div v-if="isExpanded" class="gallery__overlay" @click.self="closePreview">
      <div class="gallery__overlay-shell">
        <button class="gallery__overlay-close" @click="closePreview" aria-label="Cerrar galería">
          <X :size="24" />
        </button>

        <button v-if="hasMultiple" class="gallery__nav gallery__nav--overlay gallery__nav--prev" @click.stop="prev">
          <ChevronLeft :size="24" />
        </button>
        <button v-if="hasMultiple" class="gallery__nav gallery__nav--overlay gallery__nav--next" @click.stop="next">
          <ChevronRight :size="24" />
        </button>

        <img :src="activeImage" :alt="`Imagen ${current + 1}`" class="gallery__overlay-image" />

        <div class="gallery__overlay-footer">
          <div class="gallery__dots gallery__dots--overlay">
            <span
              v-for="(_, i) in images"
              :key="i"
              class="gallery__dot"
              :class="{ 'gallery__dot--active': i === current }"
              @click.stop="current = i"
            ></span>
          </div>
          <div class="gallery__counter gallery__counter--overlay" v-if="hasMultiple">
            {{ current + 1 }} / {{ images.length }}
          </div>
        </div>
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
  cursor: zoom-in;
}

.gallery__image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: opacity var(--transition-base);
}

.gallery__expand {
  position: absolute;
  top: var(--spacing-3);
  left: var(--spacing-3);
  width: 36px;
  height: 36px;
  border-radius: var(--radius-full);
  background: rgba(255, 255, 255, 0.9);
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

.gallery__viewport:hover .gallery__expand,
.gallery__viewport:focus-within .gallery__expand {
  opacity: 1;
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

.gallery__nav--overlay {
  opacity: 1;
  z-index: 2;
}

.gallery__dots {
  position: absolute;
  bottom: var(--spacing-3);
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  gap: var(--spacing-1);
}

.gallery__dots--overlay {
  position: static;
  transform: none;
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

.gallery__counter--overlay {
  position: static;
  background: rgba(255, 255, 255, 0.14);
  backdrop-filter: blur(8px);
  border: 1px solid rgba(255, 255, 255, 0.16);
}

.gallery__overlay {
  position: fixed;
  inset: 0;
  z-index: 1200;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: var(--spacing-4);
  background: rgba(7, 12, 24, 0.94);
}

.gallery__overlay-shell {
  position: relative;
  width: min(100%, 960px);
  max-height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.gallery__overlay-image {
  width: 100%;
  max-height: min(84vh, 860px);
  object-fit: contain;
  border-radius: var(--radius-lg);
  background: rgba(255, 255, 255, 0.04);
}

.gallery__overlay-close {
  position: absolute;
  top: var(--spacing-3);
  right: var(--spacing-3);
  width: 40px;
  height: 40px;
  border-radius: var(--radius-full);
  border: none;
  background: rgba(255, 255, 255, 0.14);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  backdrop-filter: blur(8px);
  cursor: pointer;
  z-index: 2;
}

.gallery__overlay-footer {
  position: absolute;
  bottom: var(--spacing-3);
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  align-items: center;
  gap: var(--spacing-2);
  padding: var(--spacing-2) var(--spacing-3);
  border-radius: var(--radius-full);
  background: rgba(0, 0, 0, 0.35);
  backdrop-filter: blur(8px);
}

@media (max-width: 768px) {
  .gallery__nav {
    opacity: 1;
  }

  .gallery__expand {
    opacity: 1;
  }

  .gallery__overlay {
    padding: var(--spacing-2);
  }
}
</style>
