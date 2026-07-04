import { useAuthStore } from '../stores/auth'
import { useApi } from './useApi'

export function useAuth() {
  const authStore = useAuthStore()
  const api = useApi()

  async function checkAuth() {
    if (!authStore.token) return false

    try {
      const response = await api.get('/auth/me')
      authStore.user = response.data
      return true
    } catch {
      authStore.logout()
      return false
    }
  }

  function requireAuth(to: any) {
    if (!authStore.isAuthenticated) {
      return { name: 'login', query: { redirect: to.fullPath } }
    }
    return true
  }

  function requireGuide(to: any) {
    if (!authStore.isAuthenticated) {
      return { name: 'login', query: { redirect: to.fullPath } }
    }
    if (!authStore.isGuide) {
      return { name: 'explore' } // Redirect non-guides to home
    }
    return true
  }

  return {
    checkAuth,
    requireAuth,
    requireGuide,
  }
}
