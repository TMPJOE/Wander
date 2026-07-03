import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { useApi } from '../composables/useApi';

export interface User {
  id: number;
  email: string;
  username: string;
  first_name: string;
  last_name: string;
  role: 'traveler' | 'guide' | 'admin';
  avatar_url?: string;
  bio?: string;
  phone?: string;
  languages?: string[];
}

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null);
  const token = ref<string | null>(localStorage.getItem('wander_token'));
  const loading = ref(false);
  const error = ref<string | null>(null);

  const isAuthenticated = computed(() => !!token.value);
  const isGuide = computed(() => user.value?.role === 'guide' || user.value?.role === 'admin');
  const fullName = computed(() => {
    if (!user.value) return '';
    return `${user.value.first_name} ${user.value.last_name}`.trim();
  });

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

  async function login(email: string, password: string) {
    const api = useApi();
    loading.value = true;
    error.value = null;
    try {
      const response = await api.post('/auth/login', { email, password });
      setAuth(response.data.token, response.data.user);
      return true;
    } catch (e: any) {
      error.value = e.response?.data?.error || 'Error al iniciar sesión';
      return false;
    } finally {
      loading.value = false;
    }
  }

  async function register(data: {
    email: string;
    username: string;
    password: string;
    first_name: string;
    last_name: string;
    role?: string;
  }) {
    const api = useApi();
    loading.value = true;
    error.value = null;
    try {
      const response = await api.post('/auth/register', data);
      setAuth(response.data.token, response.data.user);
      return true;
    } catch (e: any) {
      error.value = e.response?.data?.error || 'Error al registrarse';
      return false;
    } finally {
      loading.value = false;
    }
  }

  async function fetchMe() {
    const api = useApi();
    try {
      const response = await api.get('/users/me');
      user.value = response.data;
    } catch {
      logout();
    }
  }

  return {
    user,
    token,
    loading,
    error,
    isAuthenticated,
    isGuide,
    fullName,
    setAuth,
    logout,
    login,
    register,
    fetchMe,
  };
});
