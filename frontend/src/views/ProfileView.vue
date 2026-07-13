<script setup lang="ts">
import { ref, onMounted, watch, computed, nextTick } from 'vue'
import { loadStripe } from '@stripe/stripe-js'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useFavoritesStore } from '../stores/favorites'
import { useBookingsStore } from '../stores/bookings'
import { useApi } from '../composables/useApi'
import TourCard from '../components/TourCard.vue'
import ReviewCard from '../components/ReviewCard.vue'
import EmptyState from '../components/EmptyState.vue'
import StarRating from '../components/StarRating.vue'
import { normalizeTourImages } from '../utils/tourImages'
import {
  LogOut,
  Settings,
  Heart,
  Compass,
  Save,
  Calendar,
  Users,
  CreditCard,
  MapPin,
  Edit2,
  Star,
  User,
  Bell,
  HelpCircle,
  ChevronRight,
} from '@lucide/vue'

const authStore = useAuthStore()
const favoritesStore = useFavoritesStore()
const api = useApi()
const router = useRouter()

const activeTab = ref('favorites')
const adventures = ref<any[]>([])
const adventuresLoading = ref(false)
const bookingsStore = useBookingsStore()
const completedBookings = ref<any[]>([])
const bookingReviewForm = ref<{ [k: number]: { rating: number; title: string; comment: string } }>(
  {},
)
const activeBookingForm = ref<number | null>(null)
const isEditing = ref(false)
const saving = ref(false)
const showPaymentMethods = ref(false)
const isAddingCard = ref(false)
const cardElement = ref<HTMLElement | null>(null)
let stripe: any = null
let elements: any = null
let card: any = null

const stripePromise = loadStripe(import.meta.env.VITE_STRIPE_PUBLISHABLE_KEY)

async function showAddCard() {
  isAddingCard.value = true
  await nextTick()
  if (!stripe) {
    stripe = await stripePromise
    elements = stripe.elements()
  }
  if (!card && cardElement.value) {
    card = elements.create('card', {
      style: {
        base: {
          iconColor: '#913b1f',
          color: '#0c1a2e',
          fontWeight: '500',
          fontFamily: 'Inter, Roboto, sans-serif',
          fontSize: '16px',
          fontSmoothing: 'antialiased',
          '::placeholder': {
            color: '#a0aec0',
          },
        },
        invalid: {
          iconColor: '#e53e3e',
          color: '#e53e3e',
        },
      },
    })
    card.mount(cardElement.value)
  }
}

function cancelAddCard() {
  isAddingCard.value = false
}

async function saveCard() {
  if (!stripe || !card) return
  const { paymentMethod, error } = await stripe.createPaymentMethod({
    type: 'card',
    card: card,
  })
  
  if (error) {
    alert(error.message)
  } else {
    // Simulated behavior for demonstration
    alert('Tarjeta añadida con éxito (Simulado)')
    isAddingCard.value = false
    if (card) {
      card.clear()
    }
  }
}

const profileForm = ref({
  first_name: '',
  last_name: '',
  bio: '',
  phone: '',
})

// Helper to find review for a booking
function getReviewForBooking(bookingId: number) {
  return adventures.value.find((r) => r.booking_id === bookingId || r.tour_id === bookingId)
}

// Check if booking has a review
function hasReview(bookingId: number) {
  return !!getReviewForBooking(bookingId)
}

onMounted(async () => {
  if (authStore.user) {
    initForm()
  }
  try {
    await favoritesStore.fetchFavorites()
  } catch {
    // Ignore — favorites tab will show empty state
  }
})

function initForm() {
  if (!authStore.user) return
  profileForm.value = {
    first_name: authStore.user.first_name || '',
    last_name: authStore.user.last_name || '',
    bio: authStore.user.bio || '',
    phone: authStore.user.phone || '',
  }
}

async function fetchAdventures() {
  adventuresLoading.value = true
  try {
    const res = await api.get('/reviews/me')
    adventures.value = res.data || []
  } catch (e) {
    console.error('Failed to fetch adventures', e)
  } finally {
    adventuresLoading.value = false
  }
}

