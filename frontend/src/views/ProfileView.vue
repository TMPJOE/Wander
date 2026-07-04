<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useFavoritesStore } from '../stores/favorites'
import { useBookingsStore } from '../stores/bookings'
import { useApi } from '../composables/useApi'
import TourCard from '../components/TourCard.vue'
import ReviewCard from '../components/ReviewCard.vue'
import EmptyState from '../components/EmptyState.vue'
import StarRating from '../components/StarRating.vue'
import { LogOut, Settings, Heart, Save } from '@lucide/vue'

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

const profileForm = ref({
  first_name: '',
  last_name: '',
  bio: '',
  phone: '',
})

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
        bookingReviewForm.value[b.id] = bookingReviewForm.value[b.id] || {
          rating: 0,
          title: '',
          comment: '',
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
  try {
    await api.post(`/tours/${booking.tour_id}/reviews`, form)
    // refresh lists
    await fetchAdventures()
    await bookingsStore.fetchMyBookings()
    completedBookings.value = (bookingsStore.bookings || []).filter(
      (b: any) => b.status === 'completed',
    )
    activeBookingForm.value = null
  } catch (e) {
    console.error('Failed to submit review from profile', e)
    alert('Error al enviar reseña. Asegúrate de haber completado el tour y no haber reseñado ya.')
  }
}

function formFor(id: number) {
  if (!bookingReviewForm.value[id]) {
    bookingReviewForm.value[id] = { rating: 0, title: '', comment: '' }
  }
  return bookingReviewForm.value[id]
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
  <div class="page">
    <div class="header">
      <h1 class="title">Mi Perfil</h1>
      <button class="logout-btn" @click="handleLogout">
        <LogOut :size="20" />
      </button>
    </div>

    <!-- Profile Info -->
    <div class="profile-header container" v-if="authStore.user">
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
    <div class="tabs container">
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
        <Heart :size="18" />
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
    <div class="container py-4" v-if="activeTab === 'favorites'">
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
        <button class="btn btn-outline mt-3" @click="router.push('/')">Explorar tours</button>
      </div>
    </div>

    <!-- Tab Content: Adventures -->
    <div class="container py-4" v-if="activeTab === 'adventures'">
      <div v-if="adventuresLoading" class="flex flex-col gap-4">
        <div v-for="i in 3" :key="i" class="skeleton h-28 rounded-lg"></div>
      </div>

      <!-- Completed bookings: allow creating reviews here -->
      <div v-if="!adventuresLoading && completedBookings.length" class="grid gap-3 mb-4">
        <div v-for="booking in completedBookings" :key="booking.id" class="card p-4">
          <div class="flex justify-between items-center">
            <div>
              <div class="text-sm font-semibold">{{ booking.tour_title }}</div>
              <div class="text-xs text-secondary">
                {{ new Date(booking.schedule_start).toLocaleDateString('es-MX') }}
              </div>
            </div>
            <div>
              <button
                v-if="activeBookingForm !== booking.id"
                class="btn btn-outline btn-sm"
                @click="activeBookingForm = booking.id"
              >
                Escribir reseña
              </button>
              <button v-else class="btn btn-ghost btn-sm" @click="activeBookingForm = null">
                Cerrar
              </button>
            </div>
          </div>

          <div v-if="activeBookingForm === booking.id" class="mt-3">
            <div class="form-group mb-2">
              <label class="form-label">Calificación</label>
              <StarRating
                :rating="formFor(booking.id).rating"
                interactive
                @rate="(v) => (formFor(booking.id).rating = v)"
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
              <button class="btn btn-primary btn-sm" @click="submitReviewForBooking(booking)">
                Publicar reseña
              </button>
            </div>
          </div>
        </div>
      </div>

      <div v-if="!adventuresLoading && adventures.length" class="grid gap-4">
        <ReviewCard v-for="review in adventures" :key="review.id" :review="review" />
      </div>

      <EmptyState
        v-else
        :icon="Heart"
        title="Sin aventuras aún"
        message="Aún no has escrito reseñas. Comparte tu experiencia después de un tour completado."
      >
        <button class="btn btn-primary" @click="router.push('/')">Explorar tours</button>
      </EmptyState>
    </div>

    <!-- Tab Content: Settings -->
    <div class="container py-4" v-if="activeTab === 'settings'">
      <div class="card settings-card">
        <div class="flex justify-between items-center mb-4">
          <h3 class="text-lg font-semibold">Información Personal</h3>
          <button v-if="!isEditing" class="btn btn-ghost btn-sm" @click="isEditing = true">
            Editar
          </button>
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
    </div>
  </div>
</template>

<style scoped>
.container {
  padding: 0 var(--content-padding);
}
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-6) var(--spacing-4) var(--spacing-4);
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
</style>
