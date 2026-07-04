<script setup lang="ts">
import StarRating from './StarRating.vue'

defineProps<{
  review: {
    id: number
    rating: number
    title?: string
    comment: string
    created_at: string
    tour_title?: string
    user_name?: string
    user_avatar?: string
  }
}>()

function timeAgo(dateStr: string): string {
  const diff = Date.now() - new Date(dateStr).getTime()
  const mins = Math.floor(diff / 60000)
  if (mins < 60) return `Hace ${mins} min`
  const hours = Math.floor(mins / 60)
  if (hours < 24) return `Hace ${hours}h`
  const days = Math.floor(hours / 24)
  if (days < 30) return `Hace ${days}d`
  const months = Math.floor(days / 30)
  return `Hace ${months} mes${months > 1 ? 'es' : ''}`
}
</script>

<template>
  <div class="review-card">
    <div class="review-card__header">
      <img
        v-if="review.user_avatar"
        :src="review.user_avatar"
        :alt="review.user_name"
        class="review-card__avatar"
      />
      <div v-else class="review-card__avatar review-card__avatar--placeholder">
        {{ (review.user_name || 'U').charAt(0).toUpperCase() }}
      </div>
      <div class="review-card__info">
        <span class="review-card__name">{{ review.user_name || 'Usuario' }}</span>
        <span class="review-card__date">{{ timeAgo(review.created_at) }}</span>
      </div>
      <StarRating :rating="review.rating" :size="14" />
    </div>
    <p v-if="review.tour_title" class="review-card__tour-title">{{ review.tour_title }}</p>
    <p v-if="review.title" class="review-card__title">{{ review.title }}</p>
    <p class="review-card__comment">{{ review.comment }}</p>
  </div>
</template>

<style scoped>
.review-card {
  padding: var(--spacing-4);
  border-bottom: 1px solid var(--color-divider);
}

.review-card:last-child {
  border-bottom: none;
}

.review-card__header {
  display: flex;
  align-items: center;
  gap: var(--spacing-3);
  margin-bottom: var(--spacing-3);
}

.review-card__avatar {
  width: 36px;
  height: 36px;
  border-radius: var(--radius-full);
  object-fit: cover;
  flex-shrink: 0;
}

.review-card__avatar--placeholder {
  background: var(--color-secondary-50);
  color: var(--color-secondary);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: var(--font-weight-semibold);
  font-size: var(--font-size-sm);
}

.review-card__info {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.review-card__tour-title {
  font-size: var(--font-size-xs);
  color: var(--color-text-secondary);
  margin-bottom: var(--spacing-1);
  text-transform: uppercase;
  letter-spacing: var(--letter-spacing-wide);
}

.review-card__title {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-semibold);
  margin-bottom: var(--spacing-2);
}

.review-card__name {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-semibold);
}

.review-card__date {
  font-size: var(--font-size-xs);
  color: var(--color-text-light);
}

.review-card__comment {
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
  line-height: var(--line-height-relaxed);
}
</style>