watch(activeTab, async (tab) => {
  if (tab === 'adventures') {
    if (adventures.value.length === 0) await fetchAdventures()
    // ensure we have bookings to allow creating reviews from profile
    try {
      await bookingsStore.fetchMyBookings()
      completedBookings.value = (bookingsStore.bookings || []).filter(
        (b: any) => b.status === 'completed',
      )
      // initialize per-booking form state
      completedBookings.value.forEach((b: any) => {
        const existingReview = getReviewForBooking(b.id)
        bookingReviewForm.value[b.id] = {
          rating: existingReview?.rating || 0,
          title: existingReview?.title || '',
          comment: existingReview?.comment || '',
        }
      })
    } catch (e) {
      console.error('Failed to fetch bookings for adventures tab', e)
    }
  }
})

async function submitReviewForBooking(booking: any) {
  const form = bookingReviewForm.value[booking.id]
  if (!form || !form.rating || !form.comment.trim()) {
    alert('Calificación y comentario son obligatorios.')
    return
  }

  const existingReview = getReviewForBooking(booking.id)

  try {
    if (existingReview) {
      // Update existing review
      await api.put(`/reviews/${existingReview.id}`, form)
    } else {
      // Create new review
      await api.post(`/tours/${booking.tour_id}/reviews`, {
        ...form,
        booking_id: booking.id,
      })
    }

    // refresh lists
    await fetchAdventures()
    await bookingsStore.fetchMyBookings()
    completedBookings.value = (bookingsStore.bookings || []).filter(
      (b: any) => b.status === 'completed',
    )
    activeBookingForm.value = null
  } catch (e) {
    console.error('Failed to submit review from profile', e)
    alert('Error al enviar reseña.')
  }
}

function openReviewForm(booking: any) {
  activeBookingForm.value = booking.id
  // Pre-fill form with existing review if any
  const existingReview = getReviewForBooking(booking.id)
  if (existingReview) {
    bookingReviewForm.value[booking.id] = {
      rating: existingReview.rating || 0,
      title: existingReview.title || '',
      comment: existingReview.comment || '',
    }
  }
}

function formFor(id: number) {
  if (!bookingReviewForm.value[id]) {
    bookingReviewForm.value[id] = { rating: 0, title: '', comment: '' }
  }
  return bookingReviewForm.value[id]
}

function onRateFor(booking: any, value: number) {
  formFor(booking.id).rating = value
}

const PLACEHOLDER_IMG =
  'https://images.unsplash.com/photo-1501785888041-af3ef285b470?w=400&h=300&fit=crop'

function bookingImage(booking: any): string {
  if (!booking.tour_image) return PLACEHOLDER_IMG
  const imgs = normalizeTourImages(booking.tour_image)
  return imgs[0] || PLACEHOLDER_IMG
}

function bookingDate(booking: any): string {
  return new Date(booking.schedule_start).toLocaleDateString('es-MX', {
    weekday: 'short',
    day: 'numeric',
    month: 'short',
    year: 'numeric',
  })
}

function bookingTime(booking: any): string {
  return new Date(booking.schedule_start).toLocaleTimeString('es-MX', {
    hour: '2-digit',
    minute: '2-digit',
  })
}

async function saveProfile() {
  saving.value = true
  try {
    const res = await api.put('/users/me', profileForm.value)
    authStore.user = res.data
    isEditing.value = false
  } catch (e) {
    console.error(e)
    alert('Error al guardar el perfil')
  } finally {
    saving.value = false
  }
}

function handleLogout() {
  authStore.logout()
  router.push('/login')
}

const handleSettingsClick = () => {
  activeTab.value = 'settings'
  initForm()
}
</script>

