<script setup lang="ts">
import { onMounted, ref, watch, onBeforeUnmount } from 'vue'
import { useFavoritesStore } from '../stores/favorites'
import { Search, SlidersHorizontal, MapPin } from '@lucide/vue'
import { useToursStore } from '../stores/tours'
import { useCategoriesStore } from '../stores/categories'
import { applyCategoryTheme, clearCategoryTheme } from '../utils/categoryColors'
import TourCard from '../components/TourCard.vue'
import CategoryPill from '../components/CategoryPill.vue'
import FilterDrawer from '../components/FilterDrawer.vue'
import type { FilterValues } from '../components/FilterDrawer.vue'

const toursStore = useToursStore()
const categoriesStore = useCategoriesStore()
const favoritesStore = useFavoritesStore()

const searchQuery = ref('')
const activeCategory = ref('')
const showFilters = ref(false)
const currentFilters = ref<FilterValues>({
  difficulty: '',
  min_price: '',
  max_price: '',
  location: '',
})

onMounted(async () => {
  await Promise.all([categoriesStore.fetchCategories(), toursStore.fetchTours()])
  // Ensure favorites are available so cards render correct liked state
  try {
    await favoritesStore.fetchFavorites()
  } catch {
    /* ignore */
  }
})

function selectCategory(slug: string) {
  activeCategory.value = activeCategory.value === slug ? '' : slug
  applyTheme()
  fetchWithFilters()
}

function applyTheme() {
  if (activeCategory.value) {
    applyCategoryTheme(activeCategory.value, document.documentElement)
  } else {
    clearCategoryTheme(document.documentElement)
  }
}

onBeforeUnmount(() => {
  clearCategoryTheme(document.documentElement)
})

function handleSearch() {
  fetchWithFilters()
}

function applyFilters(filters: FilterValues) {
  currentFilters.value = filters
  fetchWithFilters()
}

function fetchWithFilters() {
  const params: Record<string, string> = {}
  if (searchQuery.value) params.q = searchQuery.value
  if (activeCategory.value) params.category = activeCategory.value
  if (currentFilters.value.difficulty) params.difficulty = currentFilters.value.difficulty
  if (currentFilters.value.min_price) params.min_price = currentFilters.value.min_price
  if (currentFilters.value.max_price) params.max_price = currentFilters.value.max_price
  if (currentFilters.value.location) params.location = currentFilters.value.location
  toursStore.fetchTours(params)
}
</script>

<template>
  <div class="page explore" ref="pageEl">
    <!-- Hero Section -->
    <header class="explore__hero">
      <div class="container">
        <div class="explore__greeting">
          <MapPin :size="20" class="explore__pin" />
          <span>Explora tu próxima aventura</span>
        </div>
        <h1 class="explore__headline">
          Descubre experiencias
          <span class="explore__accent">únicas</span>
        </h1>
      </div>
    </header>

    <!-- Search Bar -->
    <div class="container">
      <div class="explore__search">
        <div class="search-bar">
          <Search :size="18" class="search-bar__icon" />
          <input
            v-model="searchQuery"
            type="text"
            class="search-bar__input"
            placeholder="Buscar tours, destinos..."
            @keydown.enter="handleSearch"
          />
        </div>
        <button class="filter-btn" @click="showFilters = true">
          <SlidersHorizontal :size="18" />
        </button>
      </div>
    </div>

    <!-- Categories -->
    <section class="explore__categories">
      <div class="categories-scroll hide-scrollbar">
        <CategoryPill
          name="Todos"
          slug=""
          icon="default"
          :active="activeCategory === ''"
          @select="selectCategory"
        />
        <CategoryPill
          v-for="cat in categoriesStore.categories"
          :key="cat.id"
          :name="cat.name"
          :slug="cat.slug"
          :icon="cat.icon"
          :active="activeCategory === cat.slug"
          @select="selectCategory"
        />
      </div>
    </section>

    <!-- Tours Grid -->
    <section class="container">
      <div class="section-header">
        <h2 class="section-title">
          {{ activeCategory ? 'Resultados' : 'Tours populares' }}
        </h2>
        <span
          v-if="toursStore.tours.length"
          class="text-muted"
          style="font-size: var(--font-size-sm)"
        >
          {{ toursStore.tours.length }} tour{{ toursStore.tours.length !== 1 ? 's' : '' }}
        </span>
      </div>

      <!-- Loading State -->
      <div v-if="toursStore.loading" class="tour-grid">
        <div v-for="i in 4" :key="i" class="tour-skeleton">
          <div class="skeleton" style="aspect-ratio: 4/3"></div>
          <div style="padding: var(--spacing-3)">
            <div class="skeleton" style="height: 12px; width: 60%; margin-bottom: 8px"></div>
            <div class="skeleton" style="height: 16px; width: 90%; margin-bottom: 8px"></div>
            <div class="skeleton" style="height: 12px; width: 40%"></div>
          </div>
        </div>
      </div>

      <!-- Tour Cards -->
      <div v-else-if="toursStore.tours.length" class="tour-grid">
        <TourCard v-for="tour in toursStore.tours" :key="tour.id" :tour="tour" :allow-like="true" />
      </div>

      <!-- Empty State -->
      <div v-else class="explore__empty">
        <p class="explore__empty-text">No se encontraron tours 😕</p>
        <p class="explore__empty-hint">Intenta con otros filtros o categorías</p>
      </div>
    </section>

    <!-- Filter Drawer -->
    <FilterDrawer :open="showFilters" @close="showFilters = false" @apply="applyFilters" />
  </div>
