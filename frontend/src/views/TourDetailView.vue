<script setup lang="ts">
import { onMounted, computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useToursStore } from '../stores/tours'
import { useAuthStore } from '../stores/auth'
import { useApi } from '../composables/useApi'
import {
  ArrowLeft,
  MapPin,
  Clock,
  Users,
  Star,
  Share2,
  Heart,
  ChevronRight,
  CheckCircle2,
  Languages,
  Footprints,
  Landmark,
  Utensils
} from '@lucide/vue'
import ImageGallery from '../components/ImageGallery.vue'
// StarRating removed: creation is not allowed on tour page
import ReviewCard from '../components/ReviewCard.vue'
import GuideCard from '../components/GuideCard.vue'
import { normalizeTourImages } from '../utils/tourImages'

const route = useRoute()
const router = useRouter()
const toursStore = useToursStore()
const authStore = useAuthStore()
const api = useApi()

const reviews = ref<any[]>([])
const reviewsLoading = ref(false)
const schedules = ref<any[]>([])
// review creation removed from tour detail

const tourId = computed(() => route.params.id as string)

const tour = computed(() => toursStore.currentTour)

const selectedScheduleId = ref<number | null>(null)

function toggleSchedule(id: number) {
  selectedScheduleId.value = selectedScheduleId.value === id ? null : id
}

function getItinerary(startTimeStr: string) {
  const start = new Date(startTimeStr)
  const formatTime = (d: Date) => d.toLocaleTimeString('es-MX', { hour: 'numeric', minute: '2-digit', hour12: true })
  
  const t1 = new Date(start)
  const t2 = new Date(start.getTime() + 60 * 60000)
  const t3 = new Date(start.getTime() + 150 * 60000)
  const t4 = new Date(start.getTime() + 180 * 60000)
  
  return [
    { time: formatTime(t1), title: 'Encuentro en Plaza Nueva', desc: 'Introducción a la historia de la ciudad y comienzo del recorrido.', icon: MapPin },
    { time: formatTime(t2), title: 'Exploración del Barrio de Santa Cruz', desc: 'Caminata por el antiguo barrio judío, conociendo sus leyendas y arquitectura singular.', icon: Footprints },
    { time: formatTime(t3), title: 'Visita Exterior Catedral y Giralda', desc: 'Contexto histórico de los monumentos más emblemáticos.', icon: Landmark },
    { time: formatTime(t4), title: 'Degustación de Tapas', desc: 'Finalizamos con una selección de tapas locales y bebida en una taberna histórica.', icon: Utensils }
  ]
}

const meetingClue = computed(() => {
  if (!tour.value) return ''
  const clues = [
    'El guía llevará un paraguas naranja para facilitar su ubicación.',
    'Busca al guía con una gorra azul y una credencial de Wander.',
    'Nuestro guía te estará esperando con una mochila amarilla visible.',
    'Identificarás al guía por su chaleco reflectante o chaqueta roja.',
    'El guía tendrá un pequeño cartel con tu nombre o el nombre del tour.'
  ]
  // Use tour ID to reliably pick the same clue for the same tour
  const index = tour.value.id % clues.length
  return clues[index]
})

const images = computed<string[]>(() => {
  if (!tour.value) return []
  const imgs = normalizeTourImages(tour.value.images)
  return imgs.length
    ? imgs
    : ['https://images.unsplash.com/photo-1501785888041-af3ef285b470?w=800&h=500&fit=crop']
})

const whatIncluded = computed<string[]>(() => {
  if (!tour.value) return []
  const wi =
    typeof tour.value.what_included === 'string'
      ? JSON.parse(tour.value.what_included || '[]')
      : tour.value.what_included || []
  return wi
})

const formattedDuration = computed(() => {
  if (!tour.value) return ''
  const h = Math.floor(tour.value.duration_minutes / 60)
  const m = tour.value.duration_minutes % 60
  if (h === 0) return `${m} min`
  return m ? `${h}h ${m}m` : `${h} hora${h > 1 ? 's' : ''}`
})