<template>
  <div class="profile-page bg-surface">
    <div class="header px-content">
      <h1 class="title">Mi Perfil</h1>
      <button class="logout-btn" @click="handleLogout">
        <LogOut :size="20" />
      </button>
    </div>

    <!-- Profile Info -->
    <div class="profile-header px-content" v-if="authStore.user">
      <div class="profile-avatar-wrap">
        <img
          v-if="authStore.user.avatar_url"
          :src="authStore.user.avatar_url"
          class="profile-avatar"
        />
        <div v-else class="profile-avatar profile-avatar--placeholder">
          {{ (authStore.user.first_name || 'U').charAt(0) }}
        </div>
      </div>
      <div class="profile-meta">
        <h2 class="profile-name">{{ authStore.fullName }}</h2>
        <span class="profile-role">
          {{ authStore.user.role === 'guide' ? 'Guía Local' : 'Viajero' }}
        </span>
      </div>
    </div>

    <!-- Tabs -->
    <div class="tabs px-content">
      <button
        class="tab-btn"
        :class="{ 'tab-btn--active': activeTab === 'favorites' }"
        @click="activeTab = 'favorites'"
      >
        <Heart :size="18" />
        Favoritos
      </button>
      <button
        class="tab-btn"
        :class="{ 'tab-btn--active': activeTab === 'adventures' }"
        @click="activeTab = 'adventures'"
      >
        <Compass :size="18" />
        Mis Aventuras
      </button>
      <button
        class="tab-btn"
        :class="{ 'tab-btn--active': activeTab === 'settings' }"
        @click="handleSettingsClick"
      >
        <Settings :size="18" />
        Ajustes
      </button>
    </div>

    <!-- Tab Content: Favorites -->
    <div class="px-content py-4" v-if="activeTab === 'favorites'">
      <div v-if="favoritesStore.loading" class="grid-2 gap-4">
        <div v-for="i in 2" :key="i" class="skeleton aspect-video rounded-xl"></div>
      </div>
      <div v-else-if="favoritesStore.favorites.length" class="grid-2 gap-4">
        <TourCard
          v-for="fav in favoritesStore.favorites"
          :key="fav.id"
          :tour="{ ...fav, is_favorited: true }"
          :allow-like="false"
        />
      </div>
      <div v-else class="empty-favorites">
        <Heart :size="48" class="empty-icon" />
        <p class="empty-text">Aún no tienes tours favoritos</p>
        <button class="btn btn-primary mt-3" @click="router.push('/')">Explorar tours</button>
      </div>
    </div>

    <!-- Tab Content: Adventures -->
    <div class="px-content py-4" v-if="activeTab === 'adventures'">
      <div v-if="adventuresLoading" class="flex flex-col gap-4">
        <div v-for="i in 3" :key="i" class="skeleton h-28 rounded-lg"></div>
      </div>

      <!-- Completed bookings list -->
      <div v-if="!adventuresLoading && completedBookings.length" class="adventures-list mb-4">
        <div v-for="booking in completedBookings" :key="booking.id" class="adv-card card">
          <div class="adv-card__header">
            <span class="badge badge-primary">Completada</span>
            <span class="adv-card__id">Reserva #{{ booking.id }}</span>
          </div>

          <div class="adv-card__body">
            <img :src="bookingImage(booking)" :alt="booking.tour_title" class="adv-card__image" />
            <div class="adv-card__info">
              <h3 class="adv-card__title">{{ booking.tour_title }}</h3>
              <p class="adv-card__location">
                <MapPin :size="14" />
                {{ booking.tour_location || booking.guide_name }}
              </p>
              <div class="adv-card__details">
                <div class="detail-item">
                  <Calendar :size="14" class="detail-icon" />
                  <span>{{ bookingDate(booking) }} • {{ bookingTime(booking) }}</span>
                </div>
                <div class="detail-item">
                  <Users :size="14" class="detail-icon" />
                  <span
                    >{{ booking.guest_count }} persona{{ booking.guest_count > 1 ? 's' : '' }}</span
                  >
                </div>
                <div class="detail-item">
                  <CreditCard :size="14" class="detail-icon" />
                  <span class="font-semibold"
                    >${{ booking.total_price.toLocaleString('es-MX') }}</span
                  >
                </div>
              </div>
            </div>
          </div>

          <!-- Display existing review in card -->
          <div v-if="hasReview(booking.id)" class="adv-card__review-display">
            <div class="review-header">
              <div class="review-stars">
                <Star
                  v-for="i in 5"
                  :key="i"
                  :size="16"
                  :class="{ 'star-filled': i <= getReviewForBooking(booking.id).rating }"
                />
              </div>
              <span class="review-title">{{ getReviewForBooking(booking.id).title }}</span>
            </div>
            <p class="review-comment">{{ getReviewForBooking(booking.id).comment }}</p>
          </div>

          <div class="adv-card__footer">
            <button
              v-if="activeBookingForm !== booking.id"
              class="btn btn-outline btn-md"
              @click="openReviewForm(booking)"
            >
              <Edit2 :size="16" />
              {{ hasReview(booking.id) ? 'Editar reseña' : 'Escribir reseña' }}
            </button>
            <button v-else class="btn btn-ghost btn-md" @click="activeBookingForm = null">
              Cerrar
            </button>
          </div>

          <div v-if="activeBookingForm === booking.id" class="adv-card__review-form">
            <div class="form-group mb-2">
              <label class="adv-card__title">Calificación</label>
              <StarRating
                :rating="formFor(booking.id).rating"
                interactive
                @rate="(v: number) => onRateFor(booking, v)"
              />
            </div>
            <div class="form-group mb-2">
              <label class="form-label">Título</label>
              <input v-model="formFor(booking.id).title" type="text" class="form-input" />
            </div>
            <div class="form-group mb-2">
              <label class="form-label">Comentario</label>
              <textarea
                v-model="formFor(booking.id).comment"
                class="form-input form-textarea"
                rows="3"
              ></textarea>
            </div>
            <div class="flex justify-end">
              <button 
                class="btn btn_xl btn-review-submit" 
                style="background-color: var(--color-star); border-color: var(--color-star); color: white;"
                @click="submitReviewForBooking(booking)"
              >
                <Save :size="14" />
                {{ hasReview(booking.id) ? 'Actualizar' : 'Publicar' }} reseña
              </button>
            </div>
          </div>
        </div>
      </div>

      <EmptyState
        v-if="!adventuresLoading && !completedBookings.length"
        :icon="Compass"
        title="Sin aventuras aún"
        message="Completa un tour para comenzar a compartir tu experiencia."
      >
        <button class="btn btn-primary" @click="router.push('/')">Explorar tours</button>
      </EmptyState>
    </div>

    <!-- Tab Content: Settings -->
    <div class="px-content py-4" v-if="activeTab === 'settings'">
      <div class="card settings-card">
        <div class="flex justify-between items-center mb-4">
          <h3 class="text-lg font-semibold">Información Personal</h3>
        </div>

        <form v-if="isEditing" @submit.prevent="saveProfile" class="flex-col gap-4">
          <div class="grid-2 gap-3">
            <div class="form-group">
              <label class="form-label">Nombre</label>
              <input v-model="profileForm.first_name" type="text" class="form-input" required />
            </div>
            <div class="form-group">
              <label class="form-label">Apellido</label>
              <input v-model="profileForm.last_name" type="text" class="form-input" required />
            </div>
          </div>
          <div class="form-group">
            <label class="form-label">Teléfono</label>
            <input v-model="profileForm.phone" type="tel" class="form-input" />
          </div>
          <div class="form-group">
            <label class="form-label">Biografía</label>
            <textarea v-model="profileForm.bio" class="form-input form-textarea"></textarea>
          </div>
          <div class="flex gap-2 justify-end mt-2">
            <button type="button" class="btn btn-ghost" @click="isEditing = false">Cancelar</button>
            <button type="submit" class="btn btn-primary" :disabled="saving">
              <Save :size="16" />
              {{ saving ? 'Guardando...' : 'Guardar' }}
            </button>
          </div>
        </form>

        <div v-else class="info-grid">
          <div class="info-item">
            <span class="info-label">Email</span>
            <span class="info-value">{{ authStore.user?.email }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">Teléfono</span>
            <span class="info-value">{{ authStore.user?.phone || 'No especificado' }}</span>
          </div>
          <div class="info-item col-span-2">
            <span class="info-label">Biografía</span>
            <span class="info-value">{{ authStore.user?.bio || 'Sin biografía' }}</span>
          </div>
        </div>
      </div>

      <!-- Account Settings Menu -->
      <div class="mt-6">
        <h3 class="text-lg font-semibold mb-3">Account Settings</h3>
        <div class="card p-0 overflow-hidden settings-menu-list">
          <button class="settings-menu-item" @click="isEditing = true">
            <div class="settings-menu-icon">
              <User :size="18" />
            </div>
            <span>Editar Perfil</span>
            <ChevronRight :size="18" class="ml-auto text-muted" />
          </button>
          
          <button class="settings-menu-item" @click="showPaymentMethods = !showPaymentMethods">
            <div class="settings-menu-icon">
              <CreditCard :size="18" />
            </div>
            <span>Métodos de Pago</span>
            <ChevronRight :size="18" class="ml-auto text-muted" :class="{ 'rotate-90': showPaymentMethods }" style="transition: transform 0.2s;" />
          </button>
          
          <div v-if="showPaymentMethods" class="payment-methods-panel">
            <h4 class="text-sm font-semibold text-secondary mb-3">Tarjetas Guardadas</h4>
            <div class="saved-card">
              <div class="flex items-center gap-3">
                <CreditCard :size="20" class="text-primary" />
                <div>
                  <div class="font-semibold text-sm">•••• •••• •••• 4242</div>
                  <div class="text-xs text-muted">Expira 12/28</div>
                </div>
              </div>
              <button class="btn btn-ghost btn-sm text-error">Eliminar</button>
            </div>
            <button v-if="!isAddingCard" class="btn btn-outline w-full text-sm mt-3 flex justify-center items-center gap-2" @click="showAddCard">
              <span>+ Agregar nueva tarjeta</span>
            </button>
            
            <div v-if="isAddingCard" class="mt-4 p-4 border border-border-light rounded-lg bg-white">
              <h5 class="text-sm font-semibold mb-3">Detalles de la tarjeta</h5>
              <div ref="cardElement" class="p-3 border border-border-light rounded-md mb-4 bg-surface"></div>
              <div class="flex gap-2 justify-end">
                <button class="btn btn-ghost btn-sm" @click="cancelAddCard">Cancelar</button>
                <button class="btn btn-primary btn-sm" @click="saveCard">Guardar Tarjeta</button>
              </div>
            </div>
          </div>
          
          <button class="settings-menu-item">
            <div class="settings-menu-icon">
              <Bell :size="18" />
            </div>
            <span>Notificaciones</span>
            <ChevronRight :size="18" class="ml-auto text-muted" />
          </button>
          
          <button class="settings-menu-item">
            <div class="settings-menu-icon">
              <HelpCircle :size="18" />
            </div>
            <span>Ayuda</span>
            <ChevronRight :size="18" class="ml-auto text-muted" />
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>

.profile-page {
  flex: 1;
  width: 100%;
  min-height: 100vh;
  min-height: 100dvh;
  padding-bottom: calc(var(--bottom-nav-height) + var(--spacing-4));
}

.px-content {
  padding-left: var(--content-padding);
  padding-right: var(--content-padding);
}
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: var(--spacing-6);
  padding-bottom: var(--spacing-4);
}

