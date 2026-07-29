<script setup lang="ts">
import { useRouter } from 'vue-router'
import wanderLogo from '../assets/wander-logo.svg'
import adventureImg from '../assets/categories/adventure.png'
import culturalImg from '../assets/categories/cultural.png'
import foodImg from '../assets/categories/food.png'
import hikingImg from '../assets/categories/hiking.png'
import historicalImg from '../assets/categories/historical.png'
import natureImg from '../assets/categories/nature.png'
import nightlifeImg from '../assets/categories/nightlife.png'
import photographyImg from '../assets/categories/Photography.png'
import waterImg from '../assets/categories/water.png'

const router = useRouter()

const categories = [
  { name: 'Aventura', slug: 'adventure', image: adventureImg },
  { name: 'Cultural', slug: 'cultural', image: culturalImg },
  { name: 'Gastronomía', slug: 'food', image: foodImg },
  { name: 'Senderismo', slug: 'hiking', image: hikingImg },
  { name: 'Histórico', slug: 'historical', image: historicalImg },
  { name: 'Naturaleza', slug: 'nature', image: natureImg },
  { name: 'Vida nocturna', slug: 'nightlife', image: nightlifeImg },
  { name: 'Fotografía', slug: 'Photography', image: photographyImg },
  { name: 'Agua', slug: 'water', image: waterImg },
]

function goToExplore(categorySlug?: string) {
  if (categorySlug) {
    router.push(`/explore?category=${categorySlug}`)
  } else {
    router.push('/explore')
  }
}
</script>

<template>
  <div class="landing-page">
    <!-- Hero Section -->
    <header class="landing-hero">
      <div class="container">
        <img :src="wanderLogo" alt="Wander Logo" class="landing-hero__logo" />
        <h1 class="landing-hero__title">Descubre y ofrece experiencias únicas</h1>
        <p class="landing-hero__subtitle">
          Explora tours locales auténticos o conviértete en guía y comparte tu pasión con viajeros de todo el mundo.
        </p>
        <div class="landing-hero__actions">
          <RouterLink to="/login" class="btn btn-primary btn-lg">Entrar</RouterLink>
          <RouterLink to="/register" class="btn btn-outline btn-lg">Registrarse</RouterLink>
        </div>
      </div>
    </header>

    <!-- Category Gallery -->
    <section class="landing-categories container">
      <h2 class="landing-section-title">Explora por categoría</h2>
      <div class="categories-grid">
        <button
          v-for="cat in categories"
          :key="cat.slug"
          class="category-card"
          @click="goToExplore(cat.slug)"
        >
          <img :src="cat.image" :alt="cat.name" class="category-card__image" />
          <span class="category-card__name">{{ cat.name }}</span>
        </button>
      </div>
    </section>

    <!-- Demo Disclaimer -->
    <section class="landing-disclaimer container">
      <div class="disclaimer-card">
        <h3 class="disclaimer-card__title">⚠️ Proyecto Demo — No Production</h3>
        <ul class="disclaimer-card__list">
          <li>No se realiza ningún pago real. Stripe está en modo prueba.</li>
          <li>Los tours mostrados son datos de demostración y pueden no estar disponibles.</li>
          <li>Esta aplicación es un proyecto de portafolio con fines educativos.</li>
        </ul>
      </div>
    </section>

    <!-- Footer -->
    <footer class="landing-footer">
      <div class="container">
        <div class="landing-footer__content">
          <img :src="wanderLogo" alt="Wander Logo" class="landing-footer__logo" />
          <p class="landing-footer__text">&copy; {{ new Date().getFullYear() }} Wander. Todos los derechos reservados.</p>
          <p class="landing-footer__disclaimer">Proyecto demo — no production.</p>
        </div>
      </div>
    </footer>
  </div>
</template>

<style scoped>
.landing-page {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.container {
  max-width: var(--max-width);
  margin: 0 auto;
  padding: 0 var(--content-padding);
}

/* Hero Section */
.landing-hero {
  text-align: center;
  padding: var(--spacing-12) 0 var(--spacing-8);
  background: linear-gradient(
    160deg,
    var(--color-primary-50) 0%,
    var(--color-background) 40%,
    var(--color-secondary-50) 100%
  );
}

.landing-hero__logo {
  width: 96px;
  height: 96px;
  object-fit: contain;
  margin-bottom: var(--spacing-4);
}

.landing-hero__title {
  font-size: var(--font-size-3xl);
  font-weight: var(--font-weight-extrabold);
  color: var(--color-text);
  margin-bottom: var(--spacing-3);
  letter-spacing: var(--letter-spacing-tight);
}

.landing-hero__subtitle {
  font-size: var(--font-size-lg);
  color: var(--color-text-light);
  max-width: 500px;
  margin: 0 auto var(--spacing-6);
  line-height: var(--line-height-relaxed);
}

.landing-hero__actions {
  display: flex;
  gap: var(--spacing-3);
  justify-content: center;
}

/* Categories Section */
.landing-categories {
  padding: var(--spacing-10) var(--content-padding);
}

.landing-section-title {
  font-size: var(--font-size-2xl);
  font-weight: var(--font-weight-bold);
  text-align: center;
  margin-bottom: var(--spacing-6);
  color: var(--color-text);
}

.categories-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
  gap: var(--spacing-4);
}

.category-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--spacing-2);
  padding: var(--spacing-3);
  background: var(--color-surface);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-xl);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.category-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-md);
  border-color: var(--color-primary);
}

.category-card__image {
  width: 80px;
  height: 80px;
  object-fit: contain;
  border-radius: var(--radius-md);
}

.category-card__name {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  color: var(--color-text);
  text-align: center;
}

/* Disclaimer Section */
.landing-disclaimer {
  padding: 0 var(--content-padding) var(--spacing-8);
}

.disclaimer-card {
  background: var(--color-warning-bg);
  border: 1px solid var(--color-warning);
  border-radius: var(--radius-lg);
  padding: var(--spacing-5);
}

.disclaimer-card__title {
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  color: var(--color-warning-dark);
  margin-bottom: var(--spacing-3);
}

.disclaimer-card__list {
  list-style: disc;
  padding-left: var(--spacing-5);
  color: var(--color-warning-dark);
  font-size: var(--font-size-sm);
  line-height: var(--line-height-relaxed);
}

.disclaimer-card__list li {
  margin-bottom: var(--spacing-1);
}

/* Footer */
.landing-footer {
  margin-top: auto;
  padding: var(--spacing-6) 0;
  background: var(--color-surface);
  border-top: 1px solid var(--color-border-light);
}

.landing-footer__content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--spacing-2);
  text-align: center;
}

.landing-footer__logo {
  width: 40px;
  height: 40px;
  object-fit: contain;
}

.landing-footer__text {
  font-size: var(--font-size-sm);
  color: var(--color-text-light);
}

.landing-footer__disclaimer {
  font-size: var(--font-size-xs);
  color: var(--color-text-muted);
  font-style: italic;
}

/* Responsive */
@media (max-width: 640px) {
  .landing-hero__title {
    font-size: var(--font-size-2xl);
  }

  .landing-hero__subtitle {
    font-size: var(--font-size-base);
  }

  .landing-hero__actions {
    flex-direction: column;
    align-items: center;
  }

  .categories-grid {
    grid-template-columns: repeat(auto-fill, minmax(110px, 1fr));
  }

  .category-card__image {
    width: 64px;
    height: 64px;
  }
}
</style>