const difficultyLabels: Record<string, { label: string; class: string }> = {
  easy: { label: 'Fácil', class: 'badge-success' },
  moderate: { label: 'Moderado', class: 'badge-warning' },
  challenging: { label: 'Desafiante', class: 'badge-primary' },
  extreme: { label: 'Extremo', class: 'badge-error' },
}

onMounted(async () => {
  await toursStore.fetchTourById(tourId.value)
  fetchReviews()
  fetchSchedules()
})

async function fetchReviews() {
  reviewsLoading.value = true
  try {
    const res = await api.get(`/tours/${tourId.value}/reviews`)
    reviews.value = res.data || []
  } catch {
    /* empty */
  } finally {
    reviewsLoading.value = false
  }
}

// review creation removed from tour detail

async function fetchSchedules() {
  try {
    const res = await api.get(`/tours/${tourId.value}/schedules`)
    schedules.value = (res.data || []).filter((s: any) => new Date(s.start_time) > new Date())
  } catch {
    /* empty */
  }
}

function goBook() {
  if (!authStore.isAuthenticated) {
    router.push({ name: 'login', query: { redirect: `/tours/${tourId.value}/book` } })
    return
  }
  router.push(`/tours/${tourId.value}/book`)
}

function messageGuide() {
  if (!authStore.isAuthenticated || !tour.value) return
  router.push(`/messages/${tour.value.guide_id}`)
}
</script>

