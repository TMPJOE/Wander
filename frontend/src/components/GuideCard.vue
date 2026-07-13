<script setup lang="ts">
import { MessageCircle, Star, CheckCircle } from '@lucide/vue';
import { useRouter } from 'vue-router';
import { useApi } from '../composables/useApi';
import { ref, onMounted } from 'vue';

const router = useRouter();
const api = useApi();

const props = defineProps<{
  guide: {
    id: number;
    first_name: string;
    last_name: string;
    avatar_url?: string;
    bio?: string;
    languages?: string[];
  };
}>();

defineEmits<{
  message: [guideId: number];
}>();

const rating = ref(0);

onMounted(async () => {
  try {
    const toursRes = await api.get(`/tours?guide_id=${props.guide.id}`);
    const tours = toursRes.data || [];
    let totalRating = 0;
    let totalReviews = 0;
    tours.forEach((t: any) => {
      totalRating += (t.avg_rating || 0) * (t.review_count || 0);
      totalReviews += (t.review_count || 0);
    });
    rating.value = totalReviews > 0 ? totalRating / totalReviews : 0;
  } catch (e) {
    console.error(e);
  }
});
</script>

<template>
  <div class="guide-card" @click="router.push(`/guide/profile/${guide.id}`)">
    <div class="guide-card__top">
      <div class="guide-card__avatar-container">
        <img
          v-if="guide.avatar_url"
          :src="guide.avatar_url"
          :alt="`${guide.first_name} ${guide.last_name}`"
          class="guide-card__avatar"
        />
        <div v-else class="guide-card__avatar guide-card__avatar--placeholder">
          {{ guide.first_name[0] }}{{ guide.last_name[0] }}
        </div>
        <div class="guide-card__verified">
          <CheckCircle :size="12" class="text-white" />
        </div>
      </div>
      <div class="guide-card__info">
        <h4 class="guide-card__name">{{ guide.first_name }} {{ guide.last_name }}</h4>
        <div class="guide-card__rating">
          <Star :size="14" class="text-warning fill-current" />
          <span class="rating-val">{{ rating.toFixed(1) }}</span>
        </div>
      </div>
      <button class="guide-card__msg-btn" @click.stop="$emit('message', guide.id)">
        <MessageCircle :size="18" :stroke-width="1.8" />
      </button>
    </div>
    <p v-if="guide.bio" class="guide-card__bio">{{ guide.bio }}</p>
    <div v-if="guide.languages?.length" class="guide-card__langs">
      <span v-for="lang in guide.languages" :key="lang" class="guide-card__lang">
        {{ lang }}
      </span>
    </div>
  </div>
</template>

<style scoped>
.guide-card {
  background: white;
  border-radius: var(--radius-lg);
  padding: var(--spacing-4);
  display: flex;
  flex-direction: column;
  gap: var(--spacing-3);
  border: 1px solid var(--color-border-light);
  box-shadow: 0 4px 15px rgba(0,0,0,0.02);
  transition: all 0.2s ease;
  cursor: pointer;
}

.guide-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0,0,0,0.04);
}

.guide-card__top {
  display: flex;
  align-items: center;
  gap: var(--spacing-3);
}

.guide-card__avatar {
  width: 48px;
  height: 48px;
  border-radius: var(--radius-full);
  object-fit: cover;
  flex-shrink: 0;
  border: 2px solid var(--color-surface);
}

.guide-card__avatar--placeholder {
  background: var(--color-secondary);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: var(--font-weight-semibold);
  font-size: var(--font-size-sm);
}

.guide-card__info {
  flex: 1;
}

.guide-card__name {
  font-size: var(--font-size-base);
  font-weight: var(--font-weight-semibold);
}

.guide-card__label {
  font-size: var(--font-size-xs);
  color: var(--color-secondary);
  font-weight: var(--font-weight-medium);
}

.guide-card__msg-btn {
  width: 40px;
  height: 40px;
  border-radius: var(--radius-full);
  background: var(--color-surface);
  color: var(--color-secondary);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: var(--shadow-sm);
  transition: all var(--transition-fast);
}

.guide-card__msg-btn:hover {
  background: var(--color-secondary);
  color: white;
}

.guide-card__bio {
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
  line-height: var(--line-height-relaxed);
}

.guide-card__langs {
  display: flex;
  gap: var(--spacing-2);
  flex-wrap: wrap;
}

.guide-card__lang {
  font-size: var(--font-size-xs);
  padding: 2px var(--spacing-2);
  background: var(--color-surface);
  border-radius: var(--radius-sm);
  color: var(--color-text-secondary);
}
</style>
