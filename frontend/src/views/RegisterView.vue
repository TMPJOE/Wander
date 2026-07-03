<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import { Eye, EyeOff, MapPin } from '@lucide/vue';

const router = useRouter();
const authStore = useAuthStore();

const form = ref({
  first_name: '',
  last_name: '',
  username: '',
  email: '',
  password: '',
  role: 'traveler' as 'traveler' | 'guide',
});
const showPassword = ref(false);

async function handleSubmit() {
  const success = await authStore.register({
    email: form.value.email,
    username: form.value.username,
    password: form.value.password,
    first_name: form.value.first_name,
    last_name: form.value.last_name,
    role: form.value.role,
  });
  if (success) {
    router.push('/');
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
      <p class="auth-brand__tagline">Únete a la comunidad</p>
    </div>

    <form class="auth-form" @submit.prevent="handleSubmit">
      <h2 class="auth-form__title">Crear Cuenta</h2>

      <div v-if="authStore.error" class="auth-error">
        {{ authStore.error }}
      </div>

      <!-- Role Selector -->
      <div class="role-selector">
        <button
          type="button"
          class="role-option"
          :class="{ 'role-option--active': form.role === 'traveler' }"
          @click="form.role = 'traveler'"
        >
          <span class="role-option__emoji">🧳</span>
          <span class="role-option__label">Viajero</span>
        </button>
        <button
          type="button"
          class="role-option"
          :class="{ 'role-option--active': form.role === 'guide' }"
          @click="form.role = 'guide'"
        >
          <span class="role-option__emoji">🗺️</span>
          <span class="role-option__label">Guía</span>
        </button>
      </div>

      <div class="form-row">
        <div class="form-group">
          <label class="form-label" for="reg-fname">Nombre</label>
          <input
            id="reg-fname"
            v-model="form.first_name"
            type="text"
            class="form-input"
            placeholder="María"
            required
          />
        </div>
        <div class="form-group">
          <label class="form-label" for="reg-lname">Apellido</label>
          <input
            id="reg-lname"
            v-model="form.last_name"
            type="text"
            class="form-input"
            placeholder="García"
            required
          />
        </div>
      </div>

      <div class="form-group">
        <label class="form-label" for="reg-username">Nombre de usuario</label>
        <input
          id="reg-username"
          v-model="form.username"
          type="text"
          class="form-input"
          placeholder="maria_garcia"
          required
        />
      </div>

      <div class="form-group">
        <label class="form-label" for="reg-email">Correo electrónico</label>
        <input
          id="reg-email"
          v-model="form.email"
          type="email"
          class="form-input"
          placeholder="tu@email.com"
          required
          autocomplete="email"
        />
      </div>

      <div class="form-group">
        <label class="form-label" for="reg-password">Contraseña</label>
        <div class="password-field">
          <input
            id="reg-password"
            v-model="form.password"
            :type="showPassword ? 'text' : 'password'"
            class="form-input"
            placeholder="Mínimo 8 caracteres"
            required
            minlength="8"
            autocomplete="new-password"
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
        {{ authStore.loading ? 'Creando cuenta...' : 'Crear Cuenta' }}
      </button>

      <p class="auth-footer">
        ¿Ya tienes cuenta?
        <RouterLink to="/login" class="auth-link">Inicia sesión</RouterLink>
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
  margin-bottom: var(--spacing-8);
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
}

.auth-brand__tagline {
  font-size: var(--font-size-sm);
  color: var(--color-text-light);
  margin-top: var(--spacing-1);
}

.auth-form {
  width: 100%;
  max-width: 420px;
  background: var(--color-surface);
  border-radius: var(--radius-2xl);
  padding: var(--spacing-8);
  box-shadow: var(--shadow-lg);
  display: flex;
  flex-direction: column;
  gap: var(--spacing-4);
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

.role-selector {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--spacing-3);
}

.role-option {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--spacing-1);
  padding: var(--spacing-4) var(--spacing-3);
  border: 2px solid var(--color-border);
  border-radius: var(--radius-lg);
  cursor: pointer;
  transition: all var(--transition-fast);
  background: transparent;
}

.role-option:hover {
  border-color: var(--color-primary-light);
}

.role-option--active {
  border-color: var(--color-primary);
  background: var(--color-primary-50);
}

.role-option__emoji {
  font-size: 1.5rem;
}

.role-option__label {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--spacing-3);
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