<template>
  <div class="page tour-detail" v-if="tour">
    <div class="layered-card">
      <!-- Gallery -->
      <div class="tour-detail__gallery">
        <ImageGallery :images="images" />
        <button class="tour-detail__back" @click="router.back()">
          <ArrowLeft :size="20" />
        </button>
      </div>

      <div class="container tour-detail__content">
        <!-- Header -->
        <div class="tour-detail__header">
          <span v-if="tour.category_name" class="badge badge-secondary">{{
            tour.category_name
          }}</span>
          <h1 class="tour-detail__title">{{ tour.title }}</h1>

          <div class="tour-detail__meta">
            <span class="tour-detail__meta-item">
              <MapPin :size="14" :stroke-width="2" />
              {{ tour.location }}
            </span>
            <span v-if="tour.avg_rating > 0" class="tour-detail__meta-item">
              <Star :size="14" :stroke-width="0" fill="var(--color-star)" />
              {{ tour.avg_rating.toFixed(1) }} ({{ tour.review_count }} reseñas)
            </span>
          </div>
        </div>

        <!-- Quick Stats -->
        <div class="tour-detail__stats">
          <div class="stat-item">
            <Clock :size="18" :stroke-width="1.8" />
            <div>
              <span class="stat-item__value">{{ formattedDuration }}</span>
              <span class="stat-item__label">Duración</span>
            </div>
          </div>
          <div class="stat-item">
            <Users :size="18" :stroke-width="1.8" />
            <div>
              <span class="stat-item__value">{{ tour.max_guests }}</span>
              <span class="stat-item__label">Máx. personas</span>
            </div>
          </div>
          <div class="stat-item">
            <Languages :size="18" :stroke-width="1.8" />
            <div>
              <span class="stat-item__value">{{ (tour.languages || []).join(', ') || '—' }}</span>
              <span class="stat-item__label">Idioma</span>
            </div>
          </div>
        </div>

        <!-- Difficulty -->
        <div v-if="tour.difficulty" class="tour-detail__section">
          <span
            class="badge"
            :class="difficultyLabels[tour.difficulty]?.class || 'badge-secondary'"
          >
            {{ difficultyLabels[tour.difficulty]?.label || tour.difficulty }}
          </span>
        </div>

        <!-- Description -->
        <div class="tour-detail__section">
          <h2 class="tour-detail__section-title">Descripción</h2>
          <p class="tour-detail__description">{{ tour.description }}</p>
        </div>

        <!-- What's Included -->
        <div v-if="whatIncluded.length" class="tour-detail__section">
          <h2 class="tour-detail__section-title">¿Qué incluye?</h2>
          <ul class="included-list">
            <li v-for="(item, i) in whatIncluded" :key="i" class="included-item">
              <CheckCircle2 :size="16" :stroke-width="2" class="included-icon" />
              {{ item }}
            </li>
          </ul>
        </div>



        <!-- Guide -->
        <div class="tour-detail__section">
          <h2 class="tour-detail__section-title">Tu guía</h2>
          <GuideCard
            :guide="{
              id: tour.guide_id,
              first_name: tour.guide_name?.split(' ')[0] || 'Guía',
              last_name: tour.guide_name?.split(' ').slice(1).join(' ') || '',
              avatar_url: tour.guide_avatar,
              languages: tour.languages,
            }"
            @message="messageGuide"
          />
        </div>

        <!-- Available Dates -->
        <div v-if="schedules.length" class="tour-detail__section">
          <h2 class="tour-detail__section-title">Fechas disponibles</h2>
          <div class="schedule-preview">
            <div v-for="s in schedules.slice(0, 3)" :key="s.id" class="schedule-item-wrap">
              <button class="schedule-chip" :class="{ 'schedule-chip--active': selectedScheduleId === s.id }" @click="toggleSchedule(s.id)">
                <span class="schedule-chip__date">
                  {{
                    new Date(s.start_time).toLocaleDateString('es-MX', {
                      day: 'numeric',
                      month: 'short',
                    })
                  }}
                </span>
                <span class="schedule-chip__time">
                  {{
                    new Date(s.start_time).toLocaleTimeString('es-MX', {
                      hour: '2-digit',
                      minute: '2-digit',
                    })
                  }}
                </span>
                <span class="schedule-chip__spots">{{ s.available_spots }} lugares</span>
              </button>
              
              <div v-if="selectedScheduleId === s.id" class="itinerary-card">
                <h3 class="itinerary-card__title">Itinerario detallado</h3>
                <div class="timeline">
                  <div v-for="(item, idx) in getItinerary(s.start_time)" :key="idx" class="timeline-item">
                    <div class="timeline-marker-icon">
                      <component :is="item.icon" :size="16" color="white" />
                    </div>
                    <div class="timeline-content">
                      <h4 class="timeline-heading">
                        <span class="timeline-time">{{ item.time }} - </span>
                        <span class="timeline-title">{{ item.title }}</span>
                      </h4>
                      <p class="timeline-desc">{{ item.desc }}</p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            <p
              v-if="schedules.length > 3"
              class="text-muted"
              style="font-size: var(--font-size-xs); margin-top: var(--spacing-2)"
            >
              +{{ schedules.length - 3 }} fechas más
            </p>
          </div>
        </div>

        <!-- Meeting Point -->
        <div v-if="tour.meeting_point" class="tour-detail__section">
          <div class="itinerary-card meeting-point-card">
            <h3 class="itinerary-card__title meeting-point-card__title">
              <MapPin :size="20" :stroke-width="2" class="meeting-point__icon" />
              Punto de encuentro
            </h3>
            <p class="meeting-point__desc">{{ tour.meeting_point }}. {{ meetingClue }}</p>
            <div class="meeting-point__map">
              <iframe 
                width="100%" 
                height="100%" 
                style="border:0;" 
                loading="lazy" 
                allowfullscreen 
                :src="`https://maps.google.com/maps?q=${encodeURIComponent(tour.meeting_point + ', ' + tour.location)}&t=&z=15&ie=UTF8&iwloc=&output=embed`">
              </iframe>
            </div>
          </div>
        </div>

        <!-- Reviews -->
        <div class="tour-detail__section">
          <div class="section-header" style="margin-bottom: var(--spacing-3)">
            <h2 class="tour-detail__section-title" style="margin-bottom: 0">
              Reseñas
              <span v-if="reviews.length" class="tour-detail__review-count"
                >({{ reviews.length }})</span
              >
            </h2>
          </div>
          <div v-if="reviews.length" class="reviews-list">
            <ReviewCard v-for="review in reviews.slice(0, 5)" :key="review.id" :review="review" />
          </div>
          <p v-else class="text-muted" style="font-size: var(--font-size-sm)">
            Aún no hay reseñas para este tour.
          </p>
        </div>
      </div>
    </div>

    <!-- Sticky Book Bar -->
    <div class="book-bar">
      <div class="book-bar__price">
        <span class="book-bar__amount">${{ tour.price_per_person.toLocaleString('es-MX') }}</span>
        <span class="book-bar__unit">/ persona</span>
      </div>
      <button class="btn btn-secondary-light btn-lg" @click="goBook">
        Reservar ahora
        <ChevronRight :size="18" />
      </button>
    </div>
  </div>

  <!-- Loading State -->
  <div v-else class="page container" style="padding-top: var(--spacing-8)">
    <div class="skeleton" style="aspect-ratio: 16/10; margin-bottom: var(--spacing-4)"></div>
    <div class="skeleton" style="height: 24px; width: 30%; margin-bottom: var(--spacing-3)"></div>
    <div class="skeleton" style="height: 32px; width: 80%; margin-bottom: var(--spacing-3)"></div>
    <div class="skeleton" style="height: 16px; width: 50%"></div>
  </div>