.title {
  font-size: var(--font-size-2xl);
  font-weight: var(--font-weight-bold);
  letter-spacing: var(--letter-spacing-tight);
}

.logout-btn {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-error);
  border-radius: var(--radius-full);
  background: var(--color-error-bg);
}

.profile-header {
  display: flex;
  align-items: center;
  gap: var(--spacing-4);
  margin-bottom: var(--spacing-6);
}

.profile-avatar {
  width: 80px;
  height: 80px;
  border-radius: var(--radius-full);
  object-fit: cover;
  border: 3px solid var(--color-surface);
  box-shadow: var(--shadow-md);
}

.profile-avatar--placeholder {
  background: var(--color-primary);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--font-size-2xl);
  font-weight: var(--font-weight-bold);
}

.profile-name {
  font-size: var(--font-size-xl);
  font-weight: var(--font-weight-bold);
  line-height: var(--line-height-tight);
  margin-bottom: 2px;
}

.profile-role {
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
  font-weight: var(--font-weight-medium);
}

.tabs {
  display: flex;
  gap: var(--spacing-2);
  border-bottom: 1px solid var(--color-border-light);
  padding-bottom: 1px;
}

.tab-btn {
  display: flex;
  align-items: center;
  gap: var(--spacing-2);
  padding: var(--spacing-3) var(--spacing-4);
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-light);
  border-bottom: 2px solid transparent;
  transition: all var(--transition-fast);
}

