import axios from 'axios'
import { useAuthState } from './useAuthState'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || '/api/v1',
})

api.interceptors.request.use(
  (config) => {
    const authState = useAuthState()
    if (authState.token.value) {
      config.headers.Authorization = `Bearer ${authState.token.value}`
    }

    if (config.data instanceof FormData) {
      delete config.headers['Content-Type']
    }

    return config
  },
  (error) => {
    return Promise.reject(error)
  },
)

api.interceptors.response.use(
  (response) => {
    if (response.data && response.data.success === true && response.data.data !== undefined) {
      response.data = response.data.data
    }
    // Go returns null for empty slices — normalize to [] for frontend safety
    if (response.data === null) {
      response.data = []
    }
    return response
  },
  (error) => {
    if (error.response && error.response.status === 401) {
      // Don't auto-logout if the failing request was the initial user fetch
      const url = error.config?.url || ''
      if (!url.endsWith('/users/me')) {
        const authState = useAuthState()
        authState.logout()
      }
    }
    return Promise.reject(error)
  },
)

export function useApi() {
  return api
}