</template>

<style scoped>
.tour-detail {
  display: flex;
  flex-direction: column;
  width: 100%;
  max-width: var(--max-width);
  margin: 0 auto;
}

.tour-detail__gallery {
  position: relative;
}

.tour-detail__back {
  position: absolute;
  top: var(--spacing-4);
  left: var(--spacing-4);
  width: 40px;
  height: 40px;
  border-radius: var(--radius-full);
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: var(--shadow-md);
  z-index: 2;
  transition: all var(--transition-fast);
}

.tour-detail__back:hover {
  background: white;
  transform: scale(1.05);
}

.tour-detail__content {
  z-index: 2;
  position: relative;
  margin-top: calc(-2% - var(--bottom-nav-height));
  padding-left: var(--content-padding);
  padding-right: var(--content-padding);
  padding-top: var(--spacing-5);
  padding-bottom: calc(80px + var(--bottom-nav-height) + var(--spacing-4));
  background: var(--color-surface);
  border-top-left-radius: var(--radius-lg);
  border-top-right-radius: var(--radius-lg);
}

.tour-detail__header {
  margin-bottom: var(--spacing-5);
}

.tour-detail__title {
  font-size: var(--font-size-2xl);
  font-weight: var(--font-weight-bold);
  letter-spacing: var(--letter-spacing-tight);
  line-height: var(--line-height-tight);
  margin: var(--spacing-2) 0 var(--spacing-3);
}

.tour-detail__meta {
  display: flex;
  align-items: center;
  gap: var(--spacing-4);
  flex-wrap: wrap;
}

.tour-detail__meta-item {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
}

.tour-detail__stats {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--spacing-3);
  padding: var(--spacing-4);
  background: var(--color-surface);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-lg);
  margin-bottom: var(--spacing-6);
}

.stat-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-2);
  color: var(--color-text-secondary);
}

.stat-item__value {
  display: block;
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text);
}

.stat-item__label {
  display: block;
  font-size: var(--font-size-xs);
  color: var(--color-text-light);
}

.tour-detail__section {
  margin-bottom: var(--spacing-6);
}

.tour-detail__section-title {
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  margin-bottom: var(--spacing-3);
}

.tour-detail__review-count {
  font-weight: var(--font-weight-normal);
  color: var(--color-text-light);
}

.tour-detail__description {
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
  line-height: var(--line-height-relaxed);
  white-space: pre-line;
}

.included-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-3);
}

.included-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-2);
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
}

.included-icon {
  color: var(--color-success);
  flex-shrink: 0;
}

.meeting-point-card {
  margin-top: var(--spacing-4);
  background: var(--color-surface);
  border: 1.5px solid var(--color-primary-100);
}

.meeting-point-card__title {
  display: flex;
  align-items: center;
  gap: var(--spacing-2);
  margin-bottom: var(--spacing-3);
}

.meeting-point__icon {
  color: var(--color-primary);
  flex-shrink: 0;
}

.meeting-point__desc {
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
  line-height: var(--line-height-relaxed);
  margin-bottom: var(--spacing-4);
}

