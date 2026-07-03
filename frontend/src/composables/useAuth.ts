import { useAuthStore } from '../stores/auth';
import { useApi } from './useApi';
import { useRouter } from 'vue-router';

export function useAuth() {
  const authStore = useAuthStore();
  const api = useApi();
  const router = useRouter();

  async function checkAuth() {
    if (!authStore.token) return false;
    
    try {
      const response = await api.get('/auth/me');
      authStore.user = response.data;
      return true;
    } catch (error) {
      authStore.logout();
      return false;
    }
  }

  function requireAuth(to: any, from: any, next: any) {
    if (!authStore.isAuthenticated) {
      next({ name: 'login', query: { redirect: to.fullPath } });
    } else {
      next();
    }
  }

  function requireGuide(to: any, from: any, next: any) {
    if (!authStore.isAuthenticated) {
      next({ name: 'login', query: { redirect: to.fullPath } });
    } else if (!authStore.isGuide) {
      next({ name: 'explore' }); // Redirect non-guides to home
    } else {
      next();
    }
  }

  return {
    checkAuth,
    requireAuth,
    requireGuide,
  };
}
