import axios from 'axios';
import { useAuthStore } from '../stores/auth';

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || '/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
});

api.interceptors.request.use((config) => {
  const authStore = useAuthStore();
  if (authStore.token) {
    config.headers.Authorization = `Bearer ${authStore.token}`;
  }
  return config;
}, (error) => {
  return Promise.reject(error);
});

api.interceptors.response.use((response) => {
  if (response.data && response.data.success === true && response.data.data !== undefined) {
    response.data = response.data.data;
  }
  return response;
}, (error) => {
  if (error.response && error.response.status === 401) {
    const authStore = useAuthStore();
    authStore.logout();
  }
  return Promise.reject(error);
});

export function useApi() {
  return api;
}