.meeting-point__map {
  width: 100%;
  height: 160px;
  border-radius: var(--radius-md);
  overflow: hidden;
  border: 1.5px solid var(--color-border-light);
}

.map-placeholder {
  width: 100%;
  height: 100%;
  background-color: #e5e9ea;
  background-image: 
    linear-gradient(white 2px, transparent 2px),
    linear-gradient(90deg, white 2px, transparent 2px),
    linear-gradient(rgba(255,255,255,.3) 1px, transparent 1px),
    linear-gradient(90deg, rgba(255,255,255,.3) 1px, transparent 1px);
  background-size: 100px 100px, 100px 100px, 20px 20px, 20px 20px;
  background-position: -2px -2px, -2px -2px, -1px -1px, -1px -1px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.map-pin {
  filter: drop-shadow(0 4px 6px rgba(160, 62, 28, 0.3));
}

.schedule-preview {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-2);
}

.schedule-item-wrap {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-2);
}

.schedule-chip {
  display: flex;
  align-items: center;
  gap: var(--spacing-3);
  padding: var(--spacing-3) var(--spacing-4);
  background: var(--color-surface);
  border: 1.5px solid var(--color-border-light);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all var(--transition-fast);
  width: 100%;
  text-align: left;
}

.schedule-chip:hover, .schedule-chip--active {
  border-color: var(--color-primary-300);
  background: var(--color-primary-50);
}

.schedule-chip__date {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text);
}

.schedule-chip__time {
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
}

.schedule-chip__spots {
  margin-left: auto;
  font-size: var(--font-size-xs);
  color: var(--color-success);
  font-weight: var(--font-weight-medium);
}

.itinerary-card {
  background: var(--color-primary-50);
  border-radius: var(--radius-lg);
  padding: var(--spacing-4) var(--spacing-5);
  margin-top: var(--spacing-2);
  margin-bottom: var(--spacing-3);
  box-shadow: 0 4px 12px rgba(0,0,0,0.02);
}

.itinerary-card__title {
  font-size: var(--font-size-base);
  font-weight: var(--font-weight-bold);
  color: var(--color-primary-dark);
  margin-bottom: var(--spacing-4);
}

.timeline {
  display: flex;
  flex-direction: column;
  position: relative;
}

.timeline::before {
  content: '';
  position: absolute;
  top: 16px;
  bottom: 16px;
  left: 15px;
  width: 2px;
  background: #a9dbbc; /* Lighter green line */
}

.timeline-item {
  position: relative;
  padding-left: 48px;
  padding-bottom: var(--spacing-6);
}

.timeline-item:last-child {
  padding-bottom: 0;
}

.timeline-marker-icon {
  position: absolute;
  left: 0;
  top: -2px;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: #0b6b47; /* Dark green */
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1;
}

.timeline-heading {
  font-size: var(--font-size-sm);
  margin-bottom: 4px;
  line-height: 1.4;
}

.timeline-time {
  color: #1e293b;
  font-weight: 800;
}

.timeline-title {
  color: #0b6b47;
  font-weight: 700;
}

.timeline-desc {
  font-size: var(--font-size-sm);
  color: #5c626a;
  line-height: var(--line-height-relaxed);
}

.reviews-list {
  background: var(--color-surface);
  border-radius: var(--radius-lg);
  overflow: hidden;
}

/* Sticky Book Bar */
.book-bar {
  position: fixed;
  bottom: var(--bottom-nav-height);
  left: 0;
  right: 0;
  max-width: var(--max-width);
  margin: 0 auto;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--spacing-4) var(--spacing-5);
  background: var(--color-surface);
  border-top: 1px solid var(--color-border-light);
  z-index: var(--z-sticky);
}

@media (min-width: 1024px) {
  .book-bar {
    right: var(--nav-width);
  }
}

.book-bar__price {
  display: flex;
  align-items: baseline;
  gap: var(--spacing-1);
}

.book-bar__amount {
  font-size: var(--font-size-2xl);
  font-weight: var(--font-weight-bold);
  color: var(--color-secondary);
}

.book-bar__unit {
  font-size: var(--font-size-sm);
  color: var(--color-text-light);
}
</style>