.tab-btn--active {
  color: var(--color-primary);
  border-bottom-color: var(--color-primary);
}

.py-4 {
  padding-top: var(--spacing-4);
  padding-bottom: var(--spacing-4);
}
.grid-2 {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
}
@media (max-width: 380px) {
  .grid-2 {
    grid-template-columns: 1fr;
  }
}
.gap-3 {
  gap: var(--spacing-3);
}
.gap-4 {
  gap: var(--spacing-4);
}
.mt-3 {
  margin-top: var(--spacing-3);
}
.aspect-video {
  aspect-ratio: 16/9;
}
.rounded-xl {
  border-radius: var(--radius-xl);
}

.empty-favorites {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  padding: var(--spacing-12) 0;
}

.empty-icon {
  color: var(--color-border);
  margin-bottom: var(--spacing-3);
}

.empty-text {
  color: var(--color-text-secondary);
}

.settings-card {
  padding: var(--spacing-5);
}

.text-lg {
  font-size: var(--font-size-lg);
}
.font-semibold {
  font-weight: var(--font-weight-semibold);
}
.mb-4 {
  margin-bottom: var(--spacing-4);
}
.mt-2 {
  margin-top: var(--spacing-2);
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: var(--spacing-4);
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.col-span-2 {
  grid-column: span 2;
}

.info-label {
  font-size: var(--font-size-xs);
  color: var(--color-text-light);
  text-transform: uppercase;
  letter-spacing: var(--letter-spacing-wide);
}

.info-value {
  font-size: var(--font-size-sm);
  color: var(--color-text);
  font-weight: var(--font-weight-medium);
}

/* Mis Aventuras card (mirrors BookingCard design) */
.adventures-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-4);
}

