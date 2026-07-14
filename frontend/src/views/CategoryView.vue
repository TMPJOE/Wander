<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft } from '@lucide/vue'
import { useApi } from '../composables/useApi'
import TourCard from '../components/TourCard.vue'

const route = useRoute()
const router = useRouter()
const api = useApi()

const slug = computed(() => route.params.slug as string)
const category = ref<any>(null)
const tours = ref<any[]>([])
const loading = ref(false)

onMounted(async () => {
  loading.value = true
  try {
    const catsRes = await api.get('/categories')
    const found = (catsRes.data || []).find((c: any) => c.slug === slug.value)
    if (found) {
      category.value = found
    }

    const toursRes = await api.get('/tours', { params: { category: slug.value } })
    tours.value = toursRes.data || []
  } catch (e) {
    console.error('Failed to load category tours', e)
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="page">
    <header class="category-header container">
      <button class="back-btn" @click="router.push('/')">
        <ArrowLeft :size="20" />
      </button>
      <h1 class="category-header__title">{{ category?.name || slug }}</h1>
      <p v-if="category?.description" class="category-header__desc">{{ category.description }}</p>
    </header>

    <section class="container">
      <div v-if="loading" class="tour-grid">
        <div v-for="i in 4" :key="i" class="tour-skeleton">
          <div class="skeleton" style="aspect-ratio: 4/3"></div>
          <div style="padding: var(--spacing-3)">
            <div class="skeleton" style="height: 12px; width: 60%; margin-bottom: 8px"></div>
            <div class="skeleton" style="height: 16px; width: 90%; margin-bottom: 8px"></div>
            <div class="skeleton" style="height: 12px; width: 40%"></div>
          </div>
        </div>
      </div>

      <div v-else-if="tours.length" class="tour-grid">
        <TourCard v-for="tour in tours" :key="tour.id" :tour="tour" :allow-like="true" />
      </div>

      <div v-else class="empty">
        <p>No hay tours en esta categoría aún.</p>
      </div>
    </section>
  </div>
</template>

<style scoped>
.category-header {
  padding-top: var(--spacing-6);
  padding-bottom: var(--spacing-6);
}

.back-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border-radius: var(--radius-full);
  background: var(--color-surface);
  box-shadow: var(--shadow-sm);
  margin-bottom: var(--spacing-4);
  color: var(--color-text);
}

.category-header__title {
  font-size: var(--font-size-2xl);
  font-weight: var(--font-weight-bold);
  letter-spacing: var(--letter-spacing-tight);
}

.category-header__desc {
  font-size: var(--font-size-sm);
  color: var(--color-text-light);
  margin-top: var(--spacing-2);
}

.tour-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: var(--spacing-4);
}

.tour-skeleton {
  background: var(--color-surface);
  border-radius: var(--radius-xl);
  overflow: hidden;
}

.empty {
  text-align: center;
  padding: var(--spacing-12) var(--spacing-4);
  color: var(--color-text-light);
}
</style>
