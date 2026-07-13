<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useApi } from '../composables/useApi'
import { Search, Map, Star, Languages, Calendar, Users, MapPin } from '@lucide/vue'
import { normalizeTourImages } from '../utils/tourImages'

const route = useRoute()
const router = useRouter()
const api = useApi()

const guideId = route.params.id as string
const guide = ref<any>(null)
const tours = ref<any[]>([])
const loading = ref(true)

const defaultAvatar = 'https://images.unsplash.com/photo-1500648767791-00dcc994a43e?auto=format&fit=crop&w=300&q=80'

onMounted(async () => {
  try {
    // Fetch guide profile
    const guideRes = await api.get(`/users/${guideId}`)
    guide.value = guideRes.data

    // Fetch guide's tours
    const toursRes = await api.get(`/tours?guide_id=${guideId}`)
    tours.value = toursRes.data || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
})

const stats = computed(() => {
  if (!tours.value.length) return { avgRating: 0, reviewCount: 0, activeTours: 0 }
  
  let totalRating = 0
  let totalReviews = 0
  
  tours.value.forEach(t => {
    totalRating += (t.avg_rating || 0) * (t.review_count || 0)
    totalReviews += (t.review_count || 0)
  })
  
  const avg = totalReviews > 0 ? totalRating / totalReviews : 0
  
  return {
    avgRating: avg,
    reviewCount: totalReviews,
    activeTours: tours.value.length
  }
})

function getTourImage(tour: any) {
  const imgs = normalizeTourImages(tour.images)
  return imgs[0] || 'https://images.unsplash.com/photo-1501785888041-af3ef285b470?w=400&h=300&fit=crop'
}
</script>

<template>
  <div class="page bg-surface min-h-screen pb-8">
    <div v-if="loading" class="container py-8 flex flex-col items-center">
      <div class="skeleton avatar-skeleton mb-4"></div>
      <div class="skeleton w-48 h-6 mb-2"></div>
      <div class="skeleton w-32 h-4"></div>
    </div>

    <div v-else-if="guide" class="container guide-profile-container">
      
      <!-- Profile Card -->
      <div class="profile-card">
        <div class="profile-header">
          <Search :size="20" class="text-primary cursor-pointer" @click="router.push('/')" />
          <h2 class="profile-logo-text">LocalGuides</h2>
          <div class="w-5"></div> <!-- Spacer -->
        </div>

        <div class="profile-body">
          <img :src="guide.avatar_url || defaultAvatar" :alt="guide.first_name" class="profile-avatar" />
          <h1 class="profile-name">{{ guide.first_name }} {{ guide.last_name }}</h1>
          <p class="profile-tagline">{{ guide.bio || 'Guía Local Experto' }}</p>

          <!-- Stats Grid -->
          <div class="stats-grid">
            <div class="stat-box">
              <div class="stat-header">
                <Map :size="16" class="text-secondary" />
                <span>Tours Publicados</span>
              </div>
              <div class="stat-value">{{ stats.activeTours }}</div>
            </div>
            
            <div class="stat-box">
              <div class="stat-header">
                <Languages :size="16" class="text-success" />
                <span>Idiomas</span>
              </div>
              <div class="stat-value text-sm mt-1 flex flex-wrap gap-1 justify-center">
                <span v-if="!guide.languages?.length">ES</span>
                <span v-for="lang in guide.languages" :key="lang" class="lang-chip">{{ lang.toUpperCase() }}</span>
              </div>
            </div>
          </div>

          <!-- Full width stat -->
          <div class="stat-box full-width-stat">
            <div class="stat-header">
              <Star :size="16" class="text-warning" fill="currentColor" />
              <span>Calificación Promedio</span>
            </div>
            <div class="stat-value rating-value">
              {{ stats.avgRating.toFixed(1) }} <span class="rating-sub">/ 5.0 ({{ stats.reviewCount }} reseñas)</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Tours Section -->
      <div class="tours-section">
        <h3 class="section-title">Tours de {{ guide.first_name }}</h3>
        
        <div class="tours-list">
          <div v-if="tours.length === 0" class="text-center text-light py-4">
            Este guía aún no tiene tours activos.
          </div>
          
          <div v-for="tour in tours" :key="tour.id" class="tour-row-card">
            <div class="tour-row-img-wrapper">
              <img :src="getTourImage(tour)" alt="Tour" class="tour-row-img" />
            </div>
            <div class="tour-row-content">
              <h4 class="tour-row-title">{{ tour.title }}</h4>
              <div class="tour-row-meta">
                <span><MapPin :size="12" /> {{ tour.location }}</span>
              </div>
              <div class="tour-row-meta">
                <span><Calendar :size="12" /> {{ Math.floor(tour.duration_minutes / 60) }}h {{ tour.duration_minutes % 60 ? (tour.duration_minutes % 60) + 'm' : '' }}</span>
                <span><Users :size="12" /> Máx {{ tour.max_guests }}</span>
              </div>
              <button class="btn btn-outline-primary btn-sm mt-3 w-full" @click="router.push(`/tours/${tour.id}`)">
                Ver Detalles
              </button>
            </div>
          </div>
        </div>
      </div>

    </div>

    <div v-else class="container py-8 text-center text-light">
      Guía no encontrado.
    </div>
  </div>
</template>

<style scoped>
.bg-surface {
  background: #f8fafc;
}

.guide-profile-container {
  max-width: 480px;
  margin: 0 auto;
  padding: var(--spacing-4);
  padding-top: var(--spacing-6);
}

.profile-card {
  background: white;
  border-radius: var(--radius-xl);
  border: 1px solid var(--color-border-light);
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.03);
  overflow: hidden;
  margin-bottom: var(--spacing-6);
}

