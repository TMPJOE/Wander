import { useAuthState } from './useAuthState'
import { useApi } from './useApi'

export function useAuth() {
  const authState = useAuthState()
  const api = useApi()

  async function checkAuth() {
    if (!authState.token.value) return false

    try {
      const response = await api.get('/auth/me')
      authState.user.value = response.data
      return true
    } catch {
      authState.logout()
      return false
    }
  }

  function requireAuth(to: any) {
    if (!authState.isAuthenticated.value) {
      return { name: 'login', query: { redirect: to.fullPath } }
    }
    return true
  }

  function requireGuide(to: any) {
    if (!authState.isAuthenticated.value) {
      return { name: 'login', query: { redirect: to.fullPath } }
    }
    if (!authState.isGuide.value) {
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
