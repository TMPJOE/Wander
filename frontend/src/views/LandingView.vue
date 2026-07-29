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
      <div class="landing-hero__bg">
        <img :src="adventureImg" alt="" />
        <img :src="hikingImg" alt="" />
        <img :src="waterImg" alt="" />
        <img :src="culturalImg" alt="" />
        <img :src="natureImg" alt="" />
        <img :src="nightlifeImg" alt="" />
        <img :src="foodImg" alt="" />
        <img :src="historicalImg" alt="" />
        <img :src="photographyImg" alt="" />
      </div>
      <div class="landing-hero__overlay"></div>

      <div class="container landing-hero__content">
        <img :src="wanderLogo" alt="Wander Logo" class="landing-hero__logo" />
        <h1 class="landing-hero__title">Descubre y ofrece experiencias únicas</h1>
        <p class="landing-hero__subtitle">
          Explora tours locales auténticos o conviértete en guía y comparte tu pasión con viajeros
          de todo el mundo.
        </p>
        <div class="landing-hero__actions">
          <RouterLink to="/login" class="btn btn-primary btn-lg">Entrar</RouterLink>
          <RouterLink to="/register" class="btn btn-secondary btn-lg">Registrarse</RouterLink>
        </div>

        <div class="disclaimer-card">
          <h3 class="disclaimer-card__title">⚠️ Proyecto Demo — No Production</h3>
          <ul class="disclaimer-card__list">
            <li>No se realiza ningún pago real. Stripe está en modo prueba.</li>
            <li>Los tours mostrados son datos de demostración y pueden no estar disponibles.</li>
            <li>Esta aplicación es un proyecto de portafolio con fines educativos.</li>
          </ul>
        </div>
      </div>
    </header>

    <!-- Footer -->
    <footer class="landing-footer">
      <div class="container">
        <div class="landing-footer__content">
          <img :src="wanderLogo" alt="Wander Logo" class="landing-footer__logo" />
          <p class="landing-footer__text">
            &copy; {{ new Date().getFullYear() }} Wander. Todos los derechos reservados.
          </p>
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
  position: relative;
  overflow: hidden;
  min-height: 640px;
  flex: 1;
  display: flex;
  align-items: center;
  padding: var(--spacing-12) 0 var(--spacing-8);
}

.landing-hero__bg {
  position: absolute;
  inset: -20px; /* un poco más grande que el contenedor, así el blur no deja bordes vacíos */
  z-index: 0;
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  grid-template-rows: repeat(3, 1fr);
  gap: 6px;
}

.landing-hero__bg img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
  filter: blur(6px) brightness(0.75) saturate(1.1);
  transform: scale(1.08); /* evita que se vean bordes nítidos por el blur */
}

/* mosaico tipo bento en vez de grid uniforme */
.landing-hero__bg img:nth-child(1) {
  grid-column: span 2;
  grid-row: span 2;
}
.landing-hero__bg img:nth-child(5) {
  grid-row: span 2;
}

.landing-hero__overlay {
  position: absolute;
  inset: 0;
  z-index: 1;
  background: linear-gradient(
    180deg,
    rgba(0, 0, 0, 0.55) 0%,
    rgba(0, 0, 0, 0.35) 45%,
    rgba(0, 0, 0, 0.65) 100%
  );
}

.landing-hero__content {
  position: relative;
  z-index: 2;
  text-align: center;
}

.landing-hero__logo {
  display: block;
  width: 192px;
  height: 192px;
  object-fit: contain;
  margin: 0 auto var(--spacing-4);
  filter: drop-shadow(0 2px 8px rgba(0, 0, 0, 0.4));
}

.landing-hero__title {
  font-size: var(--font-size-3xl);
  font-weight: var(--font-weight-extrabold);
  color: var(--color-text-inverse);
  margin-bottom: var(--spacing-3);
  letter-spacing: var(--letter-spacing-tight);
  text-shadow: 0 2px 12px rgba(0, 0, 0, 0.4);
}

.landing-hero__subtitle {
  font-size: var(--font-size-lg);
  color: var(--color-text-inverse);
  max-width: 500px;
  margin: 0 auto var(--spacing-6);
  line-height: var(--line-height-relaxed);
}

.landing-hero__actions {
  display: flex;
  gap: var(--spacing-3);
  justify-content: center;
  margin-bottom: var(--spacing-8);
}

/* Disclaimer Section */
.landing-disclaimer {
  padding: 0 var(--content-padding) var(--spacing-8);
}

.disclaimer-card {
  background: var(--color-overlay);
  backdrop-filter: blur(10px);
  border-radius: var(--radius-lg);
  padding: var(--spacing-5);
}

.disclaimer-card__title {
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-inverse);
  margin-bottom: var(--spacing-3);
}

.disclaimer-card__list {
  list-style: disc;
  padding-left: var(--spacing-5);
  color: var(--color-text-inverse);
  font-size: var(--font-size-sm);
  line-height: var(--line-height-relaxed);
  text-align: left;
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
  .landing-hero {
    min-height: 480px;
  }

  .landing-hero__bg {
    grid-template-columns: repeat(2, 1fr);
    grid-template-rows: repeat(3, 1fr);
  }

  .landing-hero__bg img:nth-child(1) {
    grid-column: span 1;
    grid-row: span 1;
  }
  .landing-hero__bg img:nth-child(5) {
    grid-row: span 1;
  }
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