.btn-review-submit {
  padding-top: var(--spacing-4) !important;
  padding-right: var(--spacing-4) !important;
}

.adv-card {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-3);
}

.adv-card__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid var(--color-border-light);
  padding: var(--spacing-3);
  margin-bottom: var(--spacing-1);
}

.adv-card__id {
  font-size: var(--font-size-xs);
  color: var(--color-text-light);
  font-weight: var(--font-weight-medium);
  margin-right: var(--spacing-2);
}

.adv-card__body {
  padding-top: var(--spacing-3);
  padding-bottom: var(--spacing-3);
  padding-left: var(--spacing-8);
  padding-right: var(--spacing-8);
  display: flex;
  gap: var(--spacing-4);
}

.adv-card__image {
  width: 90px;
  height: 90px;
  border-radius: var(--radius-md);
  object-fit: cover;
  flex-shrink: 0;
}

.adv-card__info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: var(--spacing-1);
}

.adv-card__title {
  font-size: var(--font-size-base);
  font-weight: var(--font-weight-semibold);
  line-height: var(--line-height-tight);
}

.adv-card__location {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: var(--font-size-xs);
  color: var(--color-text-light);
  margin-bottom: var(--spacing-2);
}

.adv-card__details {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.detail-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
}

.detail-icon {
  color: var(--color-text-light);
}

.font-semibold {
  font-weight: var(--font-weight-semibold);
  color: var(--color-text);
}

.adv-card__footer {
  margin-top: var(--spacing-2);
  padding-right: var(--spacing-3);
  padding-bottom: var(--spacing-3);
  display: flex;
  justify-content: flex-end;
}

/* Review display in card */
.adv-card__review-display {
  margin: 0 var(--spacing-3);
  padding: var(--spacing-4);
  background: var(--color-warning-bg);
  border-radius: var(--radius-md);
}

.review-header {
  display: flex;
  align-items: center;
  gap: var(--spacing-3);
  margin-bottom: var(--spacing-2);
}

.review-stars {
  display: flex;
  gap: 2px;
}

.review-stars svg {
  color: var(--color-border);
}

.review-stars .star-filled {
  color: var(--color-warning);
  fill: var(--color-warning);
}

.review-title {
  font-weight: var(--font-weight-semibold);
  font-size: var(--font-size-sm);
  color: var(--color-text);
}

.review-comment {
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
  line-height: var(--line-height-relaxed);
  margin: 0;
}

.adv-card__review-form {
  margin: 0 var(--spacing-3) var(--spacing-3);
  padding: var(--spacing-4);
  background: var(--color-warning-bg);
  border-radius: var(--radius-md);
  display: flex;
  flex-direction: column;
}

.mt-6 {
  margin-top: var(--spacing-6);
}
.mb-3 {
  margin-bottom: var(--spacing-3);
}
.p-0 {
  padding: 0;
}
.overflow-hidden {
  overflow: hidden;
}
.ml-auto {
  margin-left: auto;
}

.settings-menu-list {
  display: flex;
  flex-direction: column;
  border: 1px solid #f0ded6;
}

.settings-menu-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-3);
  padding: var(--spacing-4);
  background: white;
  border: none;
  border-bottom: 1px solid #f0ded6;
  text-align: left;
  font-size: var(--font-size-base);
  font-weight: 500;
  color: #0c1a2e;
  cursor: pointer;
  transition: background-color 0.2s;
}

.settings-menu-item:last-child {
  border-bottom: none;
}

.settings-menu-item:hover {
  background-color: var(--color-surface);
}

.settings-menu-icon {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background-color: #f0f4fd;
  color: var(--color-primary);
  display: flex;
  align-items: center;
  justify-content: center;
}

.payment-methods-panel {
  padding: var(--spacing-4);
  background-color: #fafbfc;
  border-bottom: 1px solid #f0ded6;
}

.saved-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-3);
  background: white;
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-md);
}

.rotate-90 {
  transform: rotate(90deg);
}

.w-full {
  width: 100%;
}
</style>