.profile-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-4) var(--spacing-5);
  border-bottom: 1px solid var(--color-border-light);
}

.profile-logo-text {
  font-size: var(--font-size-base);
  font-weight: 800;
  color: var(--color-primary-dark);
}

.profile-body {
  padding: var(--spacing-6) var(--spacing-5);
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
}

.profile-avatar {
  width: 96px;
  height: 96px;
  border-radius: 50%;
  object-fit: cover;
  border: 3px solid white;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
  margin-bottom: var(--spacing-4);
}

.profile-name {
  font-size: var(--font-size-xl);
  font-weight: 800;
  color: #0c1a2e;
  margin-bottom: 4px;
}

.profile-tagline {
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
  margin-bottom: var(--spacing-6);
  max-width: 90%;
  line-height: 1.4;
}

.stats-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--spacing-3);
  width: 100%;
  margin-bottom: var(--spacing-3);
}

.stat-box {
  background: white;
  border: 1px solid #f0ded6;
  border-radius: var(--radius-lg);
  padding: var(--spacing-4);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.full-width-stat {
  width: 100%;
  align-items: flex-start;
  padding: var(--spacing-4) var(--spacing-5);
}

.stat-header {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 11px;
  font-weight: 700;
  color: var(--color-text-secondary);
  margin-bottom: var(--spacing-2);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.stat-value {
  font-size: 28px;
  font-weight: 800;
  color: #0c1a2e;
}

.rating-value {
  display: flex;
  align-items: baseline;
  gap: 8px;
}

.rating-sub {
  font-size: var(--font-size-sm);
  font-weight: 500;
  color: var(--color-text-secondary);
}

.lang-chip {
  background: var(--color-background);
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 12px;
  border: 1px solid var(--color-border);
}

.text-warning { color: #f59e0b; }
.text-success { color: #10b981; }
.text-secondary { color: #64748b; }
.text-primary { color: #a03e1c; }

.tours-section {
  width: 100%;
}

.section-title {
  font-size: var(--font-size-lg);
  font-weight: 700;
  color: #0c1a2e;
  margin-bottom: var(--spacing-4);
  padding-bottom: var(--spacing-2);
  border-bottom: 2px solid var(--color-border-light);
}

.tours-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-4);
}

.tour-row-card {
  background: white;
  border-radius: var(--radius-lg);
  border: 1px solid var(--color-border-light);
  display: flex;
  overflow: hidden;
  box-shadow: 0 4px 12px rgba(0,0,0,0.02);
}

.tour-row-img-wrapper {
  width: 100px;
  flex-shrink: 0;
}

.tour-row-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.tour-row-content {
  padding: var(--spacing-3) var(--spacing-4);
  flex: 1;
  display: flex;
  flex-direction: column;
}

.tour-row-title {
  font-size: var(--font-size-base);
  font-weight: 700;
  color: #0c1a2e;
  margin-bottom: var(--spacing-1);
  line-height: 1.2;
}

.tour-row-meta {
  display: flex;
  align-items: center;
  gap: var(--spacing-3);
  font-size: 12px;
  color: var(--color-text-secondary);
  margin-bottom: 4px;
}

.tour-row-meta span {
  display: flex;
  align-items: center;
  gap: 4px;
}

.btn-outline-primary {
  background: transparent;
  border: 1px solid var(--color-primary);
  color: var(--color-primary);
  border-radius: 20px;
  padding: 6px 0;
  font-weight: 600;
  transition: all 0.2s;
}

.btn-outline-primary:hover {
  background: var(--color-primary-50);
}

.avatar-skeleton {
  width: 96px;
  height: 96px;
  border-radius: 50%;
}
</style>
