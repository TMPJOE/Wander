import { useAuthStore } from '../stores/auth';
import { useApi } from './useApi';

export function useAuth() {
  const authStore = useAuthStore();
  const api = useApi();

  async function checkAuth() {
    if (!authStore.token) return false;
    
    try {
      const response = await api.get('/auth/me');
      authStore.user = response.data;
      return true;
    } catch {
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
