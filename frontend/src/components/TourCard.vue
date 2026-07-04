<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { Heart, MapPin, Clock, Star } from '@lucide/vue'
import { useFavoritesStore } from '../stores/favorites'
import { useAuthStore } from '../stores/auth'

const props = defineProps<{
  tour: {
    id: number
    title: string
    location: string
    price_per_person: number
    duration_minutes: number
    images: string | string[]
    avg_rating: number
    review_count: number
    difficulty?: string
    category_name?: string
    is_favorited?: boolean
    guide_name?: string
  }
  allowLike?: boolean
}>()

const router = useRouter()
const favoritesStore = useFavoritesStore()
const authStore = useAuthStore()

const isFavorited = computed(() => {
  const id = String(props.tour.id)
  return favoritesStore.favorites.some((f: any) => {
    const favId = String(f.tour_id ?? f.tour?.id ?? f.id ?? '')
    return favId === id
  })
})

const allowLike = computed(() => {
  return props.allowLike !== false
})

import { ref } from 'vue'
const justLiked = ref(false)

const imageUrl = computed(() => {
  const imgs =
    typeof props.tour.images === 'string'
      ? JSON.parse(props.tour.images || '[]')
      : props.tour.images || []
  return (
    imgs[0] || 'https://images.unsplash.com/photo-1501785888041-af3ef285b470?w=400&h=300&fit=crop'
  )
})

const formattedDuration = computed(() => {
  const h = Math.floor(props.tour.duration_minutes / 60)
  const m = props.tour.duration_minutes % 60
  if (h === 0) return `${m} min`
  return m ? `${h}h ${m}m` : `${h}h`
})

const formattedPrice = computed(() => {
  return `$${props.tour.price_per_person.toLocaleString('es-MX')}`
})

const difficultyLabel: Record<string, string> = {
  easy: 'Fácil',
  moderate: 'Moderado',
  challenging: 'Desafiante',
  extreme: 'Extremo',
}

function goToDetail() {
  router.push(`/tours/${props.tour.id}`)
}

async function toggleFavorite(e: Event) {
  e.stopPropagation()
  if (!authStore.isAuthenticated) {
    router.push({ name: 'login' })
    return
  }

  const currentlyFavorited = isFavorited.value || !!props.tour.is_favorited

  // If liking is disabled (e.g., Profile view) and the item is not already favorited,
  // prevent adding a favorite. Allow unliking always.
  if (!allowLike.value && !currentlyFavorited) return

  // Optimistic local animation for like
  if (!currentlyFavorited) {
    justLiked.value = true
    setTimeout(() => (justLiked.value = false), 700)
  }

  try {
    const ok = await favoritesStore.toggleFavorite(String(props.tour.id), currentlyFavorited)
    ok // no-op to satisfy linter
  } catch (e) {
    // rollback animation if API failed and surface error for debugging
    console.error('Failed toggling favorite', e)
    justLiked.value = false
  }
}
</script>

<template>
  <article class="tour-card" @click="goToDetail">
    <div class="tour-card__image-wrap">
      <img :src="imageUrl" :alt="tour.title" class="tour-card__image" loading="lazy" />
      <button
        :class="[
          'tour-card__fav',
          { 'tour-card__fav--liked': isFavorited, 'animate-like': justLiked },
        ]"
        @click="toggleFavorite"
      >
        <Heart
          :size="20"
          :stroke-width="isFavorited ? 0 : 2"
          :fill="isFavorited ? 'var(--color-primary)' : 'none'"
          :color="isFavorited ? 'var(--color-primary)' : 'white'"
        />
      </button>
      <span v-if="tour.difficulty" class="tour-card__badge">
        {{ difficultyLabel[tour.difficulty] || tour.difficulty }}
      </span>
    </div>

    <div class="tour-card__body">
      <div class="tour-card__meta">
        <span class="tour-card__location">
          <MapPin :size="12" :stroke-width="2" />
          {{ tour.location }}
        </span>
        <span v-if="tour.avg_rating > 0" class="tour-card__rating">
          <Star :size="12" :stroke-width="0" fill="var(--color-star)" />
          {{ tour.avg_rating.toFixed(1) }}
          <span class="tour-card__review-count">({{ tour.review_count }})</span>
        </span>
      </div>

      <h3 class="tour-card__title">{{ tour.title }}</h3>

      <div class="tour-card__footer">
        <span class="tour-card__duration">
          <Clock :size="13" :stroke-width="1.8" />
          {{ formattedDuration }}
        </span>
        <span class="tour-card__price">
          {{ formattedPrice }}
          <span class="tour-card__price-unit">/ persona</span>
        </span>
      </div>
    </div>
  </article>
</template>

<style scoped>
.tour-card {
  background: var(--color-surface);
  border-radius: var(--radius-xl);
  overflow: hidden;
  cursor: pointer;
  transition:
    transform var(--transition-fast),
    box-shadow var(--transition-fast);
  box-shadow: var(--shadow-sm);
}

.tour-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.tour-card:active {
  transform: translateY(0);
}

.tour-card__image-wrap {
  position: relative;
  aspect-ratio: 4 / 3;
  overflow: hidden;
}

.tour-card__image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform var(--transition-slow);
}

.tour-card:hover .tour-card__image {
  transform: scale(1.04);
}

.tour-card__fav {
  position: absolute;
  top: var(--spacing-3);
  right: var(--spacing-3);
  width: 36px;
  height: 36px;
  border-radius: var(--radius-full);
  background: rgba(0, 0, 0, 0.25);
  backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all var(--transition-fast);
  border: none;
  cursor: pointer;
}

.tour-card__fav:hover {
  background: rgba(0, 0, 0, 0.4);
  transform: scale(1.1);
}

.tour-card__fav--liked {
  background: rgba(208, 83, 66, 0.12);
}

@keyframes heart-pop {
  0% {
    transform: scale(1);
  }
  40% {
    transform: scale(1.35);
  }
  70% {
    transform: scale(0.95);
  }
  100% {
    transform: scale(1);
  }
}

.animate-like {
  animation: heart-pop 0.7s cubic-bezier(0.2, 0.9, 0.3, 1);
}

.tour-card__badge {
  position: absolute;
  bottom: var(--spacing-3);
  left: var(--spacing-3);
  background: rgba(0, 0, 0, 0.55);
  backdrop-filter: blur(8px);
  color: white;
  padding: 3px var(--spacing-2);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-medium);
}

.tour-card__body {
  padding: var(--spacing-3) var(--spacing-4) var(--spacing-4);
  display: flex;
  flex-direction: column;
  gap: var(--spacing-2);
}

.tour-card__meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.tour-card__location {
  display: inline-flex;
  align-items: center;
  gap: 3px;
  font-size: var(--font-size-xs);
  color: var(--color-text-light);
}

.tour-card__rating {
  display: inline-flex;
  align-items: center;
  gap: 3px;
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text);
}

.tour-card__review-count {
  color: var(--color-text-light);
  font-weight: var(--font-weight-normal);
}

.tour-card__title {
  font-size: var(--font-size-base);
  font-weight: var(--font-weight-semibold);
  line-height: var(--line-height-tight);
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.tour-card__footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: var(--spacing-1);
}

.tour-card__duration {
  display: inline-flex;
  align-items: center;
  gap: 3px;
  font-size: var(--font-size-xs);
  color: var(--color-text-light);
}

.tour-card__price {
  font-size: var(--font-size-base);
  font-weight: var(--font-weight-bold);
  color: var(--color-primary);
}

.tour-card__price-unit {
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-normal);
  color: var(--color-text-light);
}
</style>