</template>

<style scoped>
.container {
  padding: 0 var(--content-padding);
}

.explore__hero {
  padding: var(--spacing-8) 0 var(--spacing-4);
  background: linear-gradient(
    308deg,
    var(--color-primary-50) 0%,
    var(--color-background) 18%,
    var(--color-primary-100) 26%,
    var(--color-primary-100) 32%,
    var(--color-background) 70%,
    var(--color-background) 100%
  );
  border-bottom: 1px solid var(--color-primary-100);
}

.explore__greeting {
  display: flex;
  align-items: center;
  gap: var(--spacing-2);
  font-size: var(--font-size-sm);
  color: var(--color-text-light);
  margin-bottom: var(--spacing-2);
}

.explore__pin {
  color: var(--color-primary);
}

.explore__headline {
  font-size: var(--font-size-3xl);
  font-weight: var(--font-weight-extrabold);
  letter-spacing: var(--letter-spacing-tight);
  line-height: var(--line-height-tight);
}

.explore__accent {
  color: var(--color-primary);
}

.explore__search {
  display: flex;
  gap: var(--spacing-3);
  margin: var(--spacing-5) 0;
}

.search-bar {
  flex: 1;
  display: flex;
  align-items: center;
  gap: var(--spacing-2);
  background: var(--color-primary-50);
  border: 1.5px solid var(--color-primary-100);
  border-radius: var(--radius-lg);
  padding: 0 var(--spacing-4);
  transition:
    border-color var(--transition-fast),
    box-shadow var(--transition-fast);
}

.search-bar:focus-within {
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px var(--color-primary-100);
}

.search-bar__icon {
  color: var(--color-primary-dark);
  flex-shrink: 0;
}

.search-bar__input {
  flex: 1;
  border: none;
  background: none;
  padding: var(--spacing-3) 0;
  font-size: var(--font-size-sm);
  outline: none;
}

.search-bar__input::placeholder {
  color: var(--color-text-light);
}

.filter-btn {
  width: 48px;
  height: 48px;
  border-radius: var(--radius-lg);
  background: var(--color-surface);
  border: 1.5px solid var(--color-primary-100);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-primary-dark);
  transition: all var(--transition-fast);
  flex-shrink: 0;
}

.filter-btn:hover {
  border-color: var(--color-primary);
  color: var(--color-primary);
  background: var(--color-primary-50);
}

.explore__categories {
  margin-bottom: var(--spacing-6);
}

.categories-scroll {
  display: flex;
  gap: var(--spacing-2);
  overflow-x: auto;
  padding: 0 var(--content-padding);
  scroll-padding: var(--content-padding);
}

.tour-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: var(--spacing-4);
}

@media (max-width: 380px) {
  .tour-grid {
    grid-template-columns: 1fr;
  }
}

.tour-skeleton {
  background: var(--color-surface);
  border-radius: var(--radius-xl);
  overflow: hidden;
}

.explore__empty {
  text-align: center;
  padding: var(--spacing-12) var(--spacing-4);
}

.explore__empty-text {
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  margin-bottom: var(--spacing-2);
}

.explore__empty-hint {
  font-size: var(--font-size-sm);
  color: var(--color-text-light);
}
</style>
