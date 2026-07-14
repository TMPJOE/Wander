import { ref, computed } from 'vue'

export interface User {
  id: number
  email: string
  username: string
  first_name: string
  last_name: string
  role: 'traveler' | 'guide' | 'admin'
  avatar_url?: string
  bio?: string
  phone?: string
  languages?: string[]
}

const userRef = ref<User | null>(null)
const tokenRef = ref<string | null>(localStorage.getItem('wander_token'))
const loadingRef = ref(false)
const errorRef = ref<string | null>(null)

export function useAuthState() {
  const isAuthenticated = computed(() => !!tokenRef.value)
  const isGuide = computed(() => userRef.value?.role === 'guide' || userRef.value?.role === 'admin')
  const fullName = computed(() => {
    if (!userRef.value) return ''
    return `${userRef.value.first_name} ${userRef.value.last_name}`.trim()
  })

  function setAuth(newToken: string, userData: User) {
    tokenRef.value = newToken
    userRef.value = userData
    localStorage.setItem('wander_token', newToken)
  }

  function logout() {
    tokenRef.value = null
    userRef.value = null
    localStorage.removeItem('wander_token')
  }

  async function login(email: string, password: string) {
    // Lazy load useApi to prevent circular dependency
    const { useApi } = await import('./useApi')
    const api = useApi()
    loadingRef.value = true
    errorRef.value = null
    try {
      const response = await api.post('/auth/login', { email, password })
      setAuth(response.data.token, response.data.user)
      return true
    } catch (e: any) {
      errorRef.value = e.response?.data?.error || 'Error al iniciar sesión'
      return false
    } finally {
      loadingRef.value = false
    }
  }

  async function register(data: {
    email: string
    username: string
    password: string
    first_name: string
    last_name: string
    role?: string
  }) {
    const { useApi } = await import('./useApi')
    const api = useApi()
    loadingRef.value = true
    errorRef.value = null
    try {
      const response = await api.post('/auth/register', data)
      setAuth(response.data.token, response.data.user)
      return true
    } catch (e: any) {
      errorRef.value = e.response?.data?.error || 'Error al registrarse'
      return false
    } finally {
      loadingRef.value = false
    }
  }

  async function fetchMe() {
    const { useApi } = await import('./useApi')
    const api = useApi()
    try {
      const response = await api.get('/users/me')
      userRef.value = response.data
    } catch {
      logout()
    }
  }

  return {
    user: userRef,
    token: tokenRef,
    loading: loadingRef,
    error: errorRef,
    isAuthenticated,
    isGuide,
    fullName,
    setAuth,
    logout,
    login,
    register,
    fetchMe,
  }
}
