import { defineStore } from 'pinia';
import { ref, computed } from 'vue';

export interface User {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  role: 'traveler' | 'guide' | 'admin';
  avatar_url?: string;
  bio?: string;
  phone?: string;
}

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null);
  const token = ref<string | null>(localStorage.getItem('wander_token'));

  const isAuthenticated = computed(() => !!token.value);
  const isGuide = computed(() => user.value?.role === 'guide' || user.value?.role === 'admin');

  function setAuth(newToken: string, userData: User) {
    token.value = newToken;
    user.value = userData;
    localStorage.setItem('wander_token', newToken);
  }

  function logout() {
    token.value = null;
    user.value = null;
    localStorage.removeItem('wander_token');
  }

  return {
    user,
    token,
    isAuthenticated,
    isGuide,
    setAuth,
    logout,
  };
});
