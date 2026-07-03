<script setup lang="ts">
import { ref } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import { Eye, EyeOff, MapPin } from '@lucide/vue';

const router = useRouter();
const route = useRoute();
const authStore = useAuthStore();

const email = ref('');
const password = ref('');
const showPassword = ref(false);

async function handleSubmit() {
  const success = await authStore.login(email.value, password.value);
  if (success) {
    const redirect = route.query.redirect as string;
    router.push(redirect || '/');
  }
}
</script>

<template>
  <div class="auth-page">
    <div class="auth-brand">
      <div class="auth-brand__icon">
        <MapPin :size="32" :stroke-width="2" />
      </div>
      <h1 class="auth-brand__name">Wander</h1>
      <p class="auth-brand__tagline">Descubre experiencias locales únicas</p>
    </div>

    <form class="auth-form" @submit.prevent="handleSubmit">
      <h2 class="auth-form__title">Iniciar Sesión</h2>

      <div v-if="authStore.error" class="auth-error">
        {{ authStore.error }}
      </div>

      <div class="form-group">
        <label class="form-label" for="login-email">Correo electrónico</label>
        <input
          id="login-email"
          v-model="email"
          type="email"
          class="form-input"
          placeholder="tu@email.com"
          required
          autocomplete="email"
        />
      </div>

      <div class="form-group">
        <label class="form-label" for="login-password">Contraseña</label>
        <div class="password-field">
          <input
            id="login-password"
            v-model="password"
            :type="showPassword ? 'text' : 'password'"
            class="form-input"
            placeholder="Tu contraseña"
            required
            autocomplete="current-password"
          />
          <button
            type="button"
            class="password-toggle"
            @click="showPassword = !showPassword"
          >
            <EyeOff v-if="showPassword" :size="18" />
            <Eye v-else :size="18" />
          </button>
        </div>
      </div>

      <button
        type="submit"
        class="btn btn-primary btn-block btn-lg"
        :disabled="authStore.loading"
      >
        {{ authStore.loading ? 'Ingresando...' : 'Ingresar' }}
      </button>

      <p class="auth-footer">
        ¿No tienes cuenta?
        <RouterLink to="/register" class="auth-link">Regístrate</RouterLink>
      </p>
    </form>
  </div>
</template>

<style scoped>
.auth-page {
  min-height: 100vh;
  min-height: 100dvh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--spacing-6);
  background: linear-gradient(160deg, var(--color-primary-50) 0%, var(--color-background) 40%, var(--color-secondary-50) 100%);
}

.auth-brand {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: var(--spacing-10);
}

.auth-brand__icon {
  width: 56px;
  height: 56px;
  border-radius: var(--radius-xl);
  background: var(--color-primary);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: var(--spacing-3);
  box-shadow: 0 4px 16px rgba(208, 83, 66, 0.3);
}

.auth-brand__name {
  font-size: var(--font-size-3xl);
  font-weight: var(--font-weight-extrabold);
  letter-spacing: var(--letter-spacing-tight);
  color: var(--color-text);
}

.auth-brand__tagline {
  font-size: var(--font-size-sm);
  color: var(--color-text-light);
  margin-top: var(--spacing-1);
}

.auth-form {
  width: 100%;
  max-width: 380px;
  background: var(--color-surface);
  border-radius: var(--radius-2xl);
  padding: var(--spacing-8);
  box-shadow: var(--shadow-lg);
  display: flex;
  flex-direction: column;
  gap: var(--spacing-5);
}

.auth-form__title {
  font-size: var(--font-size-xl);
  font-weight: var(--font-weight-bold);
  text-align: center;
}

.auth-error {
  background: var(--color-error-bg);
  color: var(--color-error);
  padding: var(--spacing-3) var(--spacing-4);
  border-radius: var(--radius-md);
  font-size: var(--font-size-sm);
  text-align: center;
}

.password-field {
  position: relative;
}

.password-toggle {
  position: absolute;
  right: var(--spacing-3);
  top: 50%;
  transform: translateY(-50%);
  color: var(--color-text-light);
  padding: var(--spacing-1);
}

.auth-footer {
  text-align: center;
  font-size: var(--font-size-sm);
  color: var(--color-text-light);
}

.auth-link {
  color: var(--color-primary);
  font-weight: var(--font-weight-semibold);
}
</style>
